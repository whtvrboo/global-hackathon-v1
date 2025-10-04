package handlers

import (
	"context"
	"folio/api/auth"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type ListHandler struct {
	DB *pgxpool.Pool
}

type CreateListRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    *bool  `json:"is_public"`
}

type UpdateListRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IsPublic    *bool   `json:"is_public"`
}

type AddBookToListRequest struct {
	BookID string  `json:"book_id"`
	Notes  *string `json:"notes"`
}

type UpdateListItemOrderRequest struct {
	Order int `json:"order"`
}

// CreateList creates a new list
func (h *ListHandler) CreateList(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	var req CreateListRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "name is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	query := `
		INSERT INTO lists (user_id, name, description, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	var listID string
	var createdAt, updatedAt time.Time
	err := h.DB.QueryRow(ctx, query, userID, req.Name, req.Description, isPublic).Scan(&listID, &createdAt, &updatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create list",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":          listID,
		"user_id":     userID,
		"name":        req.Name,
		"description": req.Description,
		"is_public":   isPublic,
		"items_count": 0,
		"created_at":  createdAt,
		"updated_at":  updatedAt,
	})
}

// GetUserLists retrieves all lists for a given username
func (h *ListHandler) GetUserLists(c echo.Context) error {
	username := c.Param("username")
	if username == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "username is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Check if requesting user is viewing their own profile
	currentUserID := auth.GetUserID(c)
	var isOwnProfile bool
	var profileUserID string

	err := h.DB.QueryRow(ctx, "SELECT id FROM users WHERE username = $1", username).Scan(&profileUserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "user not found",
		})
	}

	isOwnProfile = (currentUserID == profileUserID)

	var query string
	if isOwnProfile {
		// Show all lists for own profile
		query = `
			SELECT id, user_id, name, description, is_public, items_count, created_at, updated_at
			FROM lists
			WHERE user_id = $1
			ORDER BY created_at DESC
		`
	} else {
		// Show only public lists for other users
		query = `
			SELECT id, user_id, name, description, is_public, items_count, created_at, updated_at
			FROM lists
			WHERE user_id = $1 AND is_public = true
			ORDER BY created_at DESC
		`
	}

	rows, err := h.DB.Query(ctx, query, profileUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch lists",
		})
	}
	defer rows.Close()

	lists := []map[string]interface{}{}
	for rows.Next() {
		var list struct {
			ID          string
			UserID      string
			Name        string
			Description *string
			IsPublic    bool
			ItemsCount  int
			CreatedAt   time.Time
			UpdatedAt   time.Time
		}

		err := rows.Scan(&list.ID, &list.UserID, &list.Name, &list.Description, &list.IsPublic, &list.ItemsCount, &list.CreatedAt, &list.UpdatedAt)
		if err != nil {
			continue
		}

		lists = append(lists, map[string]interface{}{
			"id":          list.ID,
			"user_id":     list.UserID,
			"name":        list.Name,
			"description": list.Description,
			"is_public":   list.IsPublic,
			"items_count": list.ItemsCount,
			"created_at":  list.CreatedAt,
			"updated_at":  list.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"lists": lists,
		"count": len(lists),
	})
}

