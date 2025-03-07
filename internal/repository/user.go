package repository

import (
	"fmt"

	"github.com/Comvoca-AI/comvoca-admin-back/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (UserRepository *UserRepository) GetById(id string) (entity.User, error) {
	var user entity.User

	// Use GORM to query by ID
	tx := UserRepository.db.First(&user, "id = ?", id)

	// Check for errors
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return user, fmt.Errorf("organization not found")
		}
	}
	return user, tx.Error
}

func (dao *UserRepository) Save(tx *gorm.DB, user *entity.User) error {
	return tx.Create(user).Error
}