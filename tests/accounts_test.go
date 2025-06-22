package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"internal-transfers/controllers"
	"internal-transfers/dto/accounts"
	"internal-transfers/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupAccountTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost port=5433 user=postgres password=postgres dbname=testdb sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Account{}, &models.Transaction{})
	assert.NoError(t, err)

	return db
}

func setupAccountRouter(db *gorm.DB) *gin.Engine {
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

func TestCreateAccount(t *testing.T) {
	db := setupAccountTestDB(t)
	router := setupAccountRouter(db)

	accountID := fmt.Sprintf("ACC%d", time.Now().UnixNano())

	payload := accounts.CreateAccountRequest{
		AccountID:      accountID,
		InitialBalance: 10000,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, accountID, resp["account_id"])
	assert.Equal(t, "10000", resp["balance"])
}

func TestShowAccount(t *testing.T) {
	db := setupAccountTestDB(t)
	router := setupAccountRouter(db)

	accountID := fmt.Sprintf("ACC%d", time.Now().UnixNano())
	// Seed account without transactions
	account := models.Account{
		AccountID: accountID,
		Balance:   decimal.NewFromInt(10000),
	}
	assert.NoError(t, db.Create(&account).Error)

	req := httptest.NewRequest(http.MethodGet, "/api/accounts/"+accountID, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, accountID, resp["account_id"])
	assert.Equal(t, "10000", resp["balance"])

	// Optional: check transactions are omitted
	assert.Nil(t, resp["outgoing_transactions"])
	assert.Nil(t, resp["incoming_transactions"])
}
