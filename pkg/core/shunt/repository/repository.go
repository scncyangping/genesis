package repository

import (
	"genesis/pkg/config/app/shunt"
	"genesis/pkg/core/shunt/repository/entity"
	"genesis/pkg/core/shunt/repository/mysqlRepo"
)

type UserRepo interface {
	SaveUser(entity *entity.UserEn) (string, error)
	GetUserByName(name string) (*entity.UserEn, error)
	UpdateUser(user *entity.UserEn) (*entity.UserEn, error)
	DeleteUserById(id string) (*entity.UserEn, error)
}

type Repository struct {
	UserRepo UserRepo
}

// NewRepository wire
func NewRepository() *Repository {
	return &Repository{
		UserRepo: mysqlRepo.NewUserMysqlRepository(shunt.GormDB()),
	}
}
