package services

import (
	"context"
	"errors"
	"finance-hub-api/internal/config"
	"finance-hub-api/internal/models"
	"finance-hub-api/internal/repositories"
	"finance-hub-api/internal/utils"
	"fmt"
	"time"
)

var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrPasswordsNotMatch   = errors.New("passwords do not match")
	ErrWeakPassword        = errors.New("password does not meet requirements")
	ErrUserAlreadyExists   = errors.New("user with this email already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidToken        = errors.New("invalid token")
	ErrGoogleAuthFailed    = errors.New("Google authentication failed")
	ErrEmailNotVerified    = errors.New("email not verified")
	ErrTokenExpired        = errors.New("token has expired")
	ErrTokenAlreadyUsed    = errors.New("token already used")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo      *repositories.UserRepository
	tokenRepo     *repositories.VerificationTokenRepository
	googleClient  *utils.GoogleOAuthClient
	emailService  *utils.EmailService
	jwtSecret     string
	jwtExpiresIn  string
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo *repositories.UserRepository,
	tokenRepo *repositories.VerificationTokenRepository,
	cfg *config.Config,
) *AuthService {
	googleClient := utils.NewGoogleOAuthClient(
		cfg.GoogleOAuth.ClientID,
		cfg.GoogleOAuth.ClientSecret,
		cfg.GoogleOAuth.RedirectURI,
	)

	emailService := utils.NewEmailService(cfg)

	return &AuthService{
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
		googleClient: googleClient,
		emailService: emailService,
		jwtSecret:    cfg.JWT.Secret,
		jwtExpiresIn: cfg.JWT.ExpiresIn,
	}
}

// Register registers a new user with email and password
func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Validate passwords match
	if req.Password != req.ConfirmPassword {
		return nil, ErrPasswordsNotMatch
	}

	// Validate password strength
	if !utils.ValidatePassword(req.Password) {
		return nil, ErrWeakPassword
	}

	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:         req.Email,
		PasswordHash:  &hashedPassword,
		FullName:      req.FullName,
		AuthProvider:  "email",
		EmailVerified: false,
		IsActive:      true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if err == repositories.ErrUserAlreadyExists {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Send verification email (don't fail registration if email fails)
	go func() {
		_ = s.SendVerificationEmail(context.Background(), user.ID)
	}()

	// Generate tokens
	return s.generateAuthResponse(user, true)
}

// Login authenticates a user with email and password
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Check if user has password (might be Google-only user)
	if user.PasswordHash == nil {
		return nil, errors.New("please login with Google")
	}

	// Verify password
	if err := utils.ComparePassword(*user.PasswordHash, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Generate tokens
	return s.generateAuthResponse(user, false)
}

// RefreshToken refreshes access token using refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	// Validate refresh token
	claims, err := utils.ValidateToken(refreshToken, s.jwtSecret)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Find user
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Generate new tokens
	return s.generateAuthResponse(user, false)
}

// InitiateGoogleOAuth generates Google OAuth URL
func (s *AuthService) InitiateGoogleOAuth(state string) string {
	return s.googleClient.GenerateOAuthURL(state)
}

// HandleGoogleCallback handles Google OAuth callback
func (s *AuthService) HandleGoogleCallback(ctx context.Context, code string) (*models.AuthResponse, error) {
	// Exchange code for ID token
	idToken, err := s.googleClient.ExchangeCodeForToken(ctx, code)
	if err != nil {
		fmt.Printf("Error exchanging code: %v\n", err)
		return nil, ErrGoogleAuthFailed
	}

	// Verify ID token and get user info
	return s.VerifyGoogleToken(ctx, idToken)
}

