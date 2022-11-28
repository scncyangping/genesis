package identity

import (
	"genesis/pkg/core/shunt/adapter/http/handlers/business"
	"github.com/gin-gonic/gin"
	"time"
)

func InitAuthRoute(router *gin.RouterGroup, handler *business.AuthHandler) {
	authRouter := router.Group("")
	auth(authRouter, handler)

	userRouter := router.Group("/user")

	user(userRouter, handler)
}

func auth(router *gin.RouterGroup, handler *business.AuthHandler) {
	router.Use(handler.RateLimitMiddleware(1*time.Second, 1, 1))
	router.POST("/login", handler.Login)
	router.POST("/register", handler.Register)

}

func user(router *gin.RouterGroup, handler *business.AuthHandler) {
	router.Use(handler.TokenAuthMiddleware())
	router.Use(handler.RateLimitMiddleware(1*time.Second, 3, 1))
	router.GET("/who", func(c *gin.Context) {
		un, _ := c.Get("user")
		c.JSON(200, un)
	})
}
