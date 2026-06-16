package domain

type Account struct {
	ID            int64
	AccountNumber string
	OwnerName     string
	CitizenID     string
	PhoneNumber   string
	AccountType   string
	Balance       float64
	Status        string
	CreatedAt     string
	UpdatedAt     string
}
