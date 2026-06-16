package dto

import (
	"time"

	"github.com/krizad/go-gin-api/internal/core/domain"
)

type ExampleResponse struct {
	ID        int64     `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"john@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type ExampleListResponse []ExampleResponse

type ExampleDeleteResponse struct {
	ID int64 `json:"id"`
}

func ToExampleResponse(m *domain.Example) *ExampleResponse {
	return &ExampleResponse{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToExampleListResponse(domain []domain.Example) ExampleListResponse {
	list := make(ExampleListResponse, 0, len(domain))
	for _, m := range domain {
		list = append(list, *ToExampleResponse(&m))
	}
	return list
}
