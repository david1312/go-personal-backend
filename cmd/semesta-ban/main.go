package main

import (
	"os"
	"semesta-ban/bootstrap"
	"semesta-ban/pkg/log"

	"github.com/sirupsen/logrus"

	cmd "semesta-ban/cmd/semesta-ban/commands"
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
