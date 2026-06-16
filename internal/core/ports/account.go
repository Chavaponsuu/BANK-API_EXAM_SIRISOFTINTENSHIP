package ports

import (
	"context"
	"database/sql"

	"github.com/krizad/go-gin-api/internal/core/domain"
)

type AccountService interface {
	CreateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error)
	GetAccountByNumber(ctx context.Context, accountNumber string) (*domain.Account, error)
	GetAccountList(ctx context.Context, page int, limit int) ([]*domain.Account, int, error)
	CloseAccount(ctx context.Context, accountNumber string) (*domain.Account, error)
}

type AccountRepository interface {
	// Create(ctx context.Context, account *domain.Account) (*domain.Account, error)
	CreateWithTx(ctx context.Context, tx *sql.Tx, account *domain.Account) (*domain.Account, error)
	CheckCitizenIDExists(ctx context.Context, citizenID string) (bool, error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
	GetByAccountNumber(ctx context.Context, accountNumber string) (*domain.Account, error)
	GetByAccountNumberWithLock(ctx context.Context, tx *sql.Tx, accountNumber string) (*domain.Account, error)
	GetByAccountList(ctx context.Context, page int, limit int) ([]*domain.Account, int, error)
	UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, accountID int64, newBalance float64) error
	UpdateStatus(ctx context.Context, tx *sql.Tx, accountNumber string, status string) error
	GenerateAccountNumber(ctx context.Context, tx *sql.Tx) (string, error)
}
