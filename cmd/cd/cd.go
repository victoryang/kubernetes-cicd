package main

import (
	"os"

	"github.com/victoryang/kubernetes-cicd/cmd/cd/app"
)

//go:generate go-bindata -debug -o=app/asset.go -pkg=app -ignore asset.go ../../static/...

func main() {

	command := app.NewCDManagerCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
