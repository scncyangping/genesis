// @Author: YangPing
// @Create: 2023/10/21
// @Description: 中间件配置

package middleware

import (
	"fmt"
	"genesis/app/shunt/adapter/http/handlers"
	"genesis/app/shunt/config"
	"genesis/pkg/util/jwt"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token   string
			errFlag handlers.ResultCode
		)
		token = c.Request.FormValue(config.Rt().Config().GetAuthKey())
		if token == "" {
			token = c.GetHeader(config.Rt().Config().GetAuthKey())
		}
		if token == "" {
			c.Abort()
			handlers.SendFailure(c, handlers.RequestTokenNotFound, handlers.StatusText(handlers.RequestTokenNotFound))
			return

		}
		if claims, err := jwt.ParseToken(token, config.Rt().Config().GetSecret()); err != nil {
			fmt.Println(err)
			errFlag = handlers.RequestCheckTokenError
		} else if time.Now().Unix() > claims.Base.ExpiresAt {
			errFlag = handlers.RequestCheckTokenTimeOut
		} else {
			// 设置登录信息到token里面
			c.Set("user", claims.Extra)
		}
		if errFlag > 0 {
			c.Abort()
			handlers.SendFailure(c, errFlag, handlers.StatusText(errFlag))
			return
		}
		c.Next()

	}
}

func RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.Abort()
			handlers.SendFailure(c, handlers.RateLimit, handlers.StatusText(handlers.RateLimit))
			return
		}
		c.Next()
	}
}
