package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yanggu0t/go-rdbms-practice/config"
	"github.com/yanggu0t/go-rdbms-practice/internal/handlers"
	"github.com/yanggu0t/go-rdbms-practice/internal/middleware"
	"github.com/yanggu0t/go-rdbms-practice/internal/services"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// 配置 CORS
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"}, // 允許的前端域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	// 創建服務實例
	authService := services.NewAuthService(db, cfg.JWTSecret)
	handler := handlers.NewHandler(db, cfg)

	// 設置不需要特殊權限的路由
	r.POST("/login", handler.LoginHandler)

	// 使用 AuthMiddleware
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware(authService))

	// 設置需要管理員權限的路由
	admin := authorized.Group("/admin")
	admin.Use(middleware.AuthMiddleware(authService, "admin"))

	// 設置需要管理員權限的路由
	admin.GET("/users", handler.GetAllUsersHandler)
	admin.GET("/roles", handler.GetAllRolesHandler)
	admin.GET("/permissions", handler.GetAllPermissionsHandler)
	admin.POST("/roles", handler.CreateRoleHandler)
	admin.POST("/users", handler.CreateUserHandler)
	admin.POST("/permissions", handler.CreatePermissionHandler)
	admin.PUT("/roles/:id", handler.UpdateRoleHandler)
	admin.PUT("/users/:id", handler.UpdateUserHandler)
	admin.PUT("/permissions/:id", handler.UpdatePermissionHandler)

	return r
}
