package handlers

import (
	"net/http"
	"strconv"

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
// @Summary      Deposit funds into an account
// @Description  Deposit a specified amount into the account with the given account number
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        account_number  path      string               true  "Account Number"
// @Param        request         body      dto.DepositRequest   true  "Deposit Request"
// @Success      200             {object}  dto.TransactionResponse
// @Failure      400             {object}  dto.ErrorResponse
// @Failure      404             {object}  dto.ErrorResponse
// @Failure      409             {object}  dto.ErrorResponse
// @Failure      500             {object}  dto.ErrorResponse
// @Router       /accounts/{account_number}/deposit [post]

func (t *TransactionHandler) Deposit(c *gin.Context) {
	accountNumber := c.Param("account_number")
	var req dto.TransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	deposit, err := t.service.Deposit(c.Request.Context(), accountNumber, req.Amount, req.Description)
	if err != nil {
		switch err.Error() {
		case "amount must be greater than 0":
			dto.Error(c, http.StatusBadRequest, err.Error())
		case "account not found":
			dto.Error(c, http.StatusNotFound, err.Error())
		case "account is not active", "account is already closed":
			dto.Error(c, http.StatusConflict, err.Error())
		default:
			dto.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	dto.OK(c, dto.ToTransactionResponse(deposit))
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	dto.Error(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

}

// Withdraw godoc
// @Summary      Withdraw funds from an account
// @Description  Withdraw a specified amount from the account with the given account number
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        account_number  path      string                true  "Account Number"
// @Param        request         body      dto.WithdrawRequest   true  "Withdraw Request"
// @Success      200             {object}  dto.TransactionResponse
// @Failure      400             {object}  dto.ErrorResponse
// @Failure      404             {object}  dto.ErrorResponse
// @Failure      409             {object}  dto.ErrorResponse
// @Failure      500             {object}  dto.ErrorResponse
// @Router       /accounts/{account_number}/withdraw [post]

func (t *TransactionHandler) Withdraw(c *gin.Context) {
	accountNumber := c.Param("account_number")
	var req dto.TransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	withdraw, err := t.service.Withdraw(c.Request.Context(), accountNumber, req.Amount, req.Description)
	if err != nil {
		switch err.Error() {
		case "insufficient balance":
			dto.Error(c, http.StatusBadRequest, err.Error())
		case "account not found":
			dto.Error(c, http.StatusNotFound, err.Error())
		case "account is not active", "account is already closed", "insufficient funds":
			dto.Error(c, http.StatusConflict, err.Error())
		default:
			dto.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	dto.OK(c, dto.ToTransactionResponse(withdraw))
}

func (t *TransactionHandler) GetTransactionHistory(c *gin.Context) {

	accountNumber := c.Param("account_number")
	page := 1
	limit := 1
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
		case "invalid account_number format":
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
