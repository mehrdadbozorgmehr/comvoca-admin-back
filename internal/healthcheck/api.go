package healthcheck

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

// Healthcheck responds to a healthcheck request.
func Healthcheck() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.SendString("OK - " + time.Now().Format(time.RFC3339))
	}

}
