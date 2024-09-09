// Package services
package services

import (
	"github.com/yanggu0t/go-rdbms-practice/internal/models"

	"gorm.io/gorm"
)

type RoleService struct {
	DB *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{DB: db}
}

func (s *RoleService) AddPermission(role *models.Role, permission *models.Permission) error {
	return s.DB.Model(role).Association("Permissions").Append(permission)
}

func (s *RoleService) RemovePermission(role *models.Role, permission *models.Permission) error {
	return s.DB.Model(role).Association("Permissions").Delete(permission)
}
