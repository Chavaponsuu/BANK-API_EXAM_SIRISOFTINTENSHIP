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

// MockAccountRepository is a mock implementation of AccountRepository
type MockAccountRepository struct {
	CheckCitizenIDExistsFunc       func(ctx context.Context, citizenID string) (bool, error)
	BeginTxFunc                    func(ctx context.Context) (*sql.Tx, error)
	GenerateAccountNumberFunc      func(ctx context.Context, tx *sql.Tx) (string, error)
	CreateWithTxFunc               func(ctx context.Context, tx *sql.Tx, account *domain.Account) (*domain.Account, error)
	GetByAccountNumberFunc         func(ctx context.Context, accountNumber string) (*domain.Account, error)
	GetByAccountNumberWithLockFunc func(ctx context.Context, tx *sql.Tx, accountNumber string) (*domain.Account, error)
	GetByAccountListFunc           func(ctx context.Context, page int, limit int) ([]*domain.Account, int, error)
	UpdateBalanceWithTxFunc        func(ctx context.Context, tx *sql.Tx, accountID int64, newBalance float64) error
	UpdateStatusFunc               func(ctx context.Context, tx *sql.Tx, accountNumber string, status string) error
}

func (m *MockAccountRepository) CheckCitizenIDExists(ctx context.Context, citizenID string) (bool, error) {
	if m.CheckCitizenIDExistsFunc != nil {
		return m.CheckCitizenIDExistsFunc(ctx, citizenID)
	}
	return false, nil
}

func (m *MockAccountRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	if m.BeginTxFunc != nil {
		return m.BeginTxFunc(ctx)
	}
	return nil, nil
}

func (m *MockAccountRepository) GenerateAccountNumber(ctx context.Context, tx *sql.Tx) (string, error) {
	if m.GenerateAccountNumberFunc != nil {
		return m.GenerateAccountNumberFunc(ctx, tx)
	}
	return "1234567890", nil
}

func (m *MockAccountRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, account *domain.Account) (*domain.Account, error) {
	if m.CreateWithTxFunc != nil {
		return m.CreateWithTxFunc(ctx, tx, account)
	}
	return account, nil
}

func (m *MockAccountRepository) GetByAccountNumber(ctx context.Context, accountNumber string) (*domain.Account, error) {
	if m.GetByAccountNumberFunc != nil {
		return m.GetByAccountNumberFunc(ctx, accountNumber)
	}
	return nil, nil
}

func (m *MockAccountRepository) GetByAccountNumberWithLock(ctx context.Context, tx *sql.Tx, accountNumber string) (*domain.Account, error) {
	if m.GetByAccountNumberWithLockFunc != nil {
		return m.GetByAccountNumberWithLockFunc(ctx, tx, accountNumber)
	}
	return nil, nil
}

func (m *MockAccountRepository) GetByAccountList(ctx context.Context, page int, limit int) ([]*domain.Account, int, error) {
	if m.GetByAccountListFunc != nil {
		return m.GetByAccountListFunc(ctx, page, limit)
	}
	return []*domain.Account{}, 0, nil
}

func (m *MockAccountRepository) UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, accountID int64, newBalance float64) error {
	if m.UpdateBalanceWithTxFunc != nil {
		return m.UpdateBalanceWithTxFunc(ctx, tx, accountID, newBalance)
	}
	return nil
}

func (m *MockAccountRepository) UpdateStatus(ctx context.Context, tx *sql.Tx, accountNumber string, status string) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, tx, accountNumber, status)
	}
	return nil
}

// MockTransactionRepository is a mock implementation of TransactionRepository
type MockTransactionRepository struct {
	CreateWithTxFunc      func(ctx context.Context, tx *sql.Tx, transaction *domain.Transaction) (*domain.Transaction, error)
	GetTxByAccountIDFunc  func(ctx context.Context, accountID int64, page int, limit int) ([]*domain.Transaction, int, error)
	GetAllTransactionFunc func(ctx context.Context, limit int, offset int) ([]*domain.Transaction, int, error)
}

func (m *MockTransactionRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, transaction *domain.Transaction) (*domain.Transaction, error) {
	if m.CreateWithTxFunc != nil {
		return m.CreateWithTxFunc(ctx, tx, transaction)
	}
	return transaction, nil
}

func (m *MockTransactionRepository) GenerateTransactionRef(ctx context.Context, tx *sql.Tx) (string, error) {
	return "TXN123456", nil
}

func (m *MockTransactionRepository) GetTxByAccountID(ctx context.Context, accountID int64, page int, limit int) ([]*domain.Transaction, int, error) {
	if m.GetTxByAccountIDFunc != nil {
		return m.GetTxByAccountIDFunc(ctx, accountID, page, limit)
	}
	return []*domain.Transaction{}, 0, nil
}

func (m *MockTransactionRepository) GetAllTransaction(ctx context.Context, limit int, offset int) ([]*domain.Transaction, int, error) {
	if m.GetAllTransactionFunc != nil {
		return m.GetAllTransactionFunc(ctx, limit, offset)
	}
	return []*domain.Transaction{}, 0, nil
}

// MockTx is a mock implementation of sql.Tx
type MockTx struct {
	CommitFunc   func() error
	RollbackFunc func() error
}

func (m *MockTx) Commit() error {
	if m.CommitFunc != nil {
		return m.CommitFunc()
	}
	return nil
}

func (m *MockTx) Rollback() error {
	if m.RollbackFunc != nil {
		return m.RollbackFunc()
	}
	return nil
}

