package models

type Account struct {
	ID            int64   `json:"id"`
	AccountNumber string  `json:"account_number" binding:"required,len=10"`
	OwnerName     string  `json:"owner_name" binding:"required"`
	CitizenID     string  `json:"citizen_id" binding:"required,len=13"`
	PhoneNumber   string  `json:"phone_number" binding:"required"`
	AccountType   string  `json:"account_type" binding:"required"`
	Balance       float64 `json:"balance"`
	Status        string  `json:"status" binding:"oneof=ACTIVE CLOSED"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}
