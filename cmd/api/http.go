package main

import (
	"net/http"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/arnaz06/users/internal"
	handler "github.com/arnaz06/users/internal/http"
)

const address = ":7723"

var serverCmd = &cobra.Command{
	Use:   "http",
	Short: "Turn the server on",
	Run: func(cmd *cobra.Command, args []string) {
		e := echo.New()
		e.Validator = internal.NewValidator()
		e.Use(
			handler.TimeoutMiddleware(contextTimeout),
			handler.ErrorMiddleware(),
			middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
				Validator: handler.AuthenticationMiddleware(secretKey),
				Skipper: func(c echo.Context) bool {
					if c.Path() == `/user` || c.Path() == `/user/login` {
						return true
					}
					return false
				},
			}),
		)
		handler.AddUserHandler(e, userService, secretKey, expiresTime)

		e.GET("ping", func(c echo.Context) error {
			return c.String(http.StatusOK, "pong")
		})

		log.Info("Starting HTTP server at ", address)
		err := e.Start(address)
		if err != nil {
			log.Fatalf("Failed to start server: %s", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
