// @Author: YangPing
// @Create: 2023/10/21
// @Description: 用户查询参数

package query

import "genesis/app/common/base"

type UserListQuery struct {
	PageQuery `json:"pageQuery" structs:"-"`
	base.User `structs:",flatten"`
}
