package models

type Restaurant struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ManagerID   uint       `json:"manager_id" gorm:"not null"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Description string     `json:"description" gorm:"type:text"`
	Address     string     `json:"address" gorm:"type:varchar(255);not null"`
	PhoneNumber string     `json:"phone_number" gorm:"type:varchar(255)"`
	Email       string     `json:"email" gorm:"type:varchar(255)"`
	Status      string     `json:"status" gorm:"type:varchar(10);not null;check:status IN ('open','closed')"`
	MenuItems   []MenuItem `json:"menu_items" gorm:"foreignKey:RestaurantID;constraint:OnDelete:CASCADE"`
}
