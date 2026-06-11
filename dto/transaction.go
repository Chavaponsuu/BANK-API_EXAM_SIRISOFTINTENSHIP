package dto

import "github.com/krizad/go-gin-api/models"

type TransactionRequest struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type TransactionResponse struct {
	TransactionRef string  `json:"transaction_ref"`
	TrasactionType string  `json:"transaction_type"`
	Amount         float64 `json:"amount"`
	BalanceBefore  float64 `json:"balance_before"`
	BalanceAfter   float64 `json:"balance_after"`
	Description    string  `json:"description"`
	CreatedAt      string  `json:"created_at"`
}

func ToTransactionResponse(m *models.Transaction) *TransactionResponse {
	return &TransactionResponse{
		TransactionRef: m.TransactionRef,
		TrasactionType: m.TransactionType,
		Amount:         m.Amount,
		BalanceBefore:  m.BalanceBefore,
		BalanceAfter:   m.BalanceAfter,
		Description:    m.Description,
		CreatedAt:      m.CreatedAt,
	}
}
