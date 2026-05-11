package handler

import (
	"errors"
	"net/http"

	"github.com/user/tahap2-rest-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrPhoneAlreadyRegistered) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Phone Number already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"user_id":       user.UserID,
			"first_name":    user.FirstName,
			"last_name":     user.LastName,
			"phone_number":  user.PhoneNumber,
			"address":       user.Address,
			"balance":       user.Balance,
			"created_date":  user.CreatedDate,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tokens, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Phone Number and PIN doesn't match."})
			return
		}
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Phone Number and PIN doesn't match."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"access_token":  tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
		},
	})
}
