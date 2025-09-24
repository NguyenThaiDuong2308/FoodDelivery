package models

type Order struct {
	ID           uint        `json:"id" gorm:"primaryKey"`
	CustomerID   int         `json:"customer_id" gorm:"not null"`
	RestaurantID int         `json:"restaurant_id" gorm:"not null"`
	TotalPrice   float64     `json:"total_price" gorm:"not null"`
	Status       string      `json:"status" gorm:"not null"`
	OrderItems   []OrderItem `json:"order_items" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;"`
}
