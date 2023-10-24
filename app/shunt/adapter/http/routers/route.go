// @Author: YangPing
// @Create: 2023/10/23
// @Description: 路由配置

package routers

import (
	"genesis/app/shunt/adapter/http/handlers/business"
	"genesis/app/shunt/adapter/http/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func InitRoute(router *gin.Engine, handlers *Handlers) {
	// 健康检查路由
	health(router)
	r := router.Group("/api/v1")
	withInitUserRoute(r, handlers.AuthHandler)
}

// health check
func health(router *gin.Engine) gin.IRoutes {
	return router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Application is running",
			"time":    time.Now().UTC(),
		})
	})
}

func withInitUserRoute(router *gin.RouterGroup, handler *business.AuthHandler) error {
	//router.Use(middleware.TokenAuthMiddleware())
	router = router.Group("/user")
	intUserCmdRoute(router, handler)
	intUserQueryRoute(router, handler)
	return nil
}

func intUserCmdRoute(router *gin.RouterGroup, handler *business.AuthHandler) {
	router.Use(middleware.RateLimitMiddleware(1*time.Second, 10, 10))
	router.POST("/", handler.AddUser)
	router.PUT("/", handler.UpdateUser)
	router.DELETE("/:id", handler.DeleteUserById)
	router.DELETE("/batch", handler.DeleteUserByIds)
}

func intUserQueryRoute(router *gin.RouterGroup, handler *business.AuthHandler) {
	router.Use(middleware.RateLimitMiddleware(1*time.Second, 10, 10))
	router.GET("/:id", handler.GetUserById)
	router.POST("/list", handler.QueryUser)
}
