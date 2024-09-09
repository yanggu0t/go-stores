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

func (s *AuthService) Login(username, password string) (string, error) {
	var user models.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("用戶不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("密碼錯誤")
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
		"exp":     expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(s.Secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.Secret, nil
	})
	if err != nil {
		return nil, -1, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			expTimestamp := int64(exp)
			return claims, expTimestamp, nil
		}
		return claims, -1, errors.New("無法獲取過期時間")
	}

	return nil, time.Now().Unix(), errors.New("無效的令牌")
}

func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := s.DB.Preload("Roles").Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	fmt.Printf("用戶: %+v\n", user)
	return &user, nil
}
