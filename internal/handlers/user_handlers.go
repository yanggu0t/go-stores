package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/yanggu0t/go-rdbms-practice/internal/models"
	"github.com/yanggu0t/go-rdbms-practice/internal/services"
	"github.com/yanggu0t/go-rdbms-practice/internal/utils"
)

type UserHandler struct {
	UserService   *services.UserService
	CryptoService *services.CryptoService
}

func NewUserHandler(userService *services.UserService, cryptoService *services.CryptoService) *UserHandler {
	return &UserHandler{
		UserService:   userService,
		CryptoService: cryptoService,
	}
}

type EncryptedCreateUserRequest struct {
	EncryptedData string `json:"data" binding:"required"`
}

type CreateUserData struct {
	Username     string `json:"username" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=8"`
	ProjectRoles []struct {
		ProjectID uint `json:"projects,omitempty"`
		RoleID    uint `json:"roles,omitempty"`
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var encryptedReq EncryptedCreateUserRequest
	if err := c.ShouldBindJSON(&encryptedReq); err != nil {
		utils.Response(c, http.StatusBadRequest, "error", "error_invalid_request_body", nil)
		return
	}

	// 解密數據
	decryptedData, err := h.CryptoService.Decrypt(encryptedReq.EncryptedData)
	if err != nil {
		utils.Response(c, http.StatusBadRequest, "error", "error_decrypting_data", nil)
		return
	}

	var userData CreateUserData
	if err := json.Unmarshal([]byte(decryptedData), &userData); err != nil {
		utils.Response(c, http.StatusBadRequest, "error", "error_invalid_decrypted_data", nil)
		return
	}

	// 檢查帳號格式
	if err := validateUsername(userData.Username); err != nil {
		utils.Response(c, http.StatusBadRequest, "error", "error_invalid_username", gin.H{"Error": err.Error()})
		return
	}

	// 檢查密碼格式
	if err := validatePassword(userData.Password); err != nil {
		utils.Response(c, http.StatusBadRequest, "error", "error_invalid_password", gin.H{"Error": err.Error()})
		return
	}

	if err := validateEmail(userData.Email); err != nil {
		utils.Response(c, http.StatusBadRequest, "error", "error_invalid_email", gin.H{"Error": err.Error()})
		return
	}

	// 創建用戶對象
	userID, err := models.GenerateID()
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "error_generate_user_id", nil)
		return
	}

	// 創建用戶對象
	user := models.User{
		UserID:   userID,
		Username: userData.Username,
		Email:    userData.Email,
	}

	if err := user.SetPassword(userData.Password); err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "error_encrypt_password", nil)
		return
	}

	// 開始數據庫事務
	tx := h.UserService.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := h.UserService.CreateUser(&user); err != nil {
		tx.Rollback()
		utils.Response(c, http.StatusInternalServerError, "error", "error_create_user", gin.H{"Error": err.Error()})
		return
	}

	// 處理項目和角色分配
	for _, assignment := range userData.ProjectRoles {
		if err := h.UserService.AddUserToProject(tx, user.ID, assignment.ProjectID, assignment.RoleID); err != nil {
			tx.Rollback()
			utils.Response(c, http.StatusInternalServerError, "error", "error_add_user_to_project", gin.H{"Error": err.Error()})
			return
		}
	}

	// 提交事務
	if err := tx.Commit().Error; err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "error_commit_transaction", gin.H{"Error": err.Error()})
		return
	}

	// 清除密碼，避免返回給客戶端
	user.Password = ""

	// 獲取用戶的項目和角色信息
	projects, err := h.UserService.GetUserProjects(user.ID)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "error_get_user_projects", gin.H{"Error": err.Error()})
		return
	}

	roles, err := h.UserService.GetUserRoles(user.ID)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "error_get_user_roles", gin.H{"Error": err.Error()})
		return
	}

	response := gin.H{
		"user":     user,
		"projects": projects,
		"roles":    roles,
	}

	utils.Response(c, http.StatusCreated, "success", "success_user_created", response)
}

func validateUsername(username string) error {
	if len(username) <= 7 {
		return errors.New("username must be at least 8 characters long")
	}

	// 檢查是否只包含英文字母和數字
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return errors.New("username must contain only letters and digits")
		}
	}

	// 檢查是否同時包含字母和數字
	hasLetter := false
	hasDigit := false
	for _, char := range username {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return errors.New("username must contain both letters and digits")
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 16 || len(password) > 24 {
		return errors.New("密碼長度必須在16到24個字符之間")
	}

	var (
		hasLower  = false
		hasNumber = false
	)

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		}
	}

	if !hasLower || !hasNumber {
		return errors.New("密碼必須至少包含一個小寫字母和一個數字")
	}

	// 檢查是否只包含允許的字符
	for _, char := range password {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && !unicode.IsPunct(char) && !unicode.IsSymbol(char) {
			return errors.New("密碼只能包含字母、數字和特殊字符")
		}
	}

	return nil
}

func validateEmail(email string) error {
	if !strings.Contains(email, "@") {
		return errors.New("email must contain @")
	}

	return nil
}
