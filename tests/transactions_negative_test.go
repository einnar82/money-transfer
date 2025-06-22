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

func setupTransactionNegativeTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5433 user=postgres password=postgres dbname=testdb sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)
	assert.NoError(t, db.Exec("TRUNCATE TABLE transactions, accounts RESTART IDENTITY CASCADE").Error)
	assert.NoError(t, db.AutoMigrate(&models.Account{}, &models.Transaction{}))
	return db
}

func setupTransactionNegativeRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	tc := controllers.TransactionController{DB: db}
	r.POST("/api/transactions", tc.Create)
	return r
}

func TestCreateTransaction_InvalidAmount(t *testing.T) {
	db := setupTransactionNegativeTestDB(t)
	router := setupTransactionNegativeRouter(db)

	reqBody := []byte(`{
		"source_account_id": "A1",
		"destination_account_id": "A2",
		"amount": "-100"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateTransaction_SourceNotFound(t *testing.T) {
	db := setupTransactionNegativeTestDB(t)
	router := setupTransactionNegativeRouter(db)

	reqBody := transactions.CreateTransactionRequest{
		SourceAccountID:      "A1", // does not exist
		DestinationAccountID: "A2",
		Amount:               "100",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateTransaction_DestinationNotFound(t *testing.T) {
	db := setupTransactionNegativeTestDB(t)
	router := setupTransactionNegativeRouter(db)

	db.Create(&models.Account{AccountID: "A1", Balance: decimal.NewFromInt(1000)})

	reqBody := transactions.CreateTransactionRequest{
		SourceAccountID:      "A1",
		DestinationAccountID: "A2", // not found
		Amount:               "100",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateTransaction_SameAccount(t *testing.T) {
	db := setupTransactionNegativeTestDB(t)
	router := setupTransactionNegativeRouter(db)

	db.Create(&models.Account{AccountID: "A1", Balance: decimal.NewFromInt(1000)})

	reqBody := transactions.CreateTransactionRequest{
		SourceAccountID:      "A1",
		DestinationAccountID: "A1",
		Amount:               "100",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateTransaction_InsufficientFunds(t *testing.T) {
	db := setupTransactionNegativeTestDB(t)
	router := setupTransactionNegativeRouter(db)

	db.Create(&models.Account{AccountID: "A1", Balance: decimal.NewFromInt(50)})
	db.Create(&models.Account{AccountID: "A2", Balance: decimal.NewFromInt(500)})

	reqBody := transactions.CreateTransactionRequest{
		SourceAccountID:      "A1",
		DestinationAccountID: "A2",
		Amount:               "100",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
