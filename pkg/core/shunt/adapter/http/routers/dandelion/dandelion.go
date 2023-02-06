package dandelion

import (
	"genesis/pkg/core/shunt/adapter/http/handlers/dandelion"

	"github.com/gin-gonic/gin"
)

func InitDandelionRoute(router *gin.RouterGroup, handler *dandelion.TemplateHandler) {
	router.POST("/generate", handler.TemplateGenerate)
}
