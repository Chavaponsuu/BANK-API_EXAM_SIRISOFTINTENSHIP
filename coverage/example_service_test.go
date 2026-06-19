package coverage_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/krizad/go-gin-api/internal/adapters/http/dto"
	"github.com/krizad/go-gin-api/internal/core/domain"
	services "github.com/krizad/go-gin-api/internal/core/services"
)

// MockExampleRepository is a mock implementation of ExampleRepository
type MockExampleRepository struct {
	CreateFunc     func(ctx context.Context, name, email string) (*domain.Example, error)
	GetByIDFunc    func(ctx context.Context, id int64) (*domain.Example, error)
	GetByEmailFunc func(ctx context.Context, email string) (*domain.Example, error)
	ListFunc       func(ctx context.Context) ([]domain.Example, error)
	UpdateFunc     func(ctx context.Context, id int64, name, email string) (*domain.Example, error)
	DeleteFunc     func(ctx context.Context, id int64) (bool, error)
}

func (m *MockExampleRepository) Create(ctx context.Context, name, email string) (*domain.Example, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, name, email)
	}
	return &domain.Example{ID: 1, Name: name, Email: email}, nil
}

func (m *MockExampleRepository) GetByID(ctx context.Context, id int64) (*domain.Example, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return &domain.Example{ID: id, Name: "Test", Email: "test@example.com"}, nil
}

func (m *MockExampleRepository) GetByEmail(ctx context.Context, email string) (*domain.Example, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return &domain.Example{ID: 1, Name: "Test", Email: email}, nil
}

func (m *MockExampleRepository) List(ctx context.Context) ([]domain.Example, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return []domain.Example{
		{ID: 1, Name: "Test1", Email: "test1@example.com"},
		{ID: 2, Name: "Test2", Email: "test2@example.com"},
	}, nil
}

func (m *MockExampleRepository) Update(ctx context.Context, id int64, name, email string) (*domain.Example, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, name, email)
	}
	return &domain.Example{ID: id, Name: name, Email: email}, nil
}

func (m *MockExampleRepository) Delete(ctx context.Context, id int64) (bool, error) {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return true, nil
}

func (m *MockExampleRepository) Count(ctx context.Context) (int, error) {
	return 2, nil
}

func TestExampleService_CreateExample(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		req           dto.CreateExampleRequest
		mockRepo      *MockExampleRepository
		expectedError error
	}{
		{
			name: "successful creation",
			req: dto.CreateExampleRequest{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			mockRepo: &MockExampleRepository{
				CreateFunc: func(ctx context.Context, name, email string) (*domain.Example, error) {
					return &domain.Example{ID: 1, Name: name, Email: email}, nil
				},
			},
			expectedError: nil,
		},
		{
			name: "database error",
			req: dto.CreateExampleRequest{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			mockRepo: &MockExampleRepository{
				CreateFunc: func(ctx context.Context, name, email string) (*domain.Example, error) {
					return nil, errors.New("database error")
				},
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewExampleService(tt.mockRepo)
			result, err := service.CreateExample(ctx, tt.req)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result == nil {
					t.Errorf("expected example, got nil")
				}
			}
		})
	}
}

func TestExampleService_GetExample(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		id            int64
		mockRepo      *MockExampleRepository
		expectedError error
	}{
		{
			name: "successful get",
			id:   1,
			mockRepo: &MockExampleRepository{
				GetByIDFunc: func(ctx context.Context, id int64) (*domain.Example, error) {
					return &domain.Example{ID: id, Name: "Test", Email: "test@example.com"}, nil
				},
			},
			expectedError: nil,
		},
		{
			name: "not found",
			id:   999,
			mockRepo: &MockExampleRepository{
				GetByIDFunc: func(ctx context.Context, id int64) (*domain.Example, error) {
					return nil, nil
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewExampleService(tt.mockRepo)
			result, err := service.GetExample(ctx, tt.id)

			if tt.expectedError != nil {
				if err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if tt.id == 999 && result != nil {
					t.Errorf("expected nil for not found, got %v", result)
				}
			}
		})
	}
}

func TestExampleService_GetExampleByEmail(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		email         string
		mockRepo      *MockExampleRepository
		expectedError error
	}{
		{
			name:  "successful get",
			email: "test@example.com",
			mockRepo: &MockExampleRepository{
				GetByEmailFunc: func(ctx context.Context, email string) (*domain.Example, error) {
					return &domain.Example{ID: 1, Name: "Test", Email: email}, nil
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewExampleService(tt.mockRepo)
			result, err := service.GetExampleByEmail(ctx, tt.email)

			if tt.expectedError != nil {
				if err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result == nil {
					t.Errorf("expected example, got nil")
				}
			}
		})
	}
}

