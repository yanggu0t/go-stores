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

	// i18n Middleware
	r.Use(middleware.LanguageMiddleware())

	// 創建服務實例
	userService := services.NewUserService(db)
	authService := services.NewAuthService(db, cfg.JWTSecret)
	// projectService := services.NewProjectService(db)

	// 創建處理器實例
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	// projectHandler := handlers.NewProjectHandler(projectService)

	api := r.Group("/")
	{
		// 公開路由
		api.POST("/users", userHandler.CreateUser)
		api.POST("/login", authHandler.Login)

		// 需要認證的路由
		authorized := api.Group("/")
		authorized.Use(middleware.AuthMiddleware(authService))
		// {
		// 	// 專案相關路由
		// 	projects := authorized.Group("/projects")
		// 	{
		// 		projects.GET("", projectHandler.GetAllProjects)
		// 		projects.POST("", projectHandler.CreateProject)
		// 		projects.GET("/:id", projectHandler.GetProject)
		// 		projects.PUT("/:id", projectHandler.UpdateProject)
		// 		projects.DELETE("/:id", projectHandler.DeleteProject)

		// 		// 專案內的用戶管理
		// 		projects.GET("/:id/users", projectHandler.GetProjectUsers)
		// 		projects.POST("/:id/users", projectHandler.AddUserToProject)
		// 		projects.DELETE("/:id/users/:userId", projectHandler.RemoveUserFromProject)

		// 		// 專案內的角色管理
		// 		projects.GET("/:id/roles", projectHandler.GetProjectRoles)
		// 		projects.POST("/:id/roles", projectHandler.CreateProjectRole)
		// 		projects.PUT("/:id/roles/:roleId", projectHandler.UpdateProjectRole)
		// 		projects.DELETE("/:id/roles/:roleId", projectHandler.DeleteProjectRole)

		// 		// 專案內的權限管理
		// 		projects.GET("/:id/permissions", projectHandler.GetProjectPermissions)
		// 		projects.POST("/:id/permissions", projectHandler.CreateProjectPermission)
		// 		projects.PUT("/:id/permissions/:permissionId", projectHandler.UpdateProjectPermission)
		// 		projects.DELETE("/:id/permissions/:permissionId", projectHandler.DeleteProjectPermission)
		// 	}

		// 	// 用戶管理路由
		// 	users := authorized.Group("/users")
		// 	{
		// 		users.GET("/:id", userHandler.GetUser)
		// 		users.PUT("/:id", userHandler.UpdateUser)
		// 		users.DELETE("/:id", userHandler.DeleteUser)
		// 	}

		// 	// 認證相關路由
		// 	auth := authorized.Group("/auth")
		// 	{
		// 		auth.POST("/logout", authHandler.Logout)
		// 		auth.POST("/refresh-token", authHandler.RefreshToken)
		// 	}
		// }
	}

	return r
}
