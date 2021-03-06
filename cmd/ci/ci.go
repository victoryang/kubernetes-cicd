package ci

import (
	"os"

	"github.com/victoryang/kubernetes-cicd/cmd/ci/app"
)

func main() {

	command := app.NewCIManagerCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}