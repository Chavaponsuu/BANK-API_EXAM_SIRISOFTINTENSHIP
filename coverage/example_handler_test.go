package coverage_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/krizad/go-gin-api/internal/adapters/http/dto"
	"github.com/krizad/go-gin-api/internal/core/domain"
	"github.com/krizad/go-gin-api/internal/core/services"
	handlers "github.com/krizad/go-gin-api/internal/adapters/http/handlers"
)

// MockExampleService is a mock implementation of ExampleService
type MockExampleService struct {
	CreateExampleFunc      func(ctx context.Context, req dto.CreateExampleRequest) (*domain.Example, error)
	GetExampleFunc         func(ctx context.Context, id int64) (*domain.Example, error)
	GetExampleByEmailFunc func(ctx context.Context, email string) (*domain.Example, error)
	ListExamplesFunc      func(ctx context.Context) ([]domain.Example, error)
	UpdateExampleFunc     func(ctx context.Context, id int64, req dto.UpdateExampleRequest) (*domain.Example, error)
	PatchExampleFunc      func(ctx context.Context, id int64, req dto.PatchExampleRequest) (*domain.Example, error)
	DeleteExampleFunc     func(ctx context.Context, id int64) (bool, error)
}

func (m *MockExampleService) CreateExample(ctx context.Context, req dto.CreateExampleRequest) (*domain.Example, error) {
	if m.CreateExampleFunc != nil {
		return m.CreateExampleFunc(ctx, req)
	}
	return &domain.Example{ID: 1, Name: req.Name, Email: req.Email}, nil
}

func (m *MockExampleService) GetExample(ctx context.Context, id int64) (*domain.Example, error) {
	if m.GetExampleFunc != nil {
		return m.GetExampleFunc(ctx, id)
	}
	return &domain.Example{ID: id, Name: "Test", Email: "test@example.com"}, nil
}

func (m *MockExampleService) GetExampleByEmail(ctx context.Context, email string) (*domain.Example, error) {
	if m.GetExampleByEmailFunc != nil {
		return m.GetExampleByEmailFunc(ctx, email)
	}
	return &domain.Example{ID: 1, Name: "Test", Email: email}, nil
}

func (m *MockExampleService) ListExamples(ctx context.Context) ([]domain.Example, error) {
	if m.ListExamplesFunc != nil {
		return m.ListExamplesFunc(ctx)
	}
	return []domain.Example{
		{ID: 1, Name: "Test1", Email: "test1@example.com"},
		{ID: 2, Name: "Test2", Email: "test2@example.com"},
	}, nil
}

func (m *MockExampleService) UpdateExample(ctx context.Context, id int64, req dto.UpdateExampleRequest) (*domain.Example, error) {
	if m.UpdateExampleFunc != nil {
		return m.UpdateExampleFunc(ctx, id, req)
	}
	return &domain.Example{ID: id, Name: req.Name, Email: req.Email}, nil
}

func (m *MockExampleService) PatchExample(ctx context.Context, id int64, req dto.PatchExampleRequest) (*domain.Example, error) {
	if m.PatchExampleFunc != nil {
		return m.PatchExampleFunc(ctx, id, req)
	}
	return &domain.Example{ID: id, Name: "Patched", Email: "patched@example.com"}, nil
}

func (m *MockExampleService) DeleteExample(ctx context.Context, id int64) (bool, error) {
	if m.DeleteExampleFunc != nil {
		return m.DeleteExampleFunc(ctx, id)
	}
	return true, nil
}

func setupExampleRouter(service services.ExampleService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := handlers.NewExampleHandler(service)
	
	router.GET("/examples", handler.ListExamples)
	router.GET("/examples/:id", handler.GetExample)
	router.POST("/examples", handler.CreateExample)
	router.PUT("/examples/:id", handler.UpdateExample)
	router.PATCH("/examples/:id", handler.PatchExample)
	router.DELETE("/examples/:id", handler.DeleteExample)
	router.GET("/examples/search", handler.GetExampleByEmail)
	
	return router
}

