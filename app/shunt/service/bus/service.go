// @Author: YangPing
// @Create: 2023/10/23
// @Description: 服务接口

package bus

import (
	"genesis/app/shunt/adapter/cqe/cmd"
	"genesis/app/shunt/adapter/cqe/query"
	"genesis/app/shunt/adapter/vo"
)

type UserSrv interface {
	Add(*cmd.UserSaveCmd) (string, error)
	Update(*cmd.UserUpdateCmd) error
	Query(*query.UserListQuery) (*vo.PageResult, error)
	GetById(string) (*vo.UserVO, error)
	DeleteById(string) error
	DeleteByMap(map[string]any) error
}
