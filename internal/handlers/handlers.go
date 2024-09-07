package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yanggu0t/go-rdbms-practice/config"
	"github.com/yanggu0t/go-rdbms-practice/internal/models"
	"github.com/yanggu0t/go-rdbms-practice/internal/services"
	"github.com/yanggu0t/go-rdbms-practice/internal/utils"
	"gorm.io/gorm"
)

type Handler struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func NewHandler(db *gorm.DB, cfg *config.Config) *Handler {
	return &Handler{DB: db, Cfg: cfg}
}

// ================== Search ==================

func (h *Handler) GetAllUsersHandler(c *gin.Context) {
	query, page, pageSize := utils.GetPaginationParams(c)

	db := h.DB.Model(&models.User{})
	db = utils.Search(db, query, "username", "user_id")

	var total int64
	db.Count(&total)

	var users []models.User
	if err := db.Preload("Roles").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "獲取用戶列表時出錯: "+err.Error(), nil)
		return
	}

	response := make([]gin.H, len(users))
	for i, user := range users {
		response[i] = gin.H{
			"index":    user.ID,
			"username": user.Username,
			"email":    user.Email,
			"roles":    user.Roles,
			"userId":   user.UserID,
		}
	}

	utils.Response(c, http.StatusOK, "success", "", gin.H{
		"users":      response,
		"pagination": utils.GetPaginationResponse(page, pageSize, total),
	})
}

func (h *Handler) GetAllRolesHandler(c *gin.Context) {
	query, page, pageSize := utils.GetPaginationParams(c)

	db := h.DB.Model(&models.Role{})
	db = utils.Search(db, query, "name")

	var total int64
	db.Count(&total)

	var roles []models.Role
	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&roles).Error; err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "獲取角色列表時出錯: "+err.Error(), nil)
		return
	}

	response := make([]gin.H, len(roles))
	for i, role := range roles {
		response[i] = gin.H{
			"index":       role.ID,
			"name":        role.Name,
			"description": role.Description,
			"permissions": role.Permissions,
		}
	}

	utils.Response(c, http.StatusOK, "success", "", gin.H{
		"roles":      response,
		"pagination": utils.GetPaginationResponse(page, pageSize, total),
	})
}

func (h *Handler) GetAllPermissionsHandler(c *gin.Context) {
	query, page, pageSize := utils.GetPaginationParams(c)

	db := h.DB.Model(&models.Permission{})
	db = utils.Search(db, query, "name")

	var total int64
	db.Count(&total)

	var permissions []struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Code        string `json:"code"`
	}
	if err := db.Select("id", "name", "description", "code").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&permissions).Error; err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "獲取權限列表時出錯: "+err.Error(), nil)
		return
	}

	utils.Response(c, http.StatusOK, "success", "", gin.H{
		"permissions": permissions,
		"pagination":  utils.GetPaginationResponse(page, pageSize, total),
	})
}

// ================== Create ==================

func (h *Handler) CreateRoleHandler(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		utils.Response(c, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	// 使用 GORM 創建角色
	if err := h.DB.Create(&role).Error; err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	utils.Response(c, http.StatusCreated, "success", "Role created successfully", role)
}

func (h *Handler) CreatePermissionHandler(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		utils.Response(c, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	// 使用 GORM 創建權限
	if err := h.DB.Create(&permission).Error; err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	utils.Response(c, http.StatusCreated, "success", "Permission created successfully", permission)
}

func (h *Handler) CreateUserHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Response(c, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	// 生成短 UUID
	user.UserID = uuid.New().String()[:8]

	// 加密密碼
	if err := user.SetPassword(user.Password); err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "無法加密密碼", nil)
		return
	}

	// 使用 GORM 創建用戶
	if err := h.DB.Create(&user).Error; err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	// 清除返回的密碼字段
	user.Password = ""

	utils.Response(c, http.StatusCreated, "success", "用戶創建成功", user)
}

// ================== Update ==================

type UpdateRoleRequest struct {
	RoleID      string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Permissions []uint `json:"permissions"`
}

