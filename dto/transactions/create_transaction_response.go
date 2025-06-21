package transactions

import (
	accountdto "internal-transfers/dto/accounts"
	"time"
)

type TransactionResponse struct {
	ID                   uint64                     `json:"id"`
	SourceAccountID      uint64                     `json:"source_account_id"`
	DestinationAccountID uint64                     `json:"destination_account_id"`
	Amount               string                     `json:"amount"`
	CreatedAt            time.Time                  `json:"created_at"`
	UpdatedAt            time.Time                  `json:"updated_at"`
	SourceAccount        accountdto.AccountResponse `json:"source_account"`
	DestinationAccount   accountdto.AccountResponse `json:"destination_account"`
}
