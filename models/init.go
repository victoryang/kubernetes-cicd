package models

import (
	"log"

	"github.com/victoryang/kubernetes-cicd/logger"
)

var (
	Logger		*logger.Logger
)

func InitLogger(filename string) {
	var err error
	Logger,err = logger.NewLogger(filename)
	if err!=nil {
		log.Fatal("init log failed: ", err) 
	}
}