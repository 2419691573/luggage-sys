package handlers

import (
	"net/http"

	"luggage-sys2/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "login failed",
			"error":   "invalid request",
		})
		return
	}

	user, token, err := h.authService.Login(req.Username, req.Password)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "login failed",
			"error":   "invalid username or password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login success",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
			"hotel_id": user.HotelID,
		},
		"token": token,
	})
}
