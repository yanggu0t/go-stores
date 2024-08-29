package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yanggu0t/go-rdbms-practice/internal/models"
	"github.com/yanggu0t/go-rdbms-practice/internal/services"
	"gorm.io/gorm"
)

func AuthMiddleware(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")

		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供用戶 ID"})
			c.Abort()
			return
		}

		user, err := userService.GetUserByID(userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的用戶 ID"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "服務器錯誤"})
			}
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// RoleMiddleware 是一個中間件，用於檢查用戶是否具有特定角色
func RoleMiddleware(userService *services.UserService, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用戶未登錄"})
			c.Abort()
			return
		}

		typedUser, ok := user.(*models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無效的用戶數據"})
			c.Abort()
			return
		}

		if !userService.HasRole(typedUser, role) {
			c.JSON(http.StatusForbidden, gin.H{"error": "未授權訪問"})
			c.Abort()
			return
		}

		c.Next()
	}
}
