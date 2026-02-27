package services

import (
	"finance-hub-api/internal/models"
	"finance-hub-api/internal/repositories"
	"fmt"
)

// CategoryService handles business logic for categories
type CategoryService struct {
	repo *repositories.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(userID string, req models.CreateCategoryRequest) (*models.Category, error) {
	// Validate category type
	validTypes := map[string]bool{
		"income":  true,
		"expense": true,
	}
	if !validTypes[req.Type] {
		return nil, fmt.Errorf("invalid category type")
	}

	return s.repo.Create(userID, req)
}

// GetCategory retrieves a category by ID
func (s *CategoryService) GetCategory(id, userID string) (*models.Category, error) {
	category, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}
	return category, nil
}

// GetAllCategories retrieves all categories for a user
func (s *CategoryService) GetAllCategories(userID string) ([]models.Category, error) {
	return s.repo.GetAll(userID)
}

// GetCategoriesByType retrieves categories by type
func (s *CategoryService) GetCategoriesByType(userID string, categoryType string) ([]models.Category, error) {
	validTypes := map[string]bool{
		"income":  true,
		"expense": true,
	}
	if !validTypes[categoryType] {
		return nil, fmt.Errorf("invalid category type")
	}

	return s.repo.GetByType(userID, categoryType)
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(id, userID string) error {
	// Check if category exists
	existing, err := s.repo.GetByID(id, userID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("category not found")
	}
	if existing.IsDefault {
		return fmt.Errorf("cannot delete default category")
	}

	return s.repo.Delete(id, userID)
}
