package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/krizad/go-gin-api/db"
	"github.com/krizad/go-gin-api/internal/adapters/http/dto"
	"github.com/krizad/go-gin-api/internal/core/services"
)

type ExampleHandler struct {
	service services.ExampleService
}

func NewExampleHandler(service services.ExampleService) *ExampleHandler {
	return &ExampleHandler{service: service}
}

// ListExamples godoc
//
//	@summary		List all examples
//	@description	Retrieve a list of all example records.
//	@tags			examples
//	@produce		json
//	@success		200	{object}	dto.BaseResponse{data=dto.ExampleListResponse}
//	@failure		500	{object}	dto.BaseResponse
//	@router			/examples [get]
func (h *ExampleHandler) ListExamples(c *gin.Context) {
	examples, err := h.service.ListExamples(c.Request.Context())
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	dto.OK(c, dto.ToExampleListResponse(examples))
}

// GetExample godoc
//
//	@summary		Get an example by ID
//	@description	Retrieve a single example by its ID.
//	@tags			examples
//	@produce		json
//	@param			id	path		int	true	"Example ID"
//	@success		200	{object}	dto.BaseResponse{data=dto.ExampleResponse}
//	@failure		400	{object}	dto.BaseResponse
//	@failure		404	{object}	dto.BaseResponse
//	@failure		500	{object}	dto.BaseResponse
//	@router			/examples/{id} [get]
func (h *ExampleHandler) GetExample(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	example, err := h.service.GetExample(c.Request.Context(), id)
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if example == nil {
		dto.Error(c, http.StatusNotFound, "example not found")
		return
	}
	dto.OK(c, dto.ToExampleResponse(example))
}

// CreateExample godoc
//
//	@summary		Create a new example
//	@description	Create a new example record.
//	@tags			examples
//	@accept			json
//	@produce		json
//	@param			request	body		dto.CreateExampleRequest	true	"Example payload"
//	@success		201		{object}	dto.BaseResponse{data=dto.ExampleResponse}
//	@failure		400		{object}	dto.BaseResponse
//	@failure		500		{object}	dto.BaseResponse
//	@router			/examples [post]
func (h *ExampleHandler) CreateExample(c *gin.Context) {
	var req dto.CreateExampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	example, err := h.service.CreateExample(c.Request.Context(), req)
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	dto.Created(c, dto.ToExampleResponse(example))
}

// UpdateExample godoc
//
//	@summary		Update an example
//	@description	Fully replace an existing example by ID.
//	@tags			examples
//	@accept			json
//	@produce		json
//	@param			id		path		int							true	"Example ID"
//	@param			request	body		dto.UpdateExampleRequest	true	"Example payload"
//	@success		200		{object}	dto.BaseResponse{data=dto.ExampleResponse}
//	@failure		400		{object}	dto.BaseResponse
//	@failure		404		{object}	dto.BaseResponse
//	@failure		500		{object}	dto.BaseResponse
//	@router			/examples/{id} [put]
func (h *ExampleHandler) UpdateExample(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.UpdateExampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	example, err := h.service.UpdateExample(c.Request.Context(), id, req)
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if example == nil {
		dto.Error(c, http.StatusNotFound, "example not found")
		return
	}
	dto.OK(c, dto.ToExampleResponse(example))
}

// PatchExample godoc
//
//	@summary		Partial update an example
//	@description	Partially update fields of an existing example by ID.
//	@tags			examples
//	@accept			json
//	@produce		json
//	@param			id		path		int							true	"Example ID"
//	@param			request	body		dto.PatchExampleRequest	true	"Fields to patch"
//	@success		200		{object}	dto.BaseResponse{data=dto.ExampleResponse}
//	@failure		400		{object}	dto.BaseResponse
//	@failure		404		{object}	dto.BaseResponse
//	@failure		500		{object}	dto.BaseResponse
//	@router			/examples/{id} [patch]
func (h *ExampleHandler) PatchExample(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.PatchExampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name == nil && req.Email == nil {
		dto.Error(c, http.StatusBadRequest, "at least one field must be provided")
		return
	}

	example, err := h.service.PatchExample(c.Request.Context(), id, req)
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if example == nil {
		dto.Error(c, http.StatusNotFound, "example not found")
		return
	}
	dto.OK(c, dto.ToExampleResponse(example))
}

// DeleteExample godoc
//
//	@summary		Delete an example
//	@description	Delete an example record by ID.
//	@tags			examples
//	@produce		json
//	@param			id	path		int	true	"Example ID"
//	@success		200	{object}	dto.BaseResponse{data=dto.ExampleDeleteResponse}
//	@failure		400	{object}	dto.BaseResponse
//	@failure		404	{object}	dto.BaseResponse
//	@failure		500	{object}	dto.BaseResponse
//	@router			/examples/{id} [delete]
func (h *ExampleHandler) DeleteExample(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		dto.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	deleted, err := h.service.DeleteExample(c.Request.Context(), id)
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if !deleted {
		dto.Error(c, http.StatusNotFound, "example not found")
		return
	}
	dto.OK(c, dto.ExampleDeleteResponse{ID: id})
}

// GetExampleByEmail godoc
//
//	@summary		Find example by email
//	@description	Search for an example record by email query parameter.
//	@tags			examples
//	@produce		json
//	@param			email	query		string	true	"Email address"
//	@success		200		{object}	dto.BaseResponse{data=dto.ExampleResponse}
//	@failure		400		{object}	dto.BaseResponse
//	@failure		404		{object}	dto.BaseResponse
//	@failure		500		{object}	dto.BaseResponse
//	@router			/examples/search [get]
func (h *ExampleHandler) GetExampleByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		dto.Error(c, http.StatusBadRequest, "email query parameter is required")
		return
	}

	example, err := h.service.GetExampleByEmail(c.Request.Context(), email)
	if err != nil {
		dto.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if example == nil {
		dto.Error(c, http.StatusNotFound, "example not found")
		return
	}

	dto.OK(c, dto.ToExampleResponse(example))
}

// HealthCheck godoc
//
//	@summary		API health check
//	@description	Returns the health status of the API and database connectivity.
//	@tags			health
//	@produce		json
//	@success		200	{object}	dto.BaseResponse
//	@failure		503	{object}	dto.BaseResponse
//	@router			/health [get]
func (h *ExampleHandler) HealthCheck(c *gin.Context) {
	if err := db.DB.Ping(); err != nil {
		dto.Error(c, http.StatusServiceUnavailable, "health check failed", err.Error())
		return
	}
	dto.Success(c, http.StatusOK, gin.H{"status": "healthy"}, "OK")
}
