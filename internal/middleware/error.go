package middleware

import (
	"github.com/Comvoca-AI/comvoca-admin-back/internal/errors"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/logger"
	"github.com/gofiber/fiber/v2"
	"runtime/debug"
)

// ErrorHandler creates a middleware that handles panics and errors encountered during HTTP request processing.
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic recovered: ", r)
				logger.Error(string(debug.Stack()))
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
			}
		}()

		err := c.Next()
		if err != nil {
			logger.Error("Error encountered: ", err)
			errorResponse := buildErrorResponse(err)
			return c.Status(errorResponse.Status).JSON(fiber.Map{"error": errorResponse.Message})
		}

		return nil
	}
}

// buildErrorResponse builds an error response from an error.
func buildErrorResponse(err error) errors.ErrorResponse {
	switch err.(type) {
	case errors.ErrorResponse:
		return err.(errors.ErrorResponse)

	case *fiber.Error:
		switch err.(*fiber.Error).Code {
		case fiber.StatusNotFound:
			return errors.NotFound("")
		default:
			return errors.ErrorResponse{
				Status:  err.(*fiber.Error).Code,
				Message: err.Error(),
			}
		}
	}

	return errors.InternalServerError("Internal Server DB Error")
}
