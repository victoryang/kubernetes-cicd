package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/gin-gonic/gin"

	"github.com/victoryang/kubernetes-cicd/auth"
	"github.com/victoryang/kubernetes-cicd/config"
	"github.com/victoryang/kubernetes-cicd/image"
	"github.com/victoryang/kubernetes-cicd/logger"
	"github.com/victoryang/kubernetes-cicd/orm"
	"github.com/victoryang/kubernetes-cicd/project"
)

var (
	ConfigFile string
)

func NewCDManagerCommand() *cobra.Command {

	cmd := &cobra.Command {
		Use: "kubernetes-cd-manager",
		Long: `A kubernetes cd manager`,
		//Args: cobra.MinimumNArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			conf,err := config.LoadFile(ConfigFile)
			if err!=nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return
			}

			// always do mysql init first
			orm.InitMysqlModule(conf.Database)

			// then other modules
			auth.InitAuthModule(conf.Ldap.Address, conf.Ldap.Password)
			image.InitImageModule()
			project.InitProjectModule()
			logger.InitLoggerModule(conf.Log.File, conf.Log.Level)

			kuberenetes.InitKubernetesModules()

			if err := run(conf); err!=nil {
				log.Fatal("Run app failed,", err)
			}
		},
	}

	cmd.Flags().StringVar(&ConfigFile, "config", "", "Configuration file")

	return cmd
}

func run(c *config.Config) error {

	// router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(Logger())

	//验证
	router.POST("/login.json", login)

	// hook
	router.POST("/receive_hook.json", receiveHook)

	v1 := router.Group("/")
	v1.Use(authorize())
	{
		v1.POST("/create_project.json", createProject)
		v1.POST("/update_project.json", updateProject)
		v1.GET("/project_config.json", getProjectConfig)
		v1.GET("/project_runtime.json", getProjectRuntime)
		v1.POST("/create_env.json", createEnv)
		v1.POST("/update_env_code.json", updateEnvCodeVersion)
		v1.GET("/get_env_node_num.json", getEnvNodeNum)
		v1.POST("/set_env_node_num.json", setEnvNodeNum)
		//镜像相关
		v1.GET("/images.json", getImageList)
	}

	// asset
	staticFS := assetfs.AssetFS{
		Asset:     static.Asset,
		AssetDir:  static.AssetDir,
	}
	router.StaticFS("/rolling", &staticFS)

	return router.Run(c.EndPoint)
}