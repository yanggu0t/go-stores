package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yanggu0t/go-rdbms-practice/internal/database"
	"github.com/yanggu0t/go-rdbms-practice/internal/models"
	"gorm.io/gorm"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		fmt.Println("收到的 X-User-ID 頭部:", userID)

		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供用戶 ID"})
			c.Abort()
			return
		}

		var user models.User
		err := database.DB.Where("user_id = ?", userID).Preload("Roles").First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				fmt.Printf("用戶 ID 不存在: %s\n", userID)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的用戶 ID"})
			} else {
				fmt.Printf("數據庫查詢錯誤: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "服務器錯誤"})
			}
			c.Abort()
			return
		}

		fmt.Printf("找到用戶: %+v\n", user)
		c.Set("user", user)
		c.Next()
	}
}
