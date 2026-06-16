package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krizad/go-gin-api/constants"
	"github.com/krizad/go-gin-api/internal/adapters/http/dto"
	"github.com/krizad/go-gin-api/internal/core/ports"
)

type AccountHandler struct {
	service ports.AccountService
}

func NewAccountHandler(service ports.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

// CreateAccount godoc
//
//	@summary		Create a new account
//	@description	Create a new bank account record. If initial_balance > 0, creates a DEPOSIT transaction atomically.
//	@tags			accounts
//	@accept			json
//	@produce		json
//	@param			request	body		dto.CreateAccountRequest	true	"Account payload"
//	@success		201		{object}	dto.BaseResponse{data=dto.AccountResponse}	"Account created successfully"
//	@failure		400		{object}	dto.BaseResponse	"Invalid input: citizen_id must be 13 digits, account_type must be SAVING or CURRENT, or balance cannot be negative"
//	@failure		409		{object}	dto.BaseResponse	"Citizen ID already exists"
//	@failure		500		{object}	dto.BaseResponse	"Failed to create account"
//	@router			/accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	account := dto.ToCreateAccountRequest(&req)

	createdAccount, err := h.service.CreateAccount(c.Request.Context(), account)
	if err != nil {
		constants.HandleError(c, err)
		return
	}

	dto.Created(c, dto.ToAccountResponse(createdAccount))
}

// GetAccountByNumber godoc
//
//	@summary		Get account by account number
//	@description	Retrieve account details by 10-digit account number.
//	@tags			accounts
//	@produce		json
//	@param			account_number	path		string	true	"Account Number (10 digits)"	example:"0000000001"
//	@success		200				{object}	dto.BaseResponse{data=dto.AccountResponse}	"Account found"
//	@failure		400				{object}	dto.BaseResponse	"Invalid account_number format"
//	@failure		404				{object}	dto.BaseResponse	"Account not found"
//	@failure		500				{object}	dto.BaseResponse	"Internal server error"
//	@router			/accounts/{account_number} [get]
func (h *AccountHandler) GetAccountByNumber(c *gin.Context) {
	var uri dto.AccountNumberURI
	if err := c.ShouldBindUri(&uri); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	account, err := h.service.GetAccountByNumber(c.Request.Context(), uri.AccountNumber)

	if err != nil {
		constants.HandleError(c, err)
		return
	}

	dto.OK(c, dto.ToAccountResponse(account))
}

// GetAccountList godoc
//
//	@summary		List all accounts with pagination
//	@description	Retrieve a paginated list of all accounts ordered by ID.
//	@tags			accounts
//	@produce		json
//	@param			page	query		int		false	"Page number (default: 1, min: 1)"	default(1)
//	@param			limit	query		int		false	"Items per page (default: 10, max: 100)"	default(10)
//	@success		200		{object}	dto.BaseResponse{data=[]dto.AccountResponse,meta=dto.Meta}	"List of accounts with pagination metadata"
//	@failure		400		{object}	dto.BaseResponse	"Invalid page or limit parameters"
//	@failure		500		{object}	dto.BaseResponse	"Failed to get account list"
//	@router			/accounts [get]
func (h *AccountHandler) GetAccountList(c *gin.Context) {
	// Parse query parameters with defaults
	var query dto.Pagination
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return

	}

	// Call service
	accounts, total, err := h.service.GetAccountList(c.Request.Context(), query.Page, query.Limit)
	if err != nil {
		constants.HandleError(c, err)
		return
	}

	// Convert to response DTOs
	accountResponses := make([]dto.AccountResponse, 0, len(accounts))
	for _, acc := range accounts {
		accountResponses = append(accountResponses, *dto.ToAccountResponse(acc))
	}

	// Return with pagination metadata
	dto.WithMeta(c, http.StatusOK, accountResponses, "OK", &dto.Meta{
		Page:    query.Page,
		PerPage: query.Limit,
		Total:   total,
	})
}

// CloseAccount godoc
//
//	@summary		Close an account
//	@description	Change account status from ACTIVE to CLOSED. Closed accounts cannot perform transactions (deposit/withdraw).
//	@tags			accounts
//	@produce		json
//	@param			account_number	path		string	true	"Account Number (10 digits)"	example:"0000000001"
//	@success		200				{object}	dto.BaseResponse{data=dto.CloseAccountResponse}	"Account closed successfully"
//	@failure		400				{object}	dto.BaseResponse	"Invalid account number format"
//	@failure		404				{object}	dto.BaseResponse	"Account not found"
//	@failure		409				{object}	dto.BaseResponse	"Account is already closed"
//	@failure		500				{object}	dto.BaseResponse	"Failed to close account"
//	@router			/accounts/{account_number}/close [patch]
func (h *AccountHandler) CloseAccountHandler(c *gin.Context) {

	var uri dto.AccountNumberURI
	if err := c.ShouldBindUri(&uri); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	account, err := h.service.CloseAccount(c.Request.Context(), uri.AccountNumber)
	if err != nil {
		constants.HandleError(c, err)
		return
	}

	dto.OK(c, dto.ToCloseAccountResponse(account))
}
