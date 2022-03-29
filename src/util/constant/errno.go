package constant

var Errno = map[int]string{
	400:   "系统错误",
	401:   "签名认证错误",
	10001: "钱包地址已存在，请勿重复添加",
	10002: "支付交易已存在，请勿重复创建",
	10003: "无可用钱包地址，无法发起支付",
	10004: "支付金额有误, 无法满足最小支付单位",
	10005: "无可用金额通道",
	10006: "汇率计算错误",
	10007: "订单区块已处理",
	10008: "订单不存在",
	10009: "无法解析请求参数",
}

var (
	SystemErr                  = Err(400)
	SignatureErr               = Err(401)
	WalletAddressAlreadyExists = Err(10001)
	OrderAlreadyExists         = Err(10002)
	NotAvailableWalletAddress  = Err(10003)
	PayAmountErr               = Err(10004)
	NotAvailableAmountErr      = Err(10005)
	RateAmountErr              = Err(10006)
	OrderBlockAlreadyProcess   = Err(10007)
	OrderNotExists             = Err(10008)
	ParamsMarshalErr           = Err(10009)
)

type RspError struct {
	Code int
	Msg  string
}

func (re *RspError) Error() string {
	return re.Msg
}

func Err(code int) (err error) {
	err = &RspError{
		Code: code,
		Msg:  Errno[code],
	}
	return err
}

func (re *RspError) Render() (code int, msg string) {
	return re.Code, re.Msg
}
