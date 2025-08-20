package ports

import (
	"github.com/google/uuid"
	"hexagonal-go/internal/core/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByPhoneNumber(phoneNumber string) (*domain.User, error)
	FindByID(id uuid.UUID) (*domain.User, error)
	Update(user *domain.User) error
}
