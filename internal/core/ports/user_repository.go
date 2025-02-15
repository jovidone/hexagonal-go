package ports

import (
	"hexagonal-go/internal/core/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByPhoneNumber(phoneNumber string) (*domain.User, error)
}
