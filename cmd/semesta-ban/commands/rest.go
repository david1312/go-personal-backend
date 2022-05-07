package commands

import (
	"fmt"
	"net/http"

	"github.com/semestaban/internal-api/bootstrap"
	"github.com/spf13/cobra"
)

func init() {
	registerCommand(startRestService)
}

func startRestService(dep *bootstrap.Dependency) *cobra.Command {
	return &cobra.Command{
		Use:   "rest",
		Short: "Starting REST service",
		Long:  `This command is used to start REST service`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := dep.GetConfig()
			http.HandleFunc("/hello",hello)
			fmt.Printf("server is running on %s", cfg.Host.Address)
			http.ListenAndServe(":8090",nil)
		},
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
 
	url:=r.URL
	
	fmt.Fprintf(w,"hello haha xixixixi from %v",url)
	
}