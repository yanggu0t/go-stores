package services

import (
	"github.com/yanggu0t/go-rdbms-practice/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) SetPassword(user *models.User, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

func (s *UserService) CheckPassword(user *models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (s *UserService) AddRole(user *models.User, role *models.Role) error {
	return s.DB.Model(user).Association("Roles").Append(role)
}

func (s *UserService) RemoveRole(user *models.User, role *models.Role) error {
	return s.DB.Model(user).Association("Roles").Delete(role)
}

// HasRole 檢查用戶是否擁有特定角色
func (s *UserService) HasRole(user *models.User, roleName string) bool {
	if user == nil {
		return false
	}
	for _, role := range user.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	err := s.DB.Where("user_id = ?", userID).Preload("Roles").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
