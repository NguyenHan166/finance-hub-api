package handlers

import (
	"finance-hub-api/internal/models"
	"finance-hub-api/internal/services"
	"finance-hub-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AccountHandler handles account-related HTTP requests
type AccountHandler struct {
	service *services.AccountService
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(service *services.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

// CreateAccount handles POST /accounts
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req models.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err.Error())
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID := userIDStr.(string)

	account, err := h.service.CreateAccount(userID, req)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Failed to create account", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusCreated, "Account created successfully", account)
}

// GetAccount handles GET /accounts/:id
func (h *AccountHandler) GetAccount(c *gin.Context) {
	id := c.Param("id")

	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	account, err := h.service.GetAccount(id, userID)
	if err != nil {
		response.NotFoundResponse(c, "Account")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Account retrieved successfully", account)
}

// GetAllAccounts handles GET /accounts
func (h *AccountHandler) GetAllAccounts(c *gin.Context) {
	var pagination models.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.ValidationErrorResponse(c, err.Error())
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	result, err := h.service.GetAllAccounts(userID, pagination)
	if err != nil {
		response.InternalErrorResponse(c, err)
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Accounts retrieved successfully", result)
}

// UpdateAccount handles PUT /accounts/:id
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err.Error())
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	account, err := h.service.UpdateAccount(id, userID, req)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Failed to update account", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Account updated successfully", account)
}

// DeleteAccount handles DELETE /accounts/:id
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	id := c.Param("id")

	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	if err := h.service.DeleteAccount(id, userID); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Failed to delete account", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Account deleted successfully", nil)
}
