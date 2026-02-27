package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    string    `json:"user_id" bson:"user_id"`
	FullName  *string   `json:"full_name,omitempty" bson:"full_name,omitempty"`
	AvatarURL *string   `json:"avatar_url,omitempty" bson:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// Account represents a financial account
type Account struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Name      string    `json:"name" bson:"name" binding:"required"`
	Type      string    `json:"type" bson:"type" binding:"required"` // checking, savings, credit, investment
	Balance   float64   `json:"balance" bson:"balance"`
	Currency  string    `json:"currency" bson:"currency" binding:"required"`
	BankName  *string   `json:"bank_name,omitempty" bson:"bank_name,omitempty"`
	IsActive  bool      `json:"is_active" bson:"is_active"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// Category represents a transaction category
type Category struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Name      string    `json:"name" bson:"name" binding:"required"`
	Type      string    `json:"type" bson:"type" binding:"required"` // income, expense
	Icon      *string   `json:"icon,omitempty" bson:"icon,omitempty"`
	Color     *string   `json:"color,omitempty" bson:"color,omitempty"`
	IsDefault bool      `json:"is_default" bson:"is_default"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// Transaction represents a financial transaction
type Transaction struct {
	ID              string    `json:"id" bson:"_id,omitempty"`
	UserID          string    `json:"user_id" bson:"user_id"`
	AccountID       string    `json:"account_id" bson:"account_id" binding:"required"`
	CategoryID      string    `json:"category_id" bson:"category_id" binding:"required"`
	Type            string    `json:"type" bson:"type" binding:"required"` // income, expense, transfer
	Amount          float64   `json:"amount" bson:"amount" binding:"required,gt=0"`
	Description     *string   `json:"description,omitempty" bson:"description,omitempty"`
	TransactionDate time.Time `json:"transaction_date" bson:"transaction_date" binding:"required"`
	Notes           *string   `json:"notes,omitempty" bson:"notes,omitempty"`
	AttachmentURL   *string   `json:"attachment_url,omitempty" bson:"attachment_url,omitempty"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

// Budget represents a budget
type Budget struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
	UserID     string    `json:"user_id" bson:"user_id"`
	CategoryID string    `json:"category_id" bson:"category_id" binding:"required"`
	Amount     float64   `json:"amount" bson:"amount" binding:"required,gt=0"`
	Period     string    `json:"period" bson:"period" binding:"required"` // monthly, yearly
	StartDate  time.Time `json:"start_date" bson:"start_date" binding:"required"`
	EndDate    time.Time `json:"end_date" bson:"end_date" binding:"required"`
	IsActive   bool      `json:"is_active" bson:"is_active"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
}

// DTOs (Data Transfer Objects)

// CreateAccountRequest represents request to create an account
type CreateAccountRequest struct {
	Name     string  `json:"name" binding:"required"`
	Type     string  `json:"type" binding:"required,oneof=checking savings credit investment"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency" binding:"required"`
	BankName *string `json:"bank_name,omitempty"`
}

// UpdateAccountRequest represents request to update an account
type UpdateAccountRequest struct {
	Name     *string  `json:"name,omitempty"`
	Balance  *float64 `json:"balance,omitempty"`
	BankName *string  `json:"bank_name,omitempty"`
	IsActive *bool    `json:"is_active,omitempty"`
}

// CreateTransactionRequest represents request to create a transaction
type CreateTransactionRequest struct {
	AccountID       string    `json:"account_id" binding:"required"`
	CategoryID      string    `json:"category_id" binding:"required"`
	Type            string    `json:"type" binding:"required,oneof=income expense transfer"`
	Amount          float64   `json:"amount" binding:"required,gt=0"`
	Description     *string   `json:"description,omitempty"`
	TransactionDate time.Time `json:"transaction_date" binding:"required"`
	Notes           *string   `json:"notes,omitempty"`
}

// CreateCategoryRequest represents request to create a category
type CreateCategoryRequest struct {
	Name  string  `json:"name" binding:"required"`
	Type  string  `json:"type" binding:"required,oneof=income expense"`
	Icon  *string `json:"icon,omitempty"`
	Color *string `json:"color,omitempty"`
}

// CreateBudgetRequest represents request to create a budget
type CreateBudgetRequest struct {
	CategoryID string    `json:"category_id" binding:"required"`
	Amount     float64   `json:"amount" binding:"required,gt=0"`
	Period     string    `json:"period" binding:"required,oneof=monthly yearly"`
	StartDate  time.Time `json:"start_date" binding:"required"`
	EndDate    time.Time `json:"end_date" binding:"required"`
}

// PaginationQuery represents pagination parameters
type PaginationQuery struct {
	Page  int `form:"page" binding:"min=1"`
	Limit int `form:"limit" binding:"min=1,max=100"`
}

// SetDefaults sets default values for pagination
func (p *PaginationQuery) SetDefaults() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
}

// GetOffset calculates offset for SQL query
func (p *PaginationQuery) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalItems int         `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}
