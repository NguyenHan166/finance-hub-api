package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"finance-hub-api/internal/models"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidGoogleToken = errors.New("invalid Google ID token")
	ErrTokenExpired       = errors.New("token expired")
)

// GoogleOAuthClient handles Google OAuth operations
type GoogleOAuthClient struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// NewGoogleOAuthClient creates a new Google OAuth client
func NewGoogleOAuthClient(clientID, clientSecret, redirectURI string) *GoogleOAuthClient {
	return &GoogleOAuthClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}
}

// GenerateOAuthURL generates Google OAuth authorization URL
func (g *GoogleOAuthClient) GenerateOAuthURL(state string) string {
	baseURL := "https://accounts.google.com/o/oauth2/v2/auth"
	params := url.Values{}
	params.Add("client_id", g.ClientID)
	params.Add("redirect_uri", g.RedirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "openid profile email")
	params.Add("state", state)
	params.Add("access_type", "offline")
	params.Add("prompt", "consent")

	return baseURL + "?" + params.Encode()
}

// ExchangeCodeForToken exchanges authorization code for tokens
func (g *GoogleOAuthClient) ExchangeCodeForToken(ctx context.Context, code string) (string, error) {
	tokenURL := "https://oauth2.googleapis.com/token"
	
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", g.ClientID)
	data.Set("client_secret", g.ClientSecret)
	data.Set("redirect_uri", g.RedirectURI)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to exchange code: %s", string(body))
	}

	var tokenResponse struct {
		IDToken      string `json:"id_token"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.IDToken, nil
}

// VerifyIDToken verifies Google ID token and extracts user info
func (g *GoogleOAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*models.GoogleUserInfo, error) {
	// Parse token without verification first to get the kid
	_, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return nil, ErrInvalidGoogleToken
	}

	// Verify token with Google's public keys
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(idToken, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get Google's public keys
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid not found in token header")
		}

		return getGooglePublicKey(ctx, kid)
	})

	if err != nil {
		return nil, ErrInvalidGoogleToken
	}

	// Verify issuer
	iss, ok := claims["iss"].(string)
	if !ok || (iss != "https://accounts.google.com" && iss != "accounts.google.com") {
		return nil, ErrInvalidGoogleToken
	}

	// Verify audience (client ID)
	aud, ok := claims["aud"].(string)
	if !ok || aud != g.ClientID {
		return nil, ErrInvalidGoogleToken
	}

	// Verify expiration
	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return nil, ErrTokenExpired
	}

	// Extract user info
	userInfo := &models.GoogleUserInfo{
		Sub:           getStringClaim(claims, "sub"),
		Email:         getStringClaim(claims, "email"),
		EmailVerified: getBoolClaim(claims, "email_verified"),
		Name:          getStringClaim(claims, "name"),
		GivenName:     getStringClaim(claims, "given_name"),
		FamilyName:    getStringClaim(claims, "family_name"),
		Picture:       getStringClaim(claims, "picture"),
		Locale:        getStringClaim(claims, "locale"),
	}

	return userInfo, nil
}

// getGooglePublicKey fetches Google's public keys for JWT verification
func getGooglePublicKey(ctx context.Context, kid string) (interface{}, error) {
	// Fetch Google's public keys
	keysURL := "https://www.googleapis.com/oauth2/v3/certs"
	
	req, err := http.NewRequestWithContext(ctx, "GET", keysURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var keys struct {
		Keys []struct {
			Kid string `json:"kid"`
			N   string `json:"n"`
			E   string `json:"e"`
			Kty string `json:"kty"`
			Alg string `json:"alg"`
			Use string `json:"use"`
		} `json:"keys"`
	}

	if err := json.Unmarshal(body, &keys); err != nil {
		return nil, err
	}

	// Find the key with matching kid
	for _, key := range keys.Keys {
		if key.Kid == kid {
			// Convert JWK to public key
			// For production, use a proper JWK library
			// This is a simplified version
			return []byte(key.N), nil
		}
	}

	return nil, errors.New("public key not found")
}

// Helper functions

func getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return ""
}

func getBoolClaim(claims jwt.MapClaims, key string) bool {
	if val, ok := claims[key].(bool); ok {
		return val
	}
	return false
}

// GenerateRandomState generates a random state string for CSRF protection
func GenerateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
