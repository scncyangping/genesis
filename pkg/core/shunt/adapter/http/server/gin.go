package server

import (
	"context"
	"genesis/pkg/config/app/shunt"
	"genesis/pkg/util/snowflake"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type HttpGin struct {
	Conf   *shunt.ServerConfig
	Engine *gin.Engine
	Logger *zap.SugaredLogger
}

func NewHttpGin(mod string, conf *shunt.ServerConfig, log *zap.SugaredLogger) *HttpGin {
	gin.SetMode(mod)

	g := gin.New()

	switch mod {
	case gin.ReleaseMode:
		g.Use(LogMiddleware(log), Recovery(log))
	default:
		g.Use(gin.Logger(), gin.Recovery())
	}

	return &HttpGin{Engine: g, Conf: conf, Logger: log}
}

func (e *HttpGin) Start() {
	server := &http.Server{
		Addr:           e.Conf.Addr,
		Handler:        e.Engine,
		ReadTimeout:    time.Duration(e.Conf.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(e.Conf.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	e.Logger.Infof("Server Start Success %s", e.Conf.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			e.Logger.Errorf("Server Error! %s", e.Conf.Addr)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	e.Logger.Errorf("Shutdown Server ...%s", e.Conf.Addr)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		e.Logger.Errorf("Server Shutdown: %v %v \n", e.Conf.Addr, err)
	}
	select {
	case <-ctx.Done():
		e.Logger.Errorf("Timeout of 10 seconds. %s", e.Conf.Addr)
	}

	e.Logger.Errorf("Server exiting, %s", e.Conf.Addr)
}

func LogMiddleware(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		actionId := snowflake.NextId()

		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		if raw != "" {
			path = path + "?" + raw
		}

		log.Infof("[Start-%s] %-7s | %15s | %#v",
			actionId,
			ctx.Request.Method,
			ctx.ClientIP(),
			path,
		)

		ctx.Set("ActionId", actionId)

		ctx.Next()

		// Stop timer
		latency := time.Now().Sub(start)

		if latency > time.Minute {
			// Truncate in a golang < 1.8 safe way
			latency = latency - latency%time.Second
		}
		if ctx.Writer.Status() != http.StatusOK {
			request, _ := httputil.DumpRequest(ctx.Request, true)

			log.Infof("[End - %s] %-7s | %15s | %#v | %3d | %13v \n%s",
				actionId,
				ctx.Request.Method,
				ctx.ClientIP(),
				path,
				ctx.Writer.Status(),
				latency,
				request,
			)
		} else {
			log.Infof("[End - %s] %-7s | %15s | %#v | %3d | %13v",
				actionId,
				ctx.Request.Method,
				ctx.ClientIP(),
				path,
				ctx.Writer.Status(),
				latency,
			)
		}
	}
}

func Recovery(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			var actionId = c.GetString("ActionId")
			if err := recover(); err != nil {
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {

							c.Error(err.(error))
							c.Abort()
							return
						}
					}
				}

				if httpRequest, er := httputil.DumpRequest(c.Request, true); er != nil {
					logger.Errorf("%s\n%s%v\n%v\n", actionId, string(httpRequest), err, er)
				} else {
					logger.Errorf("%s\n%s%v\n", actionId, string(httpRequest), err)
					c.JSON(200, map[string]any{"code": 500, "msg": "server error.", "data": nil})
				}
			}
		}()
		c.Next()
	}
}
