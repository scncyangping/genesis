// @Author: YangPing
// @Create: 2023/10/23
// @Description: User存储操作

package mysqlRepo

import (
	"genesis/app/common/base"
	"gorm.io/gorm"
)

type UserMysqlRepo struct {
	*base.UniversalGormRepo[base.UserGorm]
}

// NewUserMysqlRepo
// Need Wire
func NewUserMysqlRepo(db *gorm.DB) *UserMysqlRepo {
	return &UserMysqlRepo{
		UniversalGormRepo: base.NewUniversalGormRepo(base.UserGorm{}, db)}
}
