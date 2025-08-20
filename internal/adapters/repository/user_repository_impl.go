package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"hexagonal-go/internal/core/domain"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) *UserRepositoryImpl {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	return &UserRepositoryImpl{db: db}
}

//func NewUserRepositoryImpl(db *gorm.DB) *UserRepositoryImpl {
//	return &UserRepositoryImpl{db: db}
//}

func (r *UserRepositoryImpl) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) FindByPhoneNumber(phoneNumber string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("phone_number = ?", phoneNumber).First(&user).Error
	return &user, err
}

func (r *UserRepositoryImpl) FindByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("user_id = ?", id).First(&user).Error
	return &user, err
}

func (r *UserRepositoryImpl) Update(user *domain.User) error {
	return r.db.Save(user).Error
}
