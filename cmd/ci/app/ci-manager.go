package app

import (
	"fmt"

	"github.com/spf13/cobra"
	//"github.com/spf13/pflag"
)

var (
	ConfigFile string
)

func NewCIManagerCommand() *cobra.Command {

	cmd := &cobra.Command {
		Use: "kubernetes-ci-manager",
		Long: `A kubernetes ci manager`,
		//Args: cobra.MinimumNArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(ConfigFile)
		},
	}

	cmd.Flags().StringVar(&ConfigFile, "config", "", "Configuration file")

	return cmd
}
