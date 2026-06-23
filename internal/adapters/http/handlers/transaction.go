package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krizad/go-gin-api/constants"
	"github.com/krizad/go-gin-api/internal/adapters/http/dto"
	"github.com/krizad/go-gin-api/internal/core/ports"
)

type TransactionHandler struct {
	service ports.TransactionService
}

func NewTransactionHandler(service ports.TransactionService) *TransactionHandler {
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
//	@router			/accounts/{account_number}/deposit [post]
func (t *TransactionHandler) Deposit(c *gin.Context) {
	var uri dto.AccountNumberURI
	if err := c.ShouldBindUri(&uri); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	var req dto.TransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	deposit, err := t.service.Deposit(c.Request.Context(), uri.AccountNumber, req.Amount, req.Description)
	if err != nil {
		constants.HandleError(c, err)
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
//	@router			/accounts/{account_number}/withdraw [post]
func (t *TransactionHandler) Withdraw(c *gin.Context) {
	var uri dto.AccountNumberURI
	if err := c.ShouldBindUri(&uri); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	var req dto.TransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	withdraw, err := t.service.Withdraw(c.Request.Context(), uri.AccountNumber, req.Amount, req.Description)
	if err != nil {
		constants.HandleError(c, err)
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
//	@router			/accounts/{account_number}/transactions [get]
func (t *TransactionHandler) GetTransactionHistory(c *gin.Context) {

	var uri dto.AccountNumberURI
	if err := c.ShouldBindUri(&uri); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	var query dto.Pagination
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return

	}
	transactions, total, err := t.service.GetTransactionHistory(c, uri.AccountNumber, query.Page, query.Limit)

	if err != nil {
		constants.HandleError(c, err)
		return
	}
	transactionResponses := make([]dto.TransactionResponse, 0, len(transactions))
	for _, tx := range transactions {
		transactionResponses = append(transactionResponses, *dto.ToTransactionResponse(tx))
	}
	meta := &dto.Meta{
		Page:    query.Page,
		PerPage: query.Limit,
		Total:   total,
	}

	dto.WithMeta(c, http.StatusOK, transactionResponses, "OK", meta)

}



// GetAllTransactionHistory godoc
//
//	@summary		Get all transaction history
//	@description	Retrieve paginated transaction history for all accounts, ordered by newest first (created_at DESC).
//	@tags		transactions
//	@produce	json
//	@param		page	query	int	false	"Page number (default: 1, min: 1)" 	default(1)
//	@param		limit	query	int	false	"Items per page (default: 10, max: 100)" 	default(10)
//	@success	200	{object} dto.BaseResponse{data=[]dto.TransactionResponse,meta=dto.Meta}	"Transaction history with pagination"
//	@failure	400	{object} dto.BaseResponse	"Invalid pagination parameters"
//	@failure	500	{object} dto.BaseResponse	"Failed to get transaction history"
//	@router		/transactions [get]
func (t *TransactionHandler) GetAllTransactionHistory(c *gin.Context) {
	var query dto.Pagination
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return

	}
	transactions, total, err := t.service.GetAllTransaction(c, query.Page, query.Limit)

	if err != nil {
		constants.HandleError(c, err)
		return
	}
	transactionResponses := make([]dto.TransactionResponse, 0, len(transactions))
	for _, tx := range transactions {
		transactionResponses = append(transactionResponses, *dto.ToTransactionResponse(tx))
	}
	meta := &dto.Meta{
		Page:    query.Page,
		PerPage: query.Limit,
		Total:   total,
	}

	dto.WithMeta(c, http.StatusOK, transactionResponses, "OK", meta)

}
