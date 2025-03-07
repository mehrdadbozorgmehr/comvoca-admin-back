package service

import (
	"github.com/Comvoca-AI/comvoca-admin-back/internal/entity"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/repository"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/types"
	"gorm.io/gorm"
)

type OrganizationService struct {
	db  *gorm.DB
	dao *repository.OrganizationRepository
}

func NewOrganizationService(db *gorm.DB, dao *repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{db: db, dao: dao}
}

func (s *OrganizationService) SaveOrganization(tx *gorm.DB, organization *entity.Organization) error {
	return s.dao.Save(tx, organization)
}

func (s *OrganizationService) GetOrganizationById(organizationId string) (*entity.Organization, error) {
	org, err := s.dao.GetById(organizationId)
	if err != nil {
		return nil, err
	}
	return &org, nil

}

func (s *OrganizationService) UpdateOrganization(id string, dto types.OrganizationRequest) (*entity.Organization, error) {
	var organization entity.Organization
	if err := s.db.Preload("Specialities").Preload("DailySchedules").Preload("Users").
		First(&organization, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Update basic fields
	organization.Email = dto.Email
	organization.Name = dto.Name
	organization.Website = dto.Website
	organization.PhoneNumber = dto.PhoneNumber
	organization.InsuranceCompany = dto.InsuranceCompany
	organization.CallForwardingNumber = dto.CallForwardingNumber
	organization.Notifications = dto.Notifications

	// Fetch and Assign Specialities
	var specialities []entity.Speciality
	if len(dto.Specialities) > 0 {
		if err := s.db.Where("id IN ?", dto.Specialities).Find(&specialities).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.Model(&organization).Association("Specialities").Replace(specialities); err != nil {
		return nil, err
	}

	// Create and Assign DailySchedules
	var dailySchedules []entity.DailySchedule
	for _, schedule := range dto.DailySchedules {
		newSchedule := entity.DailySchedule{
			DayOfWeek: schedule.DayOfWeek,
			FromTime:  schedule.FromTime,
			ToTime:    schedule.ToTime,
		}
		if err := s.db.Create(&newSchedule).Error; err != nil {
			return nil, err
		}
		dailySchedules = append(dailySchedules, newSchedule)
	}
	if err := s.db.Model(&organization).Association("DailySchedules").Replace(dailySchedules); err != nil {
		return nil, err
	}

	// Save the organization
	if err := s.db.Save(&organization).Error; err != nil {
		return nil, err
	}

	return &organization, nil
}
