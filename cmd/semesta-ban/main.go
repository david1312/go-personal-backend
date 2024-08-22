package main

import (
	"libra-internal/bootstrap"
	"libra-internal/pkg/log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/sirupsen/logrus"

	cmd "libra-internal/cmd/semesta-ban/commands"
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
