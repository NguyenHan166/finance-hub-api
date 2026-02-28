package utils

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"finance-hub-api/internal/models"
	"fmt"
	"io"
	"math/big"
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

	fmt.Printf("Exchanging code with redirect_uri: %s\n", g.RedirectURI)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error doing request: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Token exchange failed (status %d): %s\n", resp.StatusCode, string(body))
		return "", fmt.Errorf("failed to exchange code: %s", string(body))
	}

	var tokenResponse struct {
		IDToken      string `json:"id_token"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		fmt.Printf("Error unmarshaling response: %v\n", err)
		return "", err
	}

	fmt.Printf("Token exchange successful, got ID token\n")
	return tokenResponse.IDToken, nil
}

// VerifyIDToken verifies Google ID token and extracts user info
func (g *GoogleOAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*models.GoogleUserInfo, error) {
	fmt.Println("Starting ID token verification...")
	
	// Parse token without verification first to get the kid
	_, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		fmt.Printf("Error parsing unverified token: %v\n", err)
		return nil, ErrInvalidGoogleToken
	}

	// Verify token with Google's public keys
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(idToken, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			fmt.Printf("Unexpected signing method: %v\n", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get Google's public keys
		kid, ok := token.Header["kid"].(string)
		if !ok {
			fmt.Println("kid not found in token header")
			return nil, errors.New("kid not found in token header")
		}

		fmt.Printf("Fetching public key for kid: %s\n", kid)
		return getGooglePublicKey(ctx, kid)
	})

	if err != nil {
		fmt.Printf("Error verifying token: %v\n", err)
		return nil, ErrInvalidGoogleToken
	}

	fmt.Println("Token signature verified successfully")

	// Verify issuer
	iss, ok := claims["iss"].(string)
	if !ok || (iss != "https://accounts.google.com" && iss != "accounts.google.com") {
		fmt.Printf("Invalid issuer: %v\n", iss)
		return nil, ErrInvalidGoogleToken
	}

	// Verify audience (client ID)
	aud, ok := claims["aud"].(string)
	if !ok || aud != g.ClientID {
		fmt.Printf("Invalid audience: %v (expected: %s)\n", aud, g.ClientID)
		return nil, ErrInvalidGoogleToken
	}

	// Verify expiration
	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		fmt.Printf("Token expired or invalid exp: %v\n", exp)
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

	fmt.Printf("Token verified, user: %s (%s)\n", userInfo.Name, userInfo.Email)
	return userInfo, nil
}

// getGooglePublicKey fetches Google's public keys for JWT verification
func getGooglePublicKey(ctx context.Context, kid string) (interface{}, error) {
	// Fetch Google's public keys
	keysURL := "https://www.googleapis.com/oauth2/v3/certs"
	
	req, err := http.NewRequestWithContext(ctx, "GET", keysURL, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching public keys: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
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
		fmt.Printf("Error unmarshaling keys: %v\n", err)
		return nil, err
	}

	// Find the key with matching kid
	for _, key := range keys.Keys {
		if key.Kid == kid {
			fmt.Printf("Found matching key: kid=%s, alg=%s\n", key.Kid, key.Alg)
			
			// Decode N (modulus) from base64url
			nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
			if err != nil {
				fmt.Printf("Error decoding modulus: %v\n", err)
				return nil, err
			}

			// Decode E (exponent) from base64url
			eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
			if err != nil {
				fmt.Printf("Error decoding exponent: %v\n", err)
				return nil, err
			}

			// Convert E bytes to int
			var eInt int
			if len(eBytes) == 3 {
				eInt = int(binary.BigEndian.Uint32(append([]byte{0}, eBytes...)))
			} else if len(eBytes) == 4 {
				eInt = int(binary.BigEndian.Uint32(eBytes))
			} else {
				// Fallback for other lengths
				eInt = 65537 // Common RSA exponent
			}

			// Create RSA public key
			publicKey := &rsa.PublicKey{
				N: new(big.Int).SetBytes(nBytes),
				E: eInt,
			}

			fmt.Println("Successfully created RSA public key")
			return publicKey, nil
		}
	}

	fmt.Printf("Public key not found for kid: %s\n", kid)
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
