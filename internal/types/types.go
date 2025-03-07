package types

import (
	"time"

	"github.com/Comvoca-AI/comvoca-admin-back/internal/entity"
	"github.com/google/uuid"
)

type UserRequest struct {
	Id             uuid.UUID       `json:"id"`
	Email          string          `json:"email"`
	Name           string          `json:"name"`
	Family         string          `json:"family"`
	PhoneNumber    string          `json:"phone_number"`
	Role           string          `json:"role"`
	DailySchedules []DailySchedule `json:"daily_schedules"`
}

type UserResponse struct {
	Id             uuid.UUID       `json:"id"`
	Email          string          `json:"email"`
	Name           string          `json:"name"`
	Family         string          `json:"family"`
	PhoneNumber    string          `json:"phone_number"`
	Role           string          `json:"role"`
	DailySchedules []DailySchedule `json:"daily_schedules"`
}

type DailySchedule struct {
	DayOfWeek int       `json:"day_of_week"`
	FromTime  time.Time `json:"from_time"`
	ToTime    time.Time `json:"to_time"`
}

type OrganizationRequest struct {
	Email                string                    `json:"email"`
	Name                 string                    `json:"name"`
	Website              string                    `json:"website"`
	PhoneNumber          string                    `json:"phone_number"`
	InsuranceCompany     string                    `json:"insuranceCompany"`
	Specialities         []int                     `json:"specialities"`
	DailySchedules       []DailySchedule           `json:"dailySchedules"`
	Notifications        []entity.NotificationType `json:"notifications"`
	CallForwardingNumber string                    `json:"callForwadingNumber"`
}

type OrganizationResponse struct {
	Id                   *uuid.UUID                 `json:"id"`
	Email                *string                    `json:"email"`
	Name                 *string                    `json:"name"`
	Website              *string                    `json:"website"`
	PhoneNumber          *string                    `json:"phone_number"`
	Specialities         *[]int                     `json:"specialities"`
	DailySchedules       *[]DailySchedule           `json:"daily_schedules"`
	Notifications        *[]entity.NotificationType `json:"notifications"`
	CallForwardingNumber *string                    `json:"call_forwarding_number"`
}

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type ActivateOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required"`
}

type ResendOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ConfirmTempPasswordRequest struct {
	Email             string `json:"email" validate:"required,email"`
	TemporaryPassword string `json:"temporary_password" validate:"required"`
	NewPassword       string `json:"new_password" validate:"required,min=8"`
}

type ConfirmForgotPasswordRequest struct {
	Email            string `json:"email" validate:"required,email"`
	ConfirmationCode string `json:"confirmation_code" validate:"required"`
	NewPassword      string `json:"new_password" validate:"required,min=8"`
}

type ChangePasswordRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
	OldPassword string `json:"old_password" validate:"required,min=8"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
