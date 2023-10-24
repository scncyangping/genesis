// @Author: YangPing
// @Create: 2023/10/23
// @Description: ES

package eso

import (
	"context"
	"encoding/json"
	es2 "genesis/pkg/plugin/es"
	"github.com/olivere/elastic"
	"strings"
)

type ES struct {
	client  *elastic.Client
	useType bool
}

func NewESClient(config es2.Config) (*ES, error) {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(config.GetAddress()...),
		elastic.SetBasicAuth(config.GetUserName(), config.GetPassword()),
	)
	if err != nil {
		return nil, err
	}
	return &ES{
		client:  client,
		useType: config.GetVersion() > 6,
	}, nil
}

func (e *ES) Index(index string, doc interface{}) (string, error) {
	client := e.client.Index().
		Index(index).
		BodyJson(doc)
	if !e.useType {
		client.Type(index)
	}
	_, err := client.Do(context.Background())
	if err != nil {
		return "", err
	}
	return "", nil
}

func (e *ES) Search(w *es2.WithEsSearch) ([]byte, error) {
	w.Query["from"] = w.From
	w.Query["size"] = w.Size

	client := e.client.Search().
		Index(w.Index).
		Source(w.Query)
	//.
	//	From(w.From).
	//	Size(w.Size)

	if w.Sort != "" {
		sp := strings.Split(w.Sort, ":")
		if len(sp) == 2 {
			if sp[1] == "desc" {
				w.Query["sort"] = map[string]any{
					sp[0]: map[string]any{
						"order": "desc",
					},
				}
			} else {
				client.Sort(sp[0], true)
				w.Query["sort"] = map[string]any{
					sp[0]: map[string]any{
						"order": "asc",
					},
				}
			}
		}
	}

	if searchResult, err := client.Do(context.Background()); err != nil {
		return nil, err
	} else {
		if res, err := json.Marshal(searchResult); err != nil {
			return nil, err
		} else {
			return res, nil
		}
	}
}
