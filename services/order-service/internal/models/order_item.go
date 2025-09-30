package models

type OrderItem struct {
	ID         uint `json:"id" gorm:"primaryKey"`
	OrderID    int  `json:"order_id" gorm:"not null"`
	MenuItemID int  `json:"menu_item_id" gorm:"not null"`
	Quantity   int  `json:"quantity" gorm:"not null"`
}
