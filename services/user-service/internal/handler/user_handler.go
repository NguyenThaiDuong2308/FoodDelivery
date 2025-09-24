package handler

import (
	"net/http"
	"strconv"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService}
}

const AdminRole = "admin"

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := claims["user_id"].(float64)
	if id != int(userID) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Can't access to other user"})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"email":        user.Email,
		"name":         user.Name,
		"phone_number": user.PhoneNumber,
		"address":      user.Address,
		"role":         user.Role,
	})
}

type UpdateUserRequest struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := claims["user_id"].(float64)
	if id != int(userID) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Can't access to other user"})
		return
	}
	existingUser, err := h.userService.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if req.Email != "" {
		existingUser.Email = req.Email
	}
	if req.Name != "" {
		existingUser.Name = req.Name
	}
	if req.PhoneNumber != "" {
		existingUser.PhoneNumber = req.PhoneNumber
	}
	if req.Address != "" {
		existingUser.Address = req.Address
	}
	updatedUser, err := h.userService.UpdateUser(c.Request.Context(), existingUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": updatedUser})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Can't convert user_id to float64"})
		return
	}
	if id != int(userID) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Can't access to other user"})
		return
	}
	if err := h.userService.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func (h *UserHandler) GetAllUserInfo(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userRole := claims["role"]
	if userRole == AdminRole {
		users, err := h.userService.GetAllUserInfo(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"users": users,
		})
	}
}
