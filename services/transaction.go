package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/krizad/go-gin-api/models"
	"github.com/krizad/go-gin-api/repositories"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrAccountNotActive    = errors.New("account is not active")
	ErrInvalidAmount       = errors.New("amount must be greater than 0")
)

type TransactionService interface {
	Deposit(ctx context.Context, accountNumber string, amount float64, description string) (*models.Transaction, error)
	Withdraw(ctx context.Context, accountNumber string, amount float64, description string) (*models.Transaction, error)
	GetTransactionHistory(ctx context.Context, accountNumber string, page int, limit int) ([]*models.Transaction, int, error)
}

type transactionService struct {
	accountRepo     repositories.AccountRepository
	transactionRepo repositories.TransactionRepository
}

func NewTransactionService(accountRepo repositories.AccountRepository, transactionRepo repositories.TransactionRepository) TransactionService {
	return &transactionService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

// Deposit ฝากเงินเข้าบัญชี
func (s *transactionService) Deposit(ctx context.Context, accountNumber string, amount float64, description string) (*models.Transaction, error) {
	// 1. Validate amount
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	// 2. Begin DB Transaction ก่อน (เปลี่ยนลำดับ)
	tx, err := s.accountRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 3. Get account พร้อม LOCK row (FOR UPDATE)
	account, err := s.accountRepo.GetByAccountNumberWithLock(ctx, tx, accountNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	// 4. Check account status
	if account.Status != "ACTIVE" {
		return nil, ErrAccountNotActive
	}

	// 5. Calculate balances (ใช้ balance ที่ lock แล้ว)
	balanceBefore := account.Balance
	balanceAfter := balanceBefore + amount

	// 6. Update account balance
	err = s.accountRepo.UpdateBalanceWithTx(ctx, tx, account.ID, balanceAfter)
	if err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	// 7. Generate transaction reference
	txRef, err := s.transactionRepo.GenerateTransactionRef(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate transaction ref: %w", err)
	}

	// 8. Create transaction record
	transaction := &models.Transaction{
		AccountID:       account.ID,
		TransactionRef:  txRef,
		TransactionType: "DEPOSIT",
		Amount:          amount,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
		Description:     description,
	}

	createdTx, err := s.transactionRepo.CreateWithTx(ctx, tx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// 9. Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return createdTx, nil
}

// Withdraw ถอนเงินจากบัญชี
func (s *transactionService) Withdraw(ctx context.Context, accountNumber string, amount float64, description string) (*models.Transaction, error) {
	// 1. Validate amount
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	// 2. Begin DB Transaction ก่อน
	tx, err := s.accountRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 3. Get account พร้อม LOCK row (FOR UPDATE)
	account, err := s.accountRepo.GetByAccountNumberWithLock(ctx, tx, accountNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	// 4. Check account status
	if account.Status != "ACTIVE" {
		return nil, ErrAccountNotActive
	}

	// 5. Check sufficient balance (ตรวจสอบจาก balance ที่ lock แล้ว)
	if account.Balance < amount {
		return nil, ErrInsufficientBalance
	}

	// 6. Calculate balances
	balanceBefore := account.Balance
	balanceAfter := balanceBefore - amount

	// 7. Update account balance
	err = s.accountRepo.UpdateBalanceWithTx(ctx, tx, account.ID, balanceAfter)
	if err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	// 8. Generate transaction reference
	txRef, err := s.transactionRepo.GenerateTransactionRef(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate transaction ref: %w", err)
	}

	// 9. Create transaction record
	transaction := &models.Transaction{
		AccountID:       account.ID,
		TransactionRef:  txRef,
		TransactionType: "WITHDRAW",
		Amount:          amount,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
		Description:     description,
	}

	createdTx, err := s.transactionRepo.CreateWithTx(ctx, tx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// 10. Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return createdTx, nil
}

// GetTransactionHistory ดึงประวัติธุรกรรมของบัญชี
func (s *transactionService) GetTransactionHistory(ctx context.Context, accountNumber string, page int, limit int) ([]*models.Transaction, int, error) {

	if page < 1 || limit < 1 || limit > 100 {
		return nil, 0, errors.New("invalid page or limit parameters")
	}

	// 2. Get account
	account, err := s.accountRepo.GetByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, 0, errors.New("failed to get account details")
	}
	if account == nil {
		return nil, 0, ErrAccountNotFound
	}

	// 3. Get transaction history
	transactions, total, err := s.transactionRepo.GetTxByAccountID(ctx, account.ID, page, limit)
	if err != nil {
		return nil, 0, errors.New("failed to get transaction history")
	}

	return transactions, total, nil
}
