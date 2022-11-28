package service

import (
	"genesis/pkg/core/shunt/application/cqe/cmd"
	"genesis/pkg/core/shunt/application/dto"
	"genesis/pkg/core/shunt/application/service/impl"
	"genesis/pkg/core/shunt/repository"
)

// AuthSrv 接口
type AuthSrv interface {
	Login(*cmd.LoginCmd) (dto.UserDto, error)
	Register(*cmd.RegisterCmd) (string, error)
}

type AppSrvManager struct {
	AuthSrv AuthSrv
}

// NewAppSrvManager wire
func NewAppSrvManager(repository *repository.Repository) *AppSrvManager {
	return &AppSrvManager{
		AuthSrv: impl.NewAuthSrvImp(repository.UserRepo),
	}
}
