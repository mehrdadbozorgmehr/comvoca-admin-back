package service

import (
	"github.com/Comvoca-AI/comvoca-admin-back/internal/entity"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/repository"
	"gorm.io/gorm"
)

type UserService struct {
	dao *repository.UserRepository
}

func NewUserService(dao *repository.UserRepository) *UserService {
	return &UserService{dao: dao}
}

func (s *UserService) SaveUser(tx *gorm.DB, user *entity.User) error {
	return s.dao.Save(tx, user)
}

func (s *UserService) GetUserById(id string) (entity.User, error) {
	return s.dao.GetById(id)
}
