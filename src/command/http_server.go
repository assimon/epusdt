package command

import (
	"context"
	"github.com/assimon/luuu/config"
	"github.com/assimon/luuu/middleware"
	"github.com/assimon/luuu/route"
	"github.com/assimon/luuu/util/constant"
	luluHttp "github.com/assimon/luuu/util/http"
	"github.com/assimon/luuu/util/log"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "http服务",
	Long:  "http服务相关命令",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	httpCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动",
	Long:  "启动http服务",
	Run: func(cmd *cobra.Command, args []string) {
		HttpServerStart()
	},
}

func HttpServerStart() {
	var err error
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = customHTTPErrorHandler
	// 中间件注册
	MiddlewareRegister(e)
	// 路由注册
	route.RegisterRoute(e)
	// 静态目录注册
	e.Static(config.StaticPath, "static")
	httpListen := viper.GetString("http_listen")
	go func() {
		if err = e.Start(httpListen); err != http.ErrServerClosed {
			log.Sugar.Error(err)
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// MiddlewareRegister 中间件注册
func MiddlewareRegister(e *echo.Echo) {
	if config.AppDebug {
		e.Debug = true
		e.Use(echoMiddleware.Logger())
	}
	e.Use(middleware.RequestUUID())
}

// customHTTPErrorHandler 默认消息提示
func customHTTPErrorHandler(err error, e echo.Context) {
	code := http.StatusInternalServerError
	msg := "server error"
	resp := &luluHttp.Response{
		StatusCode: code,
		Message:    msg,
		RequestID:  e.Request().Header.Get(echo.HeaderXRequestID),
	}
	if he, ok := err.(*echo.HTTPError); ok {
		e.String(http.StatusOK, he.Message.(string))
		return
	}
	if he, ok := err.(*constant.RspError); ok {
		resp.StatusCode = he.Code
		resp.Message = he.Msg
	}
	_ = e.JSON(http.StatusOK, resp)
	return
}
