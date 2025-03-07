package entity

import "github.com/google/uuid"

type UserRole int

const (
	Admin UserRole = iota
	Staff
	Support
)

type User struct {
	Base
	Email          *string
	Name           *string
	Family         *string
	PhoneNumber    *string
	Role           *UserRole
	DailySchedules []DailySchedule `gorm:"many2many:user_dailey_schedules;" json:"daileySchedules,omitempty"`
	OrganizationID uuid.UUID       `gorm:"type:uuid"`
	Organization   Organization
}
