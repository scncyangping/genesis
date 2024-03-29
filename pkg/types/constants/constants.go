// @Author: YangPing
// @Create: 2023/10/23
// @Description: 全局变量、常量定义

package constants

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type RunCmdOpts struct {
	// The first returned context is closed upon receiving first signal (SIGSTOP, SIGTERM).
	// The second returned context is closed upon receiving second signal.
	// We can start graceful shutdown when first context is closed and forcefully stop when the second one is closed.
	SetupSignalHandler func(log *slog.Logger) (context.Context, context.Context)
}

var DefaultRunCmdOpts = RunCmdOpts{
	SetupSignalHandler: SetupSignalHandler,
}

var SetupSignalHandler = func(log *slog.Logger) (context.Context, context.Context) {
	gracefulCtx, gracefulCancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 3)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-c
		log.Info("received signal, stopping instance gracefully ", " signal ", s.String())
		gracefulCancel()
		s = <-c
		log.Info("received second signal, stopping instance ", " signal ", s.String())
		cancel()
		s = <-c
		log.Info("received third signal, force exit ", " signal ", s.String())
		os.Exit(1)
	}()
	return gracefulCtx, ctx
}
