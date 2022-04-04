package http_client

import (
	"github.com/go-resty/resty/v2"
	"time"
)

// GetHttpClient 获取请求客户端
func GetHttpClient(proxys ...string) *resty.Client {
	client := resty.New()
	// 如果有代理
	if len(proxys) > 0 {
		proxy := proxys[0]
		client.SetProxy(proxy)
	}
	client.SetTimeout(time.Second * 5)
	return client
}
