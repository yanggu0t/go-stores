package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yanggu0t/go-rdbms-practice/internal/models"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) CreateRoleHandler(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 檢查用戶權限
	user := c.MustGet("user").(models.User)
	if !user.HasRole("admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can create roles"})
		return
	}

	// 使用 GORM 創建角色
	if err := h.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Role created successfully", "role": role})
}

func (h *Handler) CreatePermissionHandler(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 檢查用戶權限
	user := c.MustGet("user").(models.User)
	if !user.HasRole("admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can create permissions"})
		return
	}

	// 使用 GORM 創建權限
	if err := h.DB.Create(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Permission created successfully", "permission": permission})
}

func (h *Handler) AssignPermissionToRoleHandler(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("roleID"))
	permissionID, _ := strconv.Atoi(c.Param("permissionID"))

	// 檢查用戶權限
	user := c.MustGet("user").(models.User)
	if !user.HasRole("admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can assign permissions to roles"})
		return
	}

	// 使用 GORM 分配權限給角色
	var role models.Role
	var permission models.Permission

	if err := h.DB.First(&role, roleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	if err := h.DB.First(&permission, permissionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
		return
	}

	if err := h.DB.Model(&role).Association("Permissions").Append(&permission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission assigned to role successfully"})
}

func (h *Handler) CreateUserHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成短 UUID
	user.UserID = uuid.New().String()[:8]

	// 加密密碼
	if err := user.SetPassword(user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法加密密碼"})
		return
	}

	// 使用 GORM 創建用戶
	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 清除返回的密碼字段
	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{"message": "用戶創建成功", "user": user})
}

func (h *Handler) AssignRoleToUserHandler(c *gin.Context) {
	userID := c.Param("userID")
	roleName := c.Param("roleName")

	var user models.User

	if err := h.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用戶未找到"})
		return
	}

	var role models.Role

	if err := h.DB.First(&role, "name = ?", roleName).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "角色未找到"})
		return
	}

	if err := user.AddRole(h.DB, &role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "角色成功分配給用戶"})
}
