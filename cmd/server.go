package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mrjosh/icloud-mailbox/config"
	"github.com/mrjosh/icloud-mailbox/pkg/smtp"
	"github.com/spf13/cobra"
)

type ServerCommandFlags struct {
	Host       string
	Port       int
	ConfigFile string
}

type ServerCommand struct {
	m *Meta
}

func (s *ServerCommand) cmd() *cobra.Command {
	cFlags := new(ServerCommandFlags)
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start mailbox server",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger := s.m.Logger
			logger.Info(fmt.Sprintf("loading config file [%s]", cFlags.ConfigFile))
			conf, err := config.LoadFile(cFlags.ConfigFile)
			if err != nil {
				return err
			}

			e := echo.New()

			// Set settings for production mode
			if conf.Environment == "prod" {
				logger.Infof("Running in production mode")
				e.Debug = false
				e.HideBanner = true
				e.HidePort = true
			}

			e.Use(middleware.CORS())
			e.Use(middleware.Logger())
			e.Use(middleware.Recover())
			e.Use(middleware.Gzip())
			//e.Use(middleware.BodyLimit("2M"))

			smtpConf := smtp.Config{
				Host:     conf.SMTP.Host,
				Port:     conf.SMTP.Port,
				Username: conf.SMTP.User,
				Password: conf.SMTP.Pass,
				From:     conf.SMTP.From,
			}
			smtpClient, err := smtp.New(smtpConf)
			if err != nil {
				return err
			}

			e.POST("/v1/sendMail", func(ctx echo.Context) error {

				secretID := ctx.Request().Header.Get("Authorization")
				if secretID == "" {
					return ctx.JSON(http.StatusUnauthorized, map[string]any{
						"success": false,
						"message": "headers.Authorization is required",
					})
				}

				if secretID != conf.Auth.SecretID {
					return ctx.JSON(http.StatusUnauthorized, map[string]any{
						"success": false,
						"message": "http.StatusUnauthorized",
					})
				}

				req := new(smtp.Notification)
				if err := ctx.Bind(&req); err != nil {
					if conf.Environment == "debug" {
						logger.Error(err)
					}
					return ctx.JSON(http.StatusInternalServerError, map[string]any{
						"success": false,
						"message": "ctx.Bind.Failed",
					})
				}

				if err := smtpClient.Send(ctx.Request().Context(), *req); err != nil {
					if conf.Environment == "debug" {
						logger.Error(err)
					}
					return ctx.JSON(http.StatusInternalServerError, map[string]any{
						"success": false,
						"message": "smtpClient.Send.Failed",
					})
				}

				return ctx.JSON(http.StatusOK, map[string]any{
					"success": true,
					"message": "smtpClient.Send.Succeed",
				})
			})

			msg := make(chan error)

			go func() {
				logger.Infof(
					"listening on http://%s:%d",
					cFlags.Host,
					cFlags.Port,
				)
				e.Server = &http.Server{
					Addr:         fmt.Sprintf("%s:%d", cFlags.Host, cFlags.Port),
					ReadTimeout:  10 * time.Minute, // Timeout for reading the request
					WriteTimeout: 10 * time.Minute, // Timeout for writing the response
					IdleTimeout:  10 * time.Minute, // Timeout for idle connections
				}
				msg <- e.StartServer(e.Server)
			}()

			go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
				msg <- fmt.Errorf("%s", <-c)
			}()

			return <-msg
		},
	}

	cmd.SuggestionsMinimumDistance = 1

	// Defining The Flags
	cmd.PersistentFlags().StringVarP(&cFlags.ConfigFile, "config-file", "c", "config.yaml", "Using this flag you can specify the config filep path")
	cmd.PersistentFlags().StringVarP(&cFlags.Host, "host", "H", "127.0.0.1", "Using this flag you can specify the listen address")
	cmd.PersistentFlags().IntVarP(&cFlags.Port, "port", "P", 8937, "Using this flag you can specify the listen port")

	cmd.MarkPersistentFlagRequired("config-file")
	return cmd
}

// Add the current command to cobra interface
func (d *ServerCommand) AddCommandToCobra(rootCmd *cobra.Command) {
	rootCmd.AddCommand(d.cmd())
}
