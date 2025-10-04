package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	jwtSecret         []byte
)

// InitOAuth initializes OAuth configuration
func InitOAuth() {
	jwtSecret = []byte(getEnv("JWT_SECRET", "dev-secret-change-in-production"))
	
	googleOauthConfig = &oauth2.Config{
		ClientID:     getEnv("GOOGLE_CLIENT_ID", "your-client-id"),
		ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", "your-client-secret"),
		RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/google/callback"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// GetGoogleOAuthURL returns the OAuth URL for Google login
func GetGoogleOAuthURL(state string) string {
	return googleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeGoogleCode exchanges OAuth code for token and fetches user info
func ExchangeGoogleCode(ctx context.Context, code string) (*GoogleUserInfo, error) {
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Fetch user info
	client := googleOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s", string(body))
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

// GoogleUserInfo represents user data from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	IsGuest bool   `json:"is_guest"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a JWT token for a user
func GenerateJWT(userID, email string) (string, error) {
	return GenerateJWTWithGuest(userID, email, false)
}

// GenerateJWTWithGuest creates a JWT token for a user (guest or regular)
func GenerateJWTWithGuest(userID, email string, isGuest bool) (string, error) {
	claims := JWTClaims{
		UserID:  userID,
		Email:   email,
		IsGuest: isGuest,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates and parses a JWT token
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

