package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yanggu0t/go-rdbms-practice/internal/models"
	"github.com/yanggu0t/go-rdbms-practice/internal/services"
	"github.com/yanggu0t/go-rdbms-practice/internal/utils"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.Response(c, http.StatusBadRequest, "error", "error_invalid_request_body", nil)
		return
	}

	if err := user.SetPassword(user.Password); err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "error_encrypt_password", nil)
		return
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		utils.Response(c, http.StatusInternalServerError, "error", "error_create_user", gin.H{"Error": err.Error()})
		return
	}

	// 清除密碼，避免返回給客戶端
	user.Password = ""

	utils.Response(c, http.StatusCreated, "success", "success_user_created", user)
}
