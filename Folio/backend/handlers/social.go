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

// ToggleLike toggles a like on a log (like if not liked, unlike if already liked)
func (h *SocialHandler) ToggleLike(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	logID := c.Param("id")
	if logID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "log_id is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Check if log exists
	var logExists bool
	err := h.DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM logs WHERE id = $1)", logID).Scan(&logExists)
	if err != nil || !logExists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "log not found",
		})
	}

	// Check if already liked
	var likeID string
	err = h.DB.QueryRow(ctx, 
		"SELECT id FROM log_likes WHERE user_id = $1 AND log_id = $2",
		userID, logID,
	).Scan(&likeID)

	if err != nil {
		// No existing like, create one
		query := `
			INSERT INTO log_likes (user_id, log_id, created_at)
			VALUES ($1, $2, NOW())
			RETURNING id
		`
		err = h.DB.QueryRow(ctx, query, userID, logID).Scan(&likeID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to like log",
			})
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "log liked",
			"liked":   true,
		})
	}

	// Like exists, remove it
	_, err = h.DB.Exec(ctx, "DELETE FROM log_likes WHERE id = $1", likeID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to unlike log",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "log unliked",
		"liked":   false,
	})
}

// GetLogComments retrieves all comments for a specific log
func (h *SocialHandler) GetLogComments(c echo.Context) error {
	logID := c.Param("id")
	if logID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "log_id is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	query := `
		SELECT lc.id, lc.user_id, lc.log_id, lc.content, lc.created_at, lc.updated_at,
		       u.username, u.name, u.picture
		FROM log_comments lc
		JOIN users u ON lc.user_id = u.id
		WHERE lc.log_id = $1
		ORDER BY lc.created_at ASC
	`

	rows, err := h.DB.Query(ctx, query, logID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch comments",
		})
	}
	defer rows.Close()

	comments := []map[string]interface{}{}
	for rows.Next() {
		var comment struct {
			ID        string
			UserID    string
			LogID     string
			Content   string
			CreatedAt time.Time
			UpdatedAt time.Time
			Username  string
			Name      string
			Picture   *string
		}

		err := rows.Scan(
			&comment.ID, &comment.UserID, &comment.LogID, &comment.Content,
			&comment.CreatedAt, &comment.UpdatedAt,
			&comment.Username, &comment.Name, &comment.Picture,
		)
		if err != nil {
			continue
		}

		comments = append(comments, map[string]interface{}{
			"id":         comment.ID,
			"log_id":     comment.LogID,
			"content":    comment.Content,
			"created_at": comment.CreatedAt,
			"updated_at": comment.UpdatedAt,
			"user": map[string]interface{}{
				"id":       comment.UserID,
				"username": comment.Username,
				"name":     comment.Name,
				"picture":  comment.Picture,
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"comments": comments,
		"count":    len(comments),
	})
}

type CreateCommentRequest struct {
	Content string `json:"content"`
}

// CreateComment adds a new comment to a log
func (h *SocialHandler) CreateComment(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	logID := c.Param("id")
	if logID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "log_id is required",
		})
	}

	var req CreateCommentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.Content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "content is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Check if log exists
	var logExists bool
	err := h.DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM logs WHERE id = $1)", logID).Scan(&logExists)
	if err != nil || !logExists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "log not found",
		})
	}

	query := `
		INSERT INTO log_comments (user_id, log_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	var commentID string
	var createdAt, updatedAt time.Time
	err = h.DB.QueryRow(ctx, query, userID, logID, req.Content).Scan(&commentID, &createdAt, &updatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create comment",
		})
	}

	// Get user info for response
	var username, name string
	var picture *string
	h.DB.QueryRow(ctx, "SELECT username, name, picture FROM users WHERE id = $1", userID).Scan(&username, &name, &picture)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":         commentID,
		"log_id":     logID,
		"content":    req.Content,
		"created_at": createdAt,
		"updated_at": updatedAt,
		"user": map[string]interface{}{
			"id":       userID,
			"username": username,
			"name":     name,
			"picture":  picture,
		},
	})
}

// DeleteComment allows a user to delete their own comment
func (h *SocialHandler) DeleteComment(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	commentID := c.Param("commentId")
	if commentID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "comment_id is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Check if comment exists and belongs to user
	var ownerID string
	err := h.DB.QueryRow(ctx, "SELECT user_id FROM log_comments WHERE id = $1", commentID).Scan(&ownerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "comment not found",
		})
	}

	if ownerID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "you can only delete your own comments",
		})
	}

	// Delete comment
	_, err = h.DB.Exec(ctx, "DELETE FROM log_comments WHERE id = $1", commentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to delete comment",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "comment deleted successfully",
	})
}

