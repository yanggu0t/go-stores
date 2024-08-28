package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yanggu0t/go-rdbms-practice/internal/database"
	"github.com/yanggu0t/go-rdbms-practice/internal/handlers"
	"github.com/yanggu0t/go-rdbms-practice/internal/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	handler := handlers.NewHandler(database.DB)

	// 添加身份驗證中間件
	r.Use(middleware.AuthMiddleware())

	// 設置路由
	r.POST("/roles", handler.CreateRoleHandler)
	r.POST("/permissions", handler.CreatePermissionHandler)
	r.POST("/roles/:roleID/permissions/:permissionID", handler.AssignPermissionToRoleHandler)
	r.POST("/users", handler.CreateUserHandler)
	r.POST("/users/:userID/roles/:roleName", handler.AssignRoleToUserHandler)

	return r
}
