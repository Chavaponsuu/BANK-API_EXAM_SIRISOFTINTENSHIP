package models

type Transaction struct {
	ID              int64   `json:"id"`
	AccountID       int64   `json:"account_id" binding:"required"`
	TransactionRef  string  `json:"transaction_ref"`
	TransactionType string  `json:"transaction_type" binding:"required"`
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	BalanceBefore   float64 `json:"balance_before" binding:"required,gt=0"`
	BalanceAfter    float64 `json:"balance_after" binding:"required,gt=0"`
	Description     string  `json:"description"`
	CreatedAt       string  `json:"created_at"`
}