// GetList retrieves a specific list with its books
func (h *ListHandler) GetList(c echo.Context) error {
	listID := c.Param("id")
	if listID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "list_id is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Get list details
	var list struct {
		ID          string
		UserID      string
		Name        string
		Description *string
		IsPublic    bool
		ItemsCount  int
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	query := `SELECT id, user_id, name, description, is_public, items_count, created_at, updated_at FROM lists WHERE id = $1`
	err := h.DB.QueryRow(ctx, query, listID).Scan(&list.ID, &list.UserID, &list.Name, &list.Description, &list.IsPublic, &list.ItemsCount, &list.CreatedAt, &list.UpdatedAt)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "list not found",
		})
	}

	// Check permissions
	currentUserID := auth.GetUserID(c)
	if !list.IsPublic && currentUserID != list.UserID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "you don't have permission to view this list",
		})
	}

	// Get list items with book details
	itemsQuery := `
		SELECT li.id, li.book_id, li.notes, li.item_order, li.created_at,
		       b.title, b.authors, b.cover_url, b.description
		FROM list_items li
		JOIN books b ON li.book_id = b.id
		WHERE li.list_id = $1
		ORDER BY li.item_order ASC, li.created_at DESC
	`

	rows, err := h.DB.Query(ctx, itemsQuery, listID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch list items",
		})
	}
	defer rows.Close()

	items := []map[string]interface{}{}
	for rows.Next() {
		var item struct {
			ID          string
			BookID      string
			Notes       *string
			ItemOrder   int
			CreatedAt   time.Time
			BookTitle   string
			Authors     []string
			CoverURL    *string
			Description *string
		}

		err := rows.Scan(&item.ID, &item.BookID, &item.Notes, &item.ItemOrder, &item.CreatedAt, &item.BookTitle, &item.Authors, &item.CoverURL, &item.Description)
		if err != nil {
			continue
		}

		items = append(items, map[string]interface{}{
			"id":         item.ID,
			"book_id":    item.BookID,
			"notes":      item.Notes,
			"item_order": item.ItemOrder,
			"created_at": item.CreatedAt,
			"book": map[string]interface{}{
				"id":          item.BookID,
				"title":       item.BookTitle,
				"authors":     item.Authors,
				"cover_url":   item.CoverURL,
				"description": item.Description,
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":          list.ID,
		"user_id":     list.UserID,
		"name":        list.Name,
		"description": list.Description,
		"is_public":   list.IsPublic,
		"items_count": list.ItemsCount,
		"created_at":  list.CreatedAt,
		"updated_at":  list.UpdatedAt,
		"items":       items,
	})
}

// UpdateList updates a list's details
func (h *ListHandler) UpdateList(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	listID := c.Param("id")
	if listID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "list_id is required",
		})
	}

	var req UpdateListRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Check if list exists and belongs to user
	var ownerID string
	err := h.DB.QueryRow(ctx, "SELECT user_id FROM lists WHERE id = $1", listID).Scan(&ownerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "list not found",
		})
	}

	if ownerID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "you can only update your own lists",
		})
	}

	// Build update query dynamically
	query := "UPDATE lists SET updated_at = NOW()"
	args := []interface{}{}
	argCount := 1

	if req.Name != nil {
		argCount++
		query += ", name = $" + string(rune(argCount+'0'-1))
		args = append(args, *req.Name)
	}
	if req.Description != nil {
		argCount++
		query += ", description = $" + string(rune(argCount+'0'-1))
		args = append(args, *req.Description)
	}
	if req.IsPublic != nil {
		argCount++
		query += ", is_public = $" + string(rune(argCount+'0'-1))
		args = append(args, *req.IsPublic)
	}

	query += " WHERE id = $1 RETURNING id, name, description, is_public, updated_at"
	args = append([]interface{}{listID}, args...)

	var updatedList struct {
		ID          string
		Name        string
		Description *string
		IsPublic    bool
		UpdatedAt   time.Time
	}

	err = h.DB.QueryRow(ctx, query, args...).Scan(&updatedList.ID, &updatedList.Name, &updatedList.Description, &updatedList.IsPublic, &updatedList.UpdatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update list",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":          updatedList.ID,
		"name":        updatedList.Name,
		"description": updatedList.Description,
		"is_public":   updatedList.IsPublic,
		"updated_at":  updatedList.UpdatedAt,
	})
}

// DeleteList deletes a list
func (h *ListHandler) DeleteList(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	listID := c.Param("id")
	if listID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "list_id is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Check if list exists and belongs to user
	var ownerID string
	err := h.DB.QueryRow(ctx, "SELECT user_id FROM lists WHERE id = $1", listID).Scan(&ownerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "list not found",
		})
	}

	if ownerID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "you can only delete your own lists",
		})
	}

	_, err = h.DB.Exec(ctx, "DELETE FROM lists WHERE id = $1", listID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to delete list",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "list deleted successfully",
	})
}

// AddBookToList adds a book to a list
func (h *ListHandler) AddBookToList(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	listID := c.Param("id")
	if listID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "list_id is required",
		})
	}

	var req AddBookToListRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.BookID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "book_id is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Check if list exists and belongs to user
	var ownerID string
	err := h.DB.QueryRow(ctx, "SELECT user_id FROM lists WHERE id = $1", listID).Scan(&ownerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "list not found",
		})
	}

	if ownerID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "you can only add books to your own lists",
		})
	}

	// Get the current max order for the list
	var maxOrder int
	h.DB.QueryRow(ctx, "SELECT COALESCE(MAX(item_order), -1) FROM list_items WHERE list_id = $1", listID).Scan(&maxOrder)

	query := `
		INSERT INTO list_items (list_id, book_id, notes, item_order, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (list_id, book_id) DO NOTHING
		RETURNING id, created_at
	`

	var itemID string
	var createdAt time.Time
	err = h.DB.QueryRow(ctx, query, listID, req.BookID, req.Notes, maxOrder+1).Scan(&itemID, &createdAt)
	if err != nil {
		// Check if it's a duplicate
		var exists bool
		h.DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM list_items WHERE list_id = $1 AND book_id = $2)", listID, req.BookID).Scan(&exists)
		if exists {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "book already in list",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to add book to list",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":         itemID,
		"list_id":    listID,
		"book_id":    req.BookID,
		"notes":      req.Notes,
		"item_order": maxOrder + 1,
		"created_at": createdAt,
	})
}