// VerifyGoogleToken verifies Google ID token and creates/updates user
func (s *AuthService) VerifyGoogleToken(ctx context.Context, idToken string) (*models.AuthResponse, error) {
	// Verify token
	fmt.Println("Verifying Google ID token...")
	googleUserInfo, err := s.googleClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		fmt.Printf("Error verifying ID token: %v\n", err)
		return nil, ErrGoogleAuthFailed
	}

	fmt.Printf("Google user info: email=%s, name=%s, sub=%s\n", googleUserInfo.Email, googleUserInfo.Name, googleUserInfo.Sub)

	// Try to find user by Google ID first
	user, err := s.userRepo.FindByGoogleID(ctx, googleUserInfo.Sub)
	
	isNewUser := false

	if err == repositories.ErrUserNotFound {
		fmt.Println("User not found by Google ID, checking email...")
		// Try to find by email (existing email user wants to link Google)
		user, err = s.userRepo.FindByEmail(ctx, googleUserInfo.Email)
		
		if err == repositories.ErrUserNotFound {
			// Create new user
			fmt.Println("Creating new user...")
			user = &models.User{
				Email:         googleUserInfo.Email,
				FullName:      googleUserInfo.Name,
				AvatarURL:     &googleUserInfo.Picture,
				GoogleID:      &googleUserInfo.Sub,
				AuthProvider:  "google",
				EmailVerified: googleUserInfo.EmailVerified,
				IsActive:      true,
			}

			if err := s.userRepo.Create(ctx, user); err != nil {
				fmt.Printf("Error creating user: %v\n", err)
				return nil, err
			}

			fmt.Printf("New user created with ID: %s\n", user.ID)
			isNewUser = true
		} else if err == nil {
			// Link Google to existing email user
			fmt.Println("Linking Google account to existing user...")
			if err := s.userRepo.LinkGoogleAccount(ctx, user.ID, googleUserInfo.Sub, &googleUserInfo.Picture); err != nil {
				fmt.Printf("Error linking Google account: %v\n", err)
				return nil, err
			}
			// Refresh user data
			user, _ = s.userRepo.FindByID(ctx, user.ID)
		} else {
			fmt.Printf("Error finding user by email: %v\n", err)
			return nil, err
		}
	} else if err != nil {
		fmt.Printf("Error finding user by Google ID: %v\n", err)
		return nil, err
	} else {
		fmt.Printf("Found existing user: %s\n", user.Email)
	}

	// Check if user is active
	if !user.IsActive {
		fmt.Println("User account is inactive")
		return nil, errors.New("account is inactive")
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Generate tokens
	fmt.Println("Generating auth tokens...")
	return s.generateAuthResponse(user, isNewUser)
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(ctx context.Context, userID string, req *models.ChangePasswordRequest) error {
	// Find user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Check if user has password (Google-only users don't)
	if user.PasswordHash == nil {
		return errors.New("cannot change password for Google-only accounts")
	}

	// Verify old password
	if err := utils.ComparePassword(*user.PasswordHash, req.OldPassword); err != nil {
		return errors.New("incorrect old password")
	}

	// Validate new password strength
	if !utils.ValidatePassword(req.NewPassword) {
		return ErrWeakPassword
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	return s.userRepo.UpdatePassword(ctx, userID, hashedPassword)
}

// GetUserProfile gets user profile by ID
func (s *AuthService) GetUserProfile(ctx context.Context, userID string) (*models.UserProfile, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	profile := repositories.ToUserProfile(user)
	return &profile, nil
}

// generateAuthResponse generates authentication response with tokens
func (s *AuthService) generateAuthResponse(user *models.User, isNewUser bool) (*models.AuthResponse, error) {
	// Parse JWT expiration duration
	expiresIn, err := utils.ParseTokenDuration(s.jwtExpiresIn)
	if err != nil {
		expiresIn = 24 * time.Hour // default 24 hours
	}

	// Generate access token
	accessToken, err := utils.GenerateToken(user.ID, user.Email, s.jwtSecret, expiresIn)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Create user profile
	profile := repositories.ToUserProfile(user)

	return &models.AuthResponse{
		User:         profile,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(expiresIn.Seconds()),
		IsNewUser:    isNewUser,
	}, nil
}

// SendVerificationEmail sends email verification email
func (s *AuthService) SendVerificationEmail(ctx context.Context, userID string) error {
	// Find user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Check if already verified
	if user.EmailVerified {
		return errors.New("email already verified")
	}

	// Delete old verification tokens
	_ = s.tokenRepo.DeleteByUserIDAndType(ctx, userID, "email_verification")

	// Generate verification token
	tokenStr, err := utils.GenerateVerificationToken()
	if err != nil {
		return err
	}

	// Create token record
	token := &models.VerificationToken{
		UserID:    userID,
		Token:     tokenStr,
		Type:      "email_verification",
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours
		Used:      false,
	}

	if err := s.tokenRepo.Create(ctx, token); err != nil {
		return err
	}

	// Send verification email
	if err := s.emailService.SendVerificationEmail(user.Email, user.FullName, tokenStr); err != nil {
		return err
	}

	return nil
}

// VerifyEmail verifies user's email with token
func (s *AuthService) VerifyEmail(ctx context.Context, tokenStr string) error {
	// Find token
	token, err := s.tokenRepo.FindByToken(ctx, tokenStr, "email_verification")
	if err != nil {
		if err == repositories.ErrTokenNotFound {
			return ErrInvalidToken
		} else if err == repositories.ErrTokenExpired {
			return ErrTokenExpired
		} else if err == repositories.ErrTokenUsed {
			return ErrTokenAlreadyUsed
		}
		return err
	}

	// Mark token as used
	if err := s.tokenRepo.MarkAsUsed(ctx, token.ID); err != nil {
		return err
	}

	// Update user email_verified status
	updates := map[string]interface{}{
		"email_verified": true,
	}

	if err := s.userRepo.Update(ctx, token.UserID, updates); err != nil {
		return err
	}

	return nil
}

// RequestPasswordReset sends password reset email
func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists or not for security
		if err == repositories.ErrUserNotFound {
			return nil // Still return success
		}
		return err
	}

	// Check if user has password (Google-only users can't reset password)
	if user.PasswordHash == nil {
		return errors.New("cannot reset password for Google-only accounts")
	}

	// Delete old reset tokens
	_ = s.tokenRepo.DeleteByUserIDAndType(ctx, user.ID, "password_reset")

	// Generate reset token
	tokenStr, err := utils.GenerateResetToken()
	if err != nil {
		return err
	}

	// Create token record
	token := &models.VerificationToken{
		UserID:    user.ID,
		Token:     tokenStr,
		Type:      "password_reset",
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour
		Used:      false,
	}

	if err := s.tokenRepo.Create(ctx, token); err != nil {
		return err
	}

	// Send reset email
	if err := s.emailService.SendPasswordResetEmail(user.Email, user.FullName, tokenStr); err != nil {
		return err
	}

	return nil
}

