package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Email       string `json:"email" gorm:"unique;not null"`
	Name        string `json:"name" gorm:"not null"`
	PhoneNumber string `json:"phone_number" gorm:"not null"`
	Address     string `json:"address" gorm:"not null"`
	Password    string `json:"password" gorm:"not null"`
	Role        string `json:"role" gorm:"not null;check:role IN ('customer','restaurant_admin','shipper','admin')"`
}

type PasswordResetToken struct {
	ID         uuid.UUID `json:"id" gorm:"primary_key"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	Token      string    `json:"token" gorm:"unique"`
	Expires_at time.Time `json:"expires_at" gorm:"not null"`
	Used       bool      `json:"used" gorm:"default:false"`
}
