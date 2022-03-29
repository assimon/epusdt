package controller

import (
	"errors"
	"github.com/assimon/luuu/util/http"
	"github.com/gookit/validate"
	"github.com/gookit/validate/locales/zhcn"
	"github.com/gookit/validate/locales/zhtw"
	"github.com/labstack/echo/v4"
	"sync"
)

var Ctrl = &BaseController{}

type validatorBase struct {
	once     sync.Once
	validate *validate.Validation
}

type BaseController struct {
	http.Resp
	Validator validatorBase
	Locale    string
}

func (c *BaseController) GetLocale(ctx echo.Context) string {
	c.Locale = ctx.Request().Header.Get("locale")
	return c.Locale
}

func (c *BaseController) RegisterGlobal(ctx echo.Context) {
	locale := c.GetLocale(ctx)
	switch locale {
	case "zh":
		zhcn.RegisterGlobal()
	case "zh-tw":
		zhtw.RegisterGlobal()
	default:
		zhcn.RegisterGlobal()
	}
}

func (c *BaseController) ValidateStruct(ctx echo.Context, i interface{}, scene ...string) error {
	c.RegisterGlobal(ctx)
	v := validate.Struct(i, scene...)
	if v.Validate() {
		return nil
	} else {
		return errors.New(v.Errors.One())
	}
}
