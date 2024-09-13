package services

import (
	"strconv"

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
	err := s.DB.First(&user, id).Error
	return &user, err
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (s *UserService) GetUserProjects(userID uint) ([]models.Project, error) {
	var projects []models.Project
	err := s.DB.Joins("JOIN project_user_roles ON projects.id = project_user_roles.project_id").
		Where("project_user_roles.user_id = ?", userID).
		Find(&projects).Error
	return projects, err
}

func (s *UserService) AddUserToProject(tx *gorm.DB, userID, projectID, roleID uint) error {
	projectUserRole := models.ProjectUserRole{
		ProjectID: projectID,
		UserID:    strconv.Itoa(int(userID)),
		RoleID:    roleID,
	}
	return tx.Create(&projectUserRole).Error
}

func (s *UserService) UpdateUserRole(tx *gorm.DB, userID, projectID, newRoleID uint) error {
	return tx.Model(&models.ProjectUserRole{}).
		Where("user_id = ? AND project_id = ?", userID, projectID).
		Update("role_id", newRoleID).Error
}

func (s *UserService) GetUserRoles(userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := s.DB.Joins("JOIN project_user_roles ON roles.id = project_user_roles.role_id").
		Where("project_user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}
