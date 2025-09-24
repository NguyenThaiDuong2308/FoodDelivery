package models

type Shippers struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	UserID uint   `json:"user_id" gorm:"not null"`
	Status string `json:"status" gorm:"not null"`
}
