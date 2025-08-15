package services

import (
        "errors"
        "testing"

        "golang.org/x/crypto/bcrypt"
        "hexagonal-go/internal/core/domain"
        "hexagonal-go/internal/core/ports"
)

type mockUserRepository struct {
        createFn            func(user *domain.User) error
        findByPhoneNumberFn func(phoneNumber string) (*domain.User, error)
}

var _ ports.UserRepository = (*mockUserRepository)(nil)

func (m *mockUserRepository) Create(user *domain.User) error {
        if m.createFn != nil {
                return m.createFn(user)
        }
        return nil
}

func (m *mockUserRepository) FindByPhoneNumber(phoneNumber string) (*domain.User, error) {
        if m.findByPhoneNumberFn != nil {
                return m.findByPhoneNumberFn(phoneNumber)
        }
        return nil, errors.New("not implemented")
}

func TestUserServiceRegister(t *testing.T) {
        var savedUser *domain.User
        repo := &mockUserRepository{
                createFn: func(u *domain.User) error {
                        savedUser = u
                        return nil
                },
        }
        service := NewUserService(repo)

        if err := service.Register(&domain.User{PhoneNumber: "08123", Pin: "1234"}); err != nil {
                t.Fatalf("Register returned error: %v", err)
        }
        if savedUser == nil {
                t.Fatalf("expected Create to be called")
        }
        if savedUser.Pin == "1234" {
                t.Fatalf("expected pin to be hashed")
        }
        if err := bcrypt.CompareHashAndPassword([]byte(savedUser.Pin), []byte("1234")); err != nil {
                t.Fatalf("stored pin does not match original: %v", err)
        }
}

func TestUserServiceLogin(t *testing.T) {
        hashed, _ := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
        repo := &mockUserRepository{
                findByPhoneNumberFn: func(phone string) (*domain.User, error) {
                        return &domain.User{PhoneNumber: phone, Pin: string(hashed)}, nil
                },
        }
        service := NewUserService(repo)

        if _, err := service.Login("08123", "1234"); err != nil {
                t.Fatalf("expected login to succeed, got %v", err)
        }
        if _, err := service.Login("08123", "4321"); err == nil {
                t.Fatalf("expected error for invalid pin")
        }
}
