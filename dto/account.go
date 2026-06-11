package dto

import (
	"github.com/krizad/go-gin-api/models"
)

// Request DTOs
type CreateAccountRequest struct {
	OwnerName      string  `json:"owner_name" binding:"required" example:"John Doe"`
	CitizenID      string  `json:"citizen_id" binding:"required,len=13" example:"1234567890123"`
	PhoneNumber    string  `json:"phone_number" binding:"required" example:"0812345678"`
	AccountType    string  `json:"account_type" binding:"required,oneof=SAVING CURRENT" example:"SAVING"`
	InitialBalance float64 `json:"initial_balance" binding:"gte=0" example:"1000.00"`
}

// Response DTOs
type AccountResponse struct {
	AccountNumber string  `json:"account_number" `
	OwnerName     string  `json:"owner_name" example:"John Doe"`
	AccountType   string  `json:"account_type" example:"SAVING"`
	Balance       float64 `json:"balance" example:"1000.00"`
	Status        string  `json:"status" example:"ACTIVE"`
}

type AccountListResponse []AccountResponse

type AccountDeleteResponse struct {
	ID int64 `json:"id"`
}

// Mapper functions
func ToAccountResponse(m *models.Account) *AccountResponse {
	return &AccountResponse{
		AccountNumber: m.AccountNumber,
		OwnerName:     m.OwnerName,
		AccountType:   m.AccountType,
		Balance:       m.Balance,
		Status:        m.Status,
	}
}

func ToAccountListResponse(accounts []models.Account) AccountListResponse {
	list := make(AccountListResponse, 0, len(accounts))
	for _, a := range accounts {
		list = append(list, *ToAccountResponse(&a))
	}
	return list
}

type CloseAccountResponse struct {
	AccountNumber string `json:"account_number"`
	Status        string `json:"status"`
}

func ToCloseAccountResponse(m *models.Account) *CloseAccountResponse {
	return &CloseAccountResponse{
		AccountNumber: m.AccountNumber,
		Status:        m.Status,
	}
}
