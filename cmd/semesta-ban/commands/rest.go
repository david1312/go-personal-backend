package commands

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"semesta-ban/bootstrap"
	"semesta-ban/internal/api"
	"semesta-ban/internal/api/transactions"
	"semesta-ban/pkg/helper"
	"semesta-ban/pkg/log"
	"syscall"

	"github.com/spf13/cobra"
)

func init() {
	registerCommand(startRestService)
}

//todo create script for  migration db automation
func startRestService(dep *bootstrap.Dependency) *cobra.Command {
	return &cobra.Command{
		Use:   "rest",
		Short: "Starting REST service",
		Long:  `This command is used to start REST service`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := dep.GetConfig()
			ctx := context.Background()

			client := helper.CreateHttpClient(ctx, cfg.Midtrans.Timeout, false)

			handler := api.NewServer(dep.GetDB(), client, api.ServerConfig{
				EncKey:            cfg.Key.EncryptKey,
				JWTKey:            cfg.Key.JWT,
				AnonymousKey:      cfg.Key.Anonymous,
				BaseAssetsUrl:     cfg.Assets.BaseUrl,
				UploadPath:        cfg.Assets.UploadPath,
				ProfilePicPath:    cfg.Assets.ProfilePic.Path,
				ProfilePicMaxSize: cfg.Assets.ProfilePic.MaxSize,
				MaxFileSize:       cfg.Assets.Common.MaxFileSize,
				MidtransConfig: transactions.MidtransConfig{
					MerchantId: cfg.Midtrans.MerchantId,
					ClientKey:  cfg.Midtrans.ClientKey,
					ServerKey:  cfg.Midtrans.ServerKey,
				},
			})
			// application context, which will be cancelled upon receiving termination signal
			actx, _ := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

			srv := http.Server{Addr: cfg.Host.Address, Handler: handler}
			//end testing

			//implement graceful shutdown
			errChan := make(chan error)
			defer close(errChan)
			go func() {
				log.Infof("server is running on %s at %s env", cfg.Host.Address, cfg.Env)
				err := srv.ListenAndServe()
				if err != nil && err != http.ErrServerClosed {
					errChan <- errors.New("server error: " + err.Error())
				}
			}()

			select {
			case err := <-errChan:
				log.Error(err)
				return
			case <-actx.Done():
				err := srv.Shutdown(ctx)
				if err != nil {
					log.Error("Shutdown error:", err)
					return
				}
			}
			log.Info("Server shutdown gracefully.")
		},
	}
}
