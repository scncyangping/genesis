// @Author: YangPing
// @Create: 2023/10/23
// @Description: ES

package esv8

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	es2 "genesis/pkg/plugin/es"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

func NewESClient(config es2.Config) (es2.ClientForES, error) {
	cfg := elasticsearch.Config{
		Addresses: config.GetAddress(),
		Username:  config.GetUserName(),
		Password:  config.GetPassword(),
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS12,
				InsecureSkipVerify: true, //跳过HTTPS验证
			},
		},
	}
	if es, err := elasticsearch.NewClient(cfg); err != nil {
		return nil, err
	} else {
		return &ES{Client: es}, nil
	}
}

type ES struct {
	Client *elasticsearch.Client
}

func (e *ES) Index(index string, doc any) (string, error) {
	var buf bytes.Buffer
	//doc := map[string]interface{}{
	//	"title":   "中国",
	//	"content": "中国早日统一台湾",
	//	"time":    time.Now().Unix(),
	//	"date":    time.Now(),
	//}
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		return "", err
	}
	res, err := e.Client.Index(
		index, // Index name
		&buf,  // Document body
		//e.Client.Index.WithDocumentID(idx), // Document ID
		// Document ID
		//e.Client.Index.WithRefresh("true"), // Refresh
	)
	defer res.Body.Close()

	if err != nil {
		return "", err
	}
	return res.String(), nil
}

// IndexESApi 类型允许使用更实际的方法，您可以在其中创建一个新结构，将请求配置作为字段，并使用上下文和客户端作为参数调用 Do() 方法：
func (e *ES) IndexESApi(index, idx string, doc map[string]any) {
	//index:="my_index_name_v1"
	res, err := e.Client.Info()
	fmt.Println(res, err)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	var buf bytes.Buffer
	//doc := map[string]interface{}{
	//	"title":   "中国",
	//	"content": "中国早日统一台湾",
	//	"time":    time.Now().Unix(),
	//	"date":    time.Now(),
	//}
	if err = json.NewEncoder(&buf).Encode(doc); err != nil {
		fmt.Println(err, "Error encoding doc")
		return
	}

	req := esapi.IndexRequest{
		Index:      index, // Index name
		Body:       &buf,  // Document body
		DocumentID: idx,   // Document ID
		//Refresh:    "true", // Refresh
	}

	res, err = req.Do(context.Background(), e.Client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	fmt.Println(res.String())
	log.Println(res)
}

func (e *ES) Search(w *es2.WithEsSearch) ([]byte, error) {
	res, err := e.Client.Info()
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(w.Query); err != nil {
		return nil, err
	}

	options := []func(request *esapi.SearchRequest){
		e.Client.Search.WithContext(context.Background()),
		e.Client.Search.WithIndex(w.Index),
		e.Client.Search.WithBody(&buf),
		e.Client.Search.WithFrom(w.From),
		e.Client.Search.WithSize(w.Size),
		//e.Client.Search.WithSort(w.sort), //time:desc
		//e.Client.Search.WithTrackTotalHits(true),
		//e.Client.Search.WithPretty(),
	}

	if w.Sort != "" {
		options = append(options, e.Client.Search.WithSort(w.Sort))
	}

	// Perform the search request.
	res, err = e.Client.Search(options...)

	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}

// Delete 删除 index 根据 索引名 id
func (e *ES) Delete(index, idx string) (string, error) {
	//index:="my_index_name_v1"
	res, err := e.Client.Info()
	if err != nil {
		return res.String(), err
	}
	res, err = e.Client.Delete(
		index, // Index name
		idx,   // Document ID
		//e.Client.Delete.WithRefresh("true"),
	)
	defer res.Body.Close()

	if err != nil {
		return "", err
	}
	return res.String(), nil
}

func (e *ES) DeleteByQuery(index []string, query map[string]any) error {
	res, err := e.Client.Info()
	if err != nil {
		return err
	}
	//fmt.Println(res.String())
	// search - highlight
	var buf bytes.Buffer
	//query := map[string]interface{}{
	//	"query": map[string]interface{}{
	//		"match": map[string]interface{}{
	//			"title": title,
	//		},
	//	},
	//	},
	//}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return err
	}
	// Perform the search request.
	res, err = e.Client.DeleteByQuery(
		index,
		&buf,
	)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

func (e *ES) SearchEsapiSql(query map[string]any) (string, error) {
	jsonBody, _ := json.Marshal(query)
	req := esapi.SQLQueryRequest{
		Body: bytes.NewReader(jsonBody),
	}
	res, err := req.Do(context.Background(), e.Client)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	return res.String(), nil
}

func (e *ES) SearchHttp(method, url string, query map[string]interface{}) (string, error) {
	jsonBody, _ := json.Marshal(query)
	req, err := http.NewRequest(method, url, bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-type", "application/json")

	res, err := e.Client.Perform(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	return buf.String(), nil
}