// UpdateListItemOrder updates the order of a list item
func (h *ListHandler) UpdateListItemOrder(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	listID := c.Param("id")
	itemID := c.Param("itemId")

	var req UpdateListItemOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Verify ownership
	var ownerID string
	err := h.DB.QueryRow(ctx, `
		SELECT l.user_id FROM lists l
		JOIN list_items li ON l.id = li.list_id
		WHERE l.id = $1 AND li.id = $2
	`, listID, itemID).Scan(&ownerID)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "list item not found",
		})
	}

	if ownerID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "unauthorized",
		})
	}

	_, err = h.DB.Exec(ctx, "UPDATE list_items SET item_order = $1 WHERE id = $2", req.Order, itemID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update item order",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "item order updated",
	})
}

// RemoveBookFromList removes a book from a list
func (h *ListHandler) RemoveBookFromList(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	listID := c.Param("id")
	itemID := c.Param("itemId")

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Verify ownership
	var ownerID string
	err := h.DB.QueryRow(ctx, `
		SELECT l.user_id FROM lists l
		JOIN list_items li ON l.id = li.list_id
		WHERE l.id = $1 AND li.id = $2
	`, listID, itemID).Scan(&ownerID)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "list item not found",
		})
	}

	if ownerID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "unauthorized",
		})
	}

	_, err = h.DB.Exec(ctx, "DELETE FROM list_items WHERE id = $1", itemID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to remove book from list",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "book removed from list",
	})
}

// LikeList likes a list
func (h *ListHandler) LikeList(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	listID := c.Param("id")
	if listID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "list id is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Check if list exists and is public
	var isPublic bool
	err := h.DB.QueryRow(ctx, "SELECT is_public FROM lists WHERE id = $1", listID).Scan(&isPublic)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "list not found",
		})
	}

	if !isPublic {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "cannot like private list",
		})
	}

	// Insert like (ignore if already exists)
	query := `
		INSERT INTO list_likes (list_id, user_id, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (list_id, user_id) DO NOTHING
	`

	_, err = h.DB.Exec(ctx, query, listID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to like list",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

// UnlikeList unlikes a list
func (h *ListHandler) UnlikeList(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	listID := c.Param("id")
	if listID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "list id is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM list_likes WHERE list_id = $1 AND user_id = $2`
	_, err := h.DB.Exec(ctx, query, listID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to unlike list",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

// GetListComments gets comments for a list
func (h *ListHandler) GetListComments(c echo.Context) error {
	listID := c.Param("id")
	if listID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "list id is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	query := `
		SELECT c.id, c.user_id, c.content, c.created_at, c.updated_at,
		       u.username, u.name, u.picture
		FROM list_comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.list_id = $1
		ORDER BY c.created_at ASC
	`

	rows, err := h.DB.Query(ctx, query, listID)
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
			Content   string
			CreatedAt time.Time
			UpdatedAt time.Time
			Username  string
			Name      string
			Picture   *string
		}

		err := rows.Scan(
			&comment.ID, &comment.UserID, &comment.Content,
			&comment.CreatedAt, &comment.UpdatedAt,
			&comment.Username, &comment.Name, &comment.Picture,
		)
		if err != nil {
			continue
		}

		comments = append(comments, map[string]interface{}{
			"id":         comment.ID,
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

// AddListComment adds a comment to a list
func (h *ListHandler) AddListComment(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	listID := c.Param("id")
	if listID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "list id is required",
		})
	}

	var req struct {
		Content string `json:"content"`
	}

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

	// Check if list exists and is public
	var isPublic bool
	err := h.DB.QueryRow(ctx, "SELECT is_public FROM lists WHERE id = $1", listID).Scan(&isPublic)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "list not found",
		})
	}

	if !isPublic {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "cannot comment on private list",
		})
	}

	query := `
		INSERT INTO list_comments (list_id, user_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	var commentID string
	var createdAt, updatedAt time.Time
	err = h.DB.QueryRow(ctx, query, listID, userID, req.Content).Scan(&commentID, &createdAt, &updatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to add comment",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":         commentID,
		"content":    req.Content,
		"created_at": createdAt,
		"updated_at": updatedAt,
	})
}

