package services

import (
	"finance-hub-api/internal/models"
	"finance-hub-api/internal/repositories"
	"fmt"
)

// AccountService handles business logic for accounts
type AccountService struct {
	repo *repositories.AccountRepository
}

// NewAccountService creates a new account service
func NewAccountService(repo *repositories.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

// CreateAccount creates a new account
func (s *AccountService) CreateAccount(userID string, req models.CreateAccountRequest) (*models.Account, error) {
	// Validate account type
	validTypes := map[string]bool{
		"checking":   true,
		"savings":    true,
		"credit":     true,
		"investment": true,
	}
	if !validTypes[req.Type] {
		return nil, fmt.Errorf("invalid account type")
	}

	// Create account
	return s.repo.Create(userID, req)
}

// GetAccount retrieves an account by ID
func (s *AccountService) GetAccount(id, userID string) (*models.Account, error) {
	account, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, fmt.Errorf("account not found")
	}
	return account, nil
}

// GetAllAccounts retrieves all accounts for a user
func (s *AccountService) GetAllAccounts(userID string, pagination models.PaginationQuery) (*models.PaginatedResponse, error) {
	pagination.SetDefaults()
	
	accounts, totalCount, err := s.repo.GetAll(userID, pagination)
	if err != nil {
		return nil, err
	}

	totalPages := (totalCount + pagination.Limit - 1) / pagination.Limit

	return &models.PaginatedResponse{
		Data:       accounts,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalItems: totalCount,
		TotalPages: totalPages,
	}, nil
}

// UpdateAccount updates an account
func (s *AccountService) UpdateAccount(id, userID string, req models.UpdateAccountRequest) (*models.Account, error) {
	// Check if account exists
	existing, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("account not found")
	}

	// Update account
	return s.repo.Update(id, userID, req)
}

// DeleteAccount deletes an account
func (s *AccountService) DeleteAccount(id, userID string) error {
	return s.repo.Delete(id, userID)
}
