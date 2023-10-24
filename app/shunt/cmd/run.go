// @Author: YangPing
// @Create: 2023/10/23
// @Description: 初始化入口

package cmd

import (
	"genesis/app/shunt/adapter/http/routers"
	"genesis/app/shunt/adapter/http/server"
	"genesis/app/shunt/config"
	"genesis/pkg/runtime/component"
	"genesis/pkg/runtime/leader"
	"genesis/pkg/types/constants"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"time"
)

const gracefullyShutdownDuration = 3 * time.Second

func newRunCmd(opts constants.RunCmdOpts) *cobra.Command {
	var mod string
	cmd := &cobra.Command{
		Use:   "run",
		Short: "bastion Shunt Aggregate Url.",
		Long:  `bastion Shunt Aggregate Url.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			runLog := config.Rt().Logger().With("app", "shunt")
			rt := component.NewManager(runLog.With("manager"), leader.NewNeverLeaderElector())
			gracefulCtx, ctx := opts.SetupSignalHandler(runLog)

			// step1. withGinServer
			withGinServer(mod, runLog, rt)

			// last. start
			if err := rt.Start(gracefulCtx.Done()); err != nil {
				runLog.Error("problem running", "error", err)
				return err
			}
			select {
			case <-ctx.Done():
			case <-time.After(gracefullyShutdownDuration):
			}
			runLog.Info("stopping service")
			return nil
		},
	}

	cmd.SetOut(os.Stdout)
	// sub-commands
	cmd.PersistentFlags().StringVarP(&mod, "mod", "m", "debug", "run mod")

	return cmd
}

func withGinServer(mod string, runLog *slog.Logger, rt component.Manager) {
	// step1. add http server
	ginC := server.NewGinConfig().
		WithAddr(config.Rt().Config().GetAddr()).
		WithPort(config.Rt().Config().GetPort())
	httpServer := server.NewHttpGin(mod, ginC, runLog)

	routers.InitRoute(httpServer.Engine, NewHandler(config.Rt().GormDB()))
	rt.Add(component.NewResilientComponent(runLog, httpServer))
}
