package mdb

const (
	StatusWaitPay     = 1
	StatusPaySuccess  = 2
	StatusExpired     = 3
	CallBackConfirmOk = 1
	CallBackConfirmNo = 2
)

type Orders struct {
	TradeId            string  `gorm:"column:trade_id" json:"trade_id"`                         //  epusdt订单号
	OrderId            string  `gorm:"column:order_id" json:"order_id"`                         //  客户交易id
	BlockTransactionId string  `gorm:"column:block_transaction_id" json:"block_transaction_id"` // 区块id
	Amount             float64 `gorm:"column:amount" json:"amount"`                             //  订单金额，保留4位小数
	ActualAmount       float64 `gorm:"column:actual_amount" json:"actual_amount"`               //  订单实际需要支付的金额，保留4位小数
	Token              string  `gorm:"column:token" json:"token"`                               //  所属钱包地址
	Status             int     `gorm:"column:status" json:"status"`                             //  1：等待支付，2：支付成功，3：已过期
	NotifyUrl          string  `gorm:"column:notify_url" json:"notify_url"`                     //  异步回调地址
	RedirectUrl        string  `gorm:"column:redirect_url" json:"redirect_url"`                 //  同步回调地址
	CallbackNum        int     `gorm:"column:callback_num" json:"callback_num"`                 // 回调次数
	CallBackConfirm    int     `gorm:"column:callback_confirm" json:"callback_confirm"`         // 回调是否已确认 1是 2否
	BaseModel
}

// TableName sets the insert table name for this struct type
func (o *Orders) TableName() string {
	return "orders"
}
