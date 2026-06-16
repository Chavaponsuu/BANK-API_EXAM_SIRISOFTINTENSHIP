package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"per_page,omitempty"`
	Total   int `json:"total,omitempty"`
}

type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  []any  `json:"errors,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

func Success(c *gin.Context, status int, data any, message string) {
	c.JSON(status, BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string, errors ...any) {
	resp := BaseResponse{
		Success: false,
		Message: message,
	}
	if len(errors) > 0 {
		resp.Errors = errors
	}
	c.AbortWithStatusJSON(status, resp)
}

func OK(c *gin.Context, data any) {
	Success(c, http.StatusOK, data, "OK")
}

func Created(c *gin.Context, data any) {
	Success(c, http.StatusCreated, data, "Created")
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func WithMeta(c *gin.Context, status int, data any, message string, meta *Meta) {
	c.JSON(status, BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}
