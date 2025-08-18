package services

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"hexagonal-go/internal/core/domain"
	"hexagonal-go/internal/core/ports"
	"time"
)

type userMigration struct {
	UserID      uuid.UUID `gorm:"primaryKey;type:uuid"`
	FirstName   string
	LastName    string
	PhoneNumber string `gorm:"unique;not null"`
	Address     string
	Pin         string
	Balance     float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (userMigration) TableName() string { return "users" }

type transactionMigration struct {
	TransactionID   uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID          uuid.UUID `gorm:"type:uuid;not null"`
	TransactionType string
	Amount          float64
	Remarks         string
	BalanceBefore   float64
	BalanceAfter    float64
	CreatedAt       time.Time
}

func (transactionMigration) TableName() string { return "transactions" }

type mockTransactionRepository struct {
	createFn       func(tx *domain.Transaction) error
	createWithTxFn func(dbTx *gorm.DB, tx *domain.Transaction) error
	findByUserFn   func(userID uuid.UUID) ([]domain.Transaction, error)
}

var _ ports.TransactionRepository = (*mockTransactionRepository)(nil)

func (m *mockTransactionRepository) Create(tx *domain.Transaction) error {
	if m.createFn != nil {
		return m.createFn(tx)
	}
	return nil
}

func (m *mockTransactionRepository) CreateWithTx(dbTx *gorm.DB, tx *domain.Transaction) error {
	if m.createWithTxFn != nil {
		return m.createWithTxFn(dbTx, tx)
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

type testTransactionRepo struct {
	db *gorm.DB
}

func (r *testTransactionRepo) Create(tx *domain.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *testTransactionRepo) CreateWithTx(dbTx *gorm.DB, tx *domain.Transaction) error {
	return dbTx.Create(tx).Error
}

func (r *testTransactionRepo) FindByUser(userID uuid.UUID) ([]domain.Transaction, error) {
	var txs []domain.Transaction
	err := r.db.Where("user_id = ?", userID).Find(&txs).Error
	return txs, err
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	if err := db.AutoMigrate(&userMigration{}, &transactionMigration{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestTransactionService_Deposit_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := &testTransactionRepo{db: db}
	service := NewTransactionService(repo, db)
	user := domain.User{UserID: uuid.New(), FirstName: "A", LastName: "B", PhoneNumber: "111", Address: "addr", Pin: "1234", Balance: 100}
	db.Create(&user)

	tx, err := service.Deposit(user.UserID, 50, "deposit")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	var updated domain.User
	db.First(&updated, "user_id = ?", user.UserID)
	if updated.Balance != 150 {
		t.Fatalf("expected balance 150, got %v", updated.Balance)
	}

	var count int64
	db.Model(&domain.Transaction{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 transaction, got %d", count)
	}
	if tx.TransactionType != "CREDIT" || tx.Amount != 50 {
		t.Fatalf("unexpected transaction: %+v", tx)
	}
}

func TestTransactionService_Withdraw_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := &testTransactionRepo{db: db}
	service := NewTransactionService(repo, db)
	user := domain.User{UserID: uuid.New(), FirstName: "A", LastName: "B", PhoneNumber: "111", Address: "addr", Pin: "1234", Balance: 100}
	db.Create(&user)

	tx, err := service.Withdraw(user.UserID, 40, "withdraw")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	var updated domain.User
	db.First(&updated, "user_id = ?", user.UserID)
	if updated.Balance != 60 {
		t.Fatalf("expected balance 60, got %v", updated.Balance)
	}

	var count int64
	db.Model(&domain.Transaction{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 transaction, got %d", count)
	}
	if tx.TransactionType != "DEBIT" || tx.Amount != 40 {
		t.Fatalf("unexpected transaction: %+v", tx)
	}
}

func TestTransactionService_Withdraw_InsufficientFunds(t *testing.T) {
	db := setupTestDB(t)
	repo := &testTransactionRepo{db: db}
	service := NewTransactionService(repo, db)
	user := domain.User{UserID: uuid.New(), FirstName: "A", LastName: "B", PhoneNumber: "111", Address: "addr", Pin: "1234", Balance: 20}
	db.Create(&user)

	_, err := service.Withdraw(user.UserID, 40, "withdraw")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	var updated domain.User
	db.First(&updated, "user_id = ?", user.UserID)
	if updated.Balance != 20 {
		t.Fatalf("expected balance 20, got %v", updated.Balance)
	}

	var count int64
	db.Model(&domain.Transaction{}).Count(&count)
	if count != 0 {
		t.Fatalf("expected 0 transactions, got %d", count)
	}
}

func TestTransactionService_Transfer_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := &testTransactionRepo{db: db}
	service := NewTransactionService(repo, db)
	fromUser := domain.User{UserID: uuid.New(), FirstName: "A", LastName: "B", PhoneNumber: "111", Address: "addr", Pin: "1234", Balance: 100}
	toUser := domain.User{UserID: uuid.New(), FirstName: "C", LastName: "D", PhoneNumber: "222", Address: "addr", Pin: "1234", Balance: 50}
	db.Create(&fromUser)
	db.Create(&toUser)

	_, _, err := service.Transfer(fromUser.UserID, toUser.UserID, 30, "transfer")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	var updatedFrom, updatedTo domain.User
	db.First(&updatedFrom, "user_id = ?", fromUser.UserID)
	db.First(&updatedTo, "user_id = ?", toUser.UserID)
	if updatedFrom.Balance != 70 {
		t.Fatalf("expected from balance 70, got %v", updatedFrom.Balance)
	}
	if updatedTo.Balance != 80 {
		t.Fatalf("expected to balance 80, got %v", updatedTo.Balance)
	}

	var count int64
	db.Model(&domain.Transaction{}).Count(&count)
	if count != 2 {
		t.Fatalf("expected 2 transactions, got %d", count)
	}
}

func TestTransactionService_Transfer_InsufficientFunds(t *testing.T) {
	db := setupTestDB(t)
	repo := &testTransactionRepo{db: db}
	service := NewTransactionService(repo, db)
	fromUser := domain.User{UserID: uuid.New(), FirstName: "A", LastName: "B", PhoneNumber: "111", Address: "addr", Pin: "1234", Balance: 20}
	toUser := domain.User{UserID: uuid.New(), FirstName: "C", LastName: "D", PhoneNumber: "222", Address: "addr", Pin: "1234", Balance: 50}
	db.Create(&fromUser)
	db.Create(&toUser)

	_, _, err := service.Transfer(fromUser.UserID, toUser.UserID, 30, "transfer")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	var updatedFrom, updatedTo domain.User
	db.First(&updatedFrom, "user_id = ?", fromUser.UserID)
	db.First(&updatedTo, "user_id = ?", toUser.UserID)
	if updatedFrom.Balance != 20 {
		t.Fatalf("expected from balance 20, got %v", updatedFrom.Balance)
	}
	if updatedTo.Balance != 50 {
		t.Fatalf("expected to balance 50, got %v", updatedTo.Balance)
	}

	var count int64
	db.Model(&domain.Transaction{}).Count(&count)
	if count != 0 {
		t.Fatalf("expected 0 transactions, got %d", count)
	}
}
