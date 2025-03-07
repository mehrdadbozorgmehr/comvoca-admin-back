package service

import (
	"context"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/Comvoca-AI/comvoca-admin-back/config"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/entity"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/errors"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/logger"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/validator"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type SecurityService struct {
	PublicKey           *JWK
	db                  *gorm.DB
	userService         *UserService
	organizationService *OrganizationService
}

func NewSecurityService(db *gorm.DB, userService *UserService, organizationService *OrganizationService) *SecurityService {
	publicKey, err := fetchPublicKey()
	if err != nil {
		logger.Error("Failed to fetch public key", err)
		return nil
	}
	return &SecurityService{PublicKey: publicKey, db: db, userService: userService, organizationService: organizationService}
}

func validateRequest(c *fiber.Ctx, req interface{}) error {
	v := validator.GetValidator()
	if err := v.Struct(req); err != nil {
		logger.Error("Request validation failed", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return nil
}

func calculateSecretHash(username, clientId, clientSecret string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
func getCognitoClient(c *fiber.Ctx) (error, *cognitoidentityprovider.Client, error, bool) {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(config.AppConfig.Auth.Cognito.Region))
	if err != nil {
		return nil, nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load AWS config"}), true
	}

	// Create a Cognito client
	client := cognitoidentityprovider.NewFromConfig(cfg)
	return err, client, nil, false
}

func (authService *SecurityService) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err2 := validateRequest(c, &req); err2 != nil {
		return err2
	}

	secretHash := calculateSecretHash(req.Email, config.AppConfig.Auth.Cognito.ClientId, config.AppConfig.Auth.Cognito.Secret)

	err, client, err2, done := getCognitoClient(c)
	if done {
		return err2
	}

	var signUpInput = &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(config.AppConfig.Auth.Cognito.ClientId),
		Username:   aws.String(req.Email),
		Password:   aws.String(req.Password),
		SecretHash: aws.String(secretHash),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(req.Email),
			},
		},
	}

	// Step 1: Start GORM transaction
	transactionError := authService.db.Debug().Transaction(func(tx *gorm.DB) error {
		// Step 2: Create Organization
		var org = entity.Organization{
			Email: req.Email,
			Name:  "", // Default values
		}

		if err := authService.organizationService.SaveOrganization(tx, &org); err != nil {
			return fmt.Errorf("failed to create organization: %w", err)
		}

		// Step 3: Create User linked to Organization
		var user = entity.User{
			Email:          &req.Email,
			OrganizationID: org.ID, // Link to organization
		}

		if err := authService.userService.SaveUser(tx, &user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// Step 4: Commit transaction automatically if no error
		return nil
	})

	if transactionError != nil {
		logger.Error("Failed to register user", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
	}

	_, err = client.SignUp(context.TODO(), signUpInput)
	if err != nil {
		logger.Error("Failed to register user", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

func (authService *SecurityService) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := validateRequest(c, &req); err != nil {
		return err
	}

	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(config.AppConfig.Auth.Cognito.Region))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load AWS config"})
	}

	// Create a Cognito client
	client := cognitoidentityprovider.NewFromConfig(cfg)

	// Calculate the secret hash
	secretHash := calculateSecretHash(req.Email, config.AppConfig.Auth.Cognito.ClientId, config.AppConfig.Auth.Cognito.Secret)

	// Validate the user credentials with Cognito
	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		AuthParameters: map[string]string{
			"USERNAME":    req.Email,
			"PASSWORD":    req.Password,
			"SECRET_HASH": secretHash,
		},
		ClientId: aws.String(config.AppConfig.Auth.Cognito.ClientId),
	}

	_, err = client.InitiateAuth(context.TODO(), authInput)

	if err != nil {
		logger.Error("Authentication failed", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	claims := jwt.MapClaims{
		"user_name": req.Email,
		"exp":       time.Now().Add(time.Hour * time.Duration(config.AppConfig.Auth.Cognito.TokenExpireHour)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.AppConfig.Auth.Cognito.JwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(LoginResponse{Token: tokenString})
}

func (authService *SecurityService) ActivateOTP(c *fiber.Ctx) error {
	var req ActivateOTPRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := validateRequest(c, &req); err != nil {
		return err
	}

	secretHash := calculateSecretHash(req.Email, config.AppConfig.Auth.Cognito.ClientId, config.AppConfig.Auth.Cognito.Secret)

	// Load the AWS configuration
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(config.AppConfig.Auth.Cognito.Region))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load AWS config"})
	}

	// Create a Cognito client
	client := cognitoidentityprovider.NewFromConfig(cfg)

	// Verify the OTP
	confirmSignUpInput := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(config.AppConfig.Auth.Cognito.ClientId),
		Username:         aws.String(req.Email),
		ConfirmationCode: aws.String(req.OTP),
		SecretHash:       aws.String(secretHash),
	}

	_, err = client.ConfirmSignUp(context.TODO(), confirmSignUpInput)
	if err != nil {
		logger.Error("Failed to confirm sign up", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to confirm sign up"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User activated successfully"})
}

func (authService *SecurityService) ResendOTP(c *fiber.Ctx) error {
	var req ResendOTPRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := validateRequest(c, &req); err != nil {
		return err
	}

	// Calculate the secret hash
	secretHash := calculateSecretHash(req.Email, config.AppConfig.Auth.Cognito.ClientId, config.AppConfig.Auth.Cognito.Secret)

	// Load the AWS configuration
	err, client, err2, done := getCognitoClient(c)
	if done {
		return err2
	}

	// Resend the confirmation code
	resendInput := &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId:   aws.String(config.AppConfig.Auth.Cognito.ClientId),
		Username:   aws.String(req.Email),
		SecretHash: aws.String(secretHash),
	}

	_, err = client.ResendConfirmationCode(context.TODO(), resendInput)
	if err != nil {
		logger.Error("Failed to resend confirmation code", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to resend confirmation code"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OTP resent successfully"})
}

