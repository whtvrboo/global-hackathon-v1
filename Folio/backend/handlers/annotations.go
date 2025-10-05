package handlers

import (
	"context"
	"folio/api/auth"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type AnnotationHandler struct {
	DB *pgxpool.Pool
}

type CaptureAnnotationRequest struct {
	Content    string   `json:"content"`
	Type       string   `json:"type"` // 'highlight' or 'note'
	BookID     *string  `json:"book_id"`
	PageNumber *int     `json:"page_number"`
	Tags       []string `json:"tags"`
	Context    *string  `json:"context"`
}

type UpdateAnnotationRequest struct {
	BookID     *string  `json:"book_id"`
	Content    *string  `json:"content"`
	PageNumber *int     `json:"page_number"`
	Tags       []string `json:"tags"`
}

type RecentBook struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	CoverURL    *string `json:"cover_url"`
	Status      string  `json:"status"`
	LastUpdated string  `json:"last_updated"`
}

// CaptureAnnotation creates a new annotation with intelligent book association
func (h *AnnotationHandler) CaptureAnnotation(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	var req CaptureAnnotationRequest
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

	// Default type to 'note' if not specified
	if req.Type == "" {
		req.Type = "note"
	}

	if req.Type != "note" && req.Type != "highlight" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "type must be 'note' or 'highlight'",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	var bookID *string
	isAssociated := false

	// If book_id is provided explicitly, trust it
	if req.BookID != nil && *req.BookID != "" {
		bookID = req.BookID
		isAssociated = true
	} else {
		// Run auto-detection: find books with status='reading'
		query := `
			SELECT book_id 
			FROM logs 
			WHERE user_id = $1 AND status = 'reading'
		`
		rows, err := h.DB.Query(ctx, query, userID)
		if err == nil {
			defer rows.Close()
			
			var readingBooks []string
			for rows.Next() {
				var bid string
				if err := rows.Scan(&bid); err == nil {
					readingBooks = append(readingBooks, bid)
				}
			}

			// If exactly one book is being read, auto-associate
			if len(readingBooks) == 1 {
				bookID = &readingBooks[0]
				isAssociated = true
			}
		}
	}

	// Insert the annotation
	query := `
		INSERT INTO annotations (user_id, book_id, type, content, context, page_number, tags, is_associated, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	var annotationID string
	var createdAt, updatedAt time.Time
	err := h.DB.QueryRow(ctx, query, userID, bookID, req.Type, req.Content, req.Context, req.PageNumber, req.Tags, isAssociated).
		Scan(&annotationID, &createdAt, &updatedAt)
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create annotation",
		})
	}

	// Prepare response with full annotation details
	response := map[string]interface{}{
		"id":            annotationID,
		"user_id":       userID,
		"book_id":       bookID,
		"type":          req.Type,
		"content":       req.Content,
		"context":       req.Context,
		"page_number":   req.PageNumber,
		"tags":          req.Tags,
		"is_associated": isAssociated,
		"created_at":    createdAt,
		"updated_at":    updatedAt,
	}

	// If associated with a book, fetch and include book details
	if bookID != nil {
		bookQuery := `
			SELECT id, title, authors, cover_url
			FROM books
			WHERE id = $1
		`
		var book struct {
			ID       string
			Title    string
			Authors  []string
			CoverURL *string
		}
		err := h.DB.QueryRow(ctx, bookQuery, *bookID).Scan(&book.ID, &book.Title, &book.Authors, &book.CoverURL)
		if err == nil {
			response["book"] = map[string]interface{}{
				"id":        book.ID,
				"title":     book.Title,
				"authors":   book.Authors,
				"cover_url": book.CoverURL,
			}
		}
	}

	return c.JSON(http.StatusCreated, response)
}

// GetUserRecents returns books the user is currently reading or recently finished
func (h *AnnotationHandler) GetUserRecents(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	query := `
		SELECT l.book_id, b.title, b.cover_url, l.status, l.updated_at
		FROM logs l
		JOIN books b ON l.book_id = b.id
		WHERE l.user_id = $1 AND (l.status = 'reading' OR l.status = 'read')
		ORDER BY 
			CASE 
				WHEN l.status = 'reading' THEN 1
				ELSE 2
			END,
			l.updated_at DESC
		LIMIT 5
	`

	rows, err := h.DB.Query(ctx, query, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch recent books",
		})
	}
	defer rows.Close()

	var books []RecentBook
	for rows.Next() {
		var book RecentBook
		var updatedAt time.Time
		err := rows.Scan(&book.ID, &book.Title, &book.CoverURL, &book.Status, &updatedAt)
		if err != nil {
			continue
		}
		book.LastUpdated = updatedAt.Format(time.RFC3339)
		books = append(books, book)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"books": books,
		"count": len(books),
	})
}

// GetBookAnnotations returns all annotations for a specific book
func (h *AnnotationHandler) GetBookAnnotations(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	bookID := c.Param("id")
	if bookID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "book ID is required",
		})
	}

	annotationType := c.QueryParam("type") // Optional filter: 'note' or 'highlight'

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, type, content, context, page_number, tags, is_associated, created_at, updated_at
		FROM annotations
		WHERE user_id = $1 AND book_id = $2
	`

	args := []interface{}{userID, bookID}

	if annotationType != "" {
		query += " AND type = $3"
		args = append(args, annotationType)
	}

	query += " ORDER BY page_number ASC NULLS LAST, created_at ASC"

	rows, err := h.DB.Query(ctx, query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch annotations",
		})
	}
	defer rows.Close()

	annotations := []map[string]interface{}{}
	for rows.Next() {
		var annotation struct {
			ID           string
			Type         string
			Content      string
			Context      *string
			PageNumber   *int
			Tags         []string
			IsAssociated bool
			CreatedAt    time.Time
			UpdatedAt    time.Time
		}

		err := rows.Scan(
			&annotation.ID, &annotation.Type, &annotation.Content,
			&annotation.Context, &annotation.PageNumber, &annotation.Tags,
			&annotation.IsAssociated, &annotation.CreatedAt, &annotation.UpdatedAt,
		)
		if err != nil {
			continue
		}

		annotations = append(annotations, map[string]interface{}{
			"id":            annotation.ID,
			"type":          annotation.Type,
			"content":       annotation.Content,
			"context":       annotation.Context,
			"page_number":   annotation.PageNumber,
			"tags":          annotation.Tags,
			"is_associated": annotation.IsAssociated,
			"created_at":    annotation.CreatedAt,
			"updated_at":    annotation.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"annotations": annotations,
		"count":       len(annotations),
	})
}

