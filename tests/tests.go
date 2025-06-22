package tests

import (
	"log"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	dsn := "host=localhost port=5433 user=postgres password=postgres dbname=testdb sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test DB: %v", err)
	}

	err = db.Exec("TRUNCATE TABLE transactions, accounts RESTART IDENTITY CASCADE").Error
	if err != nil {
		log.Fatalf("Failed to truncate tables: %v", err)
	}

	code := m.Run()
	os.Exit(code)
}
