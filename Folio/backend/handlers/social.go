package handlers

import (
	"context"
	"folio/api/auth"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type SocialHandler struct {
	DB *pgxpool.Pool
}

// FollowUser creates a follow relationship
func (h *SocialHandler) FollowUser(c echo.Context) error {
	followerID := auth.GetUserID(c)
	if followerID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	username := c.Param("username")
	if username == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "username is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Get user ID from username
	var followingID string
	err := h.DB.QueryRow(ctx, "SELECT id FROM users WHERE username = $1", username).Scan(&followingID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "user not found",
		})
	}

	// Can't follow yourself
	if followerID == followingID {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "cannot follow yourself",
		})
	}

	// Create follow relationship
	query := `
		INSERT INTO followers (follower_id, following_id, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (follower_id, following_id) DO NOTHING
		RETURNING id
	`

	var followID string
	err = h.DB.QueryRow(ctx, query, followerID, followingID).Scan(&followID)
	if err != nil {
		// Check if already following
		var exists bool
		h.DB.QueryRow(ctx, 
			"SELECT EXISTS(SELECT 1 FROM followers WHERE follower_id = $1 AND following_id = $2)",
			followerID, followingID,
		).Scan(&exists)
		
		if exists {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": "already following",
				"following": true,
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create follow relationship",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":   "successfully followed user",
		"following": true,
	})
}

// UnfollowUser removes a follow relationship
func (h *SocialHandler) UnfollowUser(c echo.Context) error {
	followerID := auth.GetUserID(c)
	if followerID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	username := c.Param("username")
	if username == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "username is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Get user ID from username
	var followingID string
	err := h.DB.QueryRow(ctx, "SELECT id FROM users WHERE username = $1", username).Scan(&followingID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "user not found",
		})
	}

	// Delete follow relationship
	query := `DELETE FROM followers WHERE follower_id = $1 AND following_id = $2`
	_, err = h.DB.Exec(ctx, query, followerID, followingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to unfollow user",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "successfully unfollowed user",
		"following": false,
	})
}

