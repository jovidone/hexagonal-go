package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"hexagonal-go/internal/core/services"
	"net/http"
)

type TransactionHandler struct {
	transactionService services.TransactionService
}

func NewTransactionHandler(transactionService services.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) Deposit(c *gin.Context) {
	var request struct {
		UserID  string  `json:"user_id"`
		Amount  float64 `json:"amount"`
		Remarks string  `json:"remarks"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	tx, err := h.transactionService.Deposit(userID, request.Amount, request.Remarks)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": tx})
}

func (h *TransactionHandler) Withdraw(c *gin.Context) {
	var request struct {
		UserID  string  `json:"user_id"`
		Amount  float64 `json:"amount"`
		Remarks string  `json:"remarks"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	tx, err := h.transactionService.Withdraw(userID, request.Amount, request.Remarks)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": tx})
}

func (h *TransactionHandler) Transfer(c *gin.Context) {
	var request struct {
		FromID  string  `json:"from_id"`
		ToID    string  `json:"to_id"`
		Amount  float64 `json:"amount"`
		Remarks string  `json:"remarks"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	fromID, err := uuid.Parse(request.FromID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from_id"})
		return
	}
	toID, err := uuid.Parse(request.ToID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to_id"})
		return
	}
	debitTx, creditTx, err := h.transactionService.Transfer(fromID, toID, request.Amount, request.Remarks)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": gin.H{"debit": debitTx, "credit": creditTx}})
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	txs, err := h.transactionService.GetTransactionsByUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "result": txs})
}
