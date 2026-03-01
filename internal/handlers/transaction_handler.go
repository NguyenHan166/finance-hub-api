package handlers

import (
	"finance-hub-api/internal/models"
	"finance-hub-api/internal/services"
	"finance-hub-api/pkg/response"
	"net/http"
	"strconv"

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
	var filters models.TransactionFilterQuery
	if err := c.ShouldBindQuery(&filters); err != nil {
		response.ValidationErrorResponse(c, err.Error())
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	result, err := h.service.GetAllTransactions(userID, filters)
	if err != nil {
		response.InternalErrorResponse(c, err)
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Transactions retrieved successfully", result)
}

// UpdateTransaction handles PUT /transactions/:id
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err.Error())
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	transaction, err := h.service.UpdateTransaction(id, userID, req)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Failed to update transaction", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Transaction updated successfully", transaction)
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

// BulkUpdateCategory handles PUT /transactions/bulk/category
func (h *TransactionHandler) BulkUpdateCategory(c *gin.Context) {
	var req models.BulkUpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err.Error())
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	count, err := h.service.BulkUpdateCategory(userID, req)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Failed to update categories", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Categories updated successfully", gin.H{
		"updated_count": count,
	})
}

// BulkDelete handles DELETE /transactions/bulk
func (h *TransactionHandler) BulkDelete(c *gin.Context) {
	var req models.BulkDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err.Error())
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	count, err := h.service.BulkDelete(userID, req)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Failed to delete transactions", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Transactions deleted successfully", gin.H{
		"deleted_count": count,
	})
}

// GetRecentTransactions handles GET /transactions/recent
func (h *TransactionHandler) GetRecentTransactions(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	// Get limit from query param, default to 5
	limit := 5
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	transactions, err := h.service.GetRecentTransactions(userID, limit)
	if err != nil {
		response.InternalErrorResponse(c, err)
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Recent transactions retrieved successfully", transactions)
}

// GetTransactionSummary handles GET /transactions/summary
func (h *TransactionHandler) GetTransactionSummary(c *gin.Context) {
	var filters models.TransactionFilterQuery
	if err := c.ShouldBindQuery(&filters); err != nil {
		response.ValidationErrorResponse(c, err.Error())
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	summary, err := h.service.GetTransactionSummary(userID, filters)
	if err != nil {
		response.InternalErrorResponse(c, err)
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Transaction summary retrieved successfully", summary)
}
