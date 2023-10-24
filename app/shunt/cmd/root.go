// @Author: YangPing
// @Create: 2023/10/23
// @Description: 程序入口

package cmd

import (
	config2 "genesis/app/shunt/config"
	"genesis/pkg/config"
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
		Use:   "bastion-shunt",
		Short: "bastion Shunt Aggregate Url.",
		Long:  `bastion Shunt Aggregate Url.`,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			bc := config2.DefaultBaseConfig()
			err := config.Load(args.configPath, bc)
			if err != nil {
				return errors.Wrapf(err, "could not load the configuration")
			} else {
				config2.Rt().BuildConfig(bc)
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
