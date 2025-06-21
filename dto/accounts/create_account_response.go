package accounts

type AccountResponse struct {
	ID        uint64 `json:"id"`
	AccountID string `json:"account_id"`
	Balance   string `json:"balance"`
}
