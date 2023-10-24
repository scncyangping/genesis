// @Author: YangPing
// @Create: 2023/10/23
// @Description: 实体DTO

package dto

import "genesis/app/common/base"

type UserDto struct {
	Id string `json:"id"`
	base.User
}