func (h *Handler) UpdateRoleHandler(c *gin.Context) {
	roleID := c.Param("id")
	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response(c, http.StatusBadRequest, "error", "無效的請求體: "+err.Error(), nil)
		return
	}

	var existingRole models.Role
	if err := h.DB.Preload("Permissions").First(&existingRole, roleID).Error; err != nil {
		utils.Response(c, http.StatusNotFound, "error", "角色未找到", nil)
		return
	}

	// 更新角色基本信息
	existingRole.Name = req.Name
	existingRole.Description = req.Description

	// 開始事務
	tx := h.DB.Begin()

	// 更新角色基本信息
	if err := tx.Save(&existingRole).Error; err != nil {
		tx.Rollback()
		utils.Response(c, http.StatusInternalServerError, "error", "更新角色時出錯: "+err.Error(), nil)
		return
	}

	// 清除現有權限
	if err := tx.Model(&existingRole).Association("Permissions").Clear(); err != nil {
		tx.Rollback()
		utils.Response(c, http.StatusInternalServerError, "error", "清除權限時出錯: "+err.Error(), nil)
		return
	}

	// 添加新的權限
	if len(req.Permissions) > 0 {
		var newPermissions []models.Permission
		if err := tx.Where("id IN ?", req.Permissions).Find(&newPermissions).Error; err != nil {
			tx.Rollback()
			utils.Response(c, http.StatusInternalServerError, "error", "獲取權限時出錯: "+err.Error(), nil)
			return
		}

		if err := tx.Model(&existingRole).Association("Permissions").Append(&newPermissions); err != nil {
			tx.Rollback()
			utils.Response(c, http.StatusInternalServerError, "error", "添加權限時出錯: "+err.Error(), nil)
			return
		}
	}

	// 提交事務
	if err := tx.Commit().Error; err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "提交更改時出錯: "+err.Error(), nil)
		return
	}

	// 重新加載角色以獲取更新後的權限
	if err := h.DB.Preload("Permissions").First(&existingRole, roleID).Error; err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "獲取更新後的角色信息時出錯: "+err.Error(), nil)
		return
	}

	// 準備權限信息
	var permissions []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
	for _, perm := range existingRole.Permissions {
		permissions = append(permissions, struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}{
			ID:   perm.ID,
			Name: perm.Name,
		})
	}

	// 創建包含權限信息的響應
	response := gin.H{
		"message": "角色更新成功",
		"role": gin.H{
			"id":          existingRole.ID,
			"name":        existingRole.Name,
			"description": existingRole.Description,
			"permissions": permissions,
		},
	}

	utils.Response(c, http.StatusOK, "success", "角色更新成功", response)
}

func (h *Handler) UpdateUserHandler(c *gin.Context) {
	userID := c.Param("id")
	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		utils.Response(c, http.StatusBadRequest, "error", "無效的請求體: "+err.Error(), nil)
		return
	}

	var existingUser models.User
	if err := h.DB.First(&existingUser, "user_id = ?", userID).Error; err != nil {
		utils.Response(c, http.StatusNotFound, "error", "用戶未找到", nil)
		return
	}

	existingUser.Username = updatedUser.Username
	existingUser.Email = updatedUser.Email

	if updatedUser.Password != "" {
		if err := existingUser.SetPassword(updatedUser.Password); err != nil {
			utils.Response(c, http.StatusInternalServerError, "error", "無法加密密碼", nil)
			return
		}
	}

	if err := h.DB.Save(&existingUser).Error; err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "更新用戶時出錯: "+err.Error(), nil)
		return
	}

	existingUser.Password = "" // 清除返回的密碼字段
	utils.Response(c, http.StatusOK, "success", "用戶更新成功", existingUser)
}

type UpdatePermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
}

func (h *Handler) UpdatePermissionHandler(c *gin.Context) {
	permissionID := c.Param("id")
	var req UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response(c, http.StatusBadRequest, "error", "無效的請求體: "+err.Error(), nil)
		return
	}

	var existingPermission models.Permission
	if err := h.DB.First(&existingPermission, permissionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Response(c, http.StatusNotFound, "error", "權限未找到", nil)
		} else {
			utils.Response(c, http.StatusInternalServerError, "error", "獲取權限時出錯: "+err.Error(), nil)
		}
		return
	}

	// 檢查新名稱是否與其他權限衝突（排除當前權限）
	if req.Name != existingPermission.Name {
		var count int64
		if err := h.DB.Model(&models.Permission{}).
			Where("name = ? AND id != ?", req.Name, existingPermission.ID).
			Count(&count).Error; err != nil {
			utils.Response(c, http.StatusInternalServerError, "error", "檢查權限名稱時出錯: "+err.Error(), nil)
			return
		}
		if count > 0 {
			utils.Response(c, http.StatusConflict, "error", "權限名稱已存在", nil)
			return
		}
	}

	// 檢查新權限碼是否與其他權限衝突（排除當前權限）
	if req.Code != existingPermission.Code {
		var count int64
		if err := h.DB.Model(&models.Permission{}).
			Where("code = ? AND id != ?", req.Code, existingPermission.ID).
			Count(&count).Error; err != nil {
			utils.Response(c, http.StatusInternalServerError, "error", "檢查權限碼時出錯: "+err.Error(), nil)
			return
		}
		if count > 0 {
			utils.Response(c,http.StatusConflict, "error", "權限碼已存在", nil)
			return
		}
	}

	// 更新權限信息
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}

	if err := h.DB.Model(&existingPermission).Updates(updates).Error; err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "更新權限時出錯", nil)
		return
	}

	// 重新獲取更新後的權限信息
	if err := h.DB.First(&existingPermission, permissionID).Error; err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "獲取更新後的權限信息時出錯", nil)
		return
	}

	// 創建響應
	response := gin.H{
		"message": "權限更新成功",
		"permission": gin.H{
			"id":          existingPermission.ID,
			"name":        existingPermission.Name,
			"description": existingPermission.Description,
			"code":        existingPermission.Code,
		},
	}

	c.JSON(http.StatusOK, response)
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response(c, http.StatusBadRequest, "error", "無效的請求體", nil)
		return
	}

	authService := services.NewAuthService(h.DB, h.Cfg.JWTSecret)
	token, err := authService.Login(req.Username, req.Password)
	if err != nil {
		utils.Response(c, http.StatusUnauthorized, "error", "登錄失敗", nil)
		return
	}

	c.Header("Authorization", token)
	utils.Response(c, http.StatusOK, "success", "登錄成功", nil)
}