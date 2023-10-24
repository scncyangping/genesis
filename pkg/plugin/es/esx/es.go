// @Author: YangPing
// @Create: 2023/10/23
// @Description: ESX

package esx

import (
	"crypto/tls"
	"fmt"
	es2 "genesis/pkg/plugin/es"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

func NewESClient(config es2.Config) (es2.ClientForES, error) {
	client := &http.Client{
		Timeout: time.Duration(5) * time.Second, // 超时时间
		Transport: &http.Transport{
			MaxConnsPerHost:       10,                        // 单独host最大连接设置 0 为不限制
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1, // 单独host + port 最大空闲连接 0 为不限制
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS12,
				InsecureSkipVerify: true, //跳过HTTPS验证
			},
		},
	}
	return &ES{
		client:  resty.NewWithClient(client),
		baseUrl: config.GetAddress()[0],
	}, nil
}

type ES struct {
	cluster  []string
	lock     sync.RWMutex
	baseUrl  string
	client   *resty.Client
	username string
	password string
}

func (e *ES) Index(index string, doc any) (string, error) {
	url := fmt.Sprintf("%s/%s/%s", e.baseUrl, index, index)
	resp, err := e.client.R().
		SetBasicAuth(e.username, e.password).
		SetHeader("Accept", "application/json").
		SetBody(doc).
		Post(url)
	if err != nil {
		return "", err
	}
	// 检查 HTTP 响应状态码
	if resp.StatusCode() != 201 {
		return "", errors.New(string(resp.Body()))
	}
	return "", nil
}

func (e *ES) Search(w *es2.WithEsSearch) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/_search", e.baseUrl, w.Index)
	qp := make(map[string]any)
	qp = w.Query
	if len(qp) == 1 {
		qp = map[string]any{
			"from": w.From,
			"size": w.Size,
		}
		if w.Sort != "" {
			sp := strings.Split(w.Sort, ":")
			qp["sort"] = []map[string]any{
				{
					sp[0]: sp[1],
				},
			}
		}
	}
	var res []byte
	resp, err := e.client.R().
		SetBasicAuth(e.username, e.password).
		SetHeader("Accept", "application/json").
		SetBody(qp).
		Post(url)
	if err != nil {
		return res, err
	} else {
		res = resp.Body()
	}
	return res, nil
}
