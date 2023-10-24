//go:generate go run github.com/google/wire/cmd/wire
//go:build wireinject
// +build wireinject

package cmd

import (
	"genesis/app/shunt/api/http/routers"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func NewHandler(db *gorm.DB) *routers.Handlers {
	panic(wire.Build(
		RepositorySlice,
		wire.Struct(new(routers.Handlers), "*"),
	))
	return &routers.Handlers{}
}
