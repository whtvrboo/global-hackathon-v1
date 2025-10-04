package handlers

import (
	"context"
	"folio/api/auth"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type LogHandler struct {
	DB *pgxpool.Pool
}

type CreateLogRequest struct {
	BookID     string  `json:"book_id"`
	Status     string  `json:"status"`
	Rating     *int    `json:"rating"`
	Review     *string `json:"review"`
	Notes      *string `json:"notes"`
	StartDate  *string `json:"start_date"`
	FinishDate *string `json:"finish_date"`
	IsPublic   *bool   `json:"is_public"`
}

// CreateLog creates a new reading log entry
func (h *LogHandler) CreateLog(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	var req CreateLogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	// Validate status
	validStatuses := map[string]bool{
		"want_to_read": true,
		"reading":      true,
		"read":         true,
		"dnf":          true,
	}
	if !validStatuses[req.Status] {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid status. Must be: want_to_read, reading, read, or dnf",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO logs (user_id, book_id, status, rating, review, notes, start_date, finish_date, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	var logID string
	var createdAt, updatedAt time.Time

	err := h.DB.QueryRow(ctx, query,
		userID, req.BookID, req.Status, req.Rating, req.Review,
		req.Notes, req.StartDate, req.FinishDate, isPublic,
	).Scan(&logID, &createdAt, &updatedAt)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create log",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":          logID,
		"user_id":     userID,
		"book_id":     req.BookID,
		"status":      req.Status,
		"rating":      req.Rating,
		"review":      req.Review,
		"notes":       req.Notes,
		"start_date":  req.StartDate,
		"finish_date": req.FinishDate,
		"is_public":   isPublic,
		"created_at":  createdAt,
		"updated_at":  updatedAt,
	})
}

// GetUserLogs fetches all logs for a given username
func (h *LogHandler) GetUserLogs(c echo.Context) error {
	username := c.Param("username")
	if username == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "username is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	query := `
		SELECT l.id, l.user_id, l.book_id, l.status, l.rating, l.review,
		       l.notes, l.start_date, l.finish_date, l.is_public, l.created_at,
		       b.title, b.authors, b.cover_url
		FROM logs l
		JOIN users u ON l.user_id = u.id
		JOIN books b ON l.book_id = b.id
		WHERE u.username = $1 AND l.is_public = true
		ORDER BY l.created_at DESC
		LIMIT 50
	`

	rows, err := h.DB.Query(ctx, query, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch logs",
		})
	}
	defer rows.Close()

	logs := []map[string]interface{}{}
	for rows.Next() {
		var log struct {
			ID         string
			UserID     string
			BookID     string
			Status     string
			Rating     *int
			Review     *string
			Notes      *string
			StartDate  *string
			FinishDate *string
			IsPublic   bool
			CreatedAt  time.Time
			BookTitle  string
			Authors    []string
			CoverURL   *string
		}

		err := rows.Scan(
			&log.ID, &log.UserID, &log.BookID, &log.Status, &log.Rating,
			&log.Review, &log.Notes, &log.StartDate, &log.FinishDate,
			&log.IsPublic, &log.CreatedAt, &log.BookTitle, &log.Authors, &log.CoverURL,
		)
		if err != nil {
			continue
		}

		logs = append(logs, map[string]interface{}{
			"id":          log.ID,
			"user_id":     log.UserID,
			"book_id":     log.BookID,
			"status":      log.Status,
			"rating":      log.Rating,
			"review":      log.Review,
			"notes":       log.Notes,
			"start_date":  log.StartDate,
			"finish_date": log.FinishDate,
			"is_public":   log.IsPublic,
			"created_at":  log.CreatedAt,
			"book": map[string]interface{}{
				"title":     log.BookTitle,
				"authors":   log.Authors,
				"cover_url": log.CoverURL,
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"logs":  logs,
		"count": len(logs),
	})
}

// GetFeed gets the reading feed for the current user
func (h *LogHandler) GetFeed(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	query := `
		SELECT l.id, l.user_id, l.book_id, l.status, l.rating, l.review,
		       l.created_at, u.username, u.name, u.picture,
		       b.title, b.authors, b.cover_url
		FROM logs l
		JOIN users u ON l.user_id = u.id
		JOIN books b ON l.book_id = b.id
		WHERE l.user_id IN (
			SELECT following_id FROM followers WHERE follower_id = $1
		)
		AND l.is_public = true
		ORDER BY l.created_at DESC
		LIMIT 50
	`

	rows, err := h.DB.Query(ctx, query, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch feed",
		})
	}
	defer rows.Close()

	feed := []map[string]interface{}{}
	for rows.Next() {
		var item struct {
			LogID      string
			UserID     string
			BookID     string
			Status     string
			Rating     *int
			Review     *string
			CreatedAt  time.Time
			Username   string
			Name       string
			Picture    *string
			BookTitle  string
			Authors    []string
			CoverURL   *string
		}

		err := rows.Scan(
			&item.LogID, &item.UserID, &item.BookID, &item.Status,
			&item.Rating, &item.Review, &item.CreatedAt,
			&item.Username, &item.Name, &item.Picture,
			&item.BookTitle, &item.Authors, &item.CoverURL,
		)
		if err != nil {
			continue
		}

		feed = append(feed, map[string]interface{}{
			"id":         item.LogID,
			"status":     item.Status,
			"rating":     item.Rating,
			"review":     item.Review,
			"created_at": item.CreatedAt,
			"user": map[string]interface{}{
				"id":       item.UserID,
				"username": item.Username,
				"name":     item.Name,
				"picture":  item.Picture,
			},
			"book": map[string]interface{}{
				"id":        item.BookID,
				"title":     item.BookTitle,
				"authors":   item.Authors,
				"cover_url": item.CoverURL,
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"feed":  feed,
		"count": len(feed),
	})
}

