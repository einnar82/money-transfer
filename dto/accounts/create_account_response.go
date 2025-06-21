package accounts

import "time"

type AccountResponse struct {
	ID        uint64    `json:"id"`
	AccountID string    `json:"account_id"`
	Balance   string    `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}
