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
	"github.com/krizad/go-gin-api/internal/core/domain"
	"github.com/krizad/go-gin-api/internal/core/ports"
	handlers "github.com/krizad/go-gin-api/internal/adapters/http/handlers"
)

// MockAccountService is a mock implementation of AccountService
type MockAccountService struct {
	CreateAccountFunc      func(ctx context.Context, account *domain.Account) (*domain.Account, error)
	GetAccountByNumberFunc func(ctx context.Context, accountNumber string) (*domain.Account, error)
	GetAccountListFunc     func(ctx context.Context, page int, limit int) ([]*domain.Account, int, error)
	CloseAccountFunc       func(ctx context.Context, accountNumber string) (*domain.Account, error)
}

func (m *MockAccountService) CreateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	if m.CreateAccountFunc != nil {
		return m.CreateAccountFunc(ctx, account)
	}
	return &domain.Account{ID: 1, AccountNumber: "1234567890"}, nil
}

func (m *MockAccountService) GetAccountByNumber(ctx context.Context, accountNumber string) (*domain.Account, error) {
	if m.GetAccountByNumberFunc != nil {
		return m.GetAccountByNumberFunc(ctx, accountNumber)
	}
	return &domain.Account{ID: 1, AccountNumber: accountNumber}, nil
}

func (m *MockAccountService) GetAccountList(ctx context.Context, page int, limit int) ([]*domain.Account, int, error) {
	if m.GetAccountListFunc != nil {
		return m.GetAccountListFunc(ctx, page, limit)
	}
	return []*domain.Account{{ID: 1, AccountNumber: "1234567890"}}, 1, nil
}

func (m *MockAccountService) CloseAccount(ctx context.Context, accountNumber string) (*domain.Account, error) {
	if m.CloseAccountFunc != nil {
		return m.CloseAccountFunc(ctx, accountNumber)
	}
	return &domain.Account{ID: 1, AccountNumber: accountNumber, Status: "CLOSED"}, nil
}

func setupAccountRouter(service ports.AccountService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := handlers.NewAccountHandler(service)
	
	router.POST("/accounts", handler.CreateAccount)
	router.GET("/accounts/:account_number", handler.GetAccountByNumber)
	router.GET("/accounts", handler.GetAccountList)
	router.PATCH("/accounts/:account_number/close", handler.CloseAccountHandler)
	
	return router
}

func TestAccountHandler_CreateAccount(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockService    *MockAccountService
		expectedStatus int
	}{
		{
			name: "successful account creation",
			requestBody: dto.CreateAccountRequest{
				OwnerName:     "John Doe",
				CitizenID:     "1234567890123",
				PhoneNumber:   "0812345678",
				AccountType:   "SAVING",
				InitialBalance: 1000.0,
			},
			mockService: &MockAccountService{
				CreateAccountFunc: func(ctx context.Context, account *domain.Account) (*domain.Account, error) {
					return &domain.Account{ID: 1, AccountNumber: "1234567890", Status: "ACTIVE"}, nil
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "citizen ID already exists",
			requestBody: dto.CreateAccountRequest{
				OwnerName:     "John Doe",
				CitizenID:     "1234567890123",
				PhoneNumber:   "0812345678",
				AccountType:   "SAVING",
				InitialBalance: 1000.0,
			},
			mockService: &MockAccountService{
				CreateAccountFunc: func(ctx context.Context, account *domain.Account) (*domain.Account, error) {
					return nil, constants.ErrCitizenIDExists
				},
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "invalid request body",
			requestBody: "invalid json",
			mockService: &MockAccountService{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupAccountRouter(tt.mockService)
			
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
			
			req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAccountHandler_GetAccountByNumber(t *testing.T) {
	tests := []struct {
		name           string
		accountNumber  string
		mockService    *MockAccountService
		expectedStatus int
	}{
		{
			name:          "successful get account",
			accountNumber: "1234567890",
			mockService: &MockAccountService{
				GetAccountByNumberFunc: func(ctx context.Context, accountNumber string) (*domain.Account, error) {
					return &domain.Account{ID: 1, AccountNumber: accountNumber, Status: "ACTIVE"}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "account not found",
			accountNumber: "1234567890",
			mockService: &MockAccountService{
				GetAccountByNumberFunc: func(ctx context.Context, accountNumber string) (*domain.Account, error) {
					return nil, constants.ErrAccountNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupAccountRouter(tt.mockService)
			
			req, _ := http.NewRequest("GET", "/accounts/"+tt.accountNumber, nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAccountHandler_GetAccountList(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockService    *MockAccountService
		expectedStatus int
	}{
		{
			name:        "successful get account list",
			queryParams: "?page=1&limit=10",
			mockService: &MockAccountService{
				GetAccountListFunc: func(ctx context.Context, page int, limit int) ([]*domain.Account, int, error) {
					return []*domain.Account{
						{ID: 1, AccountNumber: "1234567890"},
						{ID: 2, AccountNumber: "0987654321"},
					}, 2, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupAccountRouter(tt.mockService)
			
			req, _ := http.NewRequest("GET", "/accounts"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAccountHandler_CloseAccount(t *testing.T) {
	tests := []struct {
		name           string
		accountNumber  string
		mockService    *MockAccountService
		expectedStatus int
	}{
		{
			name:          "successful close account",
			accountNumber: "1234567890",
			mockService: &MockAccountService{
				CloseAccountFunc: func(ctx context.Context, accountNumber string) (*domain.Account, error) {
					return &domain.Account{ID: 1, AccountNumber: accountNumber, Status: "CLOSED"}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "account not found",
			accountNumber: "1234567890",
			mockService: &MockAccountService{
				CloseAccountFunc: func(ctx context.Context, accountNumber string) (*domain.Account, error) {
					return nil, constants.ErrAccountNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupAccountRouter(tt.mockService)
			
			req, _ := http.NewRequest("PATCH", "/accounts/"+tt.accountNumber+"/close", nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
