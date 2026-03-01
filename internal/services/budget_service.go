package services

import (
	"finance-hub-api/internal/models"
	"finance-hub-api/internal/repositories"
	"fmt"
	"time"
)

// BudgetService handles business logic for budgets
type BudgetService struct {
	repo            *repositories.BudgetRepository
	transactionRepo *repositories.TransactionRepository
	categoryRepo    *repositories.CategoryRepository
}

// NewBudgetService creates a new budget service
func NewBudgetService(
	repo *repositories.BudgetRepository,
	transactionRepo *repositories.TransactionRepository,
	categoryRepo *repositories.CategoryRepository,
) *BudgetService {
	return &BudgetService{
		repo:            repo,
		transactionRepo: transactionRepo,
		categoryRepo:    categoryRepo,
	}
}

// CreateOrUpdateBudget creates a new budget or updates if exists
func (s *BudgetService) CreateOrUpdateBudget(userID string, req models.CreateBudgetRequest) (*models.Budget, error) {
	// Validate scope and categoryID
	if req.Scope == "category" && req.CategoryID == nil {
		return nil, fmt.Errorf("category_id is required for category scope")
	}

	if req.Scope == "total" && req.CategoryID != nil {
		return nil, fmt.Errorf("category_id should not be set for total scope")
	}

	// Validate category exists if provided
	if req.CategoryID != nil && *req.CategoryID != "" {
		category, err := s.categoryRepo.GetByID(*req.CategoryID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate category: %v", err)
		}
		if category == nil {
			return nil, fmt.Errorf("category not found")
		}
	}

	// Check if budget already exists
	existing, err := s.repo.GetByMonthAndScope(userID, req.Month, req.Scope, req.CategoryID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		// Update existing budget
		updateReq := models.UpdateBudgetRequest{
			Limit:          &req.Limit,
			AlertEnabled:   &req.AlertEnabled,
			AlertThreshold: req.AlertThreshold,
		}
		updated, err := s.repo.Update(existing.ID, userID, updateReq)
		if err != nil {
			return nil, err
		}

		// Recalculate spent
		if err := s.UpdateBudgetSpent(updated.ID, userID); err != nil {
			return nil, err
		}

		return s.repo.GetByID(updated.ID, userID)
	}

	// Create new budget
	budget, err := s.repo.Create(userID, req)
	if err != nil {
		return nil, err
	}

	// Calculate initial spent
	if err := s.UpdateBudgetSpent(budget.ID, userID); err != nil {
		return nil, err
	}

	return s.repo.GetByID(budget.ID, userID)
}

// GetBudget retrieves a budget by ID
func (s *BudgetService) GetBudget(id, userID string) (*models.Budget, error) {
	budget, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if budget == nil {
		return nil, fmt.Errorf("budget not found")
	}

	// Update spent before returning
	if err := s.UpdateBudgetSpent(budget.ID, userID); err != nil {
		return nil, err
	}

	return s.repo.GetByID(budget.ID, userID)
}

// GetBudgetsByMonth retrieves all budgets for a specific month
func (s *BudgetService) GetBudgetsByMonth(userID, month string) ([]models.Budget, error) {
	budgets, err := s.repo.GetByMonth(userID, month)
	if err != nil {
		return nil, err
	}

	// Update spent for all budgets
	for _, budget := range budgets {
		if err := s.UpdateBudgetSpent(budget.ID, userID); err != nil {
			// Log error but continue
			continue
		}
	}

	// Fetch updated budgets
	return s.repo.GetByMonth(userID, month)
}

// UpdateBudget updates a budget
func (s *BudgetService) UpdateBudget(id, userID string, req models.UpdateBudgetRequest) (*models.Budget, error) {
	// Check if budget exists
	existing, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("budget not found")
	}

	updated, err := s.repo.Update(id, userID, req)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

// DeleteBudget deletes a budget
func (s *BudgetService) DeleteBudget(id, userID string) error {
	// Check if budget exists
	existing, err := s.repo.GetByID(id, userID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("budget not found")
	}

	return s.repo.Delete(id, userID)
}

// UpdateBudgetSpent calculates and updates the spent amount for a budget
func (s *BudgetService) UpdateBudgetSpent(budgetID, userID string) error {
	budget, err := s.repo.GetByID(budgetID, userID)
	if err != nil {
		return err
	}
	if budget == nil {
		return fmt.Errorf("budget not found")
	}

	// Parse month to get start and end dates
	startDate, endDate, err := parseMonthRange(budget.Month)
	if err != nil {
		return err
	}

	var spent float64

	if budget.Scope == "total" {
		// Sum all expense transactions for the month
		filters := models.TransactionFilterQuery{
			StartDate: startDate.Format("2006-01-02"),
			EndDate:   endDate.Format("2006-01-02"),
			Type:      "expense",
		}

		transactions, _, err := s.transactionRepo.GetAll(userID, filters)
		if err != nil {
			return err
		}

		for _, tx := range transactions {
			spent += tx.Amount
		}
	} else if budget.Scope == "category" && budget.CategoryID != nil {
		// Sum expense transactions for the specific category
		filters := models.TransactionFilterQuery{
			StartDate:  startDate.Format("2006-01-02"),
			EndDate:    endDate.Format("2006-01-02"),
			Type:       "expense",
			CategoryID: *budget.CategoryID,
		}

		transactions, _, err := s.transactionRepo.GetAll(userID, filters)
		if err != nil {
			return err
		}

		for _, tx := range transactions {
			spent += tx.Amount
		}
	}

	// Update spent in budget
	return s.repo.UpdateSpent(budgetID, userID, spent)
}

// parseMonthRange converts YYYY-MM format to start and end dates
func parseMonthRange(month string) (time.Time, time.Time, error) {
	// Parse YYYY-MM
	t, err := time.Parse("2006-01", month)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid month format: %v", err)
	}

	// Start of month
	startDate := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)

	// End of month
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	return startDate, endDate, nil
}
