package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"main.go/bootstrap"
	cmd "main.go/cmd/commands"
	"main.go/pkg/log"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// service dependency
	dependency := bootstrap.NewDependency()
	if err := cmd.Run(dependency); err != nil {
		log.Errorf("unable to execute root command: %s", err)
		os.Exit(1)
	}
}
