// @Author: YangPing
// @Create: 2023/10/23
// @Description: 全局Repository定义

package repository

import (
	"genesis/app/common/base"
)

type UserRepositoryI interface {
	base.UniversalRepositoryI[base.UserGorm]
}
