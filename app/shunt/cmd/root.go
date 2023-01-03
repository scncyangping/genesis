package cmd

import (
	"genesis/pkg/config"
	"genesis/pkg/config/app/shunt"
	"genesis/pkg/types/constants"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

func newRootCmd() *cobra.Command {
	args := struct {
		configPath string
	}{}
	cmd := &cobra.Command{
		Use:   "genesis-shunt",
		Short: "genesis Shunt Aggregate Url.",
		Long:  `genesis Shunt Aggregate Url.`,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			err := config.Load(args.configPath, shunt.ShuntConfig())

			if err != nil {
				return errors.Wrapf(err, "could not load the configuration")
			}

			cmd.SilenceUsage = true
			return nil
		},
	}

	cmd.SetOut(os.Stdout)

	// root flags
	cmd.PersistentFlags().StringVarP(&args.configPath, "config-file", "c", "", "configuration file")

	// sub-commands
	cmd.AddCommand(newRunCmd(constants.DefaultRunCmdOpts))

	return cmd
}

func DefaultRootCmd() *cobra.Command {
	return newRootCmd()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := DefaultRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
