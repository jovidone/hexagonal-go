package http

import (
	"github.com/gin-gonic/gin"
	"hexagonal-go/internal/core/domain"
	"hexagonal-go/internal/core/services"
	"hexagonal-go/internal/utils"
	"net/http"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register handler untuk endpoint /register
func (h *UserHandler) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.userService.Register(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": user,
	})
}

// Login handler untuk endpoint /login
func (h *UserHandler) Login(c *gin.Context) {
	var request struct {
		PhoneNumber string `json:"phone_number"`
		Pin         string `json:"pin"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, err := h.userService.Login(request.PhoneNumber, request.Pin)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone number or pin"})
		return
	}

	// Generate JWT token menggunakan fungsi dari utils
	token, err := utils.GenerateJWT(user.UserID.String()) // Konversi UUID ke string
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"access_token":  token,
			"refresh_token": "", // Jika ada refresh token
		},
	})
}