func TestExampleService_ListExamples(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		mockRepo      *MockExampleRepository
		expectedError error
		expectedCount int
	}{
		{
			name: "successful list",
			mockRepo: &MockExampleRepository{
				ListFunc: func(ctx context.Context) ([]domain.Example, error) {
					return []domain.Example{
						{ID: 1, Name: "Test1", Email: "test1@example.com"},
						{ID: 2, Name: "Test2", Email: "test2@example.com"},
					}, nil
				},
			},
			expectedError: nil,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewExampleService(tt.mockRepo)
			result, err := service.ListExamples(ctx)

			if tt.expectedError != nil {
				if err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if len(result) != tt.expectedCount {
					t.Errorf("expected %d items, got %d", tt.expectedCount, len(result))
				}
			}
		})
	}
}

func TestExampleService_UpdateExample(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		id            int64
		req           dto.UpdateExampleRequest
		mockRepo      *MockExampleRepository
		expectedError error
	}{
		{
			name: "successful update",
			id:   1,
			req: dto.UpdateExampleRequest{
				Name:  "Updated Name",
				Email: "updated@example.com",
			},
			mockRepo: &MockExampleRepository{
				UpdateFunc: func(ctx context.Context, id int64, name, email string) (*domain.Example, error) {
					return &domain.Example{ID: id, Name: name, Email: email}, nil
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewExampleService(tt.mockRepo)
			result, err := service.UpdateExample(ctx, tt.id, tt.req)

			if tt.expectedError != nil {
				if err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result == nil {
					t.Errorf("expected example, got nil")
				}
			}
		})
	}
}

func TestExampleService_PatchExample(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name          string
		id            int64
		req           dto.PatchExampleRequest
		mockRepo      *MockExampleRepository
		expectedError error
	}{
		{
			name: "successful patch name",
			id:   1,
			req: dto.PatchExampleRequest{
				Name:  stringPtr("Patched Name"),
				Email: nil,
			},
			mockRepo: &MockExampleRepository{
				GetByIDFunc: func(ctx context.Context, id int64) (*domain.Example, error) {
					return &domain.Example{ID: id, Name: "Original", Email: "original@example.com", CreatedAt: now, UpdatedAt: now}, nil
				},
				UpdateFunc: func(ctx context.Context, id int64, name, email string) (*domain.Example, error) {
					return &domain.Example{ID: id, Name: name, Email: email}, nil
				},
			},
			expectedError: nil,
		},
		{
			name: "successful patch email",
			id:   1,
			req: dto.PatchExampleRequest{
				Name:  nil,
				Email: stringPtr("patched@example.com"),
			},
			mockRepo: &MockExampleRepository{
				GetByIDFunc: func(ctx context.Context, id int64) (*domain.Example, error) {
					return &domain.Example{ID: id, Name: "Original", Email: "original@example.com", CreatedAt: now, UpdatedAt: now}, nil
				},
				UpdateFunc: func(ctx context.Context, id int64, name, email string) (*domain.Example, error) {
					return &domain.Example{ID: id, Name: name, Email: email}, nil
				},
			},
			expectedError: nil,
		},
		{
			name: "not found",
			id:   999,
			req: dto.PatchExampleRequest{
				Name:  stringPtr("Patched Name"),
				Email: nil,
			},
			mockRepo: &MockExampleRepository{
				GetByIDFunc: func(ctx context.Context, id int64) (*domain.Example, error) {
					return nil, nil
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewExampleService(tt.mockRepo)
			result, err := service.PatchExample(ctx, tt.id, tt.req)

			if tt.expectedError != nil {
				if err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if tt.id == 999 && result != nil {
					t.Errorf("expected nil for not found, got %v", result)
				}
			}
		})
	}
}

func TestExampleService_DeleteExample(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		id            int64
		mockRepo      *MockExampleRepository
		expectedError error
		expectedBool  bool
	}{
		{
			name: "successful delete",
			id:   1,
			mockRepo: &MockExampleRepository{
				DeleteFunc: func(ctx context.Context, id int64) (bool, error) {
					return true, nil
				},
			},
			expectedError: nil,
			expectedBool:  true,
		},
		{
			name: "not found",
			id:   999,
			mockRepo: &MockExampleRepository{
				DeleteFunc: func(ctx context.Context, id int64) (bool, error) {
					return false, nil
				},
			},
			expectedError: nil,
			expectedBool:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewExampleService(tt.mockRepo)
			result, err := service.DeleteExample(ctx, tt.id)

			if tt.expectedError != nil {
				if err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if result != tt.expectedBool {
					t.Errorf("expected %v, got %v", tt.expectedBool, result)
				}
			}
		})
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
