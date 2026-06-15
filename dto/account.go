package dto

import (
	"github.com/krizad/go-gin-api/models"
)

// Request DTOs
type CreateAccountRequest struct {
	OwnerName      string  `json:"owner_name" binding:"required"`
	CitizenID      string  `json:"citizen_id"  binding:"required,len=13,numeric"`
	PhoneNumber    string  `json:"phone_number" binding:"required"`
	AccountType    string  `json:"account_type" binding:"required,oneof=SAVING CURRENT"`
	InitialBalance float64 `json:"initial_balance" binding:"gte=0"`
}

func ToCreateAccountRequest(accountReq *CreateAccountRequest) *models.Account {
	return &models.Account{
		OwnerName:   accountReq.OwnerName,
		CitizenID:   accountReq.CitizenID,
		PhoneNumber: accountReq.PhoneNumber,
		AccountType: accountReq.AccountType,
		Balance:     accountReq.InitialBalance,
	}

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

type AccountNumberURI struct {
	AccountNumber string `uri:"account_number" binding:"required,len=10,numeric"`
}
