// @Author: YangPing
// @Create: 2023/10/21
// @Description:

package cmd

import (
	"genesis/app/shunt/adapter/http/handlers/business"
	"genesis/app/shunt/repository"
	"genesis/app/shunt/repository/mysqlRepo"
	"genesis/app/shunt/service/bus"
	"genesis/app/shunt/service/bus/impl"
	"github.com/google/wire"
)

var RepositorySlice = wire.NewSet(
	// user handler start
	mysqlRepo.NewUserMysqlRepo,
	impl.NewUserServiceImpl,
	business.NewAuthHandler,
	wire.Bind(new(repository.UserRepositoryI), new(*mysqlRepo.UserMysqlRepo)),
	wire.Bind(new(bus.UserSrv), new(*impl.UserServiceImpl)),
	// user handler end
)
