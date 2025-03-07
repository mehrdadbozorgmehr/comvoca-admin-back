package rest

import (
	"github.com/Comvoca-AI/comvoca-admin-back/internal/service"
	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication-related requests
type CallHandler struct {
	AuthService service.SecurityService
	UserService service.UserService
}

// NewAuthHandler creates a new AuthHandler
func NewCallHandler(authService service.SecurityService, userService service.UserService) *CallHandler {
	return &CallHandler{
		authService,
		userService,
	}
}

// Register registers the authentication routes
func (authHandler *CallHandler) Register(app *fiber.App) {
	app.Get("/api/v1/call/vapi-assistant-calls", authHandler.GetVapiAssistantCalls)
}

func (authHandler *CallHandler) GetVapiAssistantCalls(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "GetVapiAssistantCalls"})
}
