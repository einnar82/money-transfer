package accounts

type CreateAccountRequest struct {
	InitialBalance float64 `json:"initial_balance" binding:"required,gte=0"`
	AccountID      string  `json:"account_id" binding:"required"`
}
