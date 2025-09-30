package handler

import (
	"net/http"
	"order-service/internal/dto"
	"order-service/internal/models"
	"order-service/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (o *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	claims := c.MustGet("claims").(jwt.MapClaims)
	customerID := claims["user_id"].(float64)
	customerRole := claims["role"].(string)
	if customerRole != "customer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only customer can create order"})
		return
	}
	if int(customerID) != req.CustomerID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can't access to other user"})
		return
	}
	var orderItems []models.OrderItem
	for _, item := range req.OrderItems {
		orderItems = append(orderItems, models.OrderItem{
			MenuItemID: item.MenuItemID,
			Quantity:   item.Quantity,
		})
	}
	err := o.orderService.CreateOrder(c.Request.Context(), req.CustomerID, req.RestaurantID, orderItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order created"})

}

func (o *OrderHandler) GetOrderByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	order, err := o.orderService.GetOrderByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := claims["user_id"].(float64)
	userRole := claims["role"].(string)
	if (int(userID) == order.CustomerID && userRole == "customer") || (int(userID) == order.RestaurantID && userRole == "restaurant") {
		var orderItems []dto.OrderItemResponse
		for _, item := range order.OrderItems {
			orderItems = append(orderItems, dto.OrderItemResponse{
				ID:         item.ID,
				OrderID:    item.OrderID,
				MenuItemID: item.MenuItemID,
				Quantity:   item.Quantity,
			})
		}
		orderResp := dto.OrderResponse{
			ID:           order.ID,
			CustomerID:   order.CustomerID,
			RestaurantID: order.RestaurantID,
			TotalPrice:   order.TotalPrice,
			Status:       order.Status,
			OrderItems:   orderItems,
		}

		c.JSON(http.StatusOK, orderResp)
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"message": "Can't access to other order's user"})
}

func (o *OrderHandler) GetOrderByCustomerID(c *gin.Context) {
	customerID, _ := strconv.Atoi(c.Param("customer_id"))
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := claims["user_id"].(float64)
	if int(userID) != customerID {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Can't access to other user"})
		return
	}
	orders, err := o.orderService.GetOrderByCustomerID(c.Request.Context(), customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	var orderResponses []dto.OrderResponse
	for _, order := range *orders {
		var orderItems []dto.OrderItemResponse
		for _, item := range order.OrderItems {
			orderItems = append(orderItems, dto.OrderItemResponse{
				ID:         item.ID,
				OrderID:    item.OrderID,
				MenuItemID: item.MenuItemID,
				Quantity:   item.Quantity,
			})
		}
		orderResponses = append(orderResponses, dto.OrderResponse{
			ID:           order.ID,
			CustomerID:   order.CustomerID,
			RestaurantID: order.RestaurantID,
			TotalPrice:   order.TotalPrice,
			Status:       order.Status,
			OrderItems:   orderItems,
		})
	}
	c.JSON(http.StatusOK, orderResponses)
}

func (o *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	order, err := o.orderService.GetOrderByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := claims["user_id"].(float64)
	userRole := claims["role"].(string)
	if (int(userID) == order.CustomerID && userRole == "customer") || (int(userID) == order.RestaurantID && userRole == "restaurant") {
		var req dto.UpdateOrderStatusRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if order.ID == req.OrderID {
			err := o.orderService.UpdateOrderStatus(c.Request.Context(), req.OrderID, req.Status)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Order updated"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
	}
	c.JSON(http.StatusUnauthorized, gin.H{"message": "Can't access to other order's user"})
}
