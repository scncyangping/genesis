package mgoRepo

import (
	"genesis/pkg/core/shunt/repository/entity"
	"genesis/pkg/plugin/mongo"
	"genesis/pkg/types"
	"go.uber.org/zap"
)

type UserMgoRepository struct {
	mgo    *mongo.MgoV
	logger *zap.SugaredLogger
}

func (u *UserMgoRepository) UpdateUser(user *entity.UserEn) (*entity.UserEn, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserMgoRepository) DeleteUserById(id string) (*entity.UserEn, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserRepository(m *mongo.MgoV, log *zap.SugaredLogger) *UserMgoRepository {
	return &UserMgoRepository{mgo: m, logger: log}
}

func (u *UserMgoRepository) SaveUser(entity *entity.UserEn) (string, error) {
	return u.mgo.InsertOne(entity)
}

//
//func (u *UserMgoRepository) UpdateUser(filter any, update any) error {
//	if _, err := u.mgo.Update(filter, update); err != nil {
//		u.logger.Error("UpdateUser Error, %v", err)
//		return err
//	}
//	return nil
//}

func (u *UserMgoRepository) DeleteUser(s string) error {
	_, err := u.mgo.DeleteOne(types.B{"_id": s})
	return err
}

func (u *UserMgoRepository) GetUserByName(name string) (*entity.UserEn, error) {
	var (
		err   error
		users entity.UserEn
	)
	err = u.mgo.FindOne(types.B{"name": name}, &users)
	if err != nil {
		u.logger.Errorf("GetUserByName Error, %v", err)
		return nil, err
	}
	return &users, nil
}

func (u *UserMgoRepository) GetUserById(s string) (*entity.UserEn, error) {
	var (
		err   error
		users entity.UserEn
	)
	if err = u.mgo.FindOne(types.B{"_id": s}, &users); err != nil {
		u.logger.Errorf("GetUserById Error, %v", err)
		return nil, err
	}
	return &users, nil
}

func (u *UserMgoRepository) ListUserBy(m any) ([]*entity.UserEn, error) {
	var (
		err   error
		users []*entity.UserEn
	)
	if err = u.mgo.Find(m, &users); err != nil {
		u.logger.Errorf("GetUserBy Error, %v", m)
		return nil, err
	}
	return users, nil
}

func (u *UserMgoRepository) ListUserPageBy(skip, limit int64, sort, filter any) ([]*entity.UserEn, error) {
	var (
		err   error
		users []*entity.UserEn
	)
	if err = u.mgo.FindBy(skip, limit, sort, filter, &users); err != nil {
		u.logger.Errorf("GetUserBy Error, %v", err)
		return nil, err
	}
	return users, nil
}
