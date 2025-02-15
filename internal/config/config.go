package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"hexagonal-go/internal/core/domain"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := "host=localhost user=admin password=root dbname=hexago port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Aktifkan ekstensi uuid-ossp
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		return nil, fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}

	// Auto migrate tabel
	if err := db.AutoMigrate(&domain.User{}, &domain.Transaction{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
