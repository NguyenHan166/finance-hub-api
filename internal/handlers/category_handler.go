package handlers

import (
	"finance-hub-api/internal/models"
	"finance-hub-api/internal/services"
	"finance-hub-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CategoryHandler handles category-related HTTP requests
type CategoryHandler struct {
	service *services.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// CreateCategory handles POST /categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
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

	category, err := h.service.CreateCategory(userID, req)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Failed to create category", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusCreated, "Category created successfully", category)
}

// GetCategory handles GET /categories/:id
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")

	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	category, err := h.service.GetCategory(id, userID)
	if err != nil {
		response.NotFoundResponse(c, "Category")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", category)
}

// GetAllCategories handles GET /categories
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	// Check if filtering by type
	categoryType := c.Query("type")
	
	var categories []models.Category
	var err error
	
	if categoryType != "" {
		categories, err = h.service.GetCategoriesByType(userID, categoryType)
	} else {
		categories, err = h.service.GetAllCategories(userID)
	}

	if err != nil {
		response.InternalErrorResponse(c, err)
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Categories retrieved successfully", categories)
}

// DeleteCategory handles DELETE /categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	if err := h.service.DeleteCategory(id, userID); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Failed to delete category", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Category deleted successfully", nil)
}
