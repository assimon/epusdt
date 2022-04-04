package service

import (
	"errors"
	"github.com/assimon/luuu/config"
	"github.com/assimon/luuu/model/data"
	"github.com/assimon/luuu/model/mdb"
	"github.com/assimon/luuu/model/response"
)

// GetCheckoutCounterByTradeId 获取收银台详情，通过订单
func GetCheckoutCounterByTradeId(tradeId string) (*response.CheckoutCounterResponse, error) {
	orderInfo, err := data.GetOrderInfoByTradeId(tradeId)
	if err != nil {
		return nil, err
	}
	if orderInfo.ID <= 0 || orderInfo.Status != mdb.StatusWaitPay {
		return nil, errors.New("不存在待支付订单或已过期！")
	}
	resp := &response.CheckoutCounterResponse{
		TradeId:        orderInfo.TradeId,
		ActualAmount:   orderInfo.ActualAmount,
		Token:          orderInfo.Token,
		ExpirationTime: orderInfo.CreatedAt.AddMinutes(config.GetOrderExpirationTime()).TimestampWithMillisecond(),
		RedirectUrl:    orderInfo.RedirectUrl,
	}
	return resp, nil
}
