// Package services
package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yanggu0t/go-rdbms-practice/internal/models"
	"golang.org/x/crypto/bcrypt"
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

func (s *AuthService) Login(usernameOrEmail, password string) (string, *models.User, error) {
	var user models.User
	if err := s.DB.Preload("Roles").Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user).Error; err != nil {
		return "", nil, errors.New(ErrUserNotFound)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New(ErrPasswordIncorrect)
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	expirationUnix := expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
		"exp":     expirationUnix,
	})

	tokenString, err := token.SignedString(s.Secret)
	if err != nil {
		return "", nil, err
	}

	// 設置用戶的過期時間
	user.Expires = expirationUnix

	// 更新數據庫中的用戶資料
	if err := s.DB.Save(&user).Error; err != nil {
		return "", nil, err
	}

	return tokenString, &user, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*models.User, int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.Secret, nil
	})
	if err != nil {
		return nil, -1, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, -1, errors.New("error_invalid_user_id")
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, -1, errors.New("error_get_expiration_time")
		}

		expTimestamp := int64(exp)

		// 查找用戶資料
		user, err := s.GetUserByID(userID)
		if err != nil {
			return nil, -1, err
		}

		// 將過期時間插入到用戶資料中
		user.Expires = expTimestamp

		return user, expTimestamp, nil
	}

	return nil, -1, errors.New("error_invalid_token")
}

func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := s.DB.Preload("Roles").Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	fmt.Printf("用戶: %+v\n", user)
	return &user, nil
}
