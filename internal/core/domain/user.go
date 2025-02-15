package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UserID      uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	FirstName   string    `gorm:"not null" json:"first_name"`
	LastName    string    `gorm:"not null" json:"last_name"`
	PhoneNumber string    `gorm:"unique;not null" json:"phone_number"`
	Address     string    `gorm:"not null" json:"address"`
	Pin         string    `gorm:"not null" json:"pin"`
	Balance     float64   `gorm:"default:0" json:"balance"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
