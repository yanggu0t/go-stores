package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yanggu0t/go-rdbms-practice/internal/models"
	"gorm.io/gorm"
)

const (
	ErrUserNotFound      = "error_user_not_found"
	ErrPasswordIncorrect = "error_password_incorrect"
)

type AuthService struct {
	DB     *gorm.DB
	Secret []byte
}

func NewAuthService(db *gorm.DB, secret string) *AuthService {
	return &AuthService{
		DB:     db,
		Secret: []byte(secret),
	}
}

func (s *AuthService) Login(account, password string) (string, *models.User, error) {
	var user models.User
	if err := s.DB.Where("username = ? OR email = ?", account, account).First(&user).Error; err != nil {
		return "", nil, errors.New(ErrUserNotFound)
	}

	if !user.CheckPassword(password) {
		return "", nil, errors.New(ErrPasswordIncorrect)
	}

	// 生成新的 token
	token, expirationUnix, err := s.GenerateToken(user.UserID)
	if err != nil {
		return "", nil, err
	}

	// 創建或更新 UserSession
	session := models.UserSession{
		UserID:  user.UserID,
		Token:   token,
		Expires: expirationUnix,
	}
	if err := s.DB.Save(&session).Error; err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

func (s *AuthService) Logout(tokenString string) error {
	session, err := s.ValidateToken(tokenString)
	if err != nil {
		return err
	}

	// 刪除 session
	if err := s.DB.Delete(&session).Error; err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ValidateToken(tokenString string) (*models.UserSession, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.Secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		_, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.New("error_invalid_user_id")
		}

		// 查找 session
		var session models.UserSession
		if err := s.DB.Where("token = ?", tokenString).Preload("User").First(&session).Error; err != nil {
			return nil, err
		}

		// 檢查 token 是否過期
		if time.Now().Unix() > session.Expires {
			return nil, errors.New("error_token_expired")
		}

		return &session, nil
	}

	return nil, errors.New("error_invalid_token")
}

func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	fmt.Printf("用戶: %+v\n", user)
	return &user, nil
}

func (s *AuthService) GenerateToken(userID string) (string, int64, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	expirationUnix := expirationTime.Unix()

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationUnix,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.Secret)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expirationUnix, nil
}

func (s *AuthService) CheckProjectPermission(userID uint, projectID string, requiredPermission string) (bool, error) {
	var projectUserRoles []models.ProjectUserRole
	if err := s.DB.Where("user_id = ? AND project_id = ?", userID, projectID).Preload("Role").Preload("Role.Permissions").Find(&projectUserRoles).Error; err != nil {
		return false, err
	}

	for _, pur := range projectUserRoles {
		for _, permission := range pur.Role.Permissions {
			if permission.Code == requiredPermission {
				return true, nil
			}
		}
	}

	return false, nil
}
