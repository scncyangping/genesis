package routers

import (
	"genesis/pkg/core/shunt/adapter/http/routers/identity"
	"genesis/pkg/core/shunt/adapter/http/server"

	//_ "genesis/pkg/core/shunt/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoute(router *gin.Engine, h *server.Handlers) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1Auth := router.Group("/api/v1/auth")
	{
		identity.InitAuthRoute(v1Auth, h.AuthHandler)
	}
}
