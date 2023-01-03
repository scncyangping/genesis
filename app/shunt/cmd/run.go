package cmd

import (
	"genesis/pkg/config/app/shunt"
	"genesis/pkg/core/shunt/adapter/http/routers"
	"genesis/pkg/core/shunt/adapter/http/server"
	"genesis/pkg/runtime/component"
	"genesis/pkg/runtime/leader"
	"genesis/pkg/types/constants"
	"github.com/spf13/cobra"
	"os"
	"time"
)

const gracefullyShutdownDuration = 3 * time.Second

func newRunCmd(opts constants.RunCmdOpts) *cobra.Command {
	var mod string
	cmd := &cobra.Command{
		Use:   "run",
		Short: "genesis Shunt Aggregate Url.",
		Long:  `genesis Shunt Aggregate Url.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			log := shunt.Log()
			runLog := log.Named("Run")

			rt := component.NewManager(shunt.Log().Named("manager"), leader.NewNeverLeaderElector())

			gracefulCtx, ctx := opts.SetupSignalHandler(log)

			// step1. add http server
			httpServer := server.NewHttpGin(mod, shunt.ShuntConfig().Server, log)
			routers.InitRoute(httpServer.Engine, NewHandler())

			rt.Add(component.NewResilientComponent(log, httpServer))

			if err := rt.Start(gracefulCtx.Done()); err != nil {
				runLog.Error(err, "problem running")
				return err
			}

			runLog.Info("Stop signal received. Waiting 3 seconds for components to stop gracefully...")
			select {
			case <-ctx.Done():
			case <-time.After(gracefullyShutdownDuration):
			}
			runLog.Info("Stopping universal service mesh manager")
			return nil
		},
	}

	cmd.SetOut(os.Stdout)
	// sub-commands
	cmd.PersistentFlags().StringVarP(&mod, "mod", "m", "debug", "run mod")

	return cmd
}
