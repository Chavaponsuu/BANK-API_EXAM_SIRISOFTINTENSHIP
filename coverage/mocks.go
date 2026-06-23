package coverage

import (
	"context"
	"database/sql"

	"github.com/krizad/go-gin-api/internal/adapters/http/dto"
	"github.com/krizad/go-gin-api/internal/core/domain"
)

// ============================================================================
// PORTS MOCKS (Domain Layer - Service & Repository Interfaces)
// ============================================================================

// MockAccountService is a mock implementation of ports.AccountService
type MockAccountService struct {
	CreateAccountFunc      func(ctx context.Context, account *domain.Account) (*domain.Account, error)
	GetAccountByNumberFunc func(ctx context.Context, accountNumber string) (*domain.Account, error)
	GetAccountListFunc     func(ctx context.Context, page int, limit int) ([]*domain.Account, int, error)
	CloseAccountFunc       func(ctx context.Context, accountNumber string) (*domain.Account, error)
}

func (m *MockAccountService) CreateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	if m.CreateAccountFunc != nil {
		return m.CreateAccountFunc(ctx, account)
	}
	return &domain.Account{ID: 1, AccountNumber: "1234567890"}, nil
}

func (m *MockAccountService) GetAccountByNumber(ctx context.Context, accountNumber string) (*domain.Account, error) {
	if m.GetAccountByNumberFunc != nil {
		return m.GetAccountByNumberFunc(ctx, accountNumber)
	}
	return &domain.Account{ID: 1, AccountNumber: accountNumber}, nil
}

func (m *MockAccountService) GetAccountList(ctx context.Context, page int, limit int) ([]*domain.Account, int, error) {
	if m.GetAccountListFunc != nil {
		return m.GetAccountListFunc(ctx, page, limit)
	}
	return []*domain.Account{{ID: 1, AccountNumber: "1234567890"}}, 1, nil
}

func (m *MockAccountService) CloseAccount(ctx context.Context, accountNumber string) (*domain.Account, error) {
	if m.CloseAccountFunc != nil {
		return m.CloseAccountFunc(ctx, accountNumber)
	}
	return &domain.Account{ID: 1, AccountNumber: accountNumber, Status: "CLOSED"}, nil
}

// MockAccountRepository is a mock implementation of ports.AccountRepository
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
	return &domain.Account{ID: 1, AccountNumber: "1234567890"}, nil
}

func (m *MockAccountRepository) GetByAccountNumber(ctx context.Context, accountNumber string) (*domain.Account, error) {
	if m.GetByAccountNumberFunc != nil {
		return m.GetByAccountNumberFunc(ctx, accountNumber)
	}
	return &domain.Account{ID: 1, AccountNumber: accountNumber}, nil
}

func (m *MockAccountRepository) GetByAccountNumberWithLock(ctx context.Context, tx *sql.Tx, accountNumber string) (*domain.Account, error) {
	if m.GetByAccountNumberWithLockFunc != nil {
		return m.GetByAccountNumberWithLockFunc(ctx, tx, accountNumber)
	}
	return &domain.Account{ID: 1, AccountNumber: accountNumber}, nil
}

