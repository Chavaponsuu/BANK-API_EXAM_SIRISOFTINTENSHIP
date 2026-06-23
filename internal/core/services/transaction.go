package services

import (
	"context"
	"fmt"

	"github.com/krizad/go-gin-api/constants"
	"github.com/krizad/go-gin-api/internal/core/domain"
	"github.com/krizad/go-gin-api/internal/core/ports"
)

type transactionService struct {
	accountRepo     ports.AccountRepository
	transactionRepo ports.TransactionRepository
}

func NewTransactionService(accountRepo ports.AccountRepository, transactionRepo ports.TransactionRepository) ports.TransactionService {
	return &transactionService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}
func (s *transactionService) Deposit(ctx context.Context, accountNumber string, amount float64, description string) (*domain.Transaction, error) {
	return s.processTransaction(ctx, accountNumber, amount, "DEPOSIT", description)
}

func (s *transactionService) Withdraw(ctx context.Context, accountNumber string, amount float64, description string) (*domain.Transaction, error) {
	return s.processTransaction(ctx, accountNumber, -amount, "WITHDRAW", description)
}

func (s *transactionService) processTransaction(ctx context.Context, accountNumber string, amount float64, txType string, description string) (*domain.Transaction, error) {
	tx, err := s.accountRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	account, err := s.accountRepo.GetByAccountNumberWithLock(ctx, tx, accountNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		return nil, constants.ErrAccountNotFound
	}
	if account.Status != "ACTIVE" {
		return nil, constants.ErrAccountNotActive
	}

	balanceBefore := account.Balance
	balanceAfter := balanceBefore + amount // amount เป็น negative สำหรับ WITHDRAW

	if balanceAfter < 0 {
		return nil, constants.ErrInsufficientBalance
	}

	if err = s.accountRepo.UpdateBalanceWithTx(ctx, tx, account.ID, balanceAfter); err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}
	if txType == "WITHDRAW" {
		amount = -amount
	}
	createdTx, err := s.transactionRepo.CreateWithTx(ctx, tx, &domain.Transaction{
		AccountID:       account.ID,
		TransactionType: txType,
		Amount:          amount,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
		Description:     description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return createdTx, nil
}

// GetTransactionHistory ดึงประวัติธุรกรรมของบัญชี
func (s *transactionService) GetTransactionHistory(ctx context.Context, accountNumber string, page int, limit int) ([]*domain.Transaction, int, error) {

	account, err := s.accountRepo.GetByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, 0, constants.ErrOperationFailed
	}
	if account == nil {
		return nil, 0, constants.ErrAccountNotFound
	}
	offset := (page - 1) * limit
	transactions, total, err := s.transactionRepo.GetTxByAccountID(ctx, account.ID, limit, offset)
	if err != nil {
		return nil, 0, constants.ErrOperationFailed
	}

	return transactions, total, nil
}

func (s *transactionService) GetAllTransaction(ctx context.Context, page int, limit int) ([]*domain.Transaction, int, error) {

	offset := (page - 1) * limit

	transactions, total, err := s.transactionRepo.GetAllTransaction(ctx, limit, offset)
	if err != nil {
		return nil, 0, constants.ErrOperationFailed
	}

	return transactions, total, nil
}
