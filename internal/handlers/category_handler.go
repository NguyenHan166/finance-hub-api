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

	// Check query parameters
	categoryType := c.Query("type")
	filter := c.Query("filter")
	parentID := c.Query("parent_id")
	
	var categories []models.Category
	var err error
	
	// Handle different filtering options
	if filter == "parent" {
		// Get only parent categories (no parent_id)
		categories, err = h.service.GetParentCategories(userID)
	} else if filter == "children" && parentID != "" {
		// Get children of specific parent
		categories, err = h.service.GetChildCategories(userID, parentID)
	} else if categoryType != "" {
		// Filter by type
		categories, err = h.service.GetCategoriesByType(userID, categoryType)
	} else {
		// Get all categories
		categories, err = h.service.GetAllCategories(userID)
	}

	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Failed to retrieve categories", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Categories retrieved successfully", categories)
}

// UpdateCategory handles PUT /categories/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err.Error())
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	category, err := h.service.UpdateCategory(id, userID, req)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Failed to update category", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Category updated successfully", category)
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

// CheckCategoryUsage handles GET /categories/:id/usage
func (h *CategoryHandler) CheckCategoryUsage(c *gin.Context) {
	id := c.Param("id")

	userIDStr, _ := c.Get("user_id")
	userID := userIDStr.(string)

	usage, err := h.service.IsCategoryInUse(id, userID)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Failed to check category usage", err.Error())
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Category usage retrieved successfully", usage)
}
