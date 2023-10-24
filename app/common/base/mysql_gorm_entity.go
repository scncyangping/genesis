// @Author: YangPing
// @Create: 2023/10/21
// @Description: Mysql GORM 对应实体

package base

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

const (
	TableUser = "user"
)

type GormEnBase struct {
	Id        string    `json:"id,omitempty" gorm:"column:id;primaryKey"`
	CreatedBy string    `json:"createdBy,omitempty"  gorm:"column:created_by"`
	UpdatedBy string    `json:"updatedBy,omitempty"  gorm:"column:updated_by"`
	Remark    string    `json:"remark,omitempty"  gorm:"column:remark"`
	CreatedAt time.Time `json:"createdAt,omitempty"  gorm:"column:created_at"`
	//TODO 使用数据库默认时间
	//UpdatedAt time.Time             `json:"updatedAt,omitempty"  gorm:"column:updated_at;default:NULL"`
	UpdatedAt time.Time             `json:"updatedAt,omitempty"  gorm:"column:updated_at"`
	Deleted   soft_delete.DeletedAt `json:"deleted,omitempty"  gorm:"softDelete:flag" gorm:"column:deleted"`
}

// En 定义使用Gorm实体作为基础实现
type En struct {
	GormEnBase
}

type UserGorm struct {
	En
	User
}

func (u User) TableName() string {
	return TableUser
}
