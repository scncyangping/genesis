// @Author: YangPing
// @Create: 2023/10/23
// @Description: es插件配置

package es

type Config interface {
	GetAddress() []string
	GetUserName() string
	GetPassword() string
	GetVersion() int
}

type ClientForES interface {
	Index(index string, doc any) (string, error)
	Search(w *WithEsSearch) ([]byte, error)
}

type WithEsSearch struct {
	Index string
	Query map[string]any
	From  int
	Size  int
	Sort  string
}

func NewWithEsSearch(index string, query map[string]any) *WithEsSearch {
	return &WithEsSearch{
		Index: index,
		Query: query,
		From:  0,
		Size:  10,
	}
}

func (w *WithEsSearch) WithEsFrom(from int) *WithEsSearch {
	w.From = from
	return w
}

func (w *WithEsSearch) WithEsSize(size int) *WithEsSearch {
	w.Size = size
	return w
}

func (w *WithEsSearch) WithEsSort(sort string) *WithEsSearch {
	w.Sort = sort
	return w
}