// Implement other sql.Tx methods to satisfy the interface
func (m *MockTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (m *MockTx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (m *MockTx) QueryRow(query string, args ...interface{}) *sql.Row {
	return &sql.Row{}
}

func (m *MockTx) Prepare(query string) (*sql.Stmt, error) {
	return nil, nil
}

func (m *MockTx) Stmt(stmt *sql.Stmt) *sql.Stmt {
	return stmt
}

func TestAccountService_CreateAccount(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		account       *domain.Account
		mockRepo      *MockAccountRepository
		mockTxRepo    *MockTransactionRepository
		expectedError error
	}{
		{
			name: "citizen ID already exists",
			account: &domain.Account{
				CitizenID:   "1234567890123",
				OwnerName:   "John Doe",
				PhoneNumber: "0812345678",
				AccountType: "SAVINGS",
				Balance:     1000.0,
			},
			mockRepo: &MockAccountRepository{
				CheckCitizenIDExistsFunc: func(ctx context.Context, citizenID string) (bool, error) {
					return true, nil
				},
			},
			expectedError: constants.ErrCitizenIDExists,
		},
		{
			name: "database error on check citizen ID",
			account: &domain.Account{
				CitizenID:   "1234567890123",
				OwnerName:   "John Doe",
				PhoneNumber: "0812345678",
				AccountType: "SAVINGS",
				Balance:     1000.0,
			},
			mockRepo: &MockAccountRepository{
				CheckCitizenIDExistsFunc: func(ctx context.Context, citizenID string) (bool, error) {
					return false, errors.New("database error")
				},
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewAccountService(tt.mockRepo, tt.mockTxRepo)
			result, err := service.CreateAccount(ctx, tt.account)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				}
				if err != tt.expectedError && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result == nil {
					t.Errorf("expected account, got nil")
				}
			}
		})
	}
}

func TestAccountService_GetAccountByNumber(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		accountNumber string
		mockRepo      *MockAccountRepository
		mockTxRepo    *MockTransactionRepository
		expectedError error
	}{
		{
			name:          "successful get account",
			accountNumber: "1234567890",
			mockRepo: &MockAccountRepository{
				GetByAccountNumberFunc: func(ctx context.Context, accountNumber string) (*domain.Account, error) {
					return &domain.Account{
						ID:            1,
						AccountNumber: accountNumber,
						OwnerName:     "John Doe",
						Status:        "ACTIVE",
					}, nil
				},
			},
			mockTxRepo:    &MockTransactionRepository{},
			expectedError: nil,
		},
		{
			name:          "account not found",
			accountNumber: "1234567890",
			mockRepo: &MockAccountRepository{
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
			mockRepo: &MockAccountRepository{
				GetByAccountNumberFunc: func(ctx context.Context, accountNumber string) (*domain.Account, error) {
					return nil, errors.New("database error")
				},
			},
			mockTxRepo:    &MockTransactionRepository{},
			expectedError: constants.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewAccountService(tt.mockRepo, tt.mockTxRepo)
			result, err := service.GetAccountByNumber(ctx, tt.accountNumber)

			if tt.expectedError != nil {
				if err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result == nil {
					t.Errorf("expected account, got nil")
				}
			}
		})
	}
}

func TestAccountService_GetAccountList(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		page          int
		limit         int
		mockRepo      *MockAccountRepository
		mockTxRepo    *MockTransactionRepository
		expectedError error
	}{
		{
			name:  "successful get account list",
			page:  1,
			limit: 10,
			mockRepo: &MockAccountRepository{
				GetByAccountListFunc: func(ctx context.Context, page int, limit int) ([]*domain.Account, int, error) {
					return []*domain.Account{
						{ID: 1, AccountNumber: "1234567890"},
						{ID: 2, AccountNumber: "0987654321"},
					}, 2, nil
				},
			},
			mockTxRepo:    &MockTransactionRepository{},
			expectedError: nil,
		},
		{
			name:  "database error",
			page:  1,
			limit: 10,
			mockRepo: &MockAccountRepository{
				GetByAccountListFunc: func(ctx context.Context, page int, limit int) ([]*domain.Account, int, error) {
					return nil, 0, errors.New("database error")
				},
			},
			mockTxRepo:    &MockTransactionRepository{},
			expectedError: constants.ErrOperationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewAccountService(tt.mockRepo, tt.mockTxRepo)
			result, total, err := service.GetAccountList(ctx, tt.page, tt.limit)

			if tt.expectedError != nil {
				if err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result == nil {
					t.Errorf("expected accounts, got nil")
				}
				if total != 2 {
					t.Errorf("expected total 2, got %d", total)
				}
			}
		})
	}
}

func TestAccountService_CloseAccount(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		accountNumber string
		mockRepo      *MockAccountRepository
		mockTxRepo    *MockTransactionRepository
		expectedError error
	}{
		{
			name:          "account not found",
			accountNumber: "1234567890",
			mockRepo: &MockAccountRepository{
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
			service := services.NewAccountService(tt.mockRepo, tt.mockTxRepo)
			result, err := service.CloseAccount(ctx, tt.accountNumber)

			if tt.expectedError != nil {
				if err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result == nil {
					t.Errorf("expected account, got nil")
				}
				if result != nil && result.Status != "CLOSED" {
					t.Errorf("expected status CLOSED, got %s", result.Status)
				}
			}
		})
	}
}
