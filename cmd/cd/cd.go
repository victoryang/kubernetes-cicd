package main

import (
	"os"

	"github.com/victoryang/kubernetes-cicd/cmd/cd/app"
)

func main() {

	command := app.NewCDManagerCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
