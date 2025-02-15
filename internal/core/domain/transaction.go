package domain

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	TransactionID   uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID          uuid.UUID `gorm:"type:uuid;not null"`
	TransactionType string    `gorm:"not null"` // CREDIT or DEBIT
	Amount          float64   `gorm:"not null"`
	Remarks         string    `gorm:"not null"`
	BalanceBefore   float64   `gorm:"not null"`
	BalanceAfter    float64   `gorm:"not null"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
}
