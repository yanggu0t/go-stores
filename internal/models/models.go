package models

import (
	"crypto/rand"
	"encoding/base32"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ProjectID   string `gorm:"type:char(12);unique;not null"`
	Name        string `gorm:"not null"`
	Description string
}

type Role struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Code        string `gorm:"not null"`
}

type User struct {
	gorm.Model
	UserID   string `gorm:"type:char(8);unique;not null;primaryKey"`
	Username string `gorm:"unique;not null;required"`
	Email    string `gorm:"unique;not null;required"`
	Password string `gorm:"not null;required"`
}

type ProjectUserRole struct {
	ProjectID uint   `gorm:"primaryKey"`
	UserID    string `gorm:"primaryKey;type:char(8)"`
	RoleID    uint   `gorm:"primaryKey"`
	Project   Project
	User      User `gorm:"foreignKey:UserID"`
	Role      Role
}

type ProjectPermission struct {
	ProjectID    uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
	Project      Project
	Permission   Permission
}

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
	Role         Role
	Permission   Permission
}

type UserSession struct {
	UserID  string `gorm:"primaryKey;type:char(8)"`
	Token   string `gorm:"not null"`
	Expires int64  `gorm:"not null"`
	User    User   `gorm:"foreignKey:UserID"`
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

// GenerateID 生成一個 12 位的英文混數字 ID
func GenerateID() (string, error) {
	bytes := make([]byte, 9) // 9 bytes will give us 12 characters in base32
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	encoded := base32.StdEncoding.EncodeToString(bytes)
	id := strings.ToLower(encoded[:12]) // 轉換為小寫並取前 12 位
	return id, nil
}
