package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/krizad/go-gin-api/internal/core/domain"
	"github.com/krizad/go-gin-api/internal/core/ports"
)

type TransactionRepo struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) ports.TransactionRepository {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) GenerateTransactionRef(ctx context.Context, tx *sql.Tx) (string, error) {
	now := time.Now()
	dateStr := now.Format("20060102")

	var count int
	query := `SELECT COUNT(*) FROM transactions WHERE transaction_ref LIKE $1`
	pattern := fmt.Sprintf("TXN%s%%", dateStr)

	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, query, pattern).Scan(&count)
	} else {
		err = r.db.QueryRowContext(ctx, query, pattern).Scan(&count)
	}

	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to count transactions: %w", err)
	}

	runningNumber := count + 1
	transactionRef := fmt.Sprintf("TXN%s%04d", dateStr, runningNumber)
	return transactionRef, nil
}

func (r *TransactionRepo) CreateWithTx(ctx context.Context, tx *sql.Tx, transaction *domain.Transaction) (*domain.Transaction, error) {
	query := `INSERT INTO transactions (account_id, transaction_ref, transaction_type, amount, balance_before, balance_after, description)
	 VALUES ($1, $2, $3, $4, $5, $6, $7)
	 RETURNING id, transaction_ref , created_at`
	TransactionRef, err := r.GenerateTransactionRef(ctx, tx)
	if err != nil {
		return nil, err
	}
	err = tx.QueryRowContext(ctx, query,
		transaction.AccountID,
		TransactionRef,
		transaction.TransactionType,
		transaction.Amount,
		transaction.BalanceBefore,
		transaction.BalanceAfter,
		transaction.Description,
	).Scan(&transaction.ID, &transaction.TransactionRef, &transaction.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}

	return transaction, nil
}

// GetByAccountID ดึงประวัติธุรกรรมของบัญชี เรียงจากใหม่ไปเก่า
func (r *TransactionRepo) GetTxByAccountID(ctx context.Context, accountID int64, limit int, offset int) ([]*domain.Transaction, int, error) {
	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM transactions WHERE account_id = $1`, accountID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count transactions: %w", err)
	}

	// Get transactions - เรียงจากใหม่ไปเก่า (ORDER BY created_at DESC)
	query := `SELECT id, account_id, transaction_ref, transaction_type, amount, balance_before, balance_after, description, created_at
		FROM transactions
		WHERE account_id = $1
		ORDER BY id DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, accountID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query transactions: %w", err)
	}
	defer rows.Close()

	transactions := []*domain.Transaction{}
	for rows.Next() {
		tx := &domain.Transaction{}
		err := rows.Scan(
			&tx.ID,
			&tx.AccountID,
			&tx.TransactionRef,
			&tx.TransactionType,
			&tx.Amount,
			&tx.BalanceBefore,
			&tx.BalanceAfter,
			&tx.Description,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return transactions, total, nil
}
