package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID        uint64          `gorm:"primaryKey;column:id" json:"id"`
	AccountID string          `gorm:"type:varchar(50);unique;not null;column:account_id" json:"account_id"`
	Balance   decimal.Decimal `gorm:"type:decimal;not null;column:balance" json:"balance"`
	CreatedAt time.Time       `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at;not null" json:"updated_at"`

	OutgoingTransactions []Transaction `gorm:"foreignKey:SourceAccountID" json:"outgoing_transactions"`
	IncomingTransactions []Transaction `gorm:"foreignKey:DestinationAccountID" json:"incoming_transactions"`
}