func (authService *SecurityService) ForgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := validateRequest(c, &req); err != nil {
		return err
	}

	// Calculate the secret hash
	secretHash := calculateSecretHash(req.Email, config.AppConfig.Auth.Cognito.ClientId, config.AppConfig.Auth.Cognito.Secret)

	// Load the AWS configuration
	err, client, err2, done := getCognitoClient(c)
	if done {
		return err2
	}

	// Initiate forgot password request
	forgotPasswordInput := &cognitoidentityprovider.ForgotPasswordInput{
		ClientId:   aws.String(config.AppConfig.Auth.Cognito.ClientId),
		Username:   aws.String(req.Email),
		SecretHash: aws.String(secretHash),
	}

	_, err = client.ForgotPassword(context.TODO(), forgotPasswordInput)
	if err != nil {
		logger.Error("Failed to initiate forgot password", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to initiate forgot password"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password reset initiated"})
}

func (authService *SecurityService) ConfirmForgotPassword(c *fiber.Ctx) error {
	var req ConfirmForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := validateRequest(c, &req); err != nil {
		return err
	}

	// Calculate the secret hash
	secretHash := calculateSecretHash(req.Email, config.AppConfig.Auth.Cognito.ClientId, config.AppConfig.Auth.Cognito.Secret)

	// Load the AWS configuration
	err, client, err2, done := getCognitoClient(c)
	if done {
		return err2
	}

	// Confirm forgot password request
	confirmForgotPasswordInput := &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(config.AppConfig.Auth.Cognito.ClientId),
		Username:         aws.String(req.Email),
		ConfirmationCode: aws.String(req.ConfirmationCode),
		Password:         aws.String(req.NewPassword),
		SecretHash:       aws.String(secretHash),
	}

	_, err = client.ConfirmForgotPassword(context.TODO(), confirmForgotPasswordInput)
	if err != nil {
		logger.Error("Failed to confirm forgot password", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to confirm forgot password"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password reset successfully"})
}

func (authService *SecurityService) ChangePassword(c *fiber.Ctx) error {
	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := validateRequest(c, &req); err != nil {
		return err
	}

	err, client, err2, done := getCognitoClient(c)
	if done {
		return err2
	}

	changePasswordInput := &cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      aws.String(req.AccessToken),
		PreviousPassword: aws.String(req.OldPassword),
		ProposedPassword: aws.String(req.NewPassword),
	}

	_, err = client.ChangePassword(context.TODO(), changePasswordInput)
	if err != nil {
		logger.Error("Failed to change password", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to change password"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password changed successfully"})
}

func fetchPublicKey() (*JWK, error) {
	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", config.AppConfig.Auth.Cognito.Region, config.AppConfig.Auth.Cognito.UserPoolId)
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("Failed to close response body", err)
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jwk JWK
	err = json.Unmarshal(body, &jwk)
	if err != nil {
		return nil, err
	}

	return &jwk, nil
}

func getRSAPublicKey(jwk *JWK, kid string) (*rsa.PublicKey, error) {
	for _, key := range jwk.Keys {
		if key.Kid == kid {
			nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
			if err != nil {
				return nil, err
			}
			eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
			if err != nil {
				return nil, err
			}
			e := big.NewInt(0).SetBytes(eBytes).Int64()
			pubKey := &rsa.PublicKey{
				N: big.NewInt(0).SetBytes(nBytes),
				E: int(e),
			}
			return pubKey, nil
		}
	}
	return nil, errors.InternalServerError("Public key not found")
}

func (authService *SecurityService) ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			logger.Error("Unexpected signing method", token.Header["alg"])
			return nil, errors.BadRequest(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}
		kid, ok := token.Header["user_id"].(string)
		if !ok {
			return nil, errors.BadRequest("user_id not found in token header")
		}
		return getRSAPublicKey(authService.PublicKey, kid)
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