func TestExampleHandler_ListExamples(t *testing.T) {
	tests := []struct {
		name           string
		mockService    *MockExampleService
		expectedStatus int
	}{
		{
			name: "successful list",
			mockService: &MockExampleService{
				ListExamplesFunc: func(ctx context.Context) ([]domain.Example, error) {
					return []domain.Example{
						{ID: 1, Name: "Test1", Email: "test1@example.com"},
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupExampleRouter(tt.mockService)
			
			req, _ := http.NewRequest("GET", "/examples", nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestExampleHandler_GetExample(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockService    *MockExampleService
		expectedStatus int
	}{
		{
			name: "successful get",
			id:   "1",
			mockService: &MockExampleService{
				GetExampleFunc: func(ctx context.Context, id int64) (*domain.Example, error) {
					return &domain.Example{ID: id, Name: "Test", Email: "test@example.com"}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "not found",
			id:   "999",
			mockService: &MockExampleService{
				GetExampleFunc: func(ctx context.Context, id int64) (*domain.Example, error) {
					return nil, nil
				},
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:        "invalid id",
			id:         "invalid",
			mockService: &MockExampleService{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupExampleRouter(tt.mockService)
			
			req, _ := http.NewRequest("GET", "/examples/"+tt.id, nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestExampleHandler_CreateExample(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockService    *MockExampleService
		expectedStatus int
	}{
		{
			name: "successful creation",
			requestBody: dto.CreateExampleRequest{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			mockService: &MockExampleService{
				CreateExampleFunc: func(ctx context.Context, req dto.CreateExampleRequest) (*domain.Example, error) {
					return &domain.Example{ID: 1, Name: req.Name, Email: req.Email}, nil
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:        "invalid request body",
			requestBody: "invalid json",
			mockService: &MockExampleService{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupExampleRouter(tt.mockService)
			
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
			
			req, _ := http.NewRequest("POST", "/examples", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestExampleHandler_UpdateExample(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		requestBody    interface{}
		mockService    *MockExampleService
		expectedStatus int
	}{
		{
			name: "successful update",
			id:   "1",
			requestBody: dto.UpdateExampleRequest{
				Name:  "Updated Name",
				Email: "updated@example.com",
			},
			mockService: &MockExampleService{
				UpdateExampleFunc: func(ctx context.Context, id int64, req dto.UpdateExampleRequest) (*domain.Example, error) {
					return &domain.Example{ID: id, Name: req.Name, Email: req.Email}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupExampleRouter(tt.mockService)
			
			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}
			
			req, _ := http.NewRequest("PUT", "/examples/"+tt.id, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestExampleHandler_PatchExample(t *testing.T) {
	name := "test"
	
	tests := []struct {
		name           string
		id             string
		requestBody    interface{}
		mockService    *MockExampleService
		expectedStatus int
	}{
		{
			name: "successful patch name",
			id:   "1",
			requestBody: dto.PatchExampleRequest{
				Name:  &name,
				Email: nil,
			},
			mockService: &MockExampleService{
				PatchExampleFunc: func(ctx context.Context, id int64, req dto.PatchExampleRequest) (*domain.Example, error) {
					return &domain.Example{ID: id, Name: "Patched", Email: "patched@example.com"}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "no fields provided",
			id:   "1",
			requestBody: dto.PatchExampleRequest{
				Name:  nil,
				Email: nil,
			},
			mockService:    &MockExampleService{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupExampleRouter(tt.mockService)
			
			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}
			
			req, _ := http.NewRequest("PATCH", "/examples/"+tt.id, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestExampleHandler_DeleteExample(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockService    *MockExampleService
		expectedStatus int
	}{
		{
			name: "successful delete",
			id:   "1",
			mockService: &MockExampleService{
				DeleteExampleFunc: func(ctx context.Context, id int64) (bool, error) {
					return true, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "not found",
			id:   "999",
			mockService: &MockExampleService{
				DeleteExampleFunc: func(ctx context.Context, id int64) (bool, error) {
					return false, nil
				},
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupExampleRouter(tt.mockService)
			
			req, _ := http.NewRequest("DELETE", "/examples/"+tt.id, nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestExampleHandler_GetExampleByEmail(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockService    *MockExampleService
		expectedStatus int
	}{
		{
			name:        "successful get by email",
			queryParams: "?email=test@example.com",
			mockService: &MockExampleService{
				GetExampleByEmailFunc: func(ctx context.Context, email string) (*domain.Example, error) {
					return &domain.Example{ID: 1, Name: "Test", Email: email}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "missing email parameter",
			queryParams: "",
			mockService: &MockExampleService{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupExampleRouter(tt.mockService)
			
			req, _ := http.NewRequest("GET", "/examples/search"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			
			router.ServeHTTP(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
