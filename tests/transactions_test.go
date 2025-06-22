package tests

import (
	"bytes"
	"encoding/json"
	"internal-transfers/controllers"
	"internal-transfers/dto/transactions"
	"internal-transfers/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5433 user=postgres password=postgres dbname=testdb sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Account{}, &models.Transaction{})
	assert.NoError(t, err)

	return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	transactionController := controllers.TransactionController{DB: db}

	api := r.Group("/api")
	{
		api.POST("/transactions", transactionController.Create)
	}

	return r
}

func TestCreateTransaction(t *testing.T) {
	db := setupTestDB(t)

	account1 := models.Account{
		AccountID: "1000000001",
		Balance:   decimal.NewFromInt(10000),
	}
	account2 := models.Account{
		AccountID: "1000000002",
		Balance:   decimal.NewFromInt(10000),
	}
	db.Create(&account1)
	db.Create(&account2)

	router := setupRouter(db)

	payload := transactions.CreateTransactionRequest{
		SourceAccountID:      "1000000001",
		DestinationAccountID: "1000000002",
		Amount:               "100",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "100", resp["amount"])
	assert.Equal(t, "1000000001", resp["source_account"].(map[string]interface{})["account_id"])
	assert.Equal(t, "1000000002", resp["destination_account"].(map[string]interface{})["account_id"])

	var updated1, updated2 models.Account
	db.First(&updated1, account1.ID)
	db.First(&updated2, account2.ID)

	assert.Equal(t, "9900", updated1.Balance.String())
	assert.Equal(t, "10100", updated2.Balance.String())
}
