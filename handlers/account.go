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
//	@description	Create a new bank account record.
//	@tags			accounts
//	@accept			json
//	@produce		json
//	@param			request	body		dto.CreateAccountRequest	true	"Account payload"
//	@success		201		{object}	dto.BaseResponse{data=dto.AccountResponse}
//	@failure		400		{object}	dto.BaseResponse
//	@failure		500		{object}	dto.BaseResponse
//	@router			/accounts [post]
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
		case "balance cannot be negative":
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
//	@description	Retrieve account details by account number.
//	@tags			accounts
//	@produce		json
//	@param			account_number	path		string	true	"Account Number (10 digits)"
//	@success		200				{object}	dto.BaseResponse{data=dto.AccountResponse}
//	@failure		400				{object}	dto.BaseResponse
//	@failure		404				{object}	dto.BaseResponse
//	@failure		500				{object}	dto.BaseResponse
//	@router			/accounts/{account_number} [get]
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
//	@description	Retrieve a paginated list of all accounts.
//	@tags			accounts
//	@produce		json
//	@param			page	query		int		false	"Page number (default: 1)"
//	@param			limit	query		int		false	"Items per page (default: 10, max: 100)"
//	@success		200		{object}	dto.BaseResponse{data=[]dto.AccountResponse}
//	@failure		400		{object}	dto.BaseResponse
//	@failure		500		{object}	dto.BaseResponse
//	@router			/accounts [get]
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
		dto.Error(c, http.StatusInternalServerError, err.Error())
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
