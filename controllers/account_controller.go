package controllers

import (
	"fmt"
	"internal-transfers/dto/accounts"
	"internal-transfers/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AccountController struct {
	DB *gorm.DB
}

func (ac *AccountController) Create(context *gin.Context) {
	var req accounts.CreateAccountRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := models.Account{
		AccountID: req.AccountID,
		Balance:   decimal.NewFromFloat(req.InitialBalance),
	}

	if err := ac.DB.Create(&account).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create an account."})
		return
	}

	context.JSON(http.StatusCreated, accounts.AccountResponse{
		ID:        account.ID,
		AccountID: account.AccountID,
		Balance:   account.Balance.String(),
	})
}

func (ac *AccountController) Show(context *gin.Context) {
	accountID := context.Param("account_id")
	fmt.Println("Account ID:", accountID)

	var account models.Account
	if err := ac.DB.
		Where("account_id = ?", accountID).
		Preload("OutgoingTransactions").
		Preload("IncomingTransactions").
		First(&account).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Account not found!"})
		return
	}

	context.JSON(http.StatusCreated, accounts.AccountResponse{
		ID:        account.ID,
		AccountID: account.AccountID,
		Balance:   account.Balance.String(),
	})
}
