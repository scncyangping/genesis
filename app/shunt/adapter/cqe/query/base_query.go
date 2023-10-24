// @Author: YangPing
// @Create: 2023/10/19
// @Description: 分页查询参数基础配置

package query

const (
	_defaultOrderField = "created_at"
	_defaultPageSize   = 20
)

type SortByType string

const (
	Asc  = "ASC"
	Desc = "DESC"
)

type PageQuery struct {
	Page   int    `json:"page" binding:"required" msg:"页数不能为空"`                           // 页
	Size   int    `json:"size" binding:"required" msg:"每页数量不能为空"`                         // 每页数量
	Sort   string `json:"sort"`                                                           // 排序字段
	SortBy string `json:"sortBy" binding:"required,oneof=ASC DESC" msg:"排序方式仅能为ASC或DESC"` // 正序/倒序
}

type PageQBase struct {
	Skip   int
	Limit  int
	Sort   string `json:"sort"`   // 排序字段
	SortBy string `json:"sortBy"` // 正序/倒序
}

func (pq *PageQuery) Q() *PageQBase {
	skip := 0

	if pq.Size == 0 {
		pq.Size = _defaultPageSize
	}
	if pq.Page > 0 {
		skip = (pq.Page - 1) * pq.Size
	}
	if pq.Sort == "" {
		pq.Sort = _defaultOrderField
	}
	return &PageQBase{
		Skip:   skip,
		Limit:  pq.Size,
		Sort:   pq.Sort,
		SortBy: pq.SortBy,
	}
}
