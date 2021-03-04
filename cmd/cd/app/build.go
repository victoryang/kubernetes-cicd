package app

import (
	"context"
	"net/http"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
	"github.com/sirupsen/logrus"

	"github.com/victoryang/kubernetes-cicd/build"
)

const (
	YamlPluginSecret = "bea26a2221fd8090ea38720fc445eca6"
)

type YamlPlugin struct {

}

func NewYamlPlugin() http.Handler {

	logrus.SetLevel(logrus.DebugLevel)

	handler := config.Handler(
		&YamlPlugin{},
		YamlPluginSecret,
		logrus.StandardLogger(),
	)

	return handler
}

func (p *YamlPlugin) Find(ctx context.Context, req *config.Request) (*drone.Config, error) {

	logrus.Info("New coming request")
	logrus.Info("Repo Info", req.Repo)
	logrus.Info("Build Info", req.Build)

	bp,err := build.NewBuildPipeline(req.Repo, req.Build)
	if err!=nil {
		return nil,err
	}

	data, err := bp.Compile()
	if err!=nil {
		return nil, err
	}

	return &drone.Config {
		Data: data,
	}, nil
}