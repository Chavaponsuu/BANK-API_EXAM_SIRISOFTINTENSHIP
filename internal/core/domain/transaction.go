package domain

type Transaction struct {
	ID              int64
	AccountID       int64
	TransactionRef  string
	TransactionType string
	Amount          float64
	BalanceBefore   float64
	BalanceAfter    float64
	Description     string
	CreatedAt       string
}
