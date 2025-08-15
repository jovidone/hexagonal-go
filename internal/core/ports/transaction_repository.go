package ports

import (
	"github.com/google/uuid"
	"hexagonal-go/internal/core/domain"
)

type TransactionRepository interface {
	Create(tx *domain.Transaction) error
	FindByUser(userID uuid.UUID) ([]domain.Transaction, error)
}