// GetUnassociatedAnnotations returns all annotations not linked to a book
func (h *AnnotationHandler) GetUnassociatedAnnotations(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, type, content, context, page_number, tags, created_at, updated_at
		FROM annotations
		WHERE user_id = $1 AND is_associated = false
		ORDER BY created_at DESC
	`

	rows, err := h.DB.Query(ctx, query, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch unassociated annotations",
		})
	}
	defer rows.Close()

	annotations := []map[string]interface{}{}
	for rows.Next() {
		var annotation struct {
			ID         string
			Type       string
			Content    string
			Context    *string
			PageNumber *int
			Tags       []string
			CreatedAt  time.Time
			UpdatedAt  time.Time
		}

		err := rows.Scan(
			&annotation.ID, &annotation.Type, &annotation.Content,
			&annotation.Context, &annotation.PageNumber, &annotation.Tags,
			&annotation.CreatedAt, &annotation.UpdatedAt,
		)
		if err != nil {
			continue
		}

		annotations = append(annotations, map[string]interface{}{
			"id":          annotation.ID,
			"type":        annotation.Type,
			"content":     annotation.Content,
			"context":     annotation.Context,
			"page_number": annotation.PageNumber,
			"tags":        annotation.Tags,
			"created_at":  annotation.CreatedAt,
			"updated_at":  annotation.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"annotations": annotations,
		"count":       len(annotations),
	})
}

// UpdateAnnotation updates an annotation (primarily for associating with a book)
func (h *AnnotationHandler) UpdateAnnotation(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	annotationID := c.Param("id")
	if annotationID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "annotation ID is required",
		})
	}

	var req UpdateAnnotationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Verify ownership
	var ownerID string
	err := h.DB.QueryRow(ctx, "SELECT user_id FROM annotations WHERE id = $1", annotationID).Scan(&ownerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "annotation not found",
		})
	}

	if ownerID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "forbidden",
		})
	}

	// Build dynamic update query
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if req.BookID != nil {
		updates = append(updates, "book_id = $"+string(rune('0'+argCount)))
		args = append(args, *req.BookID)
		argCount++
		
		// If associating with a book, set is_associated to true
		updates = append(updates, "is_associated = true")
	}

	if req.Content != nil {
		updates = append(updates, "content = $"+string(rune('0'+argCount)))
		args = append(args, *req.Content)
		argCount++
	}

	if req.PageNumber != nil {
		updates = append(updates, "page_number = $"+string(rune('0'+argCount)))
		args = append(args, *req.PageNumber)
		argCount++
	}

	if req.Tags != nil {
		updates = append(updates, "tags = $"+string(rune('0'+argCount)))
		args = append(args, req.Tags)
		argCount++
	}

	if len(updates) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "no fields to update",
		})
	}

	updates = append(updates, "updated_at = NOW()")
	args = append(args, annotationID)

	query := "UPDATE annotations SET " + updates[0]
	for i := 1; i < len(updates); i++ {
		query += ", " + updates[i]
	}
	query += " WHERE id = $" + string(rune('0'+argCount)) + " RETURNING updated_at"

	var updatedAt time.Time
	err = h.DB.QueryRow(ctx, query, args...).Scan(&updatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update annotation",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":         annotationID,
		"updated_at": updatedAt,
		"message":    "annotation updated successfully",
	})
}

// DeleteAnnotation deletes an annotation
func (h *AnnotationHandler) DeleteAnnotation(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	annotationID := c.Param("id")
	if annotationID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "annotation ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Verify ownership and delete
	query := `
		DELETE FROM annotations
		WHERE id = $1 AND user_id = $2
	`

	result, err := h.DB.Exec(ctx, query, annotationID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to delete annotation",
		})
	}

	if result.RowsAffected() == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "annotation not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "annotation deleted successfully",
	})
}

// SearchAnnotations performs full-text search across user's annotations
func (h *AnnotationHandler) SearchAnnotations(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	searchQuery := c.QueryParam("q")
	if searchQuery == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "search query 'q' is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	query := `
		SELECT a.id, a.type, a.content, a.context, a.page_number, a.tags, 
		       a.book_id, a.is_associated, a.created_at,
		       b.title as book_title, b.cover_url as book_cover
		FROM annotations a
		LEFT JOIN books b ON a.book_id = b.id
		WHERE a.user_id = $1 
		  AND to_tsvector('english', a.content) @@ plainto_tsquery('english', $2)
		ORDER BY a.created_at DESC
		LIMIT 50
	`

	rows, err := h.DB.Query(ctx, query, userID, searchQuery)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to search annotations",
		})
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var annotation struct {
			ID           string
			Type         string
			Content      string
			Context      *string
			PageNumber   *int
			Tags         []string
			BookID       *string
			IsAssociated bool
			CreatedAt    time.Time
			BookTitle    *string
			BookCover    *string
		}

		err := rows.Scan(
			&annotation.ID, &annotation.Type, &annotation.Content,
			&annotation.Context, &annotation.PageNumber, &annotation.Tags,
			&annotation.BookID, &annotation.IsAssociated, &annotation.CreatedAt,
			&annotation.BookTitle, &annotation.BookCover,
		)
		if err != nil {
			continue
		}

		result := map[string]interface{}{
			"id":            annotation.ID,
			"type":          annotation.Type,
			"content":       annotation.Content,
			"context":       annotation.Context,
			"page_number":   annotation.PageNumber,
			"tags":          annotation.Tags,
			"is_associated": annotation.IsAssociated,
			"created_at":    annotation.CreatedAt,
		}

		if annotation.BookID != nil {
			result["book"] = map[string]interface{}{
				"id":        *annotation.BookID,
				"title":     annotation.BookTitle,
				"cover_url": annotation.BookCover,
			}
		}

		results = append(results, result)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"results": results,
		"count":   len(results),
		"query":   searchQuery,
	})
}
