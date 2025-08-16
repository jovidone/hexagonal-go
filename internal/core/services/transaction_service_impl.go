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

func (s *TransactionService) Transfer(fromID, toID uuid.UUID, amount float64, remarks string) (*domain.Transaction, *domain.Transaction, error) {
	var debitTx, creditTx domain.Transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var fromUser, toUser domain.User
		if err := tx.First(&fromUser, "user_id = ?", fromID).Error; err != nil {
			return err
		}
		if err := tx.First(&toUser, "user_id = ?", toID).Error; err != nil {
			return err
		}
		if fromUser.Balance < amount {
			return errors.New("insufficient balance")
		}
		fromBalanceBefore := fromUser.Balance
		toBalanceBefore := toUser.Balance
		fromUser.Balance -= amount
		toUser.Balance += amount
		if err := tx.Save(&fromUser).Error; err != nil {
			return err
		}
		if err := tx.Save(&toUser).Error; err != nil {
			return err
		}
		debitTx = domain.Transaction{
			UserID:          fromID,
			TransactionType: "DEBIT",
			Amount:          amount,
			Remarks:         remarks,
			BalanceBefore:   fromBalanceBefore,
			BalanceAfter:    fromUser.Balance,
		}
		creditTx = domain.Transaction{
			UserID:          toID,
			TransactionType: "CREDIT",
			Amount:          amount,
			Remarks:         remarks,
			BalanceBefore:   toBalanceBefore,
			BalanceAfter:    toUser.Balance,
		}
		if err := s.transactionRepo.CreateWithTx(tx, &debitTx); err != nil {
			return err
		}
		if err := s.transactionRepo.CreateWithTx(tx, &creditTx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return &debitTx, &creditTx, nil
}

func (s *TransactionService) GetTransactionsByUser(userID uuid.UUID) ([]domain.Transaction, error) {
	return s.transactionRepo.FindByUser(userID)
}
