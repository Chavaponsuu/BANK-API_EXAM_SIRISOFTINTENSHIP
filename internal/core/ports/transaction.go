package ports

import (
	"context"
	"database/sql"

	"github.com/krizad/go-gin-api/internal/core/domain"
)

type TransactionService interface {
	Deposit(ctx context.Context, accountNumber string, amount float64, description string) (*domain.Transaction, error)
	Withdraw(ctx context.Context, accountNumber string, amount float64, description string) (*domain.Transaction, error)
	GetTransactionHistory(ctx context.Context, accountNumber string, page int, limit int) ([]*domain.Transaction, int, error)
}
type TransactionRepository interface {
	CreateWithTx(ctx context.Context, tx *sql.Tx, transaction *domain.Transaction) (*domain.Transaction, error)
	GenerateTransactionRef(ctx context.Context, tx *sql.Tx) (string, error)
	// GetByAccountID(ctx context.Context, accountID int64, page int, limit int) ([]*domain.Transaction, int, error)
	GetTxByAccountID(ctx context.Context, accountID int64, page int, limit int) ([]*domain.Transaction, int, error)
}
