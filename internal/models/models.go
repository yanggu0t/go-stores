package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string
	Users       []User `gorm:"many2many:project_users;"`
	Roles       []Role `gorm:"foreignKey:ProjectID"`
}

type Role struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	ProjectID   uint
	Project     Project      `gorm:"foreignKey:ProjectID"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
	Users       []User       `gorm:"many2many:project_user_roles;"`
}

type Permission struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Code        string `gorm:"not null"`
	Roles       []Role `gorm:"many2many:role_permissions;"`
}

type User struct {
	gorm.Model
	UserID   string    `gorm:"type:char(8);unique;not null"`
	Username string    `gorm:"unique;not null"`
	Email    string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	Projects []Project `gorm:"many2many:project_users;"`
	Roles    []Role    `gorm:"many2many:project_user_roles;"`
	Expires  int64     `gorm:"not null;default:0"`
}

// ProjectUserRole 是一個中間表，用於關聯用戶、專案和角色
type ProjectUserRole struct {
	ProjectID uint `gorm:"primaryKey"`
	UserID    uint `gorm:"primaryKey"`
	RoleID    uint `gorm:"primaryKey"`
	Project   Project
	User      User
	Role      Role
}

// SetPassword 加密並設置用戶密碼
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 驗證用戶密碼
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
