package cmd

import (
	"fmt"

	"github.com/Comvoca-AI/comvoca-admin-back/config"
	"github.com/Comvoca-AI/comvoca-admin-back/initialize"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/logger"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Running the HTTP server...")

		// Create a new Fiber instance
		app := fiber.New()

		app.Use(middleware.ErrorHandler())

		app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowMethods: "*",
			AllowHeaders: "*",
		}))
		app.Static("/", "./docs")
		app.Use(func(c *fiber.Ctx) error {
			logger.Logger.Debug().
				Str("method", c.Method()).
				Str("path", c.Path()).
				Msg("request")
			return c.Next()
		})

		app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept",
		}))

		if config.AppConfig.Application.Debug {
			app.Static("/", "./docs")
			app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
				URL:         config.AppConfig.Application.BaseURL + "/swagger.json",
				DeepLinking: false,
				// Expand ("list") or Collapse ("none") tag groups by default
				DocExpansion: "none",
				// Prefill OAuth ClientId on Authorize popup
				OAuth: &swagger.OAuthConfig{
					AppName:  "OAuth Provider",
					ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
				},
				// Ability to change OAuth2 redirect uri location
				OAuth2RedirectUrl: config.AppConfig.Application.BaseURL + "/swagger/oauth2-redirect.html",
			}))
		}

		initialize.InitializeRoutes(app)

		if err := app.Listen(fmt.Sprintf(":%d", config.AppConfig.Server.Port)); err != nil {
			logger.Error("Error running the HTTP server", err)
		}
	},
}
