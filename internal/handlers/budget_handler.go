package handlers

import (
	"finance-hub-api/internal/models"
	"finance-hub-api/internal/services"
	"finance-hub-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BudgetHandler handles budget-related HTTP requests
type BudgetHandler struct {
	service *services.BudgetService
}

// NewBudgetHandler creates a new budget handler
func NewBudgetHandler(service *services.BudgetService) *BudgetHandler {
	return &BudgetHandler{service: service}
}

// CreateOrUpdateBudget handles POST /budgets
func (h *BudgetHandler) CreateOrUpdateBudget(c *gin.Context) {
	var req models.CreateBudgetRequest
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

	budget, err := h.service.CreateOrUpdateBudget(userID, req)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Failed to create or update budget", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Budget saved successfully", budget)
}

// GetBudget handles GET /budgets/:id
func (h *BudgetHandler) GetBudget(c *gin.Context) {
	id := c.Param("id")

	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID := userIDStr.(string)

	budget, err := h.service.GetBudget(id, userID)
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "Budget not found", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Budget retrieved successfully", budget)
}

// GetBudgetsByMonth handles GET /budgets?month=YYYY-MM
func (h *BudgetHandler) GetBudgetsByMonth(c *gin.Context) {
	month := c.Query("month")
	if month == "" {
		response.ValidationErrorResponse(c, "Month parameter is required (format: YYYY-MM)")
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID := userIDStr.(string)

	budgets, err := h.service.GetBudgetsByMonth(userID, month)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve budgets", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Budgets retrieved successfully", budgets)
}

// UpdateBudget handles PUT /budgets/:id
func (h *BudgetHandler) UpdateBudget(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateBudgetRequest
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

	budget, err := h.service.UpdateBudget(id, userID, req)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Failed to update budget", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Budget updated successfully", budget)
}

// DeleteBudget handles DELETE /budgets/:id
func (h *BudgetHandler) DeleteBudget(c *gin.Context) {
	id := c.Param("id")

	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID := userIDStr.(string)

	if err := h.service.DeleteBudget(id, userID); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Failed to delete budget", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Budget deleted successfully", nil)
}
