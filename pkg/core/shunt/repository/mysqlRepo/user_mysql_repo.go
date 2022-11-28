package mysqlRepo

import (
	"genesis/pkg/core/shunt/repository/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserMysqlRepository struct {
	gm *gorm.DB
}

func NewUserMysqlRepository(gm *gorm.DB) *UserMysqlRepository {
	return &UserMysqlRepository{gm: gm}
}

func (u *UserMysqlRepository) SaveUser(entity *entity.UserEn) (string, error) {
	if err := u.gm.Create(entity).Error; err != nil {
		return "", errors.Wrapf(err, "save user error")
	}
	return entity.Id, nil
}

func (u *UserMysqlRepository) GetUserByName(name string) (*entity.UserEn, error) {
	var en entity.UserEn

	if err := u.gm.Find(&en, "name=?", name).Error; err != nil {
		return nil, errors.Wrapf(err, "query user error")
	} else {
		return &en, nil
	}
}

func (u *UserMysqlRepository) UpdateUser(user *entity.UserEn) (*entity.UserEn, error) {

	if err := u.gm.Updates(&user).Error; err != nil {
		return nil, errors.Wrapf(err, "query user error")
	} else {
		return user, nil
	}
}

func (u *UserMysqlRepository) DeleteUserById(id string) (*entity.UserEn, error) {
	var en entity.UserEn

	if err := u.gm.Delete(&en, "id=?", id).Error; err != nil {
		return nil, errors.Wrapf(err, "delete user by id error")
	} else {
		return &en, nil
	}
}
