package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	//"github.com/spf13/pflag"
	"github.com/victoryang/kubernetes-cicd/config"
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
			conf,err := config.LoadFile(ConfigFile)
			if err!=nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return
			}

			fmt.Println(conf)
		},
	}

	cmd.Flags().StringVar(&ConfigFile, "config", "", "Configuration file")

	return cmd
}
