package services

import (
	"finance-hub-api/internal/models"
	"finance-hub-api/internal/repositories"
	"fmt"
)

// TransactionService handles business logic for transactions
type TransactionService struct {
	repo        *repositories.TransactionRepository
	accountRepo *repositories.AccountRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(repo *repositories.TransactionRepository, accountRepo *repositories.AccountRepository) *TransactionService {
	return &TransactionService{
		repo:        repo,
		accountRepo: accountRepo,
	}
}

// CreateTransaction creates a new transaction
func (s *TransactionService) CreateTransaction(userID string, req models.CreateTransactionRequest) (*models.Transaction, error) {
	// Validate transaction type
	validTypes := map[string]bool{
		"income":   true,
		"expense":  true,
		"transfer": true,
	}
	if !validTypes[req.Type] {
		return nil, fmt.Errorf("invalid transaction type")
	}

	// Verify account belongs to user
	account, err := s.accountRepo.GetByID(req.AccountID, userID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, fmt.Errorf("account not found")
	}

	// Check if account has sufficient balance for expenses
	if req.Type == "expense" && account.Balance < req.Amount {
		return nil, fmt.Errorf("insufficient balance")
	}

	// Create transaction
	return s.repo.Create(userID, req)
}

// GetTransaction retrieves a transaction by ID
func (s *TransactionService) GetTransaction(id, userID string) (*models.Transaction, error) {
	transaction, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, fmt.Errorf("transaction not found")
	}
	return transaction, nil
}

// GetAllTransactions retrieves all transactions for a user
func (s *TransactionService) GetAllTransactions(userID string, pagination models.PaginationQuery) (*models.PaginatedResponse, error) {
	pagination.SetDefaults()
	
	transactions, totalCount, err := s.repo.GetAll(userID, pagination)
	if err != nil {
		return nil, err
	}

	totalPages := (totalCount + pagination.Limit - 1) / pagination.Limit

	return &models.PaginatedResponse{
		Data:       transactions,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalItems: totalCount,
		TotalPages: totalPages,
	}, nil
}

// DeleteTransaction deletes a transaction
func (s *TransactionService) DeleteTransaction(id, userID string) error {
	// Check if transaction exists
	existing, err := s.repo.GetByID(id, userID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("transaction not found")
	}

	_, err = s.repo.Delete(id, userID)
	return err
}
