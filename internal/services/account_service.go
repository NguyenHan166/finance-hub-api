package services

import (
	"finance-hub-api/internal/models"
	"finance-hub-api/internal/repositories"
	"finance-hub-api/internal/utils"
	"fmt"
)

// AccountService handles business logic for accounts
type AccountService struct {
	repo         *repositories.AccountRepository
	vietQRService *utils.VietQRService
}

// NewAccountService creates a new account service
func NewAccountService(repo *repositories.AccountRepository) *AccountService {
	return &AccountService{
		repo:         repo,
		vietQRService: utils.NewVietQRService(),
	}
}

// CreateAccount creates a new account
func (s *AccountService) CreateAccount(userID string, req models.CreateAccountRequest) (*models.Account, error) {
	// Validate account type
	validTypes := map[string]bool{
		"cash":   true,
		"bank":   true,
		"credit": true,
	}
	if !validTypes[req.Type] {
		return nil, fmt.Errorf("invalid account type: must be cash, bank, or credit")
	}

	// Validate bank-specific fields
	if req.Type == "bank" {
		if req.BankCode != nil && *req.BankCode != "" {
			// Fetch bank info from VietQR if bank code is provided
			bank, err := s.vietQRService.GetBankByCode(*req.BankCode)
			if err == nil {
				// Auto-fill bank info
				req.BankName = &bank.Name
				req.BankBIN = &bank.BIN
			}
		}
	}

	// Validate credit card fields
	if req.Type == "credit" {
		if req.CreditLimit == nil || *req.CreditLimit <= 0 {
			return nil, fmt.Errorf("credit limit is required for credit card accounts")
		}
	}

	// Set default icon and color if not provided
	if req.Icon == nil {
		icon := getDefaultIcon(req.Type)
		req.Icon = &icon
	}

	if req.Color == nil {
		color := getDefaultColor(req.Type)
		req.Color = &color
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

	// If bank code is being updated, fetch bank info
	if req.BankCode != nil && *req.BankCode != "" {
		bank, err := s.vietQRService.GetBankByCode(*req.BankCode)
		if err == nil {
			// Auto-fill bank info
			req.BankName = &bank.Name
			req.BankBIN = &bank.BIN
		}
	}

	// Update account
	return s.repo.Update(id, userID, req)
}

// DeleteAccount deletes an account
func (s *AccountService) DeleteAccount(id, userID string) error {
	// Check if account exists
	existing, err := s.repo.GetByID(id, userID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("account not found")
	}

	// TODO: Check if account has transactions before deleting
	// For now, just delete
	return s.repo.Delete(id, userID)
}

// GetAccountSummary retrieves account summary statistics
func (s *AccountService) GetAccountSummary(userID string) (*models.AccountSummary, error) {
	return s.repo.GetSummary(userID)
}

// UpdateBalance updates account balance (used by transactions)
func (s *AccountService) UpdateBalance(accountID, userID string, amount float64) error {
	// Verify account exists and belongs to user
	account, err := s.repo.GetByID(accountID, userID)
	if err != nil {
		return err
	}
	if account == nil {
		return fmt.Errorf("account not found")
	}

	return s.repo.UpdateBalance(accountID, userID, amount)
}

// GetBanks retrieves list of all banks from VietQR
func (s *AccountService) GetBanks() ([]utils.VietQRBank, error) {
	return s.vietQRService.GetBanks()
}

// SearchBanks searches banks by name or code
func (s *AccountService) SearchBanks(query string) ([]utils.VietQRBank, error) {
	return s.vietQRService.SearchBanks(query)
}

// Helper functions
func getDefaultIcon(accountType string) string {
	icons := map[string]string{
		"cash":   "💵",
		"bank":   "🏦",
		"credit": "💳",
	}
	if icon, ok := icons[accountType]; ok {
		return icon
	}
	return "💰"
}

func getDefaultColor(accountType string) string {
	colors := map[string]string{
		"cash":   "#10B981", // Green
		"bank":   "#3B82F6", // Blue
		"credit": "#F59E0B", // Orange
	}
	if color, ok := colors[accountType]; ok {
		return color
	}
	return "#6B7280" // Gray
}
