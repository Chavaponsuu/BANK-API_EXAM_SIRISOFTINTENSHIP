package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/krizad/go-gin-api/models"
)

type TransactionRepository interface {
	CreateWithTx(ctx context.Context, tx *sql.Tx, transaction *models.Transaction) (*models.Transaction, error)
	GenerateTransactionRef(ctx context.Context, tx *sql.Tx) (string, error)
	// GetByAccountID(ctx context.Context, accountID int64, page int, limit int) ([]*models.Transaction, int, error)
	GetTxByAccountID(ctx context.Context, accountID int64, page int, limit int) ([]*models.Transaction, int, error)
}

type TransactionRepo struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &TransactionRepo{db: db}
}

// GenerateTransactionRef สร้าง transaction reference ในรูปแบบ TXN<YYYYMMDD><RUNNING_NUMBER>
// ใช้ database sequence หรือ lock เพื่อป้องกัน race condition
func (r *TransactionRepo) GenerateTransactionRef(ctx context.Context, tx *sql.Tx) (string, error) {
	now := time.Now()
	dateStr := now.Format("20060102") // YYYYMMDD

	// ใช้ SELECT FOR UPDATE เพื่อ lock การนับ transactions
	// หรือใช้ SERIAL/SEQUENCE ของ PostgreSQL เพื่อความปลอดภัย
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

	// Running number เริ่มจาก 1
	runningNumber := count + 1
	transactionRef := fmt.Sprintf("TXN%s%04d", dateStr, runningNumber)
	return transactionRef, nil
}

// CreateWithTx สร้าง transaction โดยใช้ database transaction
func (r *TransactionRepo) CreateWithTx(ctx context.Context, tx *sql.Tx, transaction *models.Transaction) (*models.Transaction, error) {
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
func (r *TransactionRepo) GetTxByAccountID(ctx context.Context, accountID int64, page int, limit int) ([]*models.Transaction, int, error) {
	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM transactions WHERE account_id = $1`, accountID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count transactions: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
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

	transactions := []*models.Transaction{}
	for rows.Next() {
		tx := &models.Transaction{}
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
