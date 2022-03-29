package response

type CheckoutCounterResponse struct {
	TradeId        string  `json:"trade_id"`        //  epusdt订单号
	ActualAmount   float64 `json:"actual_amount"`   //  订单实际需要支付的金额，保留4位小数
	Token          string  `json:"token"`           //  收款钱包地址
	ExpirationTime int64   `json:"expiration_time"` // 过期时间 时间戳
	RedirectUrl    string  `json:"redirect_url"`
}

type CheckStatusResponse struct {
	TradeId string `json:"trade_id"` //  epusdt订单号
	Status  int    `json:"status"`
}
