package handler

import (
	"net/http"

	"github.com/user/tahap2-rest-api/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(string)

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := h.userService.UpdateProfile(c.Request.Context(), uid, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"user_id":       user.UserID,
			"first_name":    user.FirstName,
			"last_name":     user.LastName,
			"phone_number":  user.PhoneNumber,
			"address":       user.Address,
			"balance":       user.Balance,
			"updated_date":  user.UpdatedDate,
		},
	})
}