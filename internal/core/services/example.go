package services

import (
	"context"

	"github.com/krizad/go-gin-api/internal/adapters/http/dto"
	"github.com/krizad/go-gin-api/internal/adapters/repositories"
	"github.com/krizad/go-gin-api/internal/core/domain"
)

type ExampleService interface {
	CreateExample(ctx context.Context, req dto.CreateExampleRequest) (*domain.Example, error)
	GetExample(ctx context.Context, id int64) (*domain.Example, error)
	GetExampleByEmail(ctx context.Context, email string) (*domain.Example, error)
	ListExamples(ctx context.Context) ([]domain.Example, error)
	UpdateExample(ctx context.Context, id int64, req dto.UpdateExampleRequest) (*domain.Example, error)
	PatchExample(ctx context.Context, id int64, req dto.PatchExampleRequest) (*domain.Example, error)
	DeleteExample(ctx context.Context, id int64) (bool, error)
}

type exampleService struct {
	repo repositories.ExampleRepository
}

func NewExampleService(repo repositories.ExampleRepository) ExampleService {
	return &exampleService{repo: repo}
}

func (s *exampleService) CreateExample(ctx context.Context, req dto.CreateExampleRequest) (*domain.Example, error) {
	return s.repo.Create(ctx, req.Name, req.Email)
}

func (s *exampleService) GetExample(ctx context.Context, id int64) (*domain.Example, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *exampleService) GetExampleByEmail(ctx context.Context, email string) (*domain.Example, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *exampleService) ListExamples(ctx context.Context) ([]domain.Example, error) {
	return s.repo.List(ctx)
}

func (s *exampleService) UpdateExample(ctx context.Context, id int64, req dto.UpdateExampleRequest) (*domain.Example, error) {
	return s.repo.Update(ctx, id, req.Name, req.Email)
}

func (s *exampleService) PatchExample(ctx context.Context, id int64, req dto.PatchExampleRequest) (*domain.Example, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil // not found
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Email != nil {
		existing.Email = *req.Email
	}

	return s.repo.Update(ctx, id, existing.Name, existing.Email)
}

func (s *exampleService) DeleteExample(ctx context.Context, id int64) (bool, error) {
	return s.repo.Delete(ctx, id)
}
