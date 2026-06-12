package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/krizad/go-gin-api/dto"
	"github.com/krizad/go-gin-api/services"
)

type TransactionHandler struct {
	service services.TransactionService
}

func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// Deposit godoc
//
//	@summary		Deposit money into an account
//	@description	Deposit a specified amount into the account. Creates a DEPOSIT transaction and updates account balance atomically.
//	@tags			transactions
//	@accept			json
//	@produce		json
//	@param			account_number	path		string						true	"Account Number (10 digits)"	example:"0000000001"
//	@param			request			body		dto.TransactionRequest		true	"Deposit request payload"
//	@success		200				{object}	dto.BaseResponse{data=dto.TransactionResponse}	"Deposit successful"
//	@failure		400				{object}	dto.BaseResponse	"Amount must be greater than 0"
//	@failure		404				{object}	dto.BaseResponse	"Account not found"
//	@failure		409				{object}	dto.BaseResponse	"Account is not active or already closed"
//	@failure		500				{object}	dto.BaseResponse	"Failed to update account balance"
//	@router			/api/v1/accounts/{account_number}/deposit [post]
func (t *TransactionHandler) Deposit(c *gin.Context) {
	accountNumber := c.Param("account_number")
	var req dto.TransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	deposit, err := t.service.Deposit(c.Request.Context(), accountNumber, req.Amount, req.Description)
	if err != nil {
		switch {
		case err.Error() == "amount must be greater than 0",
			strings.Contains(err.Error(), "failed to create transaction"):
			dto.Error(c, http.StatusBadRequest, err.Error())

		case err.Error() == "account not found":
			dto.Error(c, http.StatusNotFound, err.Error())

		case err.Error() == "account is not active",
			err.Error() == "account is already closed":
			dto.Error(c, http.StatusConflict, err.Error())

		default:
			dto.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	dto.OK(c, dto.ToTransactionResponse(deposit))
}

// Withdraw godoc
//
//	@summary		Withdraw money from an account
//	@description	Withdraw a specified amount from the account. Validates sufficient balance, creates a WITHDRAW transaction, and updates account balance atomically.
//	@tags			transactions
//	@accept			json
//	@produce		json
//	@param			account_number	path		string						true	"Account Number (10 digits)"	example:"0000000001"
//	@param			request			body		dto.TransactionRequest		true	"Withdraw request payload"
//	@success		200				{object}	dto.BaseResponse{data=dto.TransactionResponse}	"Withdrawal successful"
//	@failure		400				{object}	dto.BaseResponse	"Insufficient balance or amount must be greater than 0"
//	@failure		404				{object}	dto.BaseResponse	"Account not found"
//	@failure		409				{object}	dto.BaseResponse	"Account is not active or already closed"
//	@failure		500				{object}	dto.BaseResponse	"Failed to update account balance"
//	@router			/api/v1/accounts/{account_number}/withdraw [post]
func (t *TransactionHandler) Withdraw(c *gin.Context) {
	accountNumber := c.Param("account_number")
	var req dto.TransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	withdraw, err := t.service.Withdraw(c.Request.Context(), accountNumber, req.Amount, req.Description)
	if err != nil {
		fmt.Print(err.Error())
		switch {
		case err.Error() == "amount must be greater than 0",
			strings.Contains(err.Error(), "failed to create transaction"), err.Error() == "insufficient balance":
			dto.Error(c, http.StatusBadRequest, err.Error())

		case err.Error() == "account not found":
			dto.Error(c, http.StatusNotFound, err.Error())

		case err.Error() == "account is not active",
			err.Error() == "account is already closed":
			dto.Error(c, http.StatusConflict, err.Error())

		default:
			dto.Error(c, http.StatusInternalServerError, err.Error())

		}
		return

	}

	dto.OK(c, dto.ToTransactionResponse(withdraw))
}

// GetTransactionHistory godoc
//
//	@summary		Get transaction history of an account
//	@description	Retrieve paginated transaction history for an account, ordered by newest first (created_at DESC).
//	@tags			transactions
//	@produce		json
//	@param			account_number	path		string	true	"Account Number (10 digits)"	example:"0000000001"
//	@param			page			query		int		false	"Page number (default: 1, min: 1)"	default(1)
//	@param			limit			query		int		false	"Items per page (default: 10, max: 100)"	default(10)
//	@success		200				{object}	dto.BaseResponse{data=[]dto.TransactionResponse,meta=dto.Meta}	"Transaction history with pagination"
//	@failure		400				{object}	dto.BaseResponse	"Invalid account_number format or pagination parameters"
//	@failure		404				{object}	dto.BaseResponse	"Account not found"
//	@failure		500				{object}	dto.BaseResponse	"Failed to get transaction history"
//	@router			/api/v1/accounts/{account_number}/transactions [get]
func (t *TransactionHandler) GetTransactionHistory(c *gin.Context) {

	accountNumber := c.Param("account_number")
	page := -1
	limit := -1
	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}

	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	transactions, total, err := t.service.GetTransactionHistory(c, accountNumber, page, limit)

	if err != nil {
		switch err.Error() {
		case "invalid account_number format", "invalid page or limit parameters":
			dto.Error(c, http.StatusBadRequest, err.Error())
		case "account not found":
			dto.Error(c, http.StatusNotFound, err.Error())
		default:
			dto.Error(c, http.StatusInternalServerError, err.Error())

		}

		return

	}
	transactionResponses := make([]dto.TransactionResponse, 0, len(transactions))
	for _, tx := range transactions {
		transactionResponses = append(transactionResponses, *dto.ToTransactionResponse(tx))
	}
	meta := &dto.Meta{
		Page:    page,
		PerPage: limit,
		Total:   total,
	}

	dto.WithMeta(c, http.StatusOK, transactionResponses, "OK", meta)

}
