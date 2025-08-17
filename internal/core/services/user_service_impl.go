package services

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"hexagonal-go/internal/core/domain"
	"hexagonal-go/internal/core/ports"
)

type UserService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Register(user *domain.User) error {
	hashedPin, err := bcrypt.GenerateFromPassword([]byte(user.Pin), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Pin = string(hashedPin)
	return s.userRepo.Create(user)
}

func (s *UserService) Login(phoneNumber, pin string) (*domain.User, error) {
	user, err := s.userRepo.FindByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(pin)); err != nil {
		return nil, errors.New("invalid pin")
	}
	return user, nil
}

func (s *UserService) GetByID(id uuid.UUID) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}
