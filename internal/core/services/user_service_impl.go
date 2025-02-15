package services

import (
	"errors"
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
	return s.userRepo.Create(user)
}

func (s *UserService) Login(phoneNumber, pin string) (*domain.User, error) {
	user, err := s.userRepo.FindByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, err
	}
	if user.Pin != pin {
		return nil, errors.New("invalid pin")
	}
	return user, nil
}
