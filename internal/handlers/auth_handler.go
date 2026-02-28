package handlers

import (
	"finance-hub-api/internal/models"
	"finance-hub-api/internal/services"
	"finance-hub-api/internal/utils"
	"finance-hub-api/pkg/response"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err)
		return
	}

	authResp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case services.ErrPasswordsNotMatch:
			response.BadRequestResponse(c, "Passwords do not match")
		case services.ErrWeakPassword:
			response.BadRequestResponse(c, "Password must be at least 8 characters with uppercase, lowercase, and number")
		case services.ErrUserAlreadyExists:
			response.ConflictResponse(c, "User with this email already exists")
		default:
			response.InternalErrorResponse(c, err)
		}
		return
	}

	response.SuccessResponse(c, http.StatusCreated, "Registration successful", authResp)
}

// Login handles user login
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err)
		return
	}

	authResp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			response.UnauthorizedResponse(c, "Invalid email or password")
		default:
			if err.Error() == "account is inactive" {
				response.ForbiddenResponse(c, "Account is inactive")
			} else if err.Error() == "please login with Google" {
				response.BadRequestResponse(c, "This account uses Google login. Please sign in with Google.")
			} else {
				response.InternalErrorResponse(c, err)
			}
		}
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Login successful", authResp)
}

// RefreshToken handles token refresh
// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err)
		return
	}

	authResp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch err {
		case services.ErrInvalidToken:
			response.UnauthorizedResponse(c, "Invalid or expired refresh token")
		case services.ErrUserNotFound:
			response.NotFoundResponse(c, "User")
		default:
			response.InternalErrorResponse(c, err)
		}
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", authResp)
}

// InitiateGoogleOAuth initiates Google OAuth flow
// GET /api/v1/auth/google
func (h *AuthHandler) InitiateGoogleOAuth(c *gin.Context) {
	// Generate random state for CSRF protection
	state, err := utils.GenerateRandomState()
	if err != nil {
		response.InternalErrorResponse(c, err)
		return
	}

	// Store state in session or Redis (simplified: use cookie)
	c.SetCookie("oauth_state", state, 600, "/", "", false, true) // 10 minutes

	// Get redirect URI from query parameter
	redirectURI := c.Query("redirect_uri")
	if redirectURI != "" {
		c.SetCookie("oauth_redirect_uri", redirectURI, 600, "/", "", false, true)
	}

	// Generate Google OAuth URL
	authURL := h.authService.InitiateGoogleOAuth(state)

	// Redirect to Google
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// HandleGoogleCallback handles Google OAuth callback
// GET /api/v1/auth/google/callback
func (h *AuthHandler) HandleGoogleCallback(c *gin.Context) {
	// Get parameters
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	// Get redirect URI from cookie first
	redirectURI, _ := c.Cookie("oauth_redirect_uri")
	if redirectURI == "" {
		redirectURI = "http://localhost:3000/auth/callback" // Default frontend URL
	}

	// Check for OAuth errors
	if errorParam != "" {
		errorDesc := c.Query("error_description")
		redirectToFrontend(c, redirectURI, "", errorParam, errorDesc)
		return
	}

	// Verify state (CSRF protection)
	savedState, err := c.Cookie("oauth_state")
	if err != nil || savedState != state {
		redirectToFrontend(c, redirectURI, "", "invalid_state", "Invalid state parameter")
		return
	}

	// Clear cookies
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)
	c.SetCookie("oauth_redirect_uri", "", -1, "/", "", false, true)

	// Handle OAuth callback
	authResp, err := h.authService.HandleGoogleCallback(c.Request.Context(), code)
	if err != nil {
		fmt.Printf("Google OAuth callback error: %v\n", err)
		redirectToFrontend(c, redirectURI, "", "authentication_failed", err.Error())
		return
	}

	// Redirect to frontend with tokens (URL encoded)
	redirectURL := fmt.Sprintf(
		"%s?token=%s&refresh_token=%s&is_new_user=%t",
		redirectURI,
		url.QueryEscape(authResp.AccessToken),
		url.QueryEscape(authResp.RefreshToken),
		authResp.IsNewUser,
	)

	fmt.Printf("Redirecting to: %s\n", redirectURL[:100]+"...") // Log first 100 chars
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// VerifyGoogleToken verifies Google ID token directly
// POST /api/v1/auth/google/token
func (h *AuthHandler) VerifyGoogleToken(c *gin.Context) {
	var req models.GoogleTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err)
		return
	}

	authResp, err := h.authService.VerifyGoogleToken(c.Request.Context(), req.IDToken)
	if err != nil {
		switch err {
		case services.ErrGoogleAuthFailed:
			response.UnauthorizedResponse(c, "Google authentication failed")
		default:
			response.InternalErrorResponse(c, err)
		}
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Google authentication successful", authResp)
}

