package service

import (
	"context"
	"github.com/assimon/luuu/model/data"
	"github.com/assimon/luuu/model/request"
	"github.com/assimon/luuu/util/http_client"
	"github.com/assimon/luuu/util/json"
	"github.com/assimon/luuu/util/log"
	"github.com/golang-module/carbon/v2"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/stdutil"
	"github.com/shopspring/decimal"
	"net/http"
	"sync"
	"time"
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
	ctx := context.Background()
	nowTime := time.Now().Unix()
	for _, transfer := range trc20Resp.TokenTransfers {
		if transfer.ToAddress != token || transfer.FinalResult != "SUCCESS" {
			continue
		}
		x, _ := decimal.NewFromString(transfer.Quant)
		y, _ := decimal.NewFromString("1000000")
		quant := x.Div(y).String()
		result, err := data.GetExpirationTimeByAmount(ctx, token, quant)
		if err != nil {
			panic(err)
		}
		if result != "" {
			expTime := mathutil.MustInt64(result)
			// 但是过期了
			if expTime < nowTime {
				// 删掉过期
				_ = data.ClearPayCache(token, quant)
				continue
			}
		}
		// 该钱包下有无匹配金额订单
		tradeId, err := data.GetTradeIdByAmount(ctx, token, quant)
		if err != nil {
			panic(err)
		}
		if tradeId == "" {
			continue
		}
		// 到这一步就匹配到金额了
		req := &request.OrderProcessingRequest{
			Token:              token,
			TradeId:            tradeId,
			Amount:             quant,
			BlockTransactionId: transfer.TransactionID,
		}
		err = OrderProcessing(req)
		if err != nil {
			panic(err)
		}
	}
}
