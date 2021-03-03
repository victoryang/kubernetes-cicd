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

			scm.NewGitHubClient("victoryang", conf.GithubSecret)

			// always do mysql init first
			orm.InitMysqlModule(conf.Database)

			// then other modules
			auth.InitAuthModule(conf.Ldap.Address, conf.Ldap.Password)
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
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/rolling")
	})
	staticFS := assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
	}
	router.StaticFS("/rolling", &staticFS)

	return router.Run(c.EndPoint)
}
