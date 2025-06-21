package transactions

type CreateTransactionRequest struct {
	SourceAccountID      string `json:"source_account_id" binding:"required"`
	DestinationAccountID string `json:"destination_account_id" binding:"required"`
	Amount               string `json:"amount" binding:"required"`
}
