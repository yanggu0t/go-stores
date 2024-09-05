package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yanggu0t/go-rdbms-practice/internal/services"
)

func AuthMiddleware(authService *services.AuthService, requiredRole ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供授權令牌"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的授權令牌"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的用戶 ID"})
			c.Abort()
			return
		}

		user, err := authService.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的用戶"})
			c.Abort()
			return
		}

		if len(requiredRole) > 0 && requiredRole[0] != "" {
			hasRole := false
			for _, role := range user.Roles {
				if strings.EqualFold(role.Name, requiredRole[0]) {
					hasRole = true
					break
				}
			}

			if !hasRole {
				c.JSON(http.StatusForbidden, gin.H{"error": "未授權訪問"})
				c.Abort()
				return
			}
		}

		c.Set("user", user)
		c.Next()
	}
}
