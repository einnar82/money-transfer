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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAccountTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
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

	payload := accounts.CreateAccountRequest{
		AccountID:      "1000000001",
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

	assert.Equal(t, "1000000001", resp["account_id"])
	assert.Equal(t, "10000", resp["balance"])
}

func TestShowAccount(t *testing.T) {
	db := setupAccountTestDB(t)
	router := setupAccountRouter(db)

	// Seed account without transactions
	account := models.Account{
		AccountID: "1000000001",
		Balance:   decimal.NewFromInt(10000),
	}
	assert.NoError(t, db.Create(&account).Error)

	req := httptest.NewRequest(http.MethodGet, "/api/accounts/1000000001", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "1000000001", resp["account_id"])
	assert.Equal(t, "10000", resp["balance"])

	// Optional: check transactions are omitted
	assert.Nil(t, resp["outgoing_transactions"])
	assert.Nil(t, resp["incoming_transactions"])
}
