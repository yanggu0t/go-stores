package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yanggu0t/go-rdbms-practice/internal/services"
	"github.com/yanggu0t/go-rdbms-practice/internal/utils"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	// 檢查請求頭中是否存在有效的 token
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		session, err := h.AuthService.ValidateToken(tokenString)
		if err == nil {
			utils.Response(c, http.StatusBadRequest, "error", "error_already_logged_in", session.Expires)
			return
		}
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response(c, http.StatusBadRequest, "error", "error_invalid_request_body", nil)
		return
	}

	token, user, err := h.AuthService.Login(req.Account, req.Password)
	if err != nil {
		switch err.Error() {
		case services.ErrUserNotFound:
			utils.Response(c, http.StatusUnauthorized, "error", "error_user_not_found", nil)
		case services.ErrPasswordIncorrect:
			utils.Response(c, http.StatusUnauthorized, "error", "error_password_incorrect", nil)
		default:
			utils.Response(c, http.StatusInternalServerError, "error", "error_login_failed", nil)
		}
		return
	}

	c.Header("Authorization", token)

	// 準備用戶資料
	userData := gin.H{
		"userId":   user.UserID,
		"username": user.Username,
		"email":    user.Email,
	}

	utils.Response(c, http.StatusOK, "success", "success_login", userData)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.Response(c, http.StatusBadRequest, "error", "error_no_token_provided", nil)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	err := h.AuthService.Logout(tokenString)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "error_logout_failed", nil)
		return
	}

	utils.Response(c, http.StatusOK, "success", "success_logout", nil)
}

func (h *AuthHandler) VerifyToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.Response(c, http.StatusBadRequest, "error", "error_no_token_provided", nil)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	session, err := h.AuthService.ValidateToken(tokenString)
	if err != nil {
		utils.Response(c, http.StatusUnauthorized, "error", "error_invalid_token", nil)
		return
	}

	userData := gin.H{
		"userId":  session.UserID,
		"expTime": session.Expires,
	}

	utils.Response(c, http.StatusOK, "success", "success_token_verified", userData)
}
