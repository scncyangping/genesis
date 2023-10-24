// @Author: YangPing
// @Create: 2023/10/23
// @Description: Gin Server初始化

package server

import (
	"context"
	"fmt"
	"genesis/pkg/util/snowflake"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

type GinConfig struct {
	addr         string
	port         int
	readTimeout  int64
	writeTimeout int64
}

func NewGinConfig() *GinConfig {
	return &GinConfig{}
}

func (g *GinConfig) WithAddr(addr string) *GinConfig {
	g.addr = addr
	return g
}

func (g *GinConfig) WithPort(port int) *GinConfig {
	g.port = port
	return g
}

func (g *GinConfig) WithReadTimeout(readTimeout int64) *GinConfig {
	g.readTimeout = readTimeout
	return g
}

func (g *GinConfig) WithWriteTimeout(writeTimeout int64) *GinConfig {
	g.writeTimeout = writeTimeout
	return g
}

func (g *GinConfig) CheckDefault() {
	if g.readTimeout == 0 {
		g.readTimeout = 10000
	}
	if g.writeTimeout == 0 {
		g.writeTimeout = 10000
	}
}

type HttpGin struct {
	Engine *gin.Engine
	conf   *GinConfig
	logger *slog.Logger
}

func NewHttpGin(mod string, conf *GinConfig, log *slog.Logger) *HttpGin {
	conf.CheckDefault()
	gin.SetMode(mod)

	g := gin.New()

	switch mod {
	case gin.ReleaseMode:
		g.Use(RequestFilterFunc(), LogMiddleware(log), Recovery(log))
	default:
		pprof.Register(g)
		g.Use(RequestFilterFunc(), gin.Logger(), gin.Recovery())
	}

	return &HttpGin{Engine: g, conf: conf, logger: log}
}

func (e *HttpGin) NeedLeaderElection() bool {
	return false
}

func (e *HttpGin) Start(stop <-chan struct{}) error {

	errChan := make(chan error)

	httpServer := e.startHttpServer(errChan)

	select {
	case <-stop:
		if httpServer != nil {
			return httpServer.Shutdown(context.Background())
		}
	case err := <-errChan:
		return err
	}
	return nil
}

func (e *HttpGin) startHttpServer(errChan chan error) *http.Server {
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", e.conf.addr, e.conf.port),
		Handler:        e.Engine,
		ReadTimeout:    time.Duration(e.conf.readTimeout) * time.Second,
		WriteTimeout:   time.Duration(e.conf.writeTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err != nil {
				switch err {
				case http.ErrServerClosed:
					e.logger.Info("shutting down server")
				default:
					e.logger.Error(fmt.Sprintf("could not start an HTTP Server,%v", err))
					errChan <- err
				}
			}
		} else {
			e.logger.Info(fmt.Sprintf("server start success, addr: %s", e.conf.addr))
		}
	}()

	return server
}

func LogMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		actionId := ctx.GetString("ActionId")
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		if raw != "" {
			path = path + "?" + raw
		}

		log.Info(fmt.Sprintf("[Start-%s] %-7s | %15s | %#v",
			actionId,
			ctx.Request.Method,
			ctx.ClientIP(),
			path),
		)

		ctx.Next()

		// Stop timer
		latency := time.Now().Sub(start)

		if latency > time.Minute {
			// Truncate in a golang < 1.8 safe way
			latency = latency - latency%time.Second
		}
		if ctx.Writer.Status() != http.StatusOK {
			request, _ := httputil.DumpRequest(ctx.Request, true)

			log.Info(fmt.Sprintf("[End - %s] %-7s | %15s | %#v | %3d | %13v \n%s",
				actionId,
				ctx.Request.Method,
				ctx.ClientIP(),
				path,
				ctx.Writer.Status(),
				latency,
				request))
		} else {
			log.Info(fmt.Sprintf("[End - %s] %-7s | %15s | %#v | %3d | %13v",
				actionId,
				ctx.Request.Method,
				ctx.ClientIP(),
				path,
				ctx.Writer.Status(),
				latency),
			)
		}
	}
}

func Recovery(logger *slog.Logger) gin.HandlerFunc {
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
					logger.Error(fmt.Sprintf("%s\n%s%v\n%v\n", actionId, string(httpRequest), err, er))
				} else {
					logger.Error(fmt.Sprintf("%s\n%s%v\n", actionId, string(httpRequest), err))
					c.JSON(200, map[string]any{"code": 500, "msg": "server error.", "data": nil})
				}
			}
		}()
		c.Next()
	}
}

func RequestFilterFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		actionId := snowflake.NextId()
		ctx.Set("ActionId", actionId)
		ctx.Next()
	}
}
