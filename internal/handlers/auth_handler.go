package handlers

import (
	"finance-hub-api/internal/models"
	"finance-hub-api/internal/services"
	"finance-hub-api/internal/utils"
	"finance-hub-api/pkg/response"
	"fmt"
	"net/http"

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

	// Check for OAuth errors
	if errorParam != "" {
		errorDesc := c.Query("error_description")
		redirectToFrontend(c, "", "", errorParam, errorDesc)
		return
	}

	// Verify state (CSRF protection)
	savedState, err := c.Cookie("oauth_state")
	if err != nil || savedState != state {
		redirectToFrontend(c, "", "", "invalid_state", "Invalid state parameter")
		return
	}

	// Clear state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// Get redirect URI from cookie
	redirectURI, _ := c.Cookie("oauth_redirect_uri")
	if redirectURI == "" {
		redirectURI = "http://localhost:5173" // Default frontend URL
	}
	c.SetCookie("oauth_redirect_uri", "", -1, "/", "", false, true)

	// Handle OAuth callback
	authResp, err := h.authService.HandleGoogleCallback(c.Request.Context(), code)
	if err != nil {
		redirectToFrontend(c, redirectURI, "", "authentication_failed", err.Error())
		return
	}

	// Redirect to frontend with tokens
	redirectURL := fmt.Sprintf(
		"%s/auth/callback?token=%s&refresh_token=%s&is_new_user=%t",
		redirectURI,
		authResp.AccessToken,
		authResp.RefreshToken,
		authResp.IsNewUser,
	)

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

// Helper function to redirect to frontend with error
func redirectToFrontend(c *gin.Context, redirectURI, token, errorCode, errorDesc string) {
	if redirectURI == "" {
		redirectURI = "http://localhost:5173" // Default
	}

	if errorCode != "" {
		redirectURL := fmt.Sprintf("%s/login?error=%s", redirectURI, errorCode)
		if errorDesc != "" {
			redirectURL += fmt.Sprintf("&error_description=%s", errorDesc)
		}
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, redirectURI+"/auth/callback?token="+token)
	}
}
