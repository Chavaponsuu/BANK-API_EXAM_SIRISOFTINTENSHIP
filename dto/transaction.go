package dto

type DepositRequest struct {
	Amount  string `json:"amount"`
	Description string `json:"description"`
}

