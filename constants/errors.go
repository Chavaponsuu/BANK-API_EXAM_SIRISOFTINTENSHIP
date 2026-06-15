package constants

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krizad/go-gin-api/dto"
)

// Account-related errors
var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrAccountAlreadyClosed = errors.New("account is already closed")
	ErrAccountNotActive     = errors.New("account is not active")
	ErrCitizenIDExists      = errors.New("citizen_id already exists")
	ErrNegativeBalance      = errors.New("balance cannot be negative")
)

// Transaction-related errors
var (
	ErrInsufficientBalance   = errors.New("insufficient balance")
	ErrFailedToUpdateBalance = errors.New("failed to update account balance")
)

// Generic errors
var (
	ErrInternal        = errors.New("internal error")
	ErrOperationFailed = errors.New("operation failed")
)

// ErrorStatusMap maps errors to HTTP status codes
var ErrorStatusMap = map[error]int{
	ErrCitizenIDExists:      http.StatusConflict,
	ErrAccountNotFound:      http.StatusNotFound,
	ErrAccountAlreadyClosed: http.StatusConflict,
	ErrOperationFailed:      http.StatusInternalServerError,
	ErrAccountNotActive:     http.StatusConflict,
	ErrInsufficientBalance:  http.StatusBadRequest,
	ErrInternal:             http.StatusInternalServerError,
}

// HandleError handles error responses with appropriate HTTP status codes
func HandleError(c *gin.Context, err error) {
	status, ok := ErrorStatusMap[err]
	if !ok {
		status = http.StatusInternalServerError
	}
	dto.Error(c, status, err.Error())
}
