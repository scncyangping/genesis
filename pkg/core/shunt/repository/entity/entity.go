package entity

import (
	"genesis/pkg/core/shunt"
	"gorm.io/plugin/soft_delete"
	"time"
)

const (
	DBTableUser = "user" // user table
)

type En struct {
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type UserEn struct {
	En       En                    `gorm:"embedded"`
	Id       string                `bson:"_id" gorm:"column:id;primaryKey"`
	Name     string                `json:"name,omitempty" gorm:"column:name"`
	NickName string                `json:"nick_name,omitempty" gorm:"column:nick_name"`
	Phone    string                `json:"phone,omitempty" gorm:"column:phone"`
	Email    string                `json:"email,omitempty" gorm:"column:email"`
	Pwd      string                `json:"pwd,omitempty" gorm:"column:pwd"`
	Status   shuntCore.UserStatus  `json:"status,omitempty" gorm:"column:status"`
	Deleted  soft_delete.DeletedAt `gorm:"softDelete:flag" gorm:"column:deleted"`
}

func (uen UserEn) TableName() string {
	return DBTableUser
}
