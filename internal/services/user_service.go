package services

import (
	"github.com/yanggu0t/go-rdbms-practice/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.DB.Create(user).Error
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := s.DB.Preload("Projects").First(&user, id).Error
	return &user, err
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}
