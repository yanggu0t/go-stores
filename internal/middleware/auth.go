package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yanggu0t/go-rdbms-practice/internal/models"
	"github.com/yanggu0t/go-rdbms-practice/internal/services"
	"github.com/yanggu0t/go-rdbms-practice/internal/utils"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Response(c, http.StatusUnauthorized, "error", "error_no_auth_token", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		session, err := authService.ValidateToken(tokenString)
		if err != nil {
			utils.Response(c, http.StatusUnauthorized, "error", "error_invalid_auth_token", nil)
			c.Abort()
			return
		}

		c.Set("user", session.User)
		c.Set("expTime", session.Expires)
		c.Next()
	}
}

// 新增一個檢查項目權限的中間件
func ProjectPermissionMiddleware(authService *services.AuthService, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		projectID := c.Param("projectID") // 假設項目ID在URL參數中

		hasPermission, err := authService.CheckProjectPermission(user.ID, projectID, requiredPermission)
		if err != nil {
			utils.Response(c, http.StatusInternalServerError, "error", "error_checking_permission", nil)
			c.Abort()
			return
		}

		if !hasPermission {
			utils.Response(c, http.StatusForbidden, "error", "error_unauthorized_access", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
