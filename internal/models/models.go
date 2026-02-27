package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID             string    `json:"id" bson:"_id,omitempty"`
	Email          string    `json:"email" bson:"email"`
	PasswordHash   *string   `json:"-" bson:"password_hash,omitempty"` // Hidden from JSON
	FullName       string    `json:"full_name" bson:"full_name"`
	AvatarURL      *string   `json:"avatar_url,omitempty" bson:"avatar_url,omitempty"`
	GoogleID       *string   `json:"-" bson:"google_id,omitempty"` // Hidden from JSON
	AuthProvider   string    `json:"auth_provider" bson:"auth_provider"` // "email", "google"
	EmailVerified  bool      `json:"email_verified" bson:"email_verified"`
	IsActive       bool      `json:"is_active" bson:"is_active"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty" bson:"last_login_at,omitempty"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}

// UserProfile represents public user profile (safe to send to client)
type UserProfile struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	FullName      string     `json:"full_name"`
	AvatarURL     *string    `json:"avatar_url,omitempty"`
	AuthProvider  string     `json:"auth_provider"`
	EmailVerified bool       `json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
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

// Auth DTOs

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	FullName        string `json:"full_name" binding:"required"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// GoogleTokenRequest represents Google ID token verification request
type GoogleTokenRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

// AuthResponse represents successful authentication response
type AuthResponse struct {
	User         UserProfile `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"`
	IsNewUser    bool        `json:"is_new_user"`
}

// RefreshTokenRequest represents token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// GoogleUserInfo represents user info from Google
type GoogleUserInfo struct {
	Sub           string `json:"sub"`            // Google User ID
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// ChangePasswordRequest represents change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ResetPasswordRequest represents reset password request
type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ConfirmResetPasswordRequest represents confirm reset password with token
type ConfirmResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
