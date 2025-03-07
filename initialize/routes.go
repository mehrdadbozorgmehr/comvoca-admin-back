package initialize

import (
	"github.com/Comvoca-AI/comvoca-admin-back/internal/healthcheck"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/logger"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/repository"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/repository/db"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/rest"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func InitializeRoutes(app *fiber.App) {
	// Connect to the database
	customGromDb := db.ConnectDB()

	// Create the DAO
	organizationRepo := repository.NewOrganizationRepo(customGromDb)
	userRepo := repository.NewUserDAO(customGromDb)

	// Create the service
	organizationService := service.NewOrganizationService(customGromDb, organizationRepo)
	userService := service.NewUserService(userRepo)
	securityService := service.NewSecurityService(customGromDb, userService, organizationService)

	//Create the rest API
	authHandler := rest.NewAuthHandler(*securityService, *userService)
	organizationHandler := rest.NewOrganizationHandler(*organizationService)

	//Register the handlers
	authHandler.Register(app)
	organizationHandler.Register(app)

	//Register specific routes
	app.Get("/health", healthcheck.Healthcheck())
	app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))

	logger.Info("Application initialized successfully.")
}
