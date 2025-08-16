package ports

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"hexagonal-go/internal/core/domain"
)

type TransactionRepository interface {
	Create(tx *domain.Transaction) error
	CreateWithTx(dbTx *gorm.DB, tx *domain.Transaction) error
	FindByUser(userID uuid.UUID) ([]domain.Transaction, error)
}
