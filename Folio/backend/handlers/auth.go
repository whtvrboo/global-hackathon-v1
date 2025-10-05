package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"folio/api/auth"
	"net/http"
	"os"
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
		SELECT id, google_id, email, name, username, picture, bio, 
		       favorite_book_ids, banner_url, reading_goal, reading_goal_year,
		       created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user struct {
		ID              string     `json:"id"`
		GoogleID        string     `json:"google_id"`
		Email           string     `json:"email"`
		Name            string     `json:"name"`
		Username        string     `json:"username"`
		Picture         *string    `json:"picture"`
		Bio             *string    `json:"bio"`
		FavoriteBookIDs []string   `json:"favorite_book_ids"`
		BannerURL       *string    `json:"banner_url"`
		ReadingGoal     int        `json:"reading_goal"`
		ReadingGoalYear int        `json:"reading_goal_year"`
		CreatedAt       time.Time  `json:"created_at"`
		UpdatedAt       time.Time  `json:"updated_at"`
	}

	err := h.DB.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name,
		&user.Username, &user.Picture, &user.Bio,
		&user.FavoriteBookIDs, &user.BannerURL, &user.ReadingGoal, &user.ReadingGoalYear,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "user not found",
		})
	}

	return c.JSON(http.StatusOK, user)
}

type UpdateProfileRequest struct {
	Bio             *string   `json:"bio"`
	BannerURL       *string   `json:"banner_url"`
	FavoriteBookIDs *[]string `json:"favorite_book_ids"`
	ReadingGoal     *int      `json:"reading_goal"`
}

// UpdateProfile updates the current user's profile
func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Build dynamic update query
	query := "UPDATE users SET updated_at = NOW()"
	args := []interface{}{userID}
	argIdx := 2

	if req.Bio != nil {
		query += fmt.Sprintf(", bio = $%d", argIdx)
		args = append(args, *req.Bio)
		argIdx++
	}
	if req.BannerURL != nil {
		query += fmt.Sprintf(", banner_url = $%d", argIdx)
		args = append(args, *req.BannerURL)
		argIdx++
	}
	if req.FavoriteBookIDs != nil {
		query += fmt.Sprintf(", favorite_book_ids = $%d", argIdx)
		args = append(args, *req.FavoriteBookIDs)
		argIdx++
	}
	if req.ReadingGoal != nil {
		query += fmt.Sprintf(", reading_goal = $%d, reading_goal_year = EXTRACT(YEAR FROM NOW())", argIdx)
		args = append(args, *req.ReadingGoal)
		argIdx++
	}

	query += " WHERE id = $1 RETURNING id, bio, banner_url, favorite_book_ids, reading_goal, reading_goal_year, updated_at"

	var updated struct {
		ID              string
		Bio             *string
		BannerURL       *string
		FavoriteBookIDs []string
		ReadingGoal     int
		ReadingGoalYear int
		UpdatedAt       time.Time
	}

	err := h.DB.QueryRow(ctx, query, args...).Scan(
		&updated.ID, &updated.Bio, &updated.BannerURL, &updated.FavoriteBookIDs,
		&updated.ReadingGoal, &updated.ReadingGoalYear, &updated.UpdatedAt,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("failed to update profile: %v", err),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":                 updated.ID,
		"bio":                updated.Bio,
		"banner_url":         updated.BannerURL,
		"favorite_book_ids":  updated.FavoriteBookIDs,
		"reading_goal":       updated.ReadingGoal,
		"reading_goal_year":  updated.ReadingGoalYear,
		"updated_at":         updated.UpdatedAt,
	})
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
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

