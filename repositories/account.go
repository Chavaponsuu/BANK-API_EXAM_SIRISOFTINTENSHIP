package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/krizad/go-gin-api/models"
)

type AccountRepository interface {
	// Create(ctx context.Context, account *models.Account) (*models.Account, error)
	CreateWithTx(ctx context.Context, tx *sql.Tx, account *models.Account) (*models.Account, error)
	CheckCitizenIDExists(ctx context.Context, citizenID string) (bool, error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
	GetByAccountNumber(ctx context.Context, accountNumber string) (*models.Account, error)
	GetByAccountNumberWithLock(ctx context.Context, tx *sql.Tx, accountNumber string) (*models.Account, error)
	GetByAccountList(ctx context.Context, page int, limit int) ([]*models.Account, int, error)
	UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, accountID int64, newBalance float64) error
	UpdateStatus(ctx context.Context, tx *sql.Tx, accountNumber string, status string) error
	GenerateAccountNumber(ctx context.Context, tx *sql.Tx) (string, error)
}

type AccountRepo struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &AccountRepo{db: db}
}

type querier interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (r *AccountRepo) GenerateAccountNumber(ctx context.Context, tx *sql.Tx) (string, error) {
	var next int64

	err := tx.QueryRowContext(
		ctx,
		`SELECT nextval('account_number_seq')`,
	).Scan(&next)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%010d", next), nil
}

func (r *AccountRepo) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *AccountRepo) CheckCitizenIDExists(ctx context.Context, citizenID string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM accounts WHERE citizen_id = $1`, citizenID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// func (r *AccountRepo) Create(ctx context.Context, account *models.Account) (*models.Account, error) {
// 	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
// 	defer cancel()
// 	return r.createWithQuerier(ctx, r.db, account)
// }

func (r *AccountRepo) CreateWithTx(ctx context.Context, tx *sql.Tx, account *models.Account) (*models.Account, error) {
	fmt.Println("account number ====== ", account.AccountNumber)
	err := tx.QueryRowContext(ctx, `
		INSERT INTO accounts (owner_name, citizen_id, phone_number, account_type, balance, account_number, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`,
		account.OwnerName,
		account.CitizenID,
		account.PhoneNumber,
		account.AccountType,
		account.Balance,
		account.AccountNumber,
		account.Status,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}
	return account, nil
}

// func (r *AccountRepo) createWithQuerier(ctx context.Context, q querier, account *models.Account) (*models.Account, error) {
// 	var err error
// 	account.AccountNumber, err = generateAccountNumber(ctx, q)
// 	if err != nil {
// 		return nil, fmt.Errorf("generate account number: %w", err)
// 	}

// 	err = q.QueryRowContext(ctx, `
// 		INSERT INTO accounts (owner_name, citizen_id, phone_number, account_type, balance, account_number, status)
// 		VALUES ($1, $2, $3, $4, $5, $6, $7)
// 		RETURNING id, created_at, updated_at`,
// 		account.OwnerName,
// 		account.CitizenID,
// 		account.PhoneNumber,
// 		account.AccountType,
// 		account.Balance,
// 		account.AccountNumber,
// 		account.Status,
// 	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)
// 	if err != nil {
// 		return nil, fmt.Errorf("create account: %w", err)
// 	}
// 	return account, nil
// }

func (r *AccountRepo) GetByAccountNumber(ctx context.Context, accountNumber string) (*models.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return r.getByAccountNumber(ctx, r.db, accountNumber, false)
}

func (r *AccountRepo) GetByAccountNumberWithLock(ctx context.Context, tx *sql.Tx, accountNumber string) (*models.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return r.getByAccountNumber(ctx, tx, accountNumber, true)
}

func (r *AccountRepo) getByAccountNumber(ctx context.Context, q querier, accountNumber string, lock bool) (*models.Account, error) {
	query := `SELECT id, account_number, owner_name, citizen_id, phone_number, account_type, balance, status, created_at, updated_at
		FROM accounts
		WHERE account_number = $1`
	if lock {
		query += " FOR UPDATE"
	}

	account := &models.Account{}
	err := q.QueryRowContext(ctx, query, accountNumber).Scan(
		&account.ID,
		&account.AccountNumber,
		&account.OwnerName,
		&account.CitizenID,
		&account.PhoneNumber,
		&account.AccountType,
		&account.Balance,
		&account.Status,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get account by number: %w", err)
	}

	return account, nil
}
func (r *AccountRepo) GetByAccountList(ctx context.Context, limit int, offset int) ([]*models.Account, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var total int

	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM accounts
	`).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	accountList := []*models.Account{}
	query := `SELECT account_number,owner_name,account_type,balance,status FROM accounts ORDER BY account_number LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		account := &models.Account{}

		err := rows.Scan(&account.AccountNumber, &account.OwnerName, &account.AccountType, &account.Balance, &account.Status)
		if err != nil {
			return nil, 0, err
		}
		accountList = append(accountList, account)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return accountList, total, nil
}

func (r *AccountRepo) UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, accountID int64, newBalance float64) error {
	query := `UPDATE accounts SET balance = $1, updated_at = NOW() WHERE id = $2`
	_, err := tx.ExecContext(ctx, query, newBalance, accountID)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}
	return nil
}

// GetByAccountNumberWithLock ดึงข้อมูล account พร้อม lock row (FOR UPDATE)

func (r *AccountRepo) UpdateStatus(ctx context.Context, tx *sql.Tx, accountNumber string, status string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	// Update status
	query := `UPDATE accounts SET status = $1, updated_at = NOW() WHERE account_number = $2`
	result, err := tx.ExecContext(ctx, query, status, accountNumber)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
