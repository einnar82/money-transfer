package controllers

import (
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

	err = tc.DB.Transaction(func(tx *gorm.DB) error {
		var source, destination models.Account

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("account_id = ?", req.SourceAccountID).
			First(&source).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Source account not found!"})
			return err
		}

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("account_id = ?", req.DestinationAccountID).
			First(&destination).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Destination account not found!"})
			return err
		}

		if source.AccountID == destination.AccountID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot transfer to the same account!"})
			return gorm.ErrInvalidData
		}

		if source.Balance.LessThan(amount) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance!"})
			return gorm.ErrInvalidData
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

		if err := tx.Preload("SourceAccount").
			Preload("DestinationAccount").
			First(&transaction, transaction.ID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed!"})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}
