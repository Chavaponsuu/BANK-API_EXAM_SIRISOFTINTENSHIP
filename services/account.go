package services

import (
	"context"
	"fmt"

	"github.com/krizad/go-gin-api/constants"
	"github.com/krizad/go-gin-api/models"
	"github.com/krizad/go-gin-api/repositories"
)

type AccountService interface {
	CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error)
	GetAccountByNumber(ctx context.Context, accountNumber string) (*models.Account, error)
	GetAccountList(ctx context.Context, page int, limit int) ([]*models.Account, int, error)
	CloseAccount(ctx context.Context, accountNumber string) (*models.Account, error)
}

type accountService struct {
	repo            repositories.AccountRepository
	transactionRepo repositories.TransactionRepository
}

func NewAccountService(repo repositories.AccountRepository, transactionRepo repositories.TransactionRepository) AccountService {
	return &accountService{
		repo:            repo,
		transactionRepo: transactionRepo,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
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
			tx.Rollback()
		}
	}()

	account.Status = "ACTIVE"

	// if account.Balance > 0 {
	// 	return s.createAccountWithDeposit(ctx, account)
	// }
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
		_, err = s.transactionRepo.CreateWithTx(ctx, tx, &models.Transaction{
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

// // createAccountWithDeposit สร้างบัญชีและ transaction deposit พร้อมกันด้วย DB Transaction
// func (s *accountService) createAccountWithDeposit(ctx context.Context, account *models.Account) (*models.Account, error) {
// 	// เริ่ม DB Transaction
// 	tx, err := s.repo.BeginTx(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to begin transaction: %w", err)
// 	}

// 	// ใช้ defer เพื่อ rollback ถ้าเกิด error
// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	// 1. สร้างบัญชี (ใช้ CreateWithTx เพื่อใช้ transaction เดียวกัน)
// 	createdAccount, err := s.repo.CreateWithTx(ctx, tx, account)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create account: %w", err)
// 	}

// 	// 2. Generate transaction reference

// 	// 3. สร้าง deposit transaction
// 	transaction := &models.Transaction{
// 		AccountID:       createdAccount.ID,
// 		TransactionType: "DEPOSIT",
// 		Amount:          createdAccount.Balance,
// 		BalanceBefore:   0,
// 		BalanceAfter:    createdAccount.Balance,
// 		Description:     "Initial deposit",
// 	}

// 	_, err = s.transactionRepo.CreateWithTx(ctx, tx, transaction)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create deposit transaction: %w", err)
// 	}

// 	// 4. Commit transaction (ถ้าทุกอย่างสำเร็จ)
// 	if err = tx.Commit(); err != nil {
// 		return nil, fmt.Errorf("failed to commit transaction: %w", err)
// 	}

//		return createdAccount, nil
//	}
func (s *accountService) GetAccountByNumber(ctx context.Context, accountNumber string) (*models.Account, error) {
	account, err := s.repo.GetByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, constants.ErrInternal
	}

	if account == nil {
		return nil, constants.ErrAccountNotFound
	}

	return account, nil
}

func (s *accountService) GetAccountList(ctx context.Context, page int, limit int) ([]*models.Account, int, error) {
	offset := (page - 1) * limit
	accountList, total, err := s.repo.GetByAccountList(ctx, limit, offset)
	if err != nil {
		return nil, 0, constants.ErrOperationFailed
	}

	return accountList, total, nil

}

func (s *accountService) CloseAccount(ctx context.Context, accountNumber string) (*models.Account, error) {

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
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
