package commands

import (
	"context"
	"errors"
	"libra-internal/bootstrap"
	"libra-internal/internal/api"
	"libra-internal/internal/api/customers"
	"libra-internal/internal/api/transactions"
	"libra-internal/pkg/helper"
	"libra-internal/pkg/log"
	"net/http"
	"os/signal"
	"syscall"

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
			ctx := context.Background()

			client := helper.CreateHttpClient(ctx, cfg.Midtrans.Timeout, true)

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
					AuthKey:    helper.GenerateB64AuthMidtrans(cfg.Midtrans.ServerKey),
				},
				SMTPConfig: customers.SMTPConfig{
					Host:         cfg.SMTP.Host,
					Port:         cfg.SMTP.Port,
					SenderName:   cfg.SMTP.SenderName,
					AuthEmail:    cfg.SMTP.AuthEmail,
					AuthPassword: cfg.SMTP.AuthPassword,
				},
				FcmConfig: transactions.FCMConfig{
					NotifUrl:  cfg.FCM.NotifUrl,
					ClientKey: cfg.FCM.ClientKey,
				},
			})
			// application context, which will be cancelled upon receiving termination signal
			actx, _ := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

			srv := http.Server{Addr: cfg.Host.Address, Handler: handler}

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
