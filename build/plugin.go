package build

import (
	"context"
	"net/http"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
	"github.com/drone/drone-go/plugin/webhook"
	"github.com/sirupsen/logrus"
)

const (
	YamlPluginSecret = "bea26a2221fd8090ea38720fc445eca6"
	WebhookPluginSecret = "bea26a2221fd8090ea38720fc445eca6"
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

	bp,err := NewBuildPipeline(req.Repo, req.Build)
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

type WebhookPlugin struct {

}

func NewWebhookPlugin() http.Handler {

	logrus.SetLevel(logrus.DebugLevel)

	handler := webhook.Handler(
		&WebhookPlugin{},
		WebhookPluginSecret,
		logrus.StandardLogger(),
	)

	return handler
}

func (p *WebhookPlugin) Deliver(ctx context.Context, req *webhook.Request) error {
	switch req.Event {
		case "build":
			go ProcessBuildEvent(req)
		case "user":
			go ProcessUserEvent(req)
		case "repo":
			go ProcessRepoEvent(req)
		default:
	}

	return nil
}