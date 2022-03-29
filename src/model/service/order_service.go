package service

import (
	"context"
	"fmt"
	"github.com/assimon/luuu/config"
	"github.com/assimon/luuu/model/dao"
	"github.com/assimon/luuu/model/data"
	"github.com/assimon/luuu/model/mdb"
	"github.com/assimon/luuu/model/request"
	"github.com/assimon/luuu/model/response"
	"github.com/assimon/luuu/mq"
	"github.com/assimon/luuu/mq/handle"
	"github.com/assimon/luuu/telegram"
	"github.com/assimon/luuu/util/constant"
	"github.com/golang-module/carbon/v2"
	"github.com/gookit/goutil/mathutil"
	"github.com/hibiken/asynq"
	"github.com/shopspring/decimal"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var gCreateTransactionLock sync.Mutex

// CreateTransaction 创建订单
func CreateTransaction(req *request.CreateTransactionRequest) (*response.CreateTransactionResponse, error) {
	gCreateTransactionLock.Lock()
	defer gCreateTransactionLock.Unlock()
	// 汇率计算金额
	rmb := decimal.NewFromFloat(req.Amount)
	rate := decimal.NewFromFloat(config.GetUsdtRate())
	amount := rmb.Div(rate).InexactFloat64()
	actualAmountStr := fmt.Sprintf("%.4f", amount)
	actualAmountFloat, err := strconv.ParseFloat(actualAmountStr, 64)
	if err != nil {
		return nil, err
	}
	// 是否可以满足最低支付金额
	if actualAmountFloat <= 0 {
		return nil, constant.PayAmountErr
	}
	// 已经存在了的交易
	exist, err := data.GetOrderInfoByOrderId(req.OrderId)
	if err != nil {
		return nil, err
	}
	if exist.ID > 0 {
		return nil, constant.OrderAlreadyExists
	}
	// 有无可用钱包
	walletAddress, err := data.GetAvailableWalletAddress()
	if err != nil {
		return nil, err
	}
	if len(walletAddress) <= 0 {
		return nil, constant.NotAvailableWalletAddress
	}
	availableToken, availableAmountStr, err := CalculateAvailableWalletTokenAndAmount(actualAmountStr, walletAddress)
	if err != nil {
		return nil, err
	}
	if availableToken == "" || availableAmountStr == "" {
		return nil, constant.NotAvailableAmountErr
	}
	availableAmountFloat, err := strconv.ParseFloat(availableAmountStr, 64)
	if err != nil {
		return nil, err
	}
	tx := dao.Mdb.Begin()
	order := &mdb.Orders{
		TradeId:      GenerateCode(),
		OrderId:      req.OrderId,
		Amount:       req.Amount,
		ActualAmount: availableAmountFloat,
		Token:        availableToken,
		Status:       mdb.StatusWaitPay,
		NotifyUrl:    req.NotifyUrl,
		RedirectUrl:  req.RedirectUrl,
	}
	err = data.CreateOrderWithTransaction(tx, order)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	ExpirationTime := carbon.Now().AddMinutes(config.GetOrderExpirationTime()).Timestamp()
	// 锁定支付池
	err = data.LockPayCache(availableToken, order.TradeId, availableAmountStr, ExpirationTime)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	// 超时过期消息队列
	orderExpirationQueue, _ := handle.NewOrderExpirationQueue(order.TradeId)
	mq.MClient.Enqueue(orderExpirationQueue, asynq.ProcessIn(time.Minute*time.Duration(config.GetOrderExpirationTime())))
	resp := &response.CreateTransactionResponse{
		TradeId:        order.TradeId,
		OrderId:        order.OrderId,
		Amount:         order.Amount,
		ActualAmount:   order.ActualAmount,
		Token:          order.Token,
		ExpirationTime: ExpirationTime,
		PaymentUrl:     fmt.Sprintf("%s/pay/checkout-counter/%s", config.GetAppUri(), order.TradeId),
	}
	return resp, nil
}

// OrderProcessing 成功处理订单
func OrderProcessing(req *request.OrderProcessingRequest) error {
	tx := dao.Mdb.Begin()
	exist, err := data.GetOrderByBlockIdWithTransaction(tx, req.BlockTransactionId)
	if err != nil {
		return err
	}
	if exist.ID > 0 {
		tx.Rollback()
		return constant.OrderBlockAlreadyProcess
	}
	// 标记订单成功
	err = data.OrderSuccessWithTransaction(tx, req)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = data.ClearPayCache(req.Token, req.Amount)
	tx.Commit()
	order, err := data.GetOrderInfoByTradeId(req.TradeId)
	if err != nil {
		return err
	}
	// 回调队列
	orderCallbackQueue, _ := handle.NewOrderCallbackQueue(order)
	mq.MClient.Enqueue(orderCallbackQueue, asynq.MaxRetry(5))
	// 发送机器人消息
	msgTpl := `
<b>📢📢有新的交易支付成功！</b>
<pre>交易号：%s</pre>
<pre>订单号：%s</pre>
<pre>请求支付金额：%.4f cny</pre>
<pre>实际支付金额：%.4f usdt</pre>
<pre>钱包地址：%s</pre>
<pre>订单创建时间：%s</pre>
<pre>支付成功时间：%s</pre>
`
	msg := fmt.Sprintf(msgTpl, order.TradeId, order.OrderId, order.Amount, order.ActualAmount, order.Token, order.CreatedAt.ToDateTimeString(), carbon.Now().ToDateTimeString())
	telegram.SendToBot(msg)
	return nil
}

func CalculateAvailableWalletTokenAndAmount(amount string, walletAddress []mdb.WalletAddress) (string, string, error) {
	calculateAmountStr := amount
	availableAmountStr := ""
	availableToken := ""
	for i := 0; i < 100; i++ {
		token, err := CalculateAvailableWalletToken(calculateAmountStr, walletAddress)
		if err != nil {
			return "", "", err
		}
		// 这个金额没有拿到可用的钱包，重试，金额+0.0001
		if token == "" {
			x, err := decimal.NewFromString(calculateAmountStr)
			if err != nil {
				return "", "", err
			}
			y, err := decimal.NewFromString("0.0001")
			if err != nil {
				return "", "", err
			}
			calculateAmountStr = x.Add(y).String()
			continue
		}
		availableAmountStr = calculateAmountStr
		availableToken = token
		break
	}
	return availableToken, availableAmountStr, nil
}

// CalculateAvailableWalletToken 计算可用钱包token
func CalculateAvailableWalletToken(payAmount string, walletAddress []mdb.WalletAddress) (string, error) {
	nowTime := time.Now().Unix()
	ctx := context.Background()
	walletToken := ""
	for _, address := range walletAddress {
		result, err := data.GetExpirationTimeByAmount(ctx, address.Token, payAmount)
		if err != nil {
			return "", err
		}
		// 这个钱包金额被占用了
		if result != "" {
			endTime := mathutil.MustInt64(result)
			// 但是过期了
			if endTime < nowTime {
				// 删掉过期，还能继续用这个地址
				err = data.ClearPayCache(address.Token, payAmount)
				if err != nil {
					return "", err
				}
			} else {
				continue
			}
		}
		walletToken = address.Token
		break
	}
	return walletToken, nil
}

// GenerateCode 订单号生成
func GenerateCode() string {
	date := time.Now().Format("20060102")
	r := rand.Intn(1000)
	code := fmt.Sprintf("%s%d%03d", date, time.Now().UnixNano()/1e6, r)
	return code
}

// GetOrderInfoByTradeId 通过交易号获取订单
func GetOrderInfoByTradeId(tradeId string) (*mdb.Orders, error) {
	order, err := data.GetOrderInfoByTradeId(tradeId)
	if err != nil {
		return nil, err
	}
	if order.ID <= 0 {
		return nil, constant.OrderNotExists
	}
	return order, nil
}
