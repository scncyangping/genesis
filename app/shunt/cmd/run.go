package cmd

import (
	"genesis/pkg/config/app/shunt"
	"genesis/pkg/core/shunt/adapter/http/routers"
	"genesis/pkg/core/shunt/adapter/http/server"
	"github.com/spf13/cobra"
	"os"
)

func newRunCmd() *cobra.Command {
	var mod string
	cmd := &cobra.Command{
		Use:   "run",
		Short: "genesis Shunt Aggregate Url.",
		Long:  `genesis Shunt Aggregate Url.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			h := NewHandler()
			server := server.NewHttpGin(mod, shunt.ShuntConfig().Server, shunt.Log())
			routers.InitRoute(server.Engine, h)
			server.Start()

			return nil
		},
	}

	cmd.SetOut(os.Stdout)
	// sub-commands
	cmd.PersistentFlags().StringVarP(&mod, "mod", "m", "debug", "run mod")

	return cmd
}