// ResetPassword resets user password with token
func (s *AuthService) ResetPassword(ctx context.Context, tokenStr, newPassword string) error {
	// Validate password strength
	if !utils.ValidatePassword(newPassword) {
		return ErrWeakPassword
	}

	// Find token
	token, err := s.tokenRepo.FindByToken(ctx, tokenStr, "password_reset")
	if err != nil {
		if err == repositories.ErrTokenNotFound {
			return ErrInvalidToken
		} else if err == repositories.ErrTokenExpired {
			return ErrTokenExpired
		} else if err == repositories.ErrTokenUsed {
			return ErrTokenAlreadyUsed
		}
		return err
	}

	// Mark token as used
	if err := s.tokenRepo.MarkAsUsed(ctx, token.ID); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update user password
	if err := s.userRepo.UpdatePassword(ctx, token.UserID, hashedPassword); err != nil {
		return err
	}

	// Delete all other reset tokens for this user
	_ = s.tokenRepo.DeleteByUserIDAndType(ctx, token.UserID, "password_reset")

	return nil
}

// ResendVerificationEmail resends verification email
func (s *AuthService) ResendVerificationEmail(ctx context.Context, email string) error {
	// Find user
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}

	// Check if already verified
	if user.EmailVerified {
		return errors.New("email already verified")
	}

	// Send verification email
	return s.SendVerificationEmail(ctx, user.ID)
}
