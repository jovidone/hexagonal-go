package services

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"hexagonal-go/internal/core/domain"
	"hexagonal-go/internal/core/ports"
)

type mockTransactionRepository struct {
	createFn     func(tx *domain.Transaction) error
	findByUserFn func(userID uuid.UUID) ([]domain.Transaction, error)
}

var _ ports.TransactionRepository = (*mockTransactionRepository)(nil)

func (m *mockTransactionRepository) Create(tx *domain.Transaction) error {
	if m.createFn != nil {
		return m.createFn(tx)
	}
	return nil
}

func (m *mockTransactionRepository) FindByUser(userID uuid.UUID) ([]domain.Transaction, error) {
	if m.findByUserFn != nil {
		return m.findByUserFn(userID)
	}
	return nil, errors.New("not implemented")
}

func TestTransactionService_GetTransactionsByUser(t *testing.T) {
	userID := uuid.New()
	expected := []domain.Transaction{{UserID: userID}}
	repo := &mockTransactionRepository{
		findByUserFn: func(id uuid.UUID) ([]domain.Transaction, error) {
			if id != userID {
				t.Errorf("expected userID %v, got %v", userID, id)
			}
			return expected, nil
		},
	}
	service := NewTransactionService(repo, nil)
	txs, err := service.GetTransactionsByUser(userID)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !reflect.DeepEqual(txs, expected) {
		t.Fatalf("expected %v, got %v", expected, txs)
	}
}

func TestTransactionService_GetTransactionsByUserError(t *testing.T) {
	userID := uuid.New()
	repo := &mockTransactionRepository{
		findByUserFn: func(id uuid.UUID) ([]domain.Transaction, error) {
			return nil, errors.New("db error")
		},
	}
	service := NewTransactionService(repo, nil)
	_, err := service.GetTransactionsByUser(userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
