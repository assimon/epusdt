package mdb

const (
	TokenStatusEnable  = 1
	TokenStatusDisable = 2
)

// WalletAddress  钱包表
type WalletAddress struct {
	Token  string `gorm:"column:token" json:"token"`   //  钱包token
	Status int64  `gorm:"column:status" json:"status"` //  1:启用 2:禁用
	BaseModel
}

// TableName sets the insert table name for this struct type
func (w *WalletAddress) TableName() string {
	return "wallet_address"
}
