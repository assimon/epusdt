package request

import "github.com/gookit/validate"

// CreateTransactionRequest 创建交易请求
type CreateTransactionRequest struct {
	OrderId     string  `json:"order_id" validate:"required|maxLen:32"`
	Amount      float64 `json:"amount" validate:"required|isFloat|gt:0.01"`
	NotifyUrl   string  `json:"notify_url" validate:"required"`
	Signature   string  `json:"signature"  validate:"required"`
	RedirectUrl string  `json:"redirect_url"`
}

func (r CreateTransactionRequest) Translates() map[string]string {
	return validate.MS{
		"OrderId":   "订单号",
		"Amount":    "支付金额",
		"NotifyUrl": "异步回调网址",
		"Signature": "签名",
	}
}

// OrderProcessingRequest 订单处理
type OrderProcessingRequest struct {
	Token              string
	Amount             float64
	TradeId            string
	BlockTransactionId string
}
