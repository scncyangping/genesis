// @Author: YangPing
// @Create: 2023/10/23
// @Description: 基础VO定义

package vo

type PageResult struct {
	List       any   `json:"list"`
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	TotalCount int64 `json:"totalCount"`
	TotalPage  int64 `json:"totalPage"`
}

func NewPageResult(list any, page, size int, count int64) *PageResult {
	return &PageResult{
		List:       list,
		Page:       page,
		Size:       size,
		TotalCount: count,
		TotalPage:  (count + int64(size) - 1) / int64(size),
	}
}
