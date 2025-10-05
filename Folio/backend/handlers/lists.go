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
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	IsPublic        *bool   `json:"is_public"`
	HeaderImageURL  *string `json:"header_image_url"`
	ThemeColor      *string `json:"theme_color"`
}

type UpdateListRequest struct {
	Name            *string `json:"name"`
	Description     *string `json:"description"`
	IsPublic        *bool   `json:"is_public"`
	HeaderImageURL  *string `json:"header_image_url"`
	ThemeColor      *string `json:"theme_color"`
}

type AddBookToListRequest struct {
	BookID string  `json:"book_id"`
	Notes  *string `json:"notes"`
}

type UpdateListItemRequest struct {
	Notes *string `json:"notes"`
}

type UpdateListItemOrderRequest struct {
	Order int `json:"order"`
}

type ReorderListItemsRequest struct {
	ItemIDs []string `json:"item_ids"`
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

	themeColor := "#6366f1"
	if req.ThemeColor != nil {
		themeColor = *req.ThemeColor
	}

	query := `
		INSERT INTO lists (user_id, name, description, is_public, header_image_url, theme_color, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	var listID string
	var createdAt, updatedAt time.Time
	err := h.DB.QueryRow(ctx, query, userID, req.Name, req.Description, isPublic, req.HeaderImageURL, themeColor).Scan(&listID, &createdAt, &updatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create list",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":              listID,
		"user_id":         userID,
		"name":            req.Name,
		"description":     req.Description,
		"is_public":       isPublic,
		"header_image_url": req.HeaderImageURL,
		"theme_color":     themeColor,
		"items_count":     0,
		"created_at":      createdAt,
		"updated_at":      updatedAt,
	})
}

// GetMyLists retrieves all lists for the current authenticated user
func (h *ListHandler) GetMyLists(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Get all lists for the current user (both public and private)
	query := `
		SELECT id, user_id, name, description, is_public, header_image_url, theme_color, items_count, created_at, updated_at
		FROM lists
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := h.DB.Query(ctx, query, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch lists",
		})
	}
	defer rows.Close()

	lists := []map[string]interface{}{}
	for rows.Next() {
		var list struct {
			ID             string
			UserID         string
			Name           string
			Description    *string
			IsPublic       bool
			HeaderImageURL *string
			ThemeColor     string
			ItemsCount     int
			CreatedAt      time.Time
			UpdatedAt      time.Time
		}

		err := rows.Scan(&list.ID, &list.UserID, &list.Name, &list.Description, &list.IsPublic, &list.HeaderImageURL, &list.ThemeColor, &list.ItemsCount, &list.CreatedAt, &list.UpdatedAt)
		if err != nil {
			continue
		}

		lists = append(lists, map[string]interface{}{
			"id":              list.ID,
			"user_id":         list.UserID,
			"name":            list.Name,
			"description":     list.Description,
			"is_public":       list.IsPublic,
			"header_image_url": list.HeaderImageURL,
			"theme_color":     list.ThemeColor,
			"items_count":     list.ItemsCount,
			"created_at":      list.CreatedAt,
			"updated_at":      list.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"lists": lists,
		"count": len(lists),
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
			SELECT id, user_id, name, description, is_public, header_image_url, theme_color, items_count, created_at, updated_at
			FROM lists
			WHERE user_id = $1
			ORDER BY created_at DESC
		`
	} else {
		// Show only public lists for other users
		query = `
			SELECT id, user_id, name, description, is_public, header_image_url, theme_color, items_count, created_at, updated_at
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
			ID             string
			UserID         string
			Name           string
			Description    *string
			IsPublic       bool
			HeaderImageURL *string
			ThemeColor     string
			ItemsCount     int
			CreatedAt      time.Time
			UpdatedAt      time.Time
		}

		err := rows.Scan(&list.ID, &list.UserID, &list.Name, &list.Description, &list.IsPublic, &list.HeaderImageURL, &list.ThemeColor, &list.ItemsCount, &list.CreatedAt, &list.UpdatedAt)
		if err != nil {
			continue
		}

		lists = append(lists, map[string]interface{}{
			"id":              list.ID,
			"user_id":         list.UserID,
			"name":            list.Name,
			"description":     list.Description,
			"is_public":       list.IsPublic,
			"header_image_url": list.HeaderImageURL,
			"theme_color":     list.ThemeColor,
			"items_count":     list.ItemsCount,
			"created_at":      list.CreatedAt,
			"updated_at":      list.UpdatedAt,
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

	// Get list details with creator info
	var list struct {
		ID             string
		UserID         string
		Name           string
		Description    *string
		IsPublic       bool
		HeaderImageURL *string
		ThemeColor     string
		ItemsCount     int
		CreatedAt      time.Time
		UpdatedAt      time.Time
		CreatorName    string
		CreatorUsername string
		CreatorPicture *string
	}

	query := `
		SELECT l.id, l.user_id, l.name, l.description, l.is_public, l.header_image_url, l.theme_color, l.items_count, l.created_at, l.updated_at,
		       u.name, u.username, u.picture
		FROM lists l
		JOIN users u ON l.user_id = u.id
		WHERE l.id = $1
	`
	err := h.DB.QueryRow(ctx, query, listID).Scan(&list.ID, &list.UserID, &list.Name, &list.Description, &list.IsPublic, &list.HeaderImageURL, &list.ThemeColor, &list.ItemsCount, &list.CreatedAt, &list.UpdatedAt, &list.CreatorName, &list.CreatorUsername, &list.CreatorPicture)
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

	// Get users who liked this list (top 5)
	likedByQuery := `
		SELECT u.id, u.username, u.name, u.picture
		FROM list_likes ll
		JOIN users u ON ll.user_id = u.id
		WHERE ll.list_id = $1
		ORDER BY ll.created_at DESC
		LIMIT 5
	`

	likedByRows, err := h.DB.Query(ctx, likedByQuery, listID)
	likedBy := []map[string]interface{}{}
	if err == nil {
		defer likedByRows.Close()
		for likedByRows.Next() {
			var user struct {
				ID       string
				Username string
				Name     string
				Picture  *string
			}
			if err := likedByRows.Scan(&user.ID, &user.Username, &user.Name, &user.Picture); err == nil {
				likedBy = append(likedBy, map[string]interface{}{
					"id":       user.ID,
					"username": user.Username,
					"name":     user.Name,
					"picture":  user.Picture,
				})
			}
		}
	}

	// Get likes and comments count
	var likesCount, commentsCount int
	h.DB.QueryRow(ctx, "SELECT COUNT(*) FROM list_likes WHERE list_id = $1", listID).Scan(&likesCount)
	h.DB.QueryRow(ctx, "SELECT COUNT(*) FROM list_comments WHERE list_id = $1", listID).Scan(&commentsCount)

	// Check if current user liked this list
	var isLiked bool
	if currentUserID != "" {
		h.DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM list_likes WHERE list_id = $1 AND user_id = $2)", listID, currentUserID).Scan(&isLiked)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":              list.ID,
		"user_id":         list.UserID,
		"name":            list.Name,
		"description":     list.Description,
		"is_public":       list.IsPublic,
		"header_image_url": list.HeaderImageURL,
		"theme_color":     list.ThemeColor,
		"items_count":     list.ItemsCount,
		"likes_count":     likesCount,
		"comments_count":  commentsCount,
		"is_liked":        isLiked,
		"created_at":      list.CreatedAt,
		"updated_at":      list.UpdatedAt,
		"creator": map[string]interface{}{
			"id":       list.UserID,
			"name":     list.CreatorName,
			"username": list.CreatorUsername,
			"picture":  list.CreatorPicture,
		},
		"items":    items,
		"liked_by": likedBy,
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
	if req.HeaderImageURL != nil {
		argCount++
		query += ", header_image_url = $" + string(rune(argCount+'0'-1))
		args = append(args, *req.HeaderImageURL)
	}
	if req.ThemeColor != nil {
		argCount++
		query += ", theme_color = $" + string(rune(argCount+'0'-1))
		args = append(args, *req.ThemeColor)
	}

	query += " WHERE id = $1 RETURNING id, name, description, is_public, header_image_url, theme_color, updated_at"
	args = append([]interface{}{listID}, args...)

	var updatedList struct {
		ID             string
		Name           string
		Description    *string
		IsPublic       bool
		HeaderImageURL *string
		ThemeColor     string
		UpdatedAt      time.Time
	}

	err = h.DB.QueryRow(ctx, query, args...).Scan(&updatedList.ID, &updatedList.Name, &updatedList.Description, &updatedList.IsPublic, &updatedList.HeaderImageURL, &updatedList.ThemeColor, &updatedList.UpdatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update list",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":              updatedList.ID,
		"name":            updatedList.Name,
		"description":     updatedList.Description,
		"is_public":       updatedList.IsPublic,
		"header_image_url": updatedList.HeaderImageURL,
		"theme_color":     updatedList.ThemeColor,
		"updated_at":      updatedList.UpdatedAt,
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

// UpdateListItem updates a list item's details (e.g., notes)
func (h *ListHandler) UpdateListItem(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	listID := c.Param("id")
	itemID := c.Param("itemId")

	var req UpdateListItemRequest
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

	_, err = h.DB.Exec(ctx, "UPDATE list_items SET notes = $1 WHERE id = $2", req.Notes, itemID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update list item",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "list item updated successfully",
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

// ReorderListItems reorders all items in a list
func (h *ListHandler) ReorderListItems(c echo.Context) error {
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

	var req ReorderListItemsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if len(req.ItemIDs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "item_ids cannot be empty",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	// Start transaction
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to start transaction",
		})
	}
	defer tx.Rollback(ctx)

	// Check if list exists and belongs to user
	var ownerID string
	err = tx.QueryRow(ctx, "SELECT user_id FROM lists WHERE id = $1", listID).Scan(&ownerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "list not found",
		})
	}

	if ownerID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "you can only reorder your own lists",
		})
	}

	// Update item orders in a single transaction
	for i, itemID := range req.ItemIDs {
		_, err = tx.Exec(ctx, "UPDATE list_items SET item_order = $1 WHERE id = $2 AND list_id = $3", i, itemID, listID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to update item order",
			})
		}
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to commit transaction",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "items reordered successfully",
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

