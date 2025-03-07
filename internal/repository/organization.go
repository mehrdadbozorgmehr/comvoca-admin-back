package repository

import (
	"fmt"

	"github.com/Comvoca-AI/comvoca-admin-back/internal/entity"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepo(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (dao *OrganizationRepository) GetById(organizationId string) (entity.Organization, error) {
	var organization entity.Organization

	// Use GORM to query by ID
	tx := dao.db.First(&organization, "id = ?", organizationId)

	// Check for errors
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return organization, fmt.Errorf("organization not found")
		}
	}
	return organization, tx.Error
}

func (dao *OrganizationRepository) Save(tx *gorm.DB, org *entity.Organization) error {
	return tx.Create(org).Error
}

func (dao *OrganizationRepository) Update(org entity.Organization) error {
	return dao.db.Save(org).Error
}
