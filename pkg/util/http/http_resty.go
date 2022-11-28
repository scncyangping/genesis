package http

import (
	"github.com/go-resty/resty/v2"
	"net"
	"net/http"
	"runtime"
	"time"
)

var ShuntHttp *shuntHttp

type shuntHttp struct {
	client *resty.Client
}

func init() {
	client := &http.Client{
		Timeout: time.Duration(5) * time.Second, // 超时时间
		Transport: &http.Transport{
			MaxIdleConns:        100,                       // 最大空闲连接 0 为不限制
			MaxConnsPerHost:     100,                       // 单独host最大连接设置 0 为不限制
			MaxIdleConnsPerHost: runtime.GOMAXPROCS(0) + 1, // 单独host + port 最大空闲连接 0 为不限制
			Proxy:               http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	ShuntHttp = &shuntHttp{
		client: resty.NewWithClient(client),
	}
}

func (s *shuntHttp) R() *resty.Request {
	return s.client.R()
}
