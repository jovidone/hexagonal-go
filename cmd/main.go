package main

import (
	"github.com/gin-gonic/gin"
	"hexagonal-go/internal/adapters/http"
	"hexagonal-go/internal/adapters/repository"
	"hexagonal-go/internal/config"
	"hexagonal-go/internal/core/services"
)

func main() {
	// Koneksi ke database
	db, err := config.ConnectDB()
	if err != nil {
		panic("failed to connect database")
	}

	// Inisialisasi repository
	userRepo := repository.NewUserRepositoryImpl(db)
	transactionRepo := repository.NewTransactionRepositoryImpl(db)

	// Inisialisasi service
	userService := services.NewUserService(userRepo)
	transactionService := services.NewTransactionService(transactionRepo, db)

	// Inisialisasi handler
	userHandler := http.NewUserHandler(*userService)
	transactionHandler := http.NewTransactionHandler(*transactionService)

	// Setup router menggunakan Gin
	r := gin.Default()

	// Endpoint user
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// Endpoint transaction
	r.POST("/deposit", transactionHandler.Deposit)
	r.POST("/withdraw", transactionHandler.Withdraw)

	// Jalankan server pada port 8080
	r.Run(":8080")
}
