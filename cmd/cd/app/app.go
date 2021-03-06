package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/gin-gonic/gin"
	"github.com/elazarl/go-bindata-assetfs"

	"github.com/victoryang/kubernetes-cicd/auth"
	"github.com/victoryang/kubernetes-cicd/build"
	"github.com/victoryang/kubernetes-cicd/config"
	"github.com/victoryang/kubernetes-cicd/image"
	"github.com/victoryang/kubernetes-cicd/models"
	"github.com/victoryang/kubernetes-cicd/orm"
	"github.com/victoryang/kubernetes-cicd/project"
	"github.com/victoryang/kubernetes-cicd/scm"
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

			config.SetDebugMode(conf.DebugMode)

			models.InitLogger(conf.Log.File)

			scm.NewGitHubClient(conf.GithubToken)

			// always do mysql init first
			orm.InitMysqlModule(conf.Database)

			// then other modules
			auth.InitAuthModule(conf.Ldap.Address, conf.Ldap.Password)
			build.InitCDServerClientWithLocaldMode()
			image.InitImageModule()
			project.InitProjectModule()

			if err := run(conf); err!=nil {
				log.Fatal("Run app failed,", err)
			}
		},
	}

	cmd.Flags().StringVar(&ConfigFile, "config", "settings.yaml", "Configuration file")

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

	// drone-ci build hook
	apiv1 := router.Group("/api/v1")
	{
		apiv1.POST("/build/config", gin.WrapH(build.NewYamlPlugin()))
		apiv1.POST("/build/webhook", gin.WrapH(build.NewWebhookPlugin()))
	}

	root := router.Group("/")
	root.Use(authorize())
	{
		root.POST("/create_project.json", createProject)
		root.POST("/update_project.json", updateProject)
		root.GET("/project_config.json", getProjectConfig)
		root.GET("/project_runtime.json", getProjectRuntime)
		root.POST("/create_env.json", createEnv)
		root.POST("/update_env_code.json", updateEnvCodeVersion)
		root.GET("/get_env_node_num.json", getEnvNodeNum)
		root.POST("/set_env_node_num.json", setEnvNodeNum)
		//镜像相关
		root.GET("/images.json", getImageList)
	}

	// asset
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/kubernetes")
	})
	staticFS := assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
	}
	router.StaticFS("/kubernetes", &staticFS)

	return router.Run(c.EndPoint)
}
