package request

const (
	OrderByFuncDesc = "DESC"
	OrderByFuncAsc  = "OrderByFuncASC"
)

var OrderByFuncList = []string{OrderByFuncDesc, OrderByFuncAsc}

type BaseRequest struct {
	Page       int    `json:"page"`        // 页数
	PageSize   int    `json:"page_size"`   // 每页条数
	OrderField string `json:"order_field"` // 排序字段
	OrderFunc  string `json:"order_func"`  // 排序方法
}
