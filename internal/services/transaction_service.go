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
	categoryRepo *repositories.CategoryRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(
	repo *repositories.TransactionRepository, 
	accountRepo *repositories.AccountRepository,
	categoryRepo *repositories.CategoryRepository,
) *TransactionService {
	return &TransactionService{
		repo:         repo,
		accountRepo:  accountRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateTransaction creates a new transaction and updates account balance(s)
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

	// Verify source account belongs to user
	account, err := s.accountRepo.GetByID(req.AccountID, userID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, fmt.Errorf("account not found")
	}

	// Validate transfer transaction
	if req.Type == "transfer" {
		if req.ToAccountID == nil || *req.ToAccountID == "" {
			return nil, fmt.Errorf("to_account_id is required for transfer transactions")
		}
		if *req.ToAccountID == req.AccountID {
			return nil, fmt.Errorf("cannot transfer to the same account")
		}

		// Verify destination account
		toAccount, err := s.accountRepo.GetByID(*req.ToAccountID, userID)
		if err != nil {
			return nil, err
		}
		if toAccount == nil {
			return nil, fmt.Errorf("destination account not found")
		}

		// Transfers don't need a category
		req.CategoryID = nil
	} else {
		// Non-transfer transactions must have a category
		if req.CategoryID == nil || *req.CategoryID == "" {
			return nil, fmt.Errorf("category_id is required for income and expense transactions")
		}

		// Verify category exists and belongs to user
		category, err := s.categoryRepo.GetByID(*req.CategoryID, userID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, fmt.Errorf("category not found")
		}

		// Validate category type matches transaction type
		if category.Type != req.Type && category.Type != "both" {
			return nil, fmt.Errorf("category type does not match transaction type")
		}
	}

	// Check if account has sufficient balance for expenses and transfers
	if (req.Type == "expense" || req.Type == "transfer") && account.Balance < req.Amount {
		return nil, fmt.Errorf("insufficient balance in account %s", account.Name)
	}

	// Create transaction
	transaction, err := s.repo.Create(userID, req)
	if err != nil {
		return nil, err
	}

	// Update account balance(s)
	if err := s.updateAccountBalances(userID, transaction, nil); err != nil {
		// If balance update fails, we should ideally rollback the transaction
		// For now, we'll return the error
		return nil, fmt.Errorf("failed to update account balance: %v", err)
	}

	return transaction, nil
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

// GetAllTransactions retrieves all transactions for a user with filters
func (s *TransactionService) GetAllTransactions(userID string, filters models.TransactionFilterQuery) (*models.PaginatedResponse, error) {
	filters.SetDefaults()
	
	transactions, totalCount, err := s.repo.GetAll(userID, filters)
	if err != nil {
		return nil, err
	}

	totalPages := (totalCount + filters.Limit - 1) / filters.Limit

	return &models.PaginatedResponse{
		Data:       transactions,
		Page:       filters.Page,
		Limit:      filters.Limit,
		TotalItems: totalCount,
		TotalPages: totalPages,
	}, nil
}

// UpdateTransaction updates a transaction and adjusts account balance(s)
func (s *TransactionService) UpdateTransaction(id, userID string, req models.UpdateTransactionRequest) (*models.Transaction, error) {
	// Get existing transaction
	existing, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("transaction not found")
	}

	// If type is being changed, validate it
	if req.Type != nil {
		validTypes := map[string]bool{
			"income":   true,
			"expense":  true,
			"transfer": true,
		}
		if !validTypes[*req.Type] {
			return nil, fmt.Errorf("invalid transaction type")
		}

		// If changing to transfer, require ToAccountID
		if *req.Type == "transfer" && (req.ToAccountID == nil || *req.ToAccountID == "") {
			return nil, fmt.Errorf("to_account_id is required for transfer transactions")
		}

		// If changing from transfer to income/expense, require CategoryID
		if existing.Type == "transfer" && *req.Type != "transfer" && (req.CategoryID == nil || *req.CategoryID == "") {
			return nil, fmt.Errorf("category_id is required for income and expense transactions")
		}
	}

	// Validate account if being changed
	if req.AccountID != nil {
		account, err := s.accountRepo.GetByID(*req.AccountID, userID)
		if err != nil {
			return nil, err
		}
		if account == nil {
			return nil, fmt.Errorf("account not found")
		}
	}

	// Validate destination account if being changed
	if req.ToAccountID != nil && *req.ToAccountID != "" {
		toAccount, err := s.accountRepo.GetByID(*req.ToAccountID, userID)
		if err != nil {
			return nil, err
		}
		if toAccount == nil {
			return nil, fmt.Errorf("destination account not found")
		}
	}

	// Validate category if being changed
	if req.CategoryID != nil && *req.CategoryID != "" {
		category, err := s.categoryRepo.GetByID(*req.CategoryID, userID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, fmt.Errorf("category not found")
		}
	}

	// Revert old balance changes
	if err := s.revertAccountBalances(userID, existing); err != nil {
		return nil, fmt.Errorf("failed to revert account balance: %v", err)
	}

	// Update transaction
	updated, err := s.repo.Update(id, userID, req)
	if err != nil {
		// If update fails, try to restore the old balance
		_ = s.updateAccountBalances(userID, existing, nil)
		return nil, err
	}

	// Apply new balance changes
	if err := s.updateAccountBalances(userID, updated, nil); err != nil {
		// This is a critical error - transaction is updated but balance is not
		return nil, fmt.Errorf("transaction updated but failed to update account balance: %v", err)
	}

	return updated, nil
}

// DeleteTransaction deletes a transaction and reverts account balance changes
func (s *TransactionService) DeleteTransaction(id, userID string) error {
	// Check if transaction exists
	existing, err := s.repo.GetByID(id, userID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("transaction not found")
	}

	// Revert account balance changes
	if err := s.revertAccountBalances(userID, existing); err != nil {
		return fmt.Errorf("failed to revert account balance: %v", err)
	}

	// Delete transaction
	_, err = s.repo.Delete(id, userID)
	if err != nil {
		// If delete fails, try to restore the balance
		_ = s.updateAccountBalances(userID, existing, nil)
		return err
	}

	return nil
}

// BulkUpdateCategory updates category for multiple transactions
func (s *TransactionService) BulkUpdateCategory(userID string, req models.BulkUpdateCategoryRequest) (int64, error) {
	// Verify category exists
	category, err := s.categoryRepo.GetByID(req.CategoryID, userID)
	if err != nil {
		return 0, err
	}
	if category == nil {
		return 0, fmt.Errorf("category not found")
	}

	// Perform bulk update
	count, err := s.repo.BulkUpdateCategory(userID, req.TransactionIDs, req.CategoryID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// BulkDelete deletes multiple transactions and reverts their balance changes
func (s *TransactionService) BulkDelete(userID string, req models.BulkDeleteRequest) (int64, error) {
	// Get all transactions to be deleted
	var transactionsToDelete []*models.Transaction
	for _, id := range req.TransactionIDs {
		transaction, err := s.repo.GetByID(id, userID)
		if err != nil {
			continue // Skip errors, continue with others
		}
		if transaction != nil {
			transactionsToDelete = append(transactionsToDelete, transaction)
		}
	}

	// Revert balance changes for all transactions
	for _, transaction := range transactionsToDelete {
		if err := s.revertAccountBalances(userID, transaction); err != nil {
			// Log error but continue
			fmt.Printf("Warning: failed to revert balance for transaction %s: %v\n", transaction.ID, err)
		}
	}

	// Perform bulk delete
	count, err := s.repo.BulkDelete(userID, req.TransactionIDs)
	if err != nil {
		// If delete fails, try to restore all balances
		for _, transaction := range transactionsToDelete {
			_ = s.updateAccountBalances(userID, transaction, nil)
		}
		return 0, err
	}

	return count, nil
}

// GetRecentTransactions retrieves recent transactions
func (s *TransactionService) GetRecentTransactions(userID string, limit int) ([]models.Transaction, error) {
	if limit <= 0 {
		limit = 5
	}
	if limit > 50 {
		limit = 50
	}
	
	return s.repo.GetRecentTransactions(userID, limit)
}

// GetTransactionSummary retrieves transaction summary statistics
func (s *TransactionService) GetTransactionSummary(userID string, filters models.TransactionFilterQuery) (*models.TransactionSummary, error) {
	return s.repo.GetSummary(userID, filters)
}

// Helper function to update account balances based on transaction
func (s *TransactionService) updateAccountBalances(userID string, transaction *models.Transaction, previousTransaction *models.Transaction) error {
	switch transaction.Type {
	case "income":
		// Add to account balance
		return s.accountRepo.UpdateBalance(transaction.AccountID, userID, transaction.Amount)
	
	case "expense":
		// Subtract from account balance
		return s.accountRepo.UpdateBalance(transaction.AccountID, userID, -transaction.Amount)
	
	case "transfer":
		if transaction.ToAccountID == nil {
			return fmt.Errorf("to_account_id is required for transfer")
		}
		// Subtract from source account
		if err := s.accountRepo.UpdateBalance(transaction.AccountID, userID, -transaction.Amount); err != nil {
			return err
		}
		// Add to destination account
		if err := s.accountRepo.UpdateBalance(*transaction.ToAccountID, userID, transaction.Amount); err != nil {
			// Rollback source account change
			_ = s.accountRepo.UpdateBalance(transaction.AccountID, userID, transaction.Amount)
			return err
		}
	}
	
	return nil
}

// Helper function to revert account balance changes
func (s *TransactionService) revertAccountBalances(userID string, transaction *models.Transaction) error {
	switch transaction.Type {
	case "income":
		// Subtract what was added
		return s.accountRepo.UpdateBalance(transaction.AccountID, userID, -transaction.Amount)
	
	case "expense":
		// Add back what was subtracted
		return s.accountRepo.UpdateBalance(transaction.AccountID, userID, transaction.Amount)
	
	case "transfer":
		if transaction.ToAccountID == nil {
			return fmt.Errorf("to_account_id is required for transfer")
		}
		// Add back to source account
		if err := s.accountRepo.UpdateBalance(transaction.AccountID, userID, transaction.Amount); err != nil {
			return err
		}
		// Subtract from destination account
		if err := s.accountRepo.UpdateBalance(*transaction.ToAccountID, userID, -transaction.Amount); err != nil {
			// Rollback source account change
			_ = s.accountRepo.UpdateBalance(transaction.AccountID, userID, -transaction.Amount)
			return err
		}
	}
	
	return nil
}
