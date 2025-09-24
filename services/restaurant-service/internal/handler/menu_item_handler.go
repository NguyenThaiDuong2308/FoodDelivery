package handler

import (
	"log"
	"net/http"
	"restaurant-service/internal/models"
	"restaurant-service/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MenuItemHandler struct {
	menuItemService   service.MenuItemService
	restaurantService service.RestaurantService
}

func NewMenuItemHandler(menuItemService service.MenuItemService, restaurantService service.RestaurantService) *MenuItemHandler {
	return &MenuItemHandler{
		menuItemService:   menuItemService,
		restaurantService: restaurantService,
	}
}

func (h *MenuItemHandler) GetMenuListByRestaurantInfo(c *gin.Context) {
	restaurantID, _ := strconv.Atoi(c.Param("restaurant_id"))
	menuList, err := h.menuItemService.GetMenuListByRestaurantInfo(c.Request.Context(), uint(restaurantID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, menuList)
}

type AddMenuItemRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description, omitempty"`
	Price       float64 `json:"price" binding:"required"`
	Available   bool    `json:"available" binding:"required"`
}

func (h *MenuItemHandler) CreateMenuItem(c *gin.Context) {
	restaurantID, err := strconv.Atoi(c.Param("restaurant_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	log.Println(restaurantID)
	log.Println(h.restaurantService)
	restaurant, err := h.restaurantService.GetRestaurantInfoByID(c.Request.Context(), uint(restaurantID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if restaurant.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Can't add menu item for other restaurant"})
		return
	}
	var req AddMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	menuItemReq := models.MenuItem{
		RestaurantID: restaurant.ID,
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		Available:    req.Available,
	}

	menuItem, err := h.menuItemService.CreateMenuItem(c.Request.Context(), &menuItemReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": menuItem})
}

func (h *MenuItemHandler) GetMenuItemByItemID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	menuItem, err := h.menuItemService.GetMenuItemByItemID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, menuItem)
}

type UpdateMenuItemRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Available   *bool    `json:"available"`
}

func (h *MenuItemHandler) UpdateMenuItem(c *gin.Context) {
	restaurantID, _ := strconv.Atoi(c.Param("restaurant_id"))
	itemID, _ := strconv.Atoi(c.Param("id"))
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	restaurant, err := h.restaurantService.GetRestaurantInfoByID(c.Request.Context(), uint(restaurantID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if restaurant.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Can't update menu item for other restaurant"})
		return
	}

	existingMenuItem, err := h.menuItemService.GetMenuItemByItemID(c.Request.Context(), uint(itemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var req UpdateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	if req.Name != nil {
		existingMenuItem.Name = *req.Name
	}
	if req.Description != nil {
		existingMenuItem.Description = *req.Description
	}
	if req.Price != nil {
		existingMenuItem.Price = *req.Price
	}
	if req.Available != nil {
		existingMenuItem.Available = *req.Available
	}
	updatedMenuItem, err := h.menuItemService.UpdateMenuItem(c.Request.Context(), existingMenuItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": updatedMenuItem})

}

func (h *MenuItemHandler) DeleteMenuItem(c *gin.Context) {
	restaurantID, _ := strconv.Atoi(c.Param("restaurant_id"))
	itemID, _ := strconv.Atoi(c.Param("id"))
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	restaurant, err := h.restaurantService.GetRestaurantInfoByID(c.Request.Context(), uint(restaurantID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if restaurant.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Can't delete menu item for other restaurant"})
		return
	}
	if err := h.menuItemService.DeleteMenuItem(c.Request.Context(), uint(itemID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Menu item has been deleted"})
}
