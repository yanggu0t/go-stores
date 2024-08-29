package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yanggu0t/go-rdbms-practice/internal/database"
	"github.com/yanggu0t/go-rdbms-practice/internal/handlers"
	"github.com/yanggu0t/go-rdbms-practice/internal/middleware"
	"github.com/yanggu0t/go-rdbms-practice/internal/services"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 配置 CORS
	config := cors.Config{
		AllowOrigins:     []string{"*"}, // 允許的前端域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-User-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(config))

	// 创建 UserService 实例
	userService := services.NewUserService(database.DB)

	handler := handlers.NewHandler(database.DB)

	// 添加身份驗證中間件，传入 UserService
	r.Use(middleware.AuthMiddleware(userService))

	// 創建一個需要管理員權限的路由組
	admin := r.Group("/")
	admin.Use(middleware.RoleMiddleware(userService, "admin"))

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

	// 設置不需要特殊權限的路由

	return r
}
