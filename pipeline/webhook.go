package pipeline

import (
	"github.com/drone/drone-go/plugin/webhook"
)

func ProcessBuildEvent(req *webhook.Request) {

	droneInfo := ProcessRepoAndEventInfo(req.Repo, req.Build)

	switch req.Action {
	case "created":
		CDServer.CreateImage(droneInfo.Project, droneInfo.Tag)
	case "updated":
		if req.Build.Status == "success" || req.Build.Status == "failure"{
			CDServer.UpdateImage(droneInfo.Project, droneInfo.Tag, droneInfo.Env, req.Build.Status, req.Build.Error)
		}
	}

	return
}

func ProcessUserEvent(req *webhook.Request) {
	// TODO
	return
}

func ProcessRepoEvent(req *webhook.Request) {
	// TODO
	return
}