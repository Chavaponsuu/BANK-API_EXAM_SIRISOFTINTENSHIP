package dto

type CreateExampleRequest struct {
	Name  string `json:"name" binding:"required" example:"John Doe"`
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
}

type UpdateExampleRequest struct {
	Name  string `json:"name" binding:"required" example:"John Doe"`
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
}

type PatchExampleRequest struct {
	Name  *string `json:"name" example:"John Doe"`
	Email *string `json:"email" example:"john@example.com"`
}
