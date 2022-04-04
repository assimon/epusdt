package comm

import (
	"github.com/assimon/luuu/model/request"
	"github.com/assimon/luuu/model/service"
	"github.com/assimon/luuu/util/constant"
	"github.com/labstack/echo/v4"
)

// CreateTransaction 创建交易
func (c *BaseCommController) CreateTransaction(ctx echo.Context) (err error) {
	req := new(request.CreateTransactionRequest)
	if err = ctx.Bind(req); err != nil {
		return c.FailJson(ctx, constant.ParamsMarshalErr)
	}
	if err = c.ValidateStruct(ctx, req); err != nil {
		return c.FailJson(ctx, err)
	}
	resp, err := service.CreateTransaction(req)
	if err != nil {
		return c.FailJson(ctx, err)
	}
	return c.SucJson(ctx, resp)
}
