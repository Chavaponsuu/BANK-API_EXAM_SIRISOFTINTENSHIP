package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/krizad/go-gin-api/dto"
	"github.com/krizad/go-gin-api/models"
	"github.com/krizad/go-gin-api/services"
)

type AccountHandler struct {
	service services.AccountService
}

func NewAccountHandler(service services.AccountService) *AccountHandler {
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
//	@router			/api/v1/accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	account := &models.Account{
		OwnerName:   req.OwnerName,
		CitizenID:   req.CitizenID,
		PhoneNumber: req.PhoneNumber,
		AccountType: req.AccountType,
		Balance:     req.InitialBalance,
	}

	createdAccount, err := h.service.CreateAccount(c.Request.Context(), account)
	if err != nil {
		switch err.Error() {

		case "invalid input: citizen_id must be 13 digits", "invalid input: account_type must be SAVING or CURRENT", "balance cannot be negative":
			dto.Error(c, http.StatusBadRequest, err.Error())
		case "account already exists":
			dto.Error(c, http.StatusConflict, err.Error())
		default:
			dto.Error(c, http.StatusInternalServerError, "failed to create account")
		}
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
//	@router			/api/v1/accounts/{account_number} [get]
func (h *AccountHandler) GetAccountByNumber(c *gin.Context) {
	accountNumber := c.Param("account_number")

	account, err := h.service.GetAccountByNumber(c.Request.Context(), accountNumber)

	if err != nil {
		switch err.Error() {
		case "invalid account_number format":
			dto.Error(c, http.StatusBadRequest, err.Error())
		case "account not found":
			dto.Error(c, http.StatusNotFound, err.Error())
		default:
			dto.Error(c, http.StatusInternalServerError, err.Error())
		}
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
//	@router			/api/v1/accounts [get]
func (h *AccountHandler) GetAccountList(c *gin.Context) {
	// Parse query parameters with defaults
	page := 1
	limit := 10

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

	// Call service
	accounts, total, err := h.service.GetAccountList(c.Request.Context(), page, limit)
	if err != nil {
		switch err.Error() {
		case "invalid page or limit parameters":
			dto.Error(c, http.StatusBadRequest, err.Error())
		case "failed to get account list":
			dto.Error(c, http.StatusInternalServerError, err.Error())
		case "failed to count accounts":
			dto.Error(c, http.StatusInternalServerError, err.Error())
		default:
			dto.Error(c, http.StatusInternalServerError, "failed to get account list")
		}
		return
	}

	// Convert to response DTOs
	accountResponses := make([]dto.AccountResponse, 0, len(accounts))
	for _, acc := range accounts {
		accountResponses = append(accountResponses, *dto.ToAccountResponse(acc))
	}

	// Return with pagination metadata
	dto.WithMeta(c, http.StatusOK, accountResponses, "OK", &dto.Meta{
		Page:    page,
		PerPage: limit,
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
//	@router			/api/v1/accounts/{account_number}/close [patch]
func (h *AccountHandler) CloseAccountHandler(c *gin.Context) {
	accountNumber := c.Param("account_number")
	account, err := h.service.CloseAccount(c.Request.Context(), accountNumber)
	if err != nil {
		switch err.Error() {
		case "invalid account_number format":
			dto.Error(c, http.StatusBadRequest, err.Error())
		case "account not found":
			dto.Error(c, http.StatusNotFound, err.Error())
		case "account is already closed":
			dto.Error(c, http.StatusConflict, err.Error())
		default:
			dto.Error(c, http.StatusInternalServerError, err.Error())

		}
		return

	}

	dto.OK(c, dto.ToCloseAccountResponse(account))
}
