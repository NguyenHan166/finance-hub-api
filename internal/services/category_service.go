package services

import (
	"finance-hub-api/internal/models"
	"finance-hub-api/internal/repositories"
	"fmt"
)

// CategoryService handles business logic for categories
type CategoryService struct {
	repo        *repositories.CategoryRepository
	transactionRepo *repositories.TransactionRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(repo *repositories.CategoryRepository, transactionRepo *repositories.TransactionRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
		transactionRepo: transactionRepo,
	}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(userID string, req models.CreateCategoryRequest) (*models.Category, error) {
	// Validate category type
	validTypes := map[string]bool{
		"income":  true,
		"expense": true,
		"both":    true,
	}
	if !validTypes[req.Type] {
		return nil, fmt.Errorf("invalid category type")
	}

	// Validate parent exists if ParentID is provided
	if req.ParentID != nil && *req.ParentID != "" {
		parent, err := s.repo.GetByIDWithoutUserCheck(*req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate parent category: %v", err)
		}
		if parent == nil {
			return nil, fmt.Errorf("parent category not found")
		}
		// Parent must belong to the same user
		if parent.UserID != userID {
			return nil, fmt.Errorf("parent category not found")
		}
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
		"both":    true,
	}
	if !validTypes[categoryType] {
		return nil, fmt.Errorf("invalid category type")
	}

	return s.repo.GetByType(userID, categoryType)
}

// GetParentCategories retrieves all parent categories
func (s *CategoryService) GetParentCategories(userID string) ([]models.Category, error) {
	return s.repo.GetParentCategories(userID)
}

// GetChildCategories retrieves all child categories of a parent
func (s *CategoryService) GetChildCategories(userID, parentID string) ([]models.Category, error) {
	// Validate parent exists
	parent, err := s.repo.GetByID(parentID, userID)
	if err != nil {
		return nil, err
	}
	if parent == nil {
		return nil, fmt.Errorf("parent category not found")
	}

	return s.repo.GetChildCategories(userID, parentID)
}

// UpdateCategory updates a category
func (s *CategoryService) UpdateCategory(id, userID string, req models.UpdateCategoryRequest) (*models.Category, error) {
	// Check if category exists
	existing, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("category not found")
	}
	
	// Cannot update default categories
	if existing.IsDefault {
		return nil, fmt.Errorf("cannot update default category")
	}

	// Validate type if provided
	if req.Type != nil {
		validTypes := map[string]bool{
			"income":  true,
			"expense": true,
			"both":    true,
		}
		if !validTypes[*req.Type] {
			return nil, fmt.Errorf("invalid category type")
		}
	}

	// Validate parent exists if ParentID is provided
	if req.ParentID != nil && *req.ParentID != "" {
		// Check for circular reference (cannot set itself as parent)
		if *req.ParentID == id {
			return nil, fmt.Errorf("cannot set category as its own parent")
		}

		parent, err := s.repo.GetByIDWithoutUserCheck(*req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate parent category: %v", err)
		}
		if parent == nil {
			return nil, fmt.Errorf("parent category not found")
		}
		// Parent must belong to the same user
		if parent.UserID != userID {
			return nil, fmt.Errorf("parent category not found")
		}
		
		// Check if new parent is a child of current category (prevent circular hierarchy)
		if parent.ParentID != nil && *parent.ParentID == id {
			return nil, fmt.Errorf("cannot set child category as parent")
		}
	}

	return s.repo.Update(id, userID, req)
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

	// Check if category has children
	childCount, err := s.repo.CountChildCategories(userID, id)
	if err != nil {
		return fmt.Errorf("failed to check child categories: %v", err)
	}
	if childCount > 0 {
		return fmt.Errorf("cannot delete category with child categories")
	}

	// Check if category is used in transactions
	transactionCount, err := s.transactionRepo.CountByCategoryID(id, userID)
	if err != nil {
		return fmt.Errorf("failed to check category usage: %v", err)
	}
	if transactionCount > 0 {
		return fmt.Errorf("cannot delete category that is used in transactions")
	}

	return s.repo.Delete(id, userID)
}

// IsCategoryInUse checks if a category is being used
func (s *CategoryService) IsCategoryInUse(id, userID string) (*models.CategoryUsageResponse, error) {
	// Check if category exists
	_, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	// Count child categories
	childCount, err := s.repo.CountChildCategories(userID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to check child categories: %v", err)
	}

	// Count transactions using this category
	transactionCount, err := s.transactionRepo.CountByCategoryID(id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check category usage: %v", err)
	}

	isInUse := transactionCount > 0 || childCount > 0

	return &models.CategoryUsageResponse{
		CategoryID:       id,
		IsInUse:          isInUse,
		TransactionCount: transactionCount,
		ChildCount:       childCount,
	}, nil
}
