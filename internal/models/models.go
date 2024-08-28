package models

import (
	"fmt"

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
}

type User struct {
	gorm.Model
	UserID   string `gorm:"type:char(8);unique;not null"`
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Roles    []Role `gorm:"many2many:user_roles;"`
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

// AddRole 為用戶添加角色
func (u *User) AddRole(db *gorm.DB, role *Role) error {
	return db.Model(u).Association("Roles").Append(role)
}

// RemoveRole 從用戶中移除角色
func (u *User) RemoveRole(db *gorm.DB, role *Role) error {
	return db.Model(u).Association("Roles").Delete(role)
}

// HasRole 檢查用戶是否擁有特定角色
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			fmt.Printf("111111")
			return true
		}
	}
	fmt.Printf("00000")
	return false
}
