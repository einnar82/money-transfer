package controllers

import (
	"errors"
	"internal-transfers/dto/transactions"
	"internal-transfers/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionController struct {
	DB *gorm.DB
}

func (tc *TransactionController) Create(c *gin.Context) {
	var req transactions.CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.LessThanOrEqual(decimal.Zero) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction amount!"})
		return
	}

	var transaction models.Transaction

	var (
		ErrSourceNotFound      = gorm.ErrRecordNotFound
		ErrDestinationNotFound = gorm.ErrRecordNotFound
		ErrSameAccount         = errors.New("same account")
		ErrInsufficientFunds   = errors.New("insufficient funds")
	)

	err = tc.DB.Transaction(func(tx *gorm.DB) error {
		var source, destination models.Account

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("account_id = ?", req.SourceAccountID).
			First(&source).Error; err != nil {
			return ErrSourceNotFound
		}

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("account_id = ?", req.DestinationAccountID).
			First(&destination).Error; err != nil {
			return ErrDestinationNotFound
		}

		if source.AccountID == destination.AccountID {
			return ErrSameAccount
		}

		if source.Balance.LessThan(amount) {
			return ErrInsufficientFunds
		}

		source.Balance = source.Balance.Sub(amount)
		destination.Balance = destination.Balance.Add(amount)

		if err := tx.Save(&source).Error; err != nil {
			return err
		}
		if err := tx.Save(&destination).Error; err != nil {
			return err
		}

		transaction = models.Transaction{
			SourceAccountID:      source.ID,
			DestinationAccountID: destination.ID,
			Amount:               amount,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return tx.
			Preload("SourceAccount").
			Preload("DestinationAccount").
			First(&transaction, transaction.ID).Error
	})

	switch {
	case errors.Is(err, ErrSourceNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Source account not found!"})
	case errors.Is(err, ErrDestinationNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Destination account not found!"})
	case errors.Is(err, ErrSameAccount):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot transfer to the same account!"})
	case errors.Is(err, ErrInsufficientFunds):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance!"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed!"})
	default:
		c.JSON(http.StatusCreated, transaction)
	}
}