// GetPopularLists gets popular public lists
func (h *ListHandler) GetPopularLists(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	limit := c.QueryParam("limit")
	if limit == "" {
		limit = "12"
	}

	query := `
		SELECT l.id, l.user_id, l.name, l.description, l.is_public, l.header_image_url, l.theme_color, l.items_count, l.created_at, l.updated_at,
		       u.name, u.username, u.picture,
		       COALESCE(like_counts.likes_count, 0) as likes_count,
		       COALESCE(comment_counts.comments_count, 0) as comments_count
		FROM lists l
		JOIN users u ON l.user_id = u.id
		LEFT JOIN (
			SELECT list_id, COUNT(*) as likes_count
			FROM list_likes
			GROUP BY list_id
		) like_counts ON l.id = like_counts.list_id
		LEFT JOIN (
			SELECT list_id, COUNT(*) as comments_count
			FROM list_comments
			GROUP BY list_id
		) comment_counts ON l.id = comment_counts.list_id
		WHERE l.is_public = true AND l.items_count > 0
		ORDER BY (COALESCE(like_counts.likes_count, 0) + COALESCE(comment_counts.comments_count, 0)) DESC, l.created_at DESC
		LIMIT $1
	`

	rows, err := h.DB.Query(ctx, query, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch popular lists",
		})
	}
	defer rows.Close()

	lists := []map[string]interface{}{}
	for rows.Next() {
		var list struct {
			ID             string
			UserID         string
			Name           string
			Description    *string
			IsPublic       bool
			HeaderImageURL *string
			ThemeColor     string
			ItemsCount     int
			CreatedAt      time.Time
			UpdatedAt      time.Time
			CreatorName    string
			CreatorUsername string
			CreatorPicture *string
			LikesCount     int
			CommentsCount  int
		}

		err := rows.Scan(&list.ID, &list.UserID, &list.Name, &list.Description, &list.IsPublic, &list.HeaderImageURL, &list.ThemeColor, &list.ItemsCount, &list.CreatedAt, &list.UpdatedAt, &list.CreatorName, &list.CreatorUsername, &list.CreatorPicture, &list.LikesCount, &list.CommentsCount)
		if err != nil {
			continue
		}

		lists = append(lists, map[string]interface{}{
			"id":              list.ID,
			"user_id":         list.UserID,
			"name":            list.Name,
			"description":     list.Description,
			"is_public":       list.IsPublic,
			"header_image_url": list.HeaderImageURL,
			"theme_color":     list.ThemeColor,
			"items_count":     list.ItemsCount,
			"created_at":      list.CreatedAt,
			"updated_at":      list.UpdatedAt,
			"creator": map[string]interface{}{
				"id":       list.UserID,
				"name":     list.CreatorName,
				"username": list.CreatorUsername,
				"picture":  list.CreatorPicture,
			},
			"likes_count":    list.LikesCount,
			"comments_count": list.CommentsCount,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"lists": lists,
		"count": len(lists),
	})
}

