package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"folio/api/auth"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	DB *pgxpool.Pool
}

// GoogleLogin redirects to Google OAuth consent screen
func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	// Generate random state for CSRF protection
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Store state in session/cookie (simplified for now)
	c.SetCookie(&http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300, // 5 minutes
	})

	url := auth.GetGoogleOAuthURL(state)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the OAuth callback from Google
func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	// Verify state for CSRF protection
	stateCookie, err := c.Cookie("oauth_state")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing state cookie",
		})
	}

	state := c.QueryParam("state")
	if state != stateCookie.Value {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid state parameter",
		})
	}

	// Get authorization code
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing authorization code",
		})
	}

	// Exchange code for user info
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userInfo, err := auth.ExchangeGoogleCode(ctx, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("failed to exchange code: %v", err),
		})
	}

	// Upsert user in database
	userID, err := h.upsertUser(ctx, userInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("failed to upsert user: %v", err),
		})
	}

	// Generate JWT
	token, err := auth.GenerateJWT(userID, userInfo.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate token",
		})
	}

	// Redirect to frontend with token
	frontendURL := getEnv("FRONTEND_URL", "http://localhost")
	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s", frontendURL, token)
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// upsertUser inserts or updates user in database
func (h *AuthHandler) upsertUser(ctx context.Context, userInfo *auth.GoogleUserInfo) (string, error) {
	// Generate username from email (before @)
	username := userInfo.Email[:len(userInfo.Email)-len("@gmail.com")]
	if len(username) > 50 {
		username = username[:50]
	}

	query := `
		INSERT INTO users (google_id, email, name, username, picture, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (google_id) 
		DO UPDATE SET 
			email = EXCLUDED.email,
			name = EXCLUDED.name,
			picture = EXCLUDED.picture,
			updated_at = NOW()
		RETURNING id
	`

	var userID string
	err := h.DB.QueryRow(ctx, query,
		userInfo.ID,
		userInfo.Email,
		userInfo.Name,
		username,
		userInfo.Picture,
	).Scan(&userID)

	return userID, err
}

// GetMe returns the current user's profile
func (h *AuthHandler) GetMe(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, google_id, email, name, username, picture, bio, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user struct {
		ID        string     `json:"id"`
		GoogleID  string     `json:"google_id"`
		Email     string     `json:"email"`
		Name      string     `json:"name"`
		Username  string     `json:"username"`
		Picture   *string    `json:"picture"`
		Bio       *string    `json:"bio"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
	}

	err := h.DB.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name,
		&user.Username, &user.Picture, &user.Bio,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "user not found",
		})
	}

	return c.JSON(http.StatusOK, user)
}

// ConvertGuestToUser converts a guest user to a full user after OAuth
func (h *AuthHandler) ConvertGuestToUser(c echo.Context) error {
	// Get guest session ID from query parameter
	guestSessionID := c.QueryParam("guest_session_id")
	if guestSessionID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "guest_session_id is required",
		})
	}

	// Get authorization code
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing authorization code",
		})
	}

	// Exchange code for user info
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userInfo, err := auth.ExchangeGoogleCode(ctx, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("failed to exchange code: %v", err),
		})
	}

	// Find and update the guest user
	query := `
		UPDATE users 
		SET google_id = $1, email = $2, name = $3, picture = $4, 
		    is_guest = false, guest_session_id = NULL, converted_at = NOW(), updated_at = NOW()
		WHERE guest_session_id = $5 AND is_guest = true
		RETURNING id, username
	`

	var userID, username string
	err = h.DB.QueryRow(ctx, query,
		userInfo.ID,
		userInfo.Email,
		userInfo.Name,
		userInfo.Picture,
		guestSessionID,
	).Scan(&userID, &username)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "guest user not found or already converted",
		})
	}

	// Generate new JWT for the converted user
	token, err := auth.GenerateJWT(userID, userInfo.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate token",
		})
	}

	// Redirect to frontend with token
	frontendURL := getEnv("FRONTEND_URL", "http://localhost")
	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s&converted=true", frontendURL, token)
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func getEnv(key, defaultValue string) string {
	return defaultValue
}

