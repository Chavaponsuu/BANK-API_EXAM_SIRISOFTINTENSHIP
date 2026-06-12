package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/krizad/go-gin-api/models"
	"github.com/krizad/go-gin-api/repositories"
)

var (
	ErrInvalidAccountNumber = errors.New("invalid account_number format")
	ErrAccountNotFound      = errors.New("account not found")
	ErrInternal             = errors.New("internal error")
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
	// 1. business validation
	if account.Balance < 0 {
		return nil, errors.New("balance cannot be negative")
	}

	// 2. validate citizen_id format (must be 13 digits)
	citizenIDRegex := regexp.MustCompile(`^[0-9]{13}$`)
	if !citizenIDRegex.MatchString(account.CitizenID) {
		return nil, errors.New("invalid input: citizen_id must be 13 digits")
	}

	// 3. validate account_type (must be SAVING or CURRENT)
	if account.AccountType != "SAVING" && account.AccountType != "CURRENT" {
		return nil, errors.New("invalid input: account_type must be SAVING or CURRENT")
	}

	// ตรวจสอบว่า citizen_id ซ้ำหรือไม่
	exist, err := s.repo.CheckCitizenIDExists(ctx, account.CitizenID)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, errors.New("citizen_id already exists")
	}

	account.Status = "ACTIVE"

	// ถ้า initial_balance > 0 ต้องใช้ DB Transaction
	if account.Balance > 0 {
		return s.createAccountWithDeposit(ctx, account)
	}

	// ถ้า balance = 0 สร้างบัญชีอย่างเดียว
	return s.repo.Create(ctx, account)
}

// createAccountWithDeposit สร้างบัญชีและ transaction deposit พร้อมกันด้วย DB Transaction
func (s *accountService) createAccountWithDeposit(ctx context.Context, account *models.Account) (*models.Account, error) {
	// เริ่ม DB Transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// ใช้ defer เพื่อ rollback ถ้าเกิด error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. สร้างบัญชี (ใช้ CreateWithTx เพื่อใช้ transaction เดียวกัน)
	createdAccount, err := s.repo.CreateWithTx(ctx, tx, account)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	// 2. Generate transaction reference

	// 3. สร้าง deposit transaction
	transaction := &models.Transaction{
		AccountID:       createdAccount.ID,
		TransactionRef:  "",
		TransactionType: "DEPOSIT",
		Amount:          createdAccount.Balance,
		BalanceBefore:   0,
		BalanceAfter:    createdAccount.Balance,
		Description:     "Initial deposit",
	}

	_, err = s.transactionRepo.CreateWithTx(ctx, tx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create deposit transaction: %w", err)
	}

	// 4. Commit transaction (ถ้าทุกอย่างสำเร็จ)
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return createdAccount, nil
}

var accountNumberRegex = regexp.MustCompile(`^[0-9]{10}$`)

func isValidAccountNumber(acc string) bool {
	return accountNumberRegex.MatchString(acc)
}

func (s *accountService) GetAccountByNumber(ctx context.Context, accountNumber string) (*models.Account, error) {
	if accountNumber == "" || !isValidAccountNumber(accountNumber) {
		return nil, ErrInvalidAccountNumber
	}

	account, err := s.repo.GetByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, ErrInternal
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

func (s *accountService) GetAccountList(ctx context.Context, page int, limit int) ([]*models.Account, int, error) {

	if page < 1 || limit < 1 || limit > 100 {
		return nil, 0, errors.New("invalid page or limit parameters")
	}
	accountList, total, err := s.repo.GetByAccountList(ctx, page, limit)
	if err != nil {
		return nil, 0, errors.New("failed to get account list")
	}

	return accountList, total, nil

}

var (
	ErrAccountAlreadyClosed = errors.New("account is already closed")
)

func (s *accountService) CloseAccount(ctx context.Context, accountNumber string) (*models.Account, error) {
	if accountNumber == "" || !isValidAccountNumber(accountNumber) {
		return nil, ErrInvalidAccountNumber
	}

	account, err := s.repo.GetByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, ErrInternal
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Status == "CLOSED" {
		return nil, ErrAccountAlreadyClosed
	}

	err = s.repo.UpdateStatus(ctx, accountNumber, "CLOSED")
	if err != nil {
		return nil, errors.New("failed to close account")
	}
	account.Status = "CLOSED"

	return account, nil
}
