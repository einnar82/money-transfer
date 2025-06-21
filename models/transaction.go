package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID                   uint64          `gorm:"primaryKey;column:id" json:"id"`
	SourceAccountID      uint64          `gorm:"not null;column:source_account_id" json:"source_account_id"`
	DestinationAccountID uint64          `gorm:"not null;column:destination_account_id" json:"destination_account_id"`
	Amount               decimal.Decimal `gorm:"type:decimal;not null;column:amount" json:"amount"`
	CreatedAt            time.Time       `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt            time.Time       `gorm:"column:updated_at;not null" json:"updated_at"`

	SourceAccount      Account `gorm:"foreignKey:SourceAccountID;references:ID" json:"source_account"`
	DestinationAccount Account `gorm:"foreignKey:DestinationAccountID;references:ID" json:"destination_account"`
}
