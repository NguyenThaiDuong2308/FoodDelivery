package models

type MenuItem struct {
	ID           uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	RestaurantID uint    `json:"restaurant_id" gorm:"not null;index"`
	Name         string  `json:"name" gorm:"type:varchar(255);not null"`
	Description  string  `json:"description" gorm:"type:text"`
	Price        float64 `json:"price" gorm:"not null"`
	Available    bool    `json:"available" gorm:"not null"`
}
