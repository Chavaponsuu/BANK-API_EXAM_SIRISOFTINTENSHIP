package coverage_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/krizad/go-gin-api/constants"
	"github.com/krizad/go-gin-api/internal/core/domain"
	services "github.com/krizad/go-gin-api/internal/core/services"
)

func TestTransactionService_Deposit(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		accountNumber string
		amount        float64
		description   string
		mockAccountRepo *MockAccountRepository
		mockTxRepo    *MockTransactionRepository
		expectedError error
	}{
		{
			name:          "account not found",
			accountNumber: "1234567890",
			amount:        1000.0,
			description:   "Test deposit",
			mockAccountRepo: &MockAccountRepository{
				BeginTxFunc: func(ctx context.Context) (*sql.Tx, error) {
					return nil, nil
				},
				GetByAccountNumberWithLockFunc: func(ctx context.Context, tx *sql.Tx, accountNumber string) (*domain.Account, error) {
					return nil, nil
				},
			},
			mockTxRepo:    &MockTransactionRepository{},
			expectedError: constants.ErrAccountNotFound,
		},
		{
			name:          "account not active",
			accountNumber: "1234567890",
			amount:        1000.0,
			description:   "Test deposit",
			mockAccountRepo: &MockAccountRepository{
				BeginTxFunc: func(ctx context.Context) (*sql.Tx, error) {
					return nil, nil
				},
				GetByAccountNumberWithLockFunc: func(ctx context.Context, tx *sql.Tx, accountNumber string) (*domain.Account, error) {
					return &domain.Account{
						ID:            1,
						AccountNumber: accountNumber,
						Balance:       500.0,
						Status:        "CLOSED",
					}, nil
				},
			},
			mockTxRepo:    &MockTransactionRepository{},
			expectedError: constants.ErrAccountNotActive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewTransactionService(tt.mockAccountRepo, tt.mockTxRepo)
			result, err := service.Deposit(ctx, tt.accountNumber, tt.amount, tt.description)

			if tt.expectedError != nil {
				if err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result == nil {
					t.Errorf("expected transaction, got nil")
				}
				if result != nil && result.TransactionType != "DEPOSIT" {
					t.Errorf("expected transaction type DEPOSIT, got %s", result.TransactionType)
				}
			}
		})
	}
}

func TestTransactionService_Withdraw(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		accountNumber string
		amount        float64
		description   string
		mockAccountRepo *MockAccountRepository
		mockTxRepo    *MockTransactionRepository
		expectedError error
	}{
		{
			name:          "account not found",
			accountNumber: "1234567890",
			amount:        500.0,
			description:   "Test withdrawal",
			mockAccountRepo: &MockAccountRepository{
				BeginTxFunc: func(ctx context.Context) (*sql.Tx, error) {
					return nil, nil
				},
				GetByAccountNumberWithLockFunc: func(ctx context.Context, tx *sql.Tx, accountNumber string) (*domain.Account, error) {
					return nil, nil
				},
			},
			mockTxRepo:    &MockTransactionRepository{},
			expectedError: constants.ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewTransactionService(tt.mockAccountRepo, tt.mockTxRepo)
			result, err := service.Withdraw(ctx, tt.accountNumber, tt.amount, tt.description)

			if tt.expectedError != nil {
				if err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result == nil {
					t.Errorf("expected transaction, got nil")
				}
				if result != nil && result.TransactionType != "WITHDRAW" {
					t.Errorf("expected transaction type WITHDRAW, got %s", result.TransactionType)
				}
			}
		})
	}
}

func TestTransactionService_GetTransactionHistory(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		accountNumber string
		page          int
		limit         int
		mockAccountRepo *MockAccountRepository
		mockTxRepo    *MockTransactionRepository
		expectedError error
	}{
		{
			name:          "successful get transaction history",
			accountNumber: "1234567890",
			page:          1,
			limit:         10,
			mockAccountRepo: &MockAccountRepository{
				GetByAccountNumberFunc: func(ctx context.Context, accountNumber string) (*domain.Account, error) {
					return &domain.Account{
						ID:            1,
						AccountNumber: accountNumber,
						Status:        "ACTIVE",
					}, nil
				},
			},
			mockTxRepo: &MockTransactionRepository{
				GetTxByAccountIDFunc: func(ctx context.Context, accountID int64, page int, limit int) ([]*domain.Transaction, int, error) {
					return []*domain.Transaction{
						{ID: 1, AccountID: accountID, TransactionType: "DEPOSIT", Amount: 1000.0},
						{ID: 2, AccountID: accountID, TransactionType: "WITHDRAW", Amount: 500.0},
					}, 2, nil
				},
			},
			expectedError: nil,
		},
		{
			name:          "account not found",
			accountNumber: "1234567890",
			page:          1,
			limit:         10,
			mockAccountRepo: &MockAccountRepository{
				GetByAccountNumberFunc: func(ctx context.Context, accountNumber string) (*domain.Account, error) {
					return nil, nil
				},
			},
			mockTxRepo:    &MockTransactionRepository{},
			expectedError: constants.ErrAccountNotFound,
		},
		{
			name:          "database error",
			accountNumber: "1234567890",
			page:          1,
			limit:         10,
			mockAccountRepo: &MockAccountRepository{
				GetByAccountNumberFunc: func(ctx context.Context, accountNumber string) (*domain.Account, error) {
					return nil, errors.New("database error")
				},
			},
			mockTxRepo:    &MockTransactionRepository{},
			expectedError: constants.ErrOperationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewTransactionService(tt.mockAccountRepo, tt.mockTxRepo)
			result, total, err := service.GetTransactionHistory(ctx, tt.accountNumber, tt.page, tt.limit)

			if tt.expectedError != nil {
				if err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result == nil {
					t.Errorf("expected transactions, got nil")
				}
				if total != 2 {
					t.Errorf("expected total 2, got %d", total)
				}
			}
		})
	}
}
