package app

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/gin-gonic/gin"

	"github.com/victoryang/kubernetes-cicd/pipeline"
)

func NewCIManagerCommand() *cobra.Command {

	cmd := &cobra.Command {
		Use: "kubernetes-ci-manager",
		Long: `A kubernetes ci manager`,
		//Args: cobra.MinimumNArgs(1),

		Run: func(cmd *cobra.Command, args []string) {

			pipeline.InitCDServerClientWithRemoteMode()

			if err := run(); err!=nil {
				log.Fatal("Run app failed,", err)
			}
		},
	}

	//cmd.Flags().StringVar(&ConfigFile, "config", "settings.yaml", "Configuration file")

	return cmd
}

func run() error {

	// router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(Logger())

	// drone-ci build hook
	apiv1 := router.Group("/api/v1")
	{
		apiv1.POST("/build/config", gin.WrapH(NewYamlPlugin()))
		apiv1.POST("/build/webhook", gin.WrapH(NewWebhookPlugin()))
	}

	return router.Run(":5000")
}