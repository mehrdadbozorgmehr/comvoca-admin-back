package entity

type NotificationType int

const (
	Email NotificationType = iota
	SMS
	InApp
)

type Organization struct {
	Base
	Email                string
	Name                 string
	Website              string
	PhoneNumber          string
	InsuranceCompany     string             `json:"insuranceCompany"`
	Specialities         []Speciality       `gorm:"many2many:organization_specialities;" json:"specialities,omitempty"`
	DailySchedules       []DailySchedule    `gorm:"many2many:organization_dailey_schedules;" json:"daileySchedules,omitempty"`
	Notifications        []NotificationType `gorm:"type:integer[]" json:"notifications,omitempty"`
	CallForwardingNumber string
	Users                []User `gorm:"foreignKey:OrganizationID" json:"users,omitempty"`
}
