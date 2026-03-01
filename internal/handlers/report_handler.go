package handlers

import (
	"finance-hub-api/internal/models"
	"finance-hub-api/internal/services"
	"finance-hub-api/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ReportHandler handles report-related HTTP requests
type ReportHandler struct {
	service *services.ReportService
}

// NewReportHandler creates a new report handler
func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GetOverview handles GET /reports/overview?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func (h *ReportHandler) GetOverview(c *gin.Context) {
	var query models.DateRangeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationErrorResponse(c, err.Error())
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.UnauthorizedResponse(c, "User not authenticated")
		return
	}
	userID := userIDStr.(string)

	// Parse dates
	startDate, err := time.Parse("2006-01-02", query.StartDate)
	if err != nil {
		response.ValidationErrorResponse(c, "Invalid start_date format, expected YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", query.EndDate)
	if err != nil {
		response.ValidationErrorResponse(c, "Invalid end_date format, expected YYYY-MM-DD")
		return
	}

	// Set time to end of day for endDate
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	report, err := h.service.GetOverview(userID, startDate, endDate)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate overview report", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Overview report retrieved successfully", report)
}

// GetByCategory handles GET /reports/by-category?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func (h *ReportHandler) GetByCategory(c *gin.Context) {
	var query models.DateRangeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationErrorResponse(c, err.Error())
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.UnauthorizedResponse(c, "User not authenticated")
		return
	}
	userID := userIDStr.(string)

	// Parse dates
	startDate, err := time.Parse("2006-01-02", query.StartDate)
	if err != nil {
		response.ValidationErrorResponse(c, "Invalid start_date format, expected YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", query.EndDate)
	if err != nil {
		response.ValidationErrorResponse(c, "Invalid end_date format, expected YYYY-MM-DD")
		return
	}

	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	report, err := h.service.GetByCategory(userID, startDate, endDate)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate category report", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Category report retrieved successfully", report)
}

// GetByMerchant handles GET /reports/by-merchant?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func (h *ReportHandler) GetByMerchant(c *gin.Context) {
	var query models.DateRangeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationErrorResponse(c, err.Error())
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.UnauthorizedResponse(c, "User not authenticated")
		return
	}
	userID := userIDStr.(string)

	// Parse dates
	startDate, err := time.Parse("2006-01-02", query.StartDate)
	if err != nil {
		response.ValidationErrorResponse(c, "Invalid start_date format, expected YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", query.EndDate)
	if err != nil {
		response.ValidationErrorResponse(c, "Invalid end_date format, expected YYYY-MM-DD")
		return
	}

	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	report, err := h.service.GetByMerchant(userID, startDate, endDate)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate merchant report", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Merchant report retrieved successfully", report)
}

// GetWeeklySpending handles GET /reports/weekly-spending?month=YYYY-MM&category_id=xxx
func (h *ReportHandler) GetWeeklySpending(c *gin.Context) {
	month := c.Query("month")
	if month == "" {
		response.ValidationErrorResponse(c, "Month parameter is required (format: YYYY-MM)")
		return
	}

	categoryID := c.Query("category_id") // Optional
	var categoryIDPtr *string
	if categoryID != "" {
		categoryIDPtr = &categoryID
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.UnauthorizedResponse(c, "User not authenticated")
		return
	}
	userID := userIDStr.(string)

	report, err := h.service.GetWeeklySpending(userID, month, categoryIDPtr)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate weekly spending report", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Weekly spending report retrieved successfully", report)
}

// GetWeeklyCashflow handles GET /reports/weekly-cashflow?month=YYYY-MM
func (h *ReportHandler) GetWeeklyCashflow(c *gin.Context) {
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

	report, err := h.service.GetWeeklyCashflow(userID, month)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate weekly cashflow report", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Weekly cashflow report retrieved successfully", report)
}
