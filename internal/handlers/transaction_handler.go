package handlers

import (
	"finance-hub-api/internal/models"
	"finance-hub-api/internal/services"
	"finance-hub-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TransactionHandler handles transaction-related HTTP requests
type TransactionHandler struct {
	service *services.TransactionService
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// CreateTransaction handles POST /transactions
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req models.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err.Error())
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID := userIDStr.(string)

	transaction, err := h.service.CreateTransaction(userID, req)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Failed to create transaction", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusCreated, "Transaction created successfully", transaction)
}

// GetTransaction handles GET /transactions/:id
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	id := c.Param("id")

	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	transaction, err := h.service.GetTransaction(id, userID)
	if err != nil {
		response.NotFoundResponse(c, "Transaction")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Transaction retrieved successfully", transaction)
}

// GetAllTransactions handles GET /transactions
func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	var pagination models.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.ValidationErrorResponse(c, err.Error())
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	result, err := h.service.GetAllTransactions(userID, pagination)
	if err != nil {
		response.InternalErrorResponse(c, err)
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Transactions retrieved successfully", result)
}

// DeleteTransaction handles DELETE /transactions/:id
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	id := c.Param("id")

	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	if err := h.service.DeleteTransaction(id, userID); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Failed to delete transaction", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Transaction deleted successfully", nil)
}
