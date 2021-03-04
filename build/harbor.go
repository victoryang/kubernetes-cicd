package build

import (
	"path"
)

const (
	HarborBaseUrl = "hub.snowballfinance.com"
	HarborPublicProject = "cicd"
)

func BuildImageName(project string, tag string) string {

	return path.Join(HarborBaseUrl, HarborPublicProject, project) + ":" + tag
}