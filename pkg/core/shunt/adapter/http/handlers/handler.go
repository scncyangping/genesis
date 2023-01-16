package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"go.uber.org/zap"
	"genesis/pkg/config/app/shunt"
	"genesis/pkg/util/jwt"
	"time"
)

// Handler 具体业务服务聚合
type Handler struct {
	Logger *zap.SugaredLogger
}

// NewHandler wire
func NewHandler() *Handler {
	return &Handler{
		Logger: shunt.Log().Named("http-handler"),
	}
}

func (h *Handler) TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token   string
			errFlag ResultCode
		)
		token = c.Request.FormValue(shunt.Config().Jwt.AuthKey)
		if token == "" {
			token = c.GetHeader(shunt.Config().Jwt.AuthKey)
		}
		if token == "" {
			c.Abort()
			h.SendFailure(c, RequestTokenNotFound, StatusText(RequestTokenNotFound))
			return

		}
		if claims, err := jwt.ParseToken(token, shunt.Config().Jwt.Secret); err != nil {
			fmt.Println(err)
			errFlag = RequestCheckTokenError
		} else if time.Now().Unix() > claims.Base.ExpiresAt {
			errFlag = RequestCheckTokenTimeOut
		} else {
			// 设置登录信息到token里面
			c.Set("user", claims.Extra)
		}
		if errFlag > 0 {
			c.Abort()
			h.SendFailure(c, errFlag, StatusText(errFlag))
			return
		}
		c.Next()
	}
}

func (h *Handler) RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.Abort()
			h.SendFailure(c, RateLimit, StatusText(RateLimit))
			return
		}
		c.Next()
	}
}
