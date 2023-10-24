// @Author: YangPing
// @Create: 2023/10/23
// @Description: 用户VO定义

package vo

import "genesis/app/common/base"

type UserVO struct {
	Id string `json:"id"`
	base.User
}
