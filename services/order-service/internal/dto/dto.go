package dto

type CreateOrderItemRequest struct {
	MenuItemID int `json:"menu_item_id" binding:"required"`
	Quantity   int `json:"quantity" binding:"required"`
}
type CreateOrderRequest struct {
	CustomerID   int                      `json:"customer_id" binding:"required"`
	RestaurantID int                      `json:"restaurant_id" binding:"required"`
	OrderItems   []CreateOrderItemRequest `json:"order_items" binding:"required"`
}

type OrderItemResponse struct {
	ID         uint    `json:"id"`
	OrderID    int     `json:"order_id"`
	MenuItemID int     `json:"menu_item_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
}

type OrderResponse struct {
	ID uint `json:"id"`

	CustomerID   int                 `json:"customer_id"`
	RestaurantID int                 `json:"restaurant_id"`
	TotalPrice   float64             `json:"total_price"`
	Status       string              `json:"status"`
	OrderItems   []OrderItemResponse `json:"order_items"`
}

type UpdateOrderStatusRequest struct {
	OrderID uint   `json:"order_id" binding:"required"`
	Status  string `json:"status" binding:"required"`
}
