package middleware

import (
	"github.com/Comvoca-AI/comvoca-admin-back/internal/errors"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/service"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func AuthMiddleware(securityService *service.SecurityService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// List of endpoints to exclude from token validation

		excludedPaths := []string{"/api/v1/password/",
			"/api/v1/password/forgot/confirm",
			"/api/v1/password/temp/confirm",
			"/api/v1/password/forgot",
			"/api/v1/login",
			"/api/v1/otp",
			"/api/v1/otp/resend",
			"/api/v1/register/user"}

		// Check if the request path is in the excluded paths
		for _, path := range excludedPaths {
			if strings.HasPrefix(c.Path(), path) {
				return c.Next()
			}
		}

		// Get the token from the Authorization header
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return errors.Unauthorized("Missing or invalid token")
		}

		// Validate the token
		token, err := securityService.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			return errors.Unauthorized("Invalid token")
		}

		// Proceed to the next middleware/handler
		return c.Next()
	}
}
