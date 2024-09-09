package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yanggu0t/go-rdbms-practice/internal/services"
	"github.com/yanggu0t/go-rdbms-practice/internal/utils"
)

func AuthMiddleware(authService *services.AuthService, requiredRole ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Response(c, http.StatusUnauthorized, "error", "error_no_auth_token", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, expTime, err := authService.ValidateToken(tokenString)
		if err != nil {
			utils.Response(c, http.StatusUnauthorized, "error", "error_invalid_auth_token", nil)

			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			utils.Response(c, http.StatusUnauthorized, "error", "error_invalid_user_id", nil)
			c.Abort()
			return
		}

		user, err := authService.GetUserByID(userID)
		if err != nil {
			utils.Response(c, http.StatusUnauthorized, "error", "error_invalid_user_id", nil)
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
				utils.Response(c, http.StatusForbidden, "error", "error_unauthorized_access", nil)
				c.Abort()
				return
			}
		}

		c.Set("user", user)
		c.Set("expTime", expTime)
		c.Next()
	}
}
