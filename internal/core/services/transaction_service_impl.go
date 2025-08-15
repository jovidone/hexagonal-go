package services

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"hexagonal-go/internal/core/domain"
	"hexagonal-go/internal/core/ports"
)

type TransactionService struct {
	transactionRepo ports.TransactionRepository
	db              *gorm.DB
}

func NewTransactionService(transactionRepo ports.TransactionRepository, db *gorm.DB) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo, db: db}
}

func (s *TransactionService) Deposit(userID uuid.UUID, amount float64, remarks string) (*domain.Transaction, error) {
	var user domain.User
	if err := s.db.First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	balanceBefore := user.Balance
	user.Balance += amount
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}
	tx := domain.Transaction{
		UserID:          userID,
		TransactionType: "CREDIT",
		Amount:          amount,
		Remarks:         remarks,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    user.Balance,
	}
	if err := s.transactionRepo.Create(&tx); err != nil {
		return nil, err
	}
	return &tx, nil
}

func (s *TransactionService) Withdraw(userID uuid.UUID, amount float64, remarks string) (*domain.Transaction, error) {
	var user domain.User
	if err := s.db.First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	if user.Balance < amount {
		return nil, errors.New("insufficient balance")
	}
	balanceBefore := user.Balance
	user.Balance -= amount
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}
	tx := domain.Transaction{
		UserID:          userID,
		TransactionType: "DEBIT",
		Amount:          amount,
		Remarks:         remarks,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    user.Balance,
	}
	if err := s.transactionRepo.Create(&tx); err != nil {
		return nil, err
	}
	return &tx, nil
}

func (s *TransactionService) GetTransactionsByUser(userID uuid.UUID) ([]domain.Transaction, error) {
	return s.transactionRepo.FindByUser(userID)
}
