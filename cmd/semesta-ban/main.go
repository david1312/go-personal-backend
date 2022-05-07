package main

import (
	"os"

	"github.com/semestaban/internal-api/bootstrap"
	"github.com/semestaban/internal-api/pkg/log"
	"github.com/sirupsen/logrus"

	cmd "github.com/semestaban/internal-api/cmd/semesta-ban/commands"
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
