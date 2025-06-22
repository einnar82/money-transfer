package controllers

import (
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

	context.JSON(http.StatusCreated, account)
}

func (ac *AccountController) Show(c *gin.Context) {
	accountID := c.Param("account_id")

	var account models.Account

	if err := ac.DB.
		Where("account_id = ?", accountID).
		First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found!"})
		return
	}

	c.JSON(http.StatusOK, account)
}