func (m *MockAccountRepository) GetByAccountList(ctx context.Context, page int, limit int) ([]*domain.Account, int, error) {
	if m.GetByAccountListFunc != nil {
		return m.GetByAccountListFunc(ctx, page, limit)
	}
	return []*domain.Account{{ID: 1, AccountNumber: "1234567890"}}, 1, nil
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

// MockTransactionService is a mock implementation of ports.TransactionService
type MockTransactionService struct {
	DepositFunc               func(ctx context.Context, accountNumber string, amount float64, description string) (*domain.Transaction, error)
	WithdrawFunc              func(ctx context.Context, accountNumber string, amount float64, description string) (*domain.Transaction, error)
	GetTransactionHistoryFunc func(ctx context.Context, accountNumber string, page int, limit int) ([]*domain.Transaction, int, error)
	GetAllTransactionFunc     func(ctx context.Context, page int, limit int) ([]*domain.Transaction, int, error)
}

func (m *MockTransactionService) Deposit(ctx context.Context, accountNumber string, amount float64, description string) (*domain.Transaction, error) {
	if m.DepositFunc != nil {
		return m.DepositFunc(ctx, accountNumber, amount, description)
	}
	return &domain.Transaction{ID: 1, AccountID: 1, TransactionType: "DEPOSIT", Amount: amount}, nil
}

func (m *MockTransactionService) Withdraw(ctx context.Context, accountNumber string, amount float64, description string) (*domain.Transaction, error) {
	if m.WithdrawFunc != nil {
		return m.WithdrawFunc(ctx, accountNumber, amount, description)
	}
	return &domain.Transaction{ID: 1, AccountID: 1, TransactionType: "WITHDRAW", Amount: amount}, nil
}

func (m *MockTransactionService) GetTransactionHistory(ctx context.Context, accountNumber string, page int, limit int) ([]*domain.Transaction, int, error) {
	if m.GetTransactionHistoryFunc != nil {
		return m.GetTransactionHistoryFunc(ctx, accountNumber, page, limit)
	}
	return []*domain.Transaction{{ID: 1, TransactionType: "DEPOSIT", Amount: 1000.0}}, 1, nil
}

func (m *MockTransactionService) GetAllTransaction(ctx context.Context, page int, limit int) ([]*domain.Transaction, int, error) {
	if m.GetAllTransactionFunc != nil {
		return m.GetAllTransactionFunc(ctx, page, limit)
	}
	return []*domain.Transaction{{ID: 1, TransactionType: "DEPOSIT", Amount: 1000.0}, {ID: 2, TransactionType: "WITHDRAW", Amount: 500.0}}, 2, nil
}

// MockTransactionRepository is a mock implementation of ports.TransactionRepository
type MockTransactionRepository struct {
	CreateWithTxFunc           func(ctx context.Context, tx *sql.Tx, transaction *domain.Transaction) (*domain.Transaction, error)
	GenerateTransactionRefFunc func(ctx context.Context, tx *sql.Tx) (string, error)
	GetTxByAccountIDFunc       func(ctx context.Context, accountID int64, page int, limit int) ([]*domain.Transaction, int, error)
	GetAllTransactionFunc      func(ctx context.Context, limit int, offset int) ([]*domain.Transaction, int, error)
}

func (m *MockTransactionRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, transaction *domain.Transaction) (*domain.Transaction, error) {
	if m.CreateWithTxFunc != nil {
		return m.CreateWithTxFunc(ctx, tx, transaction)
	}
	return &domain.Transaction{ID: 1, TransactionType: "DEPOSIT", Amount: 1000.0}, nil
}

func (m *MockTransactionRepository) GenerateTransactionRef(ctx context.Context, tx *sql.Tx) (string, error) {
	if m.GenerateTransactionRefFunc != nil {
		return m.GenerateTransactionRefFunc(ctx, tx)
	}
	return "TXN-20250623-001", nil
}

func (m *MockTransactionRepository) GetTxByAccountID(ctx context.Context, accountID int64, page int, limit int) ([]*domain.Transaction, int, error) {
	if m.GetTxByAccountIDFunc != nil {
		return m.GetTxByAccountIDFunc(ctx, accountID, page, limit)
	}
	return []*domain.Transaction{{ID: 1, TransactionType: "DEPOSIT", Amount: 1000.0}}, 1, nil
}

func (m *MockTransactionRepository) GetAllTransaction(ctx context.Context, limit int, offset int) ([]*domain.Transaction, int, error) {
	if m.GetAllTransactionFunc != nil {
		return m.GetAllTransactionFunc(ctx, limit, offset)
	}
	return []*domain.Transaction{{ID: 1, TransactionType: "DEPOSIT", Amount: 1000.0}}, 1, nil
}

// ============================================================================
// SERVICES LAYER MOCKS
// ============================================================================

// MockExampleService is a mock implementation of services.ExampleService
type MockExampleService struct {
	CreateExampleFunc     func(ctx context.Context, req dto.CreateExampleRequest) (*domain.Example, error)
	GetExampleFunc        func(ctx context.Context, id int64) (*domain.Example, error)
	GetExampleByEmailFunc func(ctx context.Context, email string) (*domain.Example, error)
	ListExamplesFunc      func(ctx context.Context) ([]domain.Example, error)
	UpdateExampleFunc     func(ctx context.Context, id int64, req dto.UpdateExampleRequest) (*domain.Example, error)
	PatchExampleFunc      func(ctx context.Context, id int64, req dto.PatchExampleRequest) (*domain.Example, error)
	DeleteExampleFunc     func(ctx context.Context, id int64) (bool, error)
}

func (m *MockExampleService) CreateExample(ctx context.Context, req dto.CreateExampleRequest) (*domain.Example, error) {
	if m.CreateExampleFunc != nil {
		return m.CreateExampleFunc(ctx, req)
	}
	return &domain.Example{ID: 1, Name: req.Name, Email: req.Email}, nil
}

func (m *MockExampleService) GetExample(ctx context.Context, id int64) (*domain.Example, error) {
	if m.GetExampleFunc != nil {
		return m.GetExampleFunc(ctx, id)
	}
	return &domain.Example{ID: id, Name: "Test", Email: "test@example.com"}, nil
}

func (m *MockExampleService) GetExampleByEmail(ctx context.Context, email string) (*domain.Example, error) {
	if m.GetExampleByEmailFunc != nil {
		return m.GetExampleByEmailFunc(ctx, email)
	}
	return &domain.Example{ID: 1, Name: "Test", Email: email}, nil
}

func (m *MockExampleService) ListExamples(ctx context.Context) ([]domain.Example, error) {
	if m.ListExamplesFunc != nil {
		return m.ListExamplesFunc(ctx)
	}
	return []domain.Example{
		{ID: 1, Name: "Test1", Email: "test1@example.com"},
		{ID: 2, Name: "Test2", Email: "test2@example.com"},
	}, nil
}

func (m *MockExampleService) UpdateExample(ctx context.Context, id int64, req dto.UpdateExampleRequest) (*domain.Example, error) {
	if m.UpdateExampleFunc != nil {
		return m.UpdateExampleFunc(ctx, id, req)
	}
	return &domain.Example{ID: id, Name: req.Name, Email: req.Email}, nil
}

func (m *MockExampleService) PatchExample(ctx context.Context, id int64, req dto.PatchExampleRequest) (*domain.Example, error) {
	if m.PatchExampleFunc != nil {
		return m.PatchExampleFunc(ctx, id, req)
	}
	return &domain.Example{ID: id, Name: "Patched", Email: "patched@example.com"}, nil
}

func (m *MockExampleService) DeleteExample(ctx context.Context, id int64) (bool, error) {
	if m.DeleteExampleFunc != nil {
		return m.DeleteExampleFunc(ctx, id)
	}
	return true, nil
}

// ============================================================================
// REPOSITORY LAYER MOCKS
// ============================================================================

// MockExampleRepository is a mock implementation of repositories.ExampleRepository
type MockExampleRepository struct {
	CreateFunc     func(ctx context.Context, name, email string) (*domain.Example, error)
	GetByIDFunc    func(ctx context.Context, id int64) (*domain.Example, error)
	GetByEmailFunc func(ctx context.Context, email string) (*domain.Example, error)
	ListFunc       func(ctx context.Context) ([]domain.Example, error)
	UpdateFunc     func(ctx context.Context, id int64, name, email string) (*domain.Example, error)
	DeleteFunc     func(ctx context.Context, id int64) (bool, error)
	CountFunc      func(ctx context.Context) (int, error)
}

func (m *MockExampleRepository) Create(ctx context.Context, name, email string) (*domain.Example, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, name, email)
	}
	return &domain.Example{ID: 1, Name: name, Email: email}, nil
}

func (m *MockExampleRepository) GetByID(ctx context.Context, id int64) (*domain.Example, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return &domain.Example{ID: id, Name: "Test", Email: "test@example.com"}, nil
}

func (m *MockExampleRepository) GetByEmail(ctx context.Context, email string) (*domain.Example, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return &domain.Example{ID: 1, Name: "Test", Email: email}, nil
}

func (m *MockExampleRepository) List(ctx context.Context) ([]domain.Example, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return []domain.Example{
		{ID: 1, Name: "Test1", Email: "test1@example.com"},
		{ID: 2, Name: "Test2", Email: "test2@example.com"},
	}, nil
}

func (m *MockExampleRepository) Update(ctx context.Context, id int64, name, email string) (*domain.Example, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, name, email)
	}
	return &domain.Example{ID: id, Name: name, Email: email}, nil
}

func (m *MockExampleRepository) Delete(ctx context.Context, id int64) (bool, error) {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return true, nil
}

func (m *MockExampleRepository) Count(ctx context.Context) (int, error) {
	if m.CountFunc != nil {
		return m.CountFunc(ctx)
	}
	return 2, nil
}

// ============================================================================
// DATABASE TRANSACTION MOCK
// ============================================================================

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
