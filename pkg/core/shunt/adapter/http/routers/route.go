package routers

import (
	"genesis/pkg/core/shunt/adapter/http/routers/dandelion"
	"genesis/pkg/core/shunt/adapter/http/routers/identity"
	"genesis/pkg/core/shunt/adapter/http/server"

	"github.com/gin-gonic/gin"
)

func InitRoute(router *gin.Engine, h *server.Handlers) {
	v1Auth := router.Group("/api/v1/auth")
	{
		identity.InitAuthRoute(v1Auth, h.AuthHandler)
	}

	v1Dandelion := router.Group("/api/v1/dandelion")
	{
		dandelion.InitDandelionRoute(v1Dandelion, h.TemplateHandler)
	}
}
