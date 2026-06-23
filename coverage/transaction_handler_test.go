package coverage_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/krizad/go-gin-api/constants"
	"github.com/krizad/go-gin-api/internal/adapters/http/dto"
	handlers "github.com/krizad/go-gin-api/internal/adapters/http/handlers"
	"github.com/krizad/go-gin-api/internal/core/domain"
	"github.com/krizad/go-gin-api/internal/core/ports"
)

// MockTransactionService is a mock implementation of TransactionService
type MockTransactionService struct {
	DepositFunc               func(ctx context.Context, accountNumber string, amount float64, description string) (*domain.Transaction, error)
	WithdrawFunc              func(ctx context.Context, accountNumber string, amount float64, description string) (*domain.Transaction, error)
	GetTransactionHistoryFunc func(ctx context.Context, accountNumber string, page int, limit int) ([]*domain.Transaction, int, error)
	GetAllTransactionFunc     func(ctx context.Context, page int, limit int) ([]*domain.Transaction, int, error)
}

func (m *MockTransactionService) Deposit(ctx context.Context, accountNumber string, amount float64, description string) (*domain.Transaction, error) {
	if m.DepositFunc != nil {
		return m.DepositFunc(ctx, accountNumber, amount, description)
	}
	return &domain.Transaction{ID: 1, AccountID: 1, TransactionType: "DEPOSIT", Amount: amount}, nil
}

func (m *MockTransactionService) Withdraw(ctx context.Context, accountNumber string, amount float64, description string) (*domain.Transaction, error) {
	if m.WithdrawFunc != nil {
		return m.WithdrawFunc(ctx, accountNumber, amount, description)
	}
	return &domain.Transaction{ID: 1, AccountID: 1, TransactionType: "WITHDRAW", Amount: amount}, nil
}

func (m *MockTransactionService) GetTransactionHistory(ctx context.Context, accountNumber string, page int, limit int) ([]*domain.Transaction, int, error) {
	if m.GetTransactionHistoryFunc != nil {
		return m.GetTransactionHistoryFunc(ctx, accountNumber, page, limit)
	}
	return []*domain.Transaction{{ID: 1, TransactionType: "DEPOSIT", Amount: 1000.0}}, 1, nil
}

func (m *MockTransactionService) GetAllTransaction(ctx context.Context, page int, limit int) ([]*domain.Transaction, int, error) {
	if m.GetAllTransactionFunc != nil {
		return m.GetAllTransactionFunc(ctx, page, limit)
	}
	return []*domain.Transaction{{ID: 1, TransactionType: "DEPOSIT", Amount: 1000.0}, {ID: 2, TransactionType: "WITHDRAW", Amount: 500.0}}, 2, nil
}

func setupTransactionRouter(service ports.TransactionService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := handlers.NewTransactionHandler(service)

	router.POST("/accounts/:account_number/deposit", handler.Deposit)
	router.POST("/accounts/:account_number/withdraw", handler.Withdraw)
	router.GET("/accounts/:account_number/transactions", handler.GetTransactionHistory)

	return router
}

func TestTransactionHandler_Deposit(t *testing.T) {
	tests := []struct {
		name           string
		accountNumber  string
		requestBody    interface{}
		mockService    *MockTransactionService
		expectedStatus int
	}{
		{
			name:          "successful deposit",
			accountNumber: "1234567890",
			requestBody: dto.TransactionRequest{
				Amount:      1000.0,
				Description: "Test deposit",
			},
			mockService: &MockTransactionService{
				DepositFunc: func(ctx context.Context, accountNumber string, amount float64, description string) (*domain.Transaction, error) {
					return &domain.Transaction{ID: 1, AccountID: 1, TransactionType: "DEPOSIT", Amount: amount}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "account not found",
			accountNumber: "1234567890",
			requestBody: dto.TransactionRequest{
				Amount:      1000.0,
				Description: "Test deposit",
			},
			mockService: &MockTransactionService{
				DepositFunc: func(ctx context.Context, accountNumber string, amount float64, description string) (*domain.Transaction, error) {
					return nil, constants.ErrAccountNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid request body",
			accountNumber:  "1234567890",
			requestBody:    "invalid json",
			mockService:    &MockTransactionService{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTransactionRouter(tt.mockService)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
			}

			req, _ := http.NewRequest("POST", "/accounts/"+tt.accountNumber+"/deposit", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestTransactionHandler_Withdraw(t *testing.T) {
	tests := []struct {
		name           string
		accountNumber  string
		requestBody    interface{}
		mockService    *MockTransactionService
		expectedStatus int
	}{
		{
			name:          "successful withdrawal",
			accountNumber: "1234567890",
			requestBody: dto.TransactionRequest{
				Amount:      500.0,
				Description: "Test withdrawal",
			},
			mockService: &MockTransactionService{
				WithdrawFunc: func(ctx context.Context, accountNumber string, amount float64, description string) (*domain.Transaction, error) {
					return &domain.Transaction{ID: 1, AccountID: 1, TransactionType: "WITHDRAW", Amount: amount}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "insufficient balance",
			accountNumber: "1234567890",
			requestBody: dto.TransactionRequest{
				Amount:      5000.0,
				Description: "Test withdrawal",
			},
			mockService: &MockTransactionService{
				WithdrawFunc: func(ctx context.Context, accountNumber string, amount float64, description string) (*domain.Transaction, error) {
					return nil, constants.ErrInsufficientBalance
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTransactionRouter(tt.mockService)

			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			req, _ := http.NewRequest("POST", "/accounts/"+tt.accountNumber+"/withdraw", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestTransactionHandler_GetTransactionHistory(t *testing.T) {
	tests := []struct {
		name           string
		accountNumber  string
		queryParams    string
		mockService    *MockTransactionService
		expectedStatus int
	}{
		{
			name:          "successful get transaction history",
			accountNumber: "1234567890",
			queryParams:   "?page=1&limit=10",
			mockService: &MockTransactionService{
				GetTransactionHistoryFunc: func(ctx context.Context, accountNumber string, page int, limit int) ([]*domain.Transaction, int, error) {
					return []*domain.Transaction{
						{ID: 1, TransactionType: "DEPOSIT", Amount: 1000.0},
						{ID: 2, TransactionType: "WITHDRAW", Amount: 500.0},
					}, 2, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "account not found",
			accountNumber: "1234567890",
			queryParams:   "?page=1&limit=10",
			mockService: &MockTransactionService{
				GetTransactionHistoryFunc: func(ctx context.Context, accountNumber string, page int, limit int) ([]*domain.Transaction, int, error) {
					return nil, 0, constants.ErrAccountNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTransactionRouter(tt.mockService)

			req, _ := http.NewRequest("GET", "/accounts/"+tt.accountNumber+"/transactions"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
