package tests

import (
	"bytes"
	"encoding/json"
	"internal-transfers/controllers"
	"internal-transfers/dto/accounts"
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

func setupNegativeTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5433 user=postgres password=postgres dbname=testdb sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)
	err = db.Exec("TRUNCATE TABLE transactions, accounts RESTART IDENTITY CASCADE").Error
	assert.NoError(t, err)
	err = db.AutoMigrate(&models.Account{}, &models.Transaction{})
	assert.NoError(t, err)
	return db
}

func setupNegativeTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	accountController := controllers.AccountController{DB: db}
	api := r.Group("/api")
	{
		api.POST("/accounts", accountController.Create)
		api.GET("/accounts/:account_id", accountController.Show)
	}
	return r
}

func TestCreateAccount_MissingFields(t *testing.T) {
	db := setupNegativeTestDB(t)
	router := setupNegativeTestRouter(db)

	body := []byte(`{"initial_balance": 1000}`) // Missing account_id
	req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateAccount_DuplicateID(t *testing.T) {
	db := setupNegativeTestDB(t)
	router := setupNegativeTestRouter(db)

	account := models.Account{
		AccountID: "DUPLICATE_ID",
		Balance:   decimal.NewFromInt(10000),
	}
	db.Create(&account)

	payload := accounts.CreateAccountRequest{
		AccountID:      "DUPLICATE_ID",
		InitialBalance: 5000,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestShowAccount_NotFound(t *testing.T) {
	db := setupNegativeTestDB(t)
	router := setupNegativeTestRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/api/accounts/NOT_FOUND", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
