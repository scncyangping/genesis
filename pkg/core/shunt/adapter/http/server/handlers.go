package server

import (
	"genesis/pkg/core/shunt/adapter/http/handlers"
	"genesis/pkg/core/shunt/adapter/http/handlers/business"
	"genesis/pkg/core/shunt/adapter/http/handlers/dandelion"
	"genesis/pkg/core/shunt/application/service"
)

type Handlers struct {
	AuthHandler     *business.AuthHandler
	TemplateHandler *dandelion.TemplateHandler
}

// NewHandlers wire
func NewHandlers(handler *handlers.Handler, srvM *service.AppSrvManager) *Handlers {
	return &Handlers{
		AuthHandler: business.NewAuthHandler(handler, srvM.AuthSrv),
		TemplateHandler: dandelion.NewTemplateHandler(handler),
	}
}
