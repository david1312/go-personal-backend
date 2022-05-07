package commands

import (
	"io"
	"os"
	"runtime"

	"github.com/ndv6/semestaban/internal-api/bootstrap"
	"github.com/ndv6/semestaban/internal-api/pkg/log"
	"github.com/spf13/cobra"
)

type commandFn func(dep *bootstrap.Dependency) *cobra.Command

var subCommands []commandFn

func registerCommand(fn commandFn) {
	subCommands = append(subCommands, fn)
}

func Run(dep *bootstrap.Dependency) error {
	var cpu int
	var config string
	var verbose bool
	var tracerCloser io.Closer

	rootCommand := &cobra.Command{
		Use: "one-notif",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
	

			// Set number of CPU
			x := runtime.GOMAXPROCS(cpu)
			if cpu == 0 {
				cpu = x
			}
			log.Debugf("using %v cpu(s)", cpu)

			// Set configuration file
			log.Debugf("load config from: %v", config)
			cfg, err := bootstrap.LoadConfig(config)
			if err != nil {
				log.Errorf("unable to load config file %s: %s", config, err)
				os.Exit(1)
			}

			// Initialize dependency injection
			dep.SetConfig(cfg)
			if err := dep.Initialize(); err != nil {
				log.Errorf("Fail to load dependency: %v", err)
			}

		},

		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if tracerCloser != nil {
				tracerCloser.Close()
			}
		},

		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCommand.PersistentFlags().IntVar(&cpu, "cpu", 0, "set the number of CPU to use, default to 0, which means it will use all available CPU")
	rootCommand.PersistentFlags().StringVarP(&config, "config", "c", "config.yaml", "set configuration file to use")
	rootCommand.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show debug level log")

	for _, fn := range subCommands {
		rootCommand.AddCommand(fn(dep))
	}

	return rootCommand.Execute()
}
