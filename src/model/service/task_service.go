package service

import (
	"fmt"
	"github.com/assimon/luuu/model/data"
	"github.com/assimon/luuu/model/request"
	"github.com/assimon/luuu/mq"
	"github.com/assimon/luuu/mq/handle"
	"github.com/assimon/luuu/telegram"
	"github.com/assimon/luuu/util/http_client"
	"github.com/assimon/luuu/util/json"
	"github.com/assimon/luuu/util/log"
	"github.com/golang-module/carbon/v2"
	"github.com/gookit/goutil/stdutil"
	"github.com/hibiken/asynq"
	"github.com/shopspring/decimal"
	"net/http"
	"sync"
)

const UsdtTrc20ApiUri = "https://apilist.tronscan.org/api/token_trc20/transfers"

type UsdtTrc20Resp struct {
	Total          int              `json:"total"`
	RangeTotal     int              `json:"rangeTotal"`
	TokenTransfers []TokenTransfers `json:"token_transfers"`
}

type TokenInfo struct {
	TokenID      string `json:"tokenId"`
	TokenAbbr    string `json:"tokenAbbr"`
	TokenName    string `json:"tokenName"`
	TokenDecimal int    `json:"tokenDecimal"`
	TokenCanShow int    `json:"tokenCanShow"`
	TokenType    string `json:"tokenType"`
	TokenLogo    string `json:"tokenLogo"`
	TokenLevel   string `json:"tokenLevel"`
	Vip          bool   `json:"vip"`
}
type TokenTransfers struct {
	TransactionID         string    `json:"transaction_id"`
	BlockTs               int64     `json:"block_ts"`
	FromAddress           string    `json:"from_address"`
	ToAddress             string    `json:"to_address"`
	Block                 int       `json:"block"`
	ContractAddress       string    `json:"contract_address"`
	Quant                 string    `json:"quant"`
	ApprovalAmount        string    `json:"approval_amount"`
	EventType             string    `json:"event_type"`
	ContractType          string    `json:"contract_type"`
	Confirmed             bool      `json:"confirmed"`
	ContractRet           string    `json:"contractRet"`
	FinalResult           string    `json:"finalResult"`
	TokenInfo             TokenInfo `json:"tokenInfo"`
	FromAddressIsContract bool      `json:"fromAddressIsContract"`
	ToAddressIsContract   bool      `json:"toAddressIsContract"`
	Revert                bool      `json:"revert"`
}

// Trc20CallBack trc20回调
func Trc20CallBack(token string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			log.Sugar.Error(err)
		}
	}()
	client := http_client.GetHttpClient()
	startTime := carbon.Now().AddHours(-24).TimestampWithMillisecond()
	endTime := carbon.Now().TimestampWithMillisecond()
	resp, err := client.R().SetQueryParams(map[string]string{
		"limit":           "200",
		"start":           "0",
		"direction":       "in",
		"tokens":          "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
		"relatedAddress":  token,
		"start_timestamp": stdutil.ToString(startTime),
		"end_timestamp":   stdutil.ToString(endTime),
	}).Get(UsdtTrc20ApiUri)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode() != http.StatusOK {
		panic(err)
	}
	var trc20Resp UsdtTrc20Resp
	err = json.Cjson.Unmarshal(resp.Body(), &trc20Resp)
	if err != nil {
		panic(err)
	}
	if trc20Resp.Total <= 0 {
		return
	}
	for _, transfer := range trc20Resp.TokenTransfers {
		if transfer.ToAddress != token || transfer.FinalResult != "SUCCESS" {
			continue
		}
		decimalQuant, err := decimal.NewFromString(transfer.Quant)
		if err != nil {
			panic(err)
		}
		decimalDivisor := decimal.NewFromFloat(1000000)
		amount := decimalQuant.Div(decimalDivisor).InexactFloat64()
		tradeId, err := data.GetTradeIdByWalletAddressAndAmount(token, amount)
		if err != nil {
			panic(err)
		}
		if tradeId == "" {
			continue
		}
		order, err := data.GetOrderInfoByTradeId(tradeId)
		if err != nil {
			panic(err)
		}
		// 区块的确认时间必须在订单创建时间之后
		createTime := order.CreatedAt.TimestampWithMillisecond()
		if transfer.BlockTs < createTime {
			panic("Orders cannot actually be matched")
		}
		// 到这一步就完全算是支付成功了
		req := &request.OrderProcessingRequest{
			Token:              token,
			TradeId:            tradeId,
			Amount:             amount,
			BlockTransactionId: transfer.TransactionID,
		}
		err = OrderProcessing(req)
		if err != nil {
			panic(err)
		}
		// 回调队列
		orderCallbackQueue, _ := handle.NewOrderCallbackQueue(order)
		mq.MClient.Enqueue(orderCallbackQueue, asynq.MaxRetry(5))
		// 发送机器人消息
		msgTpl := `
<b>📢📢有新的交易支付成功！</b>
<pre>交易号：%s</pre>
<pre>订单号：%s</pre>
<pre>请求支付金额：%f cny</pre>
<pre>实际支付金额：%f usdt</pre>
<pre>钱包地址：%s</pre>
<pre>订单创建时间：%s</pre>
<pre>支付成功时间：%s</pre>
`
		msg := fmt.Sprintf(msgTpl, order.TradeId, order.OrderId, order.Amount, order.ActualAmount, order.Token, order.CreatedAt.ToDateTimeString(), carbon.Now().ToDateTimeString())
		telegram.SendToBot(msg)
	}
}
