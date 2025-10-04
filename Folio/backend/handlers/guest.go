package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"folio/api/auth"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type GuestHandler struct {
	DB *pgxpool.Pool
}

// CreateGuestUser creates a new guest user
func (h *GuestHandler) CreateGuestUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Generate unique guest session ID
	guestSessionID, err := generateGuestSessionID()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate guest session",
		})
	}

	// Generate a temporary username
	username := fmt.Sprintf("guest_%s", guestSessionID[:8])

	// Create guest user
	query := `
		INSERT INTO users (guest_session_id, username, name, is_guest, created_at, updated_at)
		VALUES ($1, $2, $3, true, NOW(), NOW())
		RETURNING id, username, name, is_guest, created_at
	`

	var user struct {
		ID        string    `json:"id"`
		Username  string    `json:"username"`
		Name      string    `json:"name"`
		IsGuest   bool      `json:"is_guest"`
		CreatedAt time.Time `json:"created_at"`
	}

	err = h.DB.QueryRow(ctx, query, guestSessionID, username, "Guest User").Scan(
		&user.ID, &user.Username, &user.Name, &user.IsGuest, &user.CreatedAt,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create guest user",
		})
	}

	// Generate JWT for guest user
	token, err := auth.GenerateJWTWithGuest(user.ID, "", true)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate token",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"user":  user,
		"token": token,
		"guest_session_id": guestSessionID,
	})
}

// ConvertGuestToUser converts a guest user to a full user
func (h *GuestHandler) ConvertGuestToUser(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Get the guest user
	var guestUser struct {
		ID              string
		GuestSessionID  string
		Username        string
		IsGuest         bool
	}

	query := `SELECT id, guest_session_id, username, is_guest FROM users WHERE id = $1 AND is_guest = true`
	err := h.DB.QueryRow(ctx, query, userID).Scan(
		&guestUser.ID, &guestUser.GuestSessionID, &guestUser.Username, &guestUser.IsGuest,
	)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "guest user not found",
		})
	}

	// This endpoint would be called after OAuth completion
	// The OAuth callback would need to be modified to handle guest conversion
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "guest user ready for conversion",
		"guest_session_id": guestUser.GuestSessionID,
	})
}

// GetGuestUser gets guest user info
func (h *GuestHandler) GetGuestUser(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, username, name, is_guest, guest_session_id, created_at, converted_at
		FROM users
		WHERE id = $1
	`

	var user struct {
		ID              string     `json:"id"`
		Username        string     `json:"username"`
		Name            string     `json:"name"`
		IsGuest         bool       `json:"is_guest"`
		GuestSessionID  *string    `json:"guest_session_id"`
		CreatedAt       time.Time  `json:"created_at"`
		ConvertedAt     *time.Time `json:"converted_at"`
	}

	err := h.DB.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.Username, &user.Name, &user.IsGuest,
		&user.GuestSessionID, &user.CreatedAt, &user.ConvertedAt,
	)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "user not found",
		})
	}

	return c.JSON(http.StatusOK, user)
}

// generateGuestSessionID creates a unique session ID for guest users
func generateGuestSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

