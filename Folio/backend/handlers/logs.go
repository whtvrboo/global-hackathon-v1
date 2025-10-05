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
	BookID      string  `json:"book_id"`
	Status      string  `json:"status"`
	Rating      *int    `json:"rating"`
	Review      *string `json:"review"`
	Notes       *string `json:"notes"`
	StartDate   *string `json:"start_date"`
	FinishDate  *string `json:"finish_date"`
	IsPublic    *bool   `json:"is_public"`
	SpoilerFlag *bool   `json:"spoiler_flag"`
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
		INSERT INTO logs (user_id, book_id, status, rating, review, notes, start_date, finish_date, is_public, spoiler_flag, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	spoilerFlag := false
	if req.SpoilerFlag != nil {
		spoilerFlag = *req.SpoilerFlag
	}

	var logID string
	var createdAt, updatedAt time.Time

	err := h.DB.QueryRow(ctx, query,
		userID, req.BookID, req.Status, req.Rating, req.Review,
		req.Notes, req.StartDate, req.FinishDate, isPublic, spoilerFlag,
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

	// Check if the requesting user is the same as the profile owner
	currentUserID := auth.GetUserID(c)
	var isOwnProfile bool
	if currentUserID != "" {
		var profileUserID string
		err := h.DB.QueryRow(ctx, "SELECT id FROM users WHERE username = $1", username).Scan(&profileUserID)
		if err == nil {
			isOwnProfile = (currentUserID == profileUserID)
		}
	}
	

	var query string
	var args []interface{}

	if isOwnProfile {
		// Show all logs (public and private) for own profile
		query = `
			SELECT l.id, l.user_id, l.book_id, l.status, l.rating, l.review,
			       l.notes, l.start_date, l.finish_date, l.is_public, l.created_at,
			       l.likes_count, l.comments_count,
			       b.title, b.authors, b.cover_url,
			       EXISTS(SELECT 1 FROM log_likes WHERE log_id = l.id AND user_id = $2) as is_liked
			FROM logs l
			JOIN users u ON l.user_id = u.id
			JOIN books b ON l.book_id = b.id
			WHERE u.username = $1
			ORDER BY l.created_at DESC
			LIMIT 50
		`
		args = []interface{}{username, currentUserID}
	} else if currentUserID != "" {
		// Show only public logs for authenticated users viewing other profiles
		query = `
			SELECT l.id, l.user_id, l.book_id, l.status, l.rating, l.review,
			       l.notes, l.start_date, l.finish_date, l.is_public, l.created_at,
			       l.likes_count, l.comments_count,
			       b.title, b.authors, b.cover_url,
			       EXISTS(SELECT 1 FROM log_likes WHERE log_id = l.id AND user_id = $2) as is_liked
			FROM logs l
			JOIN users u ON l.user_id = u.id
			JOIN books b ON l.book_id = b.id
			WHERE u.username = $1 AND l.is_public = true
			ORDER BY l.created_at DESC
			LIMIT 50
		`
		args = []interface{}{username, currentUserID}
	} else {
		// Show only public logs for unauthenticated users
		query = `
			SELECT l.id, l.user_id, l.book_id, l.status, l.rating, l.review,
			       l.notes, l.start_date, l.finish_date, l.is_public, l.created_at,
			       l.likes_count, l.comments_count,
			       b.title, b.authors, b.cover_url,
			       false as is_liked
			FROM logs l
			JOIN users u ON l.user_id = u.id
			JOIN books b ON l.book_id = b.id
			WHERE u.username = $1 AND l.is_public = true
			ORDER BY l.created_at DESC
			LIMIT 50
		`
		args = []interface{}{username}
	}

	rows, err := h.DB.Query(ctx, query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch logs",
		})
	}
	defer rows.Close()

	logs := []map[string]interface{}{}
	for rows.Next() {
		var log struct {
			ID            string
			UserID        string
			BookID        string
			Status        string
			Rating        *int
			Review        *string
			Notes         *string
			StartDate     *string
			FinishDate    *string
			IsPublic      bool
			CreatedAt     time.Time
			LikesCount    int
			CommentsCount int
			BookTitle     string
			Authors       []string
			CoverURL      *string
			IsLiked       bool
		}

		err := rows.Scan(
			&log.ID, &log.UserID, &log.BookID, &log.Status, &log.Rating,
			&log.Review, &log.Notes, &log.StartDate, &log.FinishDate,
			&log.IsPublic, &log.CreatedAt, &log.LikesCount, &log.CommentsCount,
			&log.BookTitle, &log.Authors, &log.CoverURL, &log.IsLiked,
		)
		if err != nil {
			continue
		}

		logs = append(logs, map[string]interface{}{
			"id":             log.ID,
			"user_id":        log.UserID,
			"book_id":        log.BookID,
			"status":         log.Status,
			"rating":         log.Rating,
			"review":         log.Review,
			"notes":          log.Notes,
			"start_date":     log.StartDate,
			"finish_date":    log.FinishDate,
			"is_public":      log.IsPublic,
			"created_at":     log.CreatedAt,
			"likes_count":    log.LikesCount,
			"comments_count": log.CommentsCount,
			"is_liked":       log.IsLiked,
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

// GetFeed gets the list-based activity feed for the current user
func (h *LogHandler) GetFeed(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// First, check if user has any followers
	var followerCount int
	err := h.DB.QueryRow(ctx, "SELECT COUNT(*) FROM followers WHERE follower_id = $1", userID).Scan(&followerCount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to check followers",
		})
	}

	var query string
	var args []interface{}

	if followerCount > 0 {
		// Show list activity from followed users
		query = `
			SELECT 
				'list_created' as event_type,
				l.id as entity_id,
				l.user_id,
				l.name as title,
				l.description,
				l.items_count,
				l.likes_count,
				l.comments_count,
				l.header_image_url,
				l.theme_color,
				l.created_at,
				u.username,
				u.name as user_name,
				u.picture,
				NULL::text[] as book_ids,
				NULL::text[] as book_titles,
				NULL::text[] as book_covers,
				EXISTS(SELECT 1 FROM list_likes WHERE list_id = l.id AND user_id = $1) as is_liked
			FROM lists l
			JOIN users u ON l.user_id = u.id
			WHERE l.user_id IN (
				SELECT following_id FROM followers WHERE follower_id = $1
			)
			AND l.is_public = true
			ORDER BY l.created_at DESC
			LIMIT 50
		`
		args = []interface{}{userID}
	} else {
		// If no followers, show popular public lists to bootstrap engagement
		query = `
			SELECT 
				'list_created' as event_type,
				l.id as entity_id,
				l.user_id,
				l.name as title,
				l.description,
				l.items_count,
				l.likes_count,
				l.comments_count,
				l.header_image_url,
				l.theme_color,
				l.created_at,
				u.username,
				u.name as user_name,
				u.picture,
				NULL::text[] as book_ids,
				NULL::text[] as book_titles,
				NULL::text[] as book_covers,
				EXISTS(SELECT 1 FROM list_likes WHERE list_id = l.id AND user_id = $1) as is_liked
			FROM lists l
			JOIN users u ON l.user_id = u.id
			WHERE l.is_public = true
			AND l.items_count > 0
			ORDER BY l.likes_count DESC, l.created_at DESC
			LIMIT 50
		`
		args = []interface{}{userID}
	}

	rows, err := h.DB.Query(ctx, query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch feed",
		})
	}
	defer rows.Close()

	feed := []map[string]interface{}{}
	for rows.Next() {
		var item struct {
			EventType      string
			EntityID       string
			UserID         string
			Title          string
			Description    *string
			ItemsCount     int
			LikesCount     int
			CommentsCount  int
			HeaderImageURL *string
			ThemeColor     *string
			CreatedAt      time.Time
			Username       string
			UserName       string
			Picture        *string
			BookIDs        []string
			BookTitles     []string
			BookCovers     []string
			IsLiked        bool
		}

		err := rows.Scan(
			&item.EventType, &item.EntityID, &item.UserID,
			&item.Title, &item.Description, &item.ItemsCount,
			&item.LikesCount, &item.CommentsCount,
			&item.HeaderImageURL, &item.ThemeColor,
			&item.CreatedAt, &item.Username, &item.UserName, &item.Picture,
			&item.BookIDs, &item.BookTitles, &item.BookCovers,
			&item.IsLiked,
		)
		if err != nil {
			continue
		}

		// Fetch preview books for this list (top 3 covers)
		previewBooks := []map[string]interface{}{}
		if item.ItemsCount > 0 {
			bookQuery := `
				SELECT b.id, b.title, b.cover_url
				FROM list_items li
				JOIN books b ON li.book_id = b.id
				WHERE li.list_id = $1
				ORDER BY li.item_order
				LIMIT 3
			`
			bookRows, err := h.DB.Query(ctx, bookQuery, item.EntityID)
			if err == nil {
				defer bookRows.Close()
				for bookRows.Next() {
					var bookID, bookTitle string
					var bookCover *string
					if err := bookRows.Scan(&bookID, &bookTitle, &bookCover); err == nil {
						previewBooks = append(previewBooks, map[string]interface{}{
							"id":        bookID,
							"title":     bookTitle,
							"cover_url": bookCover,
						})
					}
				}
			}
		}

		feed = append(feed, map[string]interface{}{
			"event_type":       item.EventType,
			"id":               item.EntityID,
			"title":            item.Title,
			"description":      item.Description,
			"items_count":      item.ItemsCount,
			"likes_count":      item.LikesCount,
			"comments_count":   item.CommentsCount,
			"header_image_url": item.HeaderImageURL,
			"theme_color":      item.ThemeColor,
			"created_at":       item.CreatedAt,
			"is_liked":         item.IsLiked,
			"preview_books":    previewBooks,
			"user": map[string]interface{}{
				"id":       item.UserID,
				"username": item.Username,
				"name":     item.UserName,
				"picture":  item.Picture,
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"feed":  feed,
		"count": len(feed),
	})
}

// GetSingleLog retrieves a specific log with full details for dedicated review view
func (h *LogHandler) GetSingleLog(c echo.Context) error {
	logID := c.Param("id")
	if logID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "log_id is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	currentUserID := auth.GetUserID(c)

	query := `
		SELECT l.id, l.user_id, l.book_id, l.status, l.rating, l.review,
		       l.notes, l.start_date, l.finish_date, l.is_public, l.spoiler_flag, l.created_at,
		       l.likes_count, l.comments_count,
		       u.username, u.name, u.picture,
		       b.title, b.authors, b.cover_url, b.description, b.pages, b.published_date,
		       EXISTS(SELECT 1 FROM log_likes WHERE log_id = l.id AND user_id = $2) as is_liked
		FROM logs l
		JOIN users u ON l.user_id = u.id
		JOIN books b ON l.book_id = b.id
		WHERE l.id = $1
	`

	var log struct {
		ID            string
		UserID        string
		BookID        string
		Status        string
		Rating        *int
		Review        *string
		Notes         *string
		StartDate     *string
		FinishDate    *string
		IsPublic      bool
		SpoilerFlag   bool
		CreatedAt     time.Time
		LikesCount    int
		CommentsCount int
		Username      string
		Name          string
		Picture       *string
		BookTitle     string
		Authors       []string
		CoverURL      *string
		Description   *string
		Pages         *int
		PublishedDate *string
		IsLiked       bool
	}

	err := h.DB.QueryRow(ctx, query, logID, currentUserID).Scan(
		&log.ID, &log.UserID, &log.BookID, &log.Status, &log.Rating, &log.Review,
		&log.Notes, &log.StartDate, &log.FinishDate, &log.IsPublic, &log.SpoilerFlag, &log.CreatedAt,
		&log.LikesCount, &log.CommentsCount,
		&log.Username, &log.Name, &log.Picture,
		&log.BookTitle, &log.Authors, &log.CoverURL, &log.Description, &log.Pages, &log.PublishedDate,
		&log.IsLiked,
	)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "log not found",
		})
	}

	// Check visibility
	if !log.IsPublic && currentUserID != log.UserID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "you don't have permission to view this log",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":             log.ID,
		"status":         log.Status,
		"rating":         log.Rating,
		"review":         log.Review,
		"notes":          log.Notes,
		"start_date":     log.StartDate,
		"finish_date":    log.FinishDate,
		"is_public":      log.IsPublic,
		"spoiler_flag":   log.SpoilerFlag,
		"created_at":     log.CreatedAt,
		"likes_count":    log.LikesCount,
		"comments_count": log.CommentsCount,
		"is_liked":       log.IsLiked,
		"user": map[string]interface{}{
			"id":       log.UserID,
			"username": log.Username,
			"name":     log.Name,
			"picture":  log.Picture,
		},
		"book": map[string]interface{}{
			"id":             log.BookID,
			"title":          log.BookTitle,
			"authors":        log.Authors,
			"cover_url":      log.CoverURL,
			"description":    log.Description,
			"pages":          log.Pages,
			"published_date": log.PublishedDate,
		},
	})
}

