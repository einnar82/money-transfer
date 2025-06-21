package routes

import (
	"internal-transfers/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	account := controllers.AccountController{DB: db}
	transaction := controllers.TransactionController{DB: db}

	api := router.Group("/api")
	{
		api.POST("/accounts", account.Create)
		api.GET("/accounts/:account_id", account.Show)
		api.POST("/transactions", transaction.Create)
	}
}
