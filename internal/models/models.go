package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string
	Code        string `gorm:"unique;not null"`
}

type User struct {
	gorm.Model
	UserID   string `gorm:"type:char(8);unique;not null"`
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Roles    []Role `gorm:"many2many:user_roles;"`
	Expires  int64  `gorm:"not null;default:0"`
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
