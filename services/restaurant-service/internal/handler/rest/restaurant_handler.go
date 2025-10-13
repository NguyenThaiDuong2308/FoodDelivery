package rest

import (
	"net/http"
	"restaurant-service/internal/models"
	"restaurant-service/internal/service"
	"strconv"

	"restaurant-service/constant"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type RestaurantHandler struct {
	restaurantService service.RestaurantService
}

func NewRestaurantHandler(restaurantService service.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{restaurantService}
}

type CreateRestaurantRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"omitempty"`
	Address     string `json:"address" validate:"required,min=5,max=255"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,len=10,numeric"`
	Email       string `json:"email" validate:"omitempty,email"`
	Status      string `json:"status" validate:"required,oneof=open closed"`
}

func (h *RestaurantHandler) CreateRestaurant(c *gin.Context) {
	var req CreateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID, ok := claims[constant.UserID].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "wrong type of userID"})
		return
	}
	userRole, ok := claims[constant.UserRole].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "wrong type of userRole"})
		return
	}
	if userRole != constant.RestaurantAdminRole {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "only restaurant_admin can create"})
		return
	}
	restaurant := &models.Restaurant{
		ManagerID:   uint(userID),
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Status:      req.Status,
	}
	err := h.restaurantService.CreateRestaurant(c.Request.Context(), restaurant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, restaurant)

}

func (h *RestaurantHandler) GetRestaurantInfoByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	restaurant, err := h.restaurantService.GetRestaurantInfoByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":           restaurant.ID,
		"manager_id":   restaurant.ManagerID,
		"name":         restaurant.Name,
		"description":  restaurant.Description,
		"address":      restaurant.Address,
		"phone_number": restaurant.PhoneNumber,
		"email":        restaurant.Email,
		"status":       restaurant.Status,
	})
}

type UpdateRestaurantRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Status      string `json:"status" binding:"required"`
}

func (h *RestaurantHandler) UpdateRestaurant(c *gin.Context) {
	var req UpdateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	existingRestaurant, err := h.restaurantService.GetRestaurantInfoByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := claims["user_id"].(float64)
	if existingRestaurant.ManagerID != uint(userID) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "can't change other restaurant status"})
		return
	}
	if req.Status != "" {
		existingRestaurant.Status = req.Status
	}
	err = h.restaurantService.UpdateRestaurant(c.Request.Context(), existingRestaurant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":           existingRestaurant.ID,
		"manager_id":   existingRestaurant.ManagerID,
		"name":         existingRestaurant.Name,
		"description":  existingRestaurant.Description,
		"address":      existingRestaurant.Address,
		"phone_number": existingRestaurant.PhoneNumber,
		"email":        existingRestaurant.Email,
		"status":       existingRestaurant.Status,
	})
}

func (h *RestaurantHandler) GetRestaurantList(c *gin.Context) {
	restaurants, err := h.restaurantService.GetRestaurantList(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"restaurants": restaurants,
	})
}
