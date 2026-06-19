package services

import (
	"context"
	"fmt"

	"github.com/krizad/go-gin-api/constants"
	"github.com/krizad/go-gin-api/internal/core/domain"
	"github.com/krizad/go-gin-api/internal/core/ports"
)

type accountService struct {
	repo            ports.AccountRepository
	transactionRepo ports.TransactionRepository
}

func NewAccountService(repo ports.AccountRepository, transactionRepo ports.TransactionRepository) ports.AccountService {
	return &accountService{
		repo:            repo,
		transactionRepo: transactionRepo,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	exist, err := s.repo.CheckCitizenIDExists(ctx, account.CitizenID)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, constants.ErrCitizenIDExists
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	account.Status = "ACTIVE"

	account.AccountNumber, err = s.repo.GenerateAccountNumber(ctx, tx)
	if err != nil {
		return nil, err
	}
	createdAccount, err := s.repo.CreateWithTx(ctx, tx, account)
	// createdAccount.AccountNumber = s.repo.

	if err != nil {
		return nil, err
	}
	if createdAccount.Balance > 0 {
		_, err = s.transactionRepo.CreateWithTx(ctx, tx, &domain.Transaction{
			AccountID:       createdAccount.ID,
			TransactionType: "DEPOSIT",
			Amount:          createdAccount.Balance,
			BalanceBefore:   0,
			BalanceAfter:    createdAccount.Balance,
			Description:     "Initial deposit",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create deposit transaction: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return createdAccount, nil
}

func (s *accountService) GetAccountByNumber(ctx context.Context, accountNumber string) (*domain.Account, error) {
	account, err := s.repo.GetByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, constants.ErrInternal
	}

	if account == nil {
		return nil, constants.ErrAccountNotFound
	}

	return account, nil
}

func (s *accountService) GetAccountList(ctx context.Context, page int, limit int) ([]*domain.Account, int, error) {
	offset := (page - 1) * limit
	accountList, total, err := s.repo.GetByAccountList(ctx, limit, offset)
	if err != nil {
		return nil, 0, constants.ErrOperationFailed
	}

	return accountList, total, nil

}

func (s *accountService) CloseAccount(ctx context.Context, accountNumber string) (*domain.Account, error) {

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	account, err := s.repo.GetByAccountNumberWithLock(ctx, tx, accountNumber)
	if err != nil {
		return nil, constants.ErrInternal
	}

	if account == nil {
		return nil, constants.ErrAccountNotFound
	}

	if account.Status == "CLOSED" {
		return nil, constants.ErrAccountAlreadyClosed
	}

	err = s.repo.UpdateStatus(ctx, tx, accountNumber, "CLOSED")
	if err != nil {
		return nil, constants.ErrOperationFailed
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	account, err = s.repo.GetByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, constants.ErrInternal
	}

	return account, nil
}
