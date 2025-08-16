package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"hexagonal-go/internal/core/domain"
)

type TransactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepositoryImpl(db *gorm.DB) *TransactionRepositoryImpl {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	return &TransactionRepositoryImpl{db: db}
}

func (r *TransactionRepositoryImpl) Create(tx *domain.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *TransactionRepositoryImpl) CreateWithTx(dbTx *gorm.DB, tx *domain.Transaction) error {
	return dbTx.Create(tx).Error
}

func (r *TransactionRepositoryImpl) FindByUser(userID uuid.UUID) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}
