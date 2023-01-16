package impl

import (
	"errors"
	"genesis/pkg/config/app/shunt"
	"genesis/pkg/core/shunt/application/cqe/cmd"
	"genesis/pkg/core/shunt/application/dto"
	"genesis/pkg/core/shunt/repository"
	"genesis/pkg/core/shunt/repository/entity"
	"genesis/pkg/util"
	"genesis/pkg/util/jwt"
	"genesis/pkg/util/snowflake"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthSrvImp struct {
	userRepo repository.UserRepo
}

func NewAuthSrvImp(userRepo repository.UserRepo) *AuthSrvImp {
	return &AuthSrvImp{
		userRepo: userRepo,
	}
}

func (a *AuthSrvImp) Login(re *cmd.LoginCmd) (dto.UserDto, error) {
	var userDto dto.UserDto

	if vo, err := a.userRepo.GetUserByName(re.Name); err != nil {
		return userDto, errors.New("name or password error")
	} else {
		if vo.Pwd == re.Pwd {
			err := util.StructCopy(&userDto, vo)
			if err != nil {
				return dto.UserDto{}, err
			}

			userDto.Token, err = jwt.GenerateToken(
				vo.Name,
				shunt.Config().Jwt.Issuer,
				shunt.Config().Jwt.Secret,
				shunt.Config().Jwt.ExpireTime)

			return userDto, err
		} else {
			return userDto, errors.New("name or password error")
		}
	}
}

func (a *AuthSrvImp) Register(re *cmd.RegisterCmd) (string, error) {
	var ue entity.UserEn

	if err := util.StructCopy(&ue, re); err != nil {
		return "", err
	}

	if v, err := a.userRepo.GetUserByName(re.Name); err != nil && err != mongo.ErrNoDocuments {
		return "", err
	} else {
		if v != nil && v.Id != "" {
			return "", errors.New("对应用户名已存在")
		}
	}

	ue.Id = snowflake.NextId()

	return a.userRepo.SaveUser(&ue)
}
