package rest

import (
	"github.com/Comvoca-AI/comvoca-admin-back/internal/logger"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/service"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/types"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/validator"
	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	AuthService service.SecurityService
	UserService service.UserService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService service.SecurityService, userService service.UserService) *AuthHandler {
	return &AuthHandler{
		authService,
		userService,
	}
}

// Register registers the authentication routes
func (authHandler *AuthHandler) Register(app *fiber.App) {
	app.Post("/api/v1/password/", authHandler.ChangePassword)
	app.Post("/api/v1/password/forgot/confirm", authHandler.ConfirmForgotPassword)
	app.Post("/api/v1/password/temp/confirm", authHandler.ConfirmTempPassword)
	app.Post("/api/v1/password/forgot", authHandler.ForgotPassword)
	app.Post("/api/v1/auth", authHandler.Authenticate)
	app.Post("/api/v1/otp", authHandler.ActivateOTP)
	app.Post("/api/v1/otp/resend", authHandler.ResendOTP)
	app.Post("/api/v1/user/register", authHandler.RegisterUser)
}

// validateRequest validates the request body
func validateRequest(c *fiber.Ctx, req interface{}) error {
	v := validator.GetValidator()
	if err := v.Struct(req); err != nil {
		logger.Error("Request validation failed", err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return nil
}

// Authenticate godoc
// @Summary Authenticate user
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body types.AuthRequest true "Authenticate request"
// @Success 200 {string} string "Authenticate successful, Change Password Required, Activate OTP Required"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Invalid Credentials"
// @Router /api/v1/login [post]
func (authHandler *AuthHandler) Authenticate(c *fiber.Ctx) error {
	return authHandler.AuthService.Login(c)
}

// ActivateOTP godoc
// @Summary Activate OTP
// @Description Activate OTP for user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body types.ActivateOTPRequest true "Activate OTP request"
// @Success 200 {string} string "User activated successfully"
// @Failure 400 {string} string "Invalid request"
// @Router /api/v1/otp [post]
func (authHandler *AuthHandler) ActivateOTP(c *fiber.Ctx) error {
	return authHandler.AuthService.ActivateOTP(c)
}

// ResendOTP godoc
// @Summary Resend OTP
// @Description Resend OTP to user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body types.ResendOTPRequest true "Resend OTP request"
// @Success 200 {string} string "OTP resent successfully"
// @Failure 400 {string} string "Invalid request"
// @Router /api/v1/otp/resend [post]
func (authHandler *AuthHandler) ResendOTP(c *fiber.Ctx) error {
	return authHandler.AuthService.ResendOTP(c)
}

// ForgotPassword godoc
// @Summary Forgot password
// @Description Initiate forgot password process
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body types.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {string} string "Password reset initiated"
// @Failure 400 {string} string "Invalid request"
// @Router /api/v1/password/forgot [post]
func (authHandler *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req types.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	if err := validateRequest(c, &req); err != nil {
		return err
	}

	return c.SendString("Password reset initiated")
}

// ConfirmForgotPassword godoc
// @Summary Confirm forgot password
// @Description Confirm forgot password with confirmation code
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body types.ConfirmForgotPasswordRequest true "Confirm forgot password request"
// @Success 200 {string} string "Password reset successfully"
// @Failure 400 {string} string "Invalid request"
// @Router /api/v1/password/forgot/confirm [post]
func (authHandler *AuthHandler) ConfirmForgotPassword(c *fiber.Ctx) error {
	var req types.ConfirmForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	if err := validateRequest(c, &req); err != nil {
		return err
	}

	return c.SendString("Password reset successfully")
}

// ConfirmTempPassword godoc
// @Summary Confirm temporary password
// @Description Confirm temporary password with new password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body types.ConfirmTempPasswordRequest true "Confirm temporary password request"
// @Success 200 {string} string "Password reset successfully"
// @Failure 400 {string} string "Invalid request"
// @Router /api/v1/password/temp/confirm [post]
func (authHandler *AuthHandler) ConfirmTempPassword(c *fiber.Ctx) error {
	var req types.ConfirmTempPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	if err := validateRequest(c, &req); err != nil {
		return err
	}

	return c.SendString("Password reset successfully")
}

// ChangePassword godoc
// @Summary Change password
// @Description Change user password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body types.ChangePasswordRequest true "Change password request"
// @Success 200 {string} string "Password changed successfully"
// @Failure 400 {string} string "Invalid request"
// @Router /api/v1/password [post]
func (authHandler *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	var req types.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	if err := validateRequest(c, &req); err != nil {
		return err
	}

	return c.SendString("Password changed successfully")
}

// RegisterUser godoc
// @Summary Register user
// @Description Register a new user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body types.RegisterRequest true "Register request"
// @Success 201 {string} string "User registered successfully"
// @Failure 400 {string} string "Invalid request"
// @Router /api/v1/register [post]
func (authHandler *AuthHandler) RegisterUser(c *fiber.Ctx) error {
	return authHandler.AuthService.Register(c)
}