// GetProfile gets current user profile
// GET /api/v1/auth/profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user ID from JWT middleware
	userID, exists := c.Get("user_id")
	if !exists {
		response.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	profile, err := h.authService.GetUserProfile(c.Request.Context(), userID.(string))
	if err != nil {
		if err == services.ErrUserNotFound {
			response.NotFoundResponse(c, "User")
		} else {
			response.InternalErrorResponse(c, err)
		}
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", profile)
}

// ChangePassword changes user password
// POST /api/v1/auth/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	// Get user ID from JWT middleware
	userID, exists := c.Get("user_id")
	if !exists {
		response.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err)
		return
	}

	err := h.authService.ChangePassword(c.Request.Context(), userID.(string), &req)
	if err != nil {
		switch err {
		case services.ErrWeakPassword:
			response.BadRequestResponse(c, "Password must be at least 8 characters with uppercase, lowercase, and number")
		case services.ErrUserNotFound:
			response.NotFoundResponse(c, "User")
		default:
			if err.Error() == "incorrect old password" {
				response.BadRequestResponse(c, "Incorrect old password")
			} else if err.Error() == "cannot change password for Google-only accounts" {
				response.BadRequestResponse(c, "Cannot change password for Google-only accounts")
			} else {
				response.InternalErrorResponse(c, err)
			}
		}
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Password changed successfully", nil)
}

// Logout handles user logout (client-side token removal)
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT is stateless, so logout is handled client-side by removing token
	// This endpoint is just for consistency and can trigger additional cleanup if needed
	response.SuccessResponse(c, http.StatusOK, "Logged out successfully", nil)
}

// SendVerificationEmail sends verification email to user
// POST /api/v1/auth/send-verification-email
func (h *AuthHandler) SendVerificationEmail(c *gin.Context) {
	// Get user ID from JWT middleware
	userID, exists := c.Get("user_id")
	if !exists {
		response.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	err := h.authService.SendVerificationEmail(c.Request.Context(), userID.(string))
	if err != nil {
		if err.Error() == "email already verified" {
			response.BadRequestResponse(c, "Email already verified")
		} else {
			response.InternalErrorResponse(c, err)
		}
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Verification email sent successfully", nil)
}

// VerifyEmail verifies user email with token
// POST /api/v1/auth/verify-email
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req models.VerifyEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err)
		return
	}

	err := h.authService.VerifyEmail(c.Request.Context(), req.Token)
	if err != nil {
		switch err {
		case services.ErrInvalidToken:
			response.BadRequestResponse(c, "Invalid verification token")
		case services.ErrTokenExpired:
			response.BadRequestResponse(c, "Verification token has expired")
		case services.ErrTokenAlreadyUsed:
			response.BadRequestResponse(c, "Verification token already used")
		default:
			response.InternalErrorResponse(c, err)
		}
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Email verified successfully", nil)
}

// RequestPasswordReset sends password reset email
// POST /api/v1/auth/forgot-password
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req models.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err)
		return
	}

	err := h.authService.RequestPasswordReset(c.Request.Context(), req.Email)
	if err != nil {
		if err.Error() == "cannot reset password for Google-only accounts" {
			response.BadRequestResponse(c, "Cannot reset password for Google-only accounts. Please sign in with Google.")
		} else {
			// Don't expose internal errors for security
			response.SuccessResponse(c, http.StatusOK, "If the email exists, a password reset link has been sent", nil)
			return
		}
		return
	}

	response.SuccessResponse(c, http.StatusOK, "If the email exists, a password reset link has been sent", nil)
}

// ResetPassword resets user password with token
// POST /api/v1/auth/reset-password
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ConfirmResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err)
		return
	}

	err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		switch err {
		case services.ErrInvalidToken:
			response.BadRequestResponse(c, "Invalid reset token")
		case services.ErrTokenExpired:
			response.BadRequestResponse(c, "Reset token has expired. Please request a new one.")
		case services.ErrTokenAlreadyUsed:
			response.BadRequestResponse(c, "Reset token already used. Please request a new one.")
		case services.ErrWeakPassword:
			response.BadRequestResponse(c, "Password must be at least 8 characters with uppercase, lowercase, and number")
		default:
			response.InternalErrorResponse(c, err)
		}
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Password reset successfully", nil)
}

// ResendVerificationEmail resends verification email
// POST /api/v1/auth/resend-verification-email
func (h *AuthHandler) ResendVerificationEmail(c *gin.Context) {
	var req models.ResetPasswordRequest // Reuse same struct (just need email)

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err)
		return
	}

	err := h.authService.ResendVerificationEmail(c.Request.Context(), req.Email)
	if err != nil {
		if err == services.ErrUserNotFound {
			response.NotFoundResponse(c, "User")
		} else if err.Error() == "email already verified" {
			response.BadRequestResponse(c, "Email already verified")
		} else {
			response.InternalErrorResponse(c, err)
		}
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Verification email sent successfully", nil)
}

// Helper function to redirect to frontend with error
func redirectToFrontend(c *gin.Context, redirectURI, token, errorCode, errorDesc string) {
	if redirectURI == "" {
		redirectURI = "http://localhost:3000/auth/callback" // Default
	}

	if errorCode != "" {
		redirectURL := fmt.Sprintf("%s?error=%s", redirectURI, errorCode)
		if errorDesc != "" {
			redirectURL += fmt.Sprintf("&error_description=%s", errorDesc)
		}
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, redirectURI+"?token="+token)
	}
}
