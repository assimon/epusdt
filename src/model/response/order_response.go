package response

// CreateTransactionResponse 创建订单成功返回
type CreateTransactionResponse struct {
	TradeId        string  `json:"trade_id"`        //  epusdt订单号
	OrderId        string  `json:"order_id"`        //  客户交易id
	Amount         float64 `json:"amount"`          //  订单金额，保留4位小数
	ActualAmount   float64 `json:"actual_amount"`   //  订单实际需要支付的金额，保留4位小数
	Token          string  `json:"token"`           //  收款钱包地址
	ExpirationTime int64   `json:"expiration_time"` // 过期时间 时间戳
	PaymentUrl     string  `json:"payment_url"`     // 收银台地址
}

// OrderNotifyResponse 订单异步回调结构体
type OrderNotifyResponse struct {
	TradeId            string  `json:"trade_id"`             //  epusdt订单号
	OrderId            string  `json:"order_id"`             //  客户交易id
	Amount             float64 `json:"amount"`               //  订单金额，保留4位小数
	ActualAmount       float64 `json:"actual_amount"`        //  订单实际需要支付的金额，保留4位小数
	Token              string  `json:"token"`                //  收款钱包地址
	BlockTransactionId string  `json:"block_transaction_id"` // 区块id
	Signature          string  `json:"signature"`            // 签名
	Status             int     `json:"status"`               //  1：等待支付，2：支付成功，3：已过期
}
