//go:generate go run github.com/google/wire/cmd/wire
//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/google/wire"
	"genesis/pkg/core/shunt/adapter/http/handlers"
	"genesis/pkg/core/shunt/adapter/http/server"
	"genesis/pkg/core/shunt/application/service"
	"genesis/pkg/core/shunt/repository"
)

var providerSet = wire.NewSet(
	repository.NewRepository,
	service.NewAppSrvManager,
	handlers.NewHandler,
	server.NewHandlers,
)

func NewHandler() *server.Handlers {
	panic(wire.Build(providerSet))
}
