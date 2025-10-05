package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type BookHandler struct {
	DB *pgxpool.Pool
}

type BookSearchResult struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Authors     []string `json:"authors"`
	Description string   `json:"description"`
	CoverURL    string   `json:"cover_url"`
	PublishedDate string `json:"published_date"`
}

// SearchBooks searches for books using local database first, then Google Books API
func (h *BookHandler) SearchBooks(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "query parameter 'q' is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	// First, search local database
	localBooks, err := h.searchLocalBooks(ctx, query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("failed to search local books: %v", err),
		})
	}

	// If we have local results, return them
	if len(localBooks) > 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"results": localBooks,
			"count":   len(localBooks),
			"source":  "local",
		})
	}

	// If no local results, search Google Books API and cache results
	googleBooks, err := h.searchGoogleBooks(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("failed to search books: %v", err),
		})
	}

	// Cache the Google Books results in our database
	if len(googleBooks) > 0 {
		go h.cacheBooks(ctx, googleBooks)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"results": googleBooks,
		"count":   len(googleBooks),
		"source":  "google",
	})
}

// GetBook gets detailed book information
func (h *BookHandler) GetBook(c echo.Context) error {
	bookID := c.Param("id")
	if bookID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "book ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Check if book exists in cache
	var book struct {
		ID            string    `json:"id"`
		Title         string    `json:"title"`
		Authors       []string  `json:"authors"`
		Description   *string   `json:"description"`
		CoverURL      *string   `json:"cover_url"`
		PublishedDate *string   `json:"published_date"`
		PageCount     *int      `json:"page_count"`
		ISBN10        *string   `json:"isbn_10"`
		ISBN13        *string   `json:"isbn_13"`
		Categories    []string  `json:"categories"`
		Language      *string   `json:"language"`
		Publisher     *string   `json:"publisher"`
		Rating        *float64  `json:"rating"`
		RatingsCount  *int      `json:"ratings_count"`
	}

	query := `
		SELECT id, title, authors, description, cover_url, published_date,
		       page_count, isbn_10, isbn_13, categories, language, publisher,
		       rating, ratings_count
		FROM books
		WHERE id = $1
	`

	err := h.DB.QueryRow(ctx, query, bookID).Scan(
		&book.ID, &book.Title, &book.Authors, &book.Description,
		&book.CoverURL, &book.PublishedDate, &book.PageCount,
		&book.ISBN10, &book.ISBN13, &book.Categories,
		&book.Language, &book.Publisher, &book.Rating, &book.RatingsCount,
	)

	if err == nil {
		// Book found in cache
		return c.JSON(http.StatusOK, book)
	}

	// Fetch from Google Books API and cache
	bookData, err := h.fetchGoogleBook(bookID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "book not found",
		})
	}

	// Cache the book for future requests
	go h.cacheGoogleBook(ctx, bookData)

	return c.JSON(http.StatusOK, bookData)
}

// searchLocalBooks searches for books in our local database
func (h *BookHandler) searchLocalBooks(ctx context.Context, query string) ([]BookSearchResult, error) {
	searchQuery := "%" + query + "%"
	
	sql := `
		SELECT id, title, author, description, cover_url, isbn, published_date, pages, categories
		FROM books 
		WHERE LOWER(title) LIKE LOWER($1) 
		   OR LOWER(author) LIKE LOWER($1)
		   OR LOWER(description) LIKE LOWER($1)
		ORDER BY 
			CASE 
				WHEN LOWER(title) LIKE LOWER($1) THEN 1
				WHEN LOWER(author) LIKE LOWER($1) THEN 2
				ELSE 3
			END,
			title
		LIMIT 20
	`
	
	rows, err := h.DB.Query(ctx, sql, searchQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []BookSearchResult
	for rows.Next() {
		var book struct {
			ID            string
			Title         string
			Author        string
			Description   string
			CoverURL      string
			ISBN          string
			PublishedDate string
			Pages         int
			Categories    string
		}

		err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.Description,
			&book.CoverURL, &book.ISBN, &book.PublishedDate, &book.Pages, &book.Categories,
		)
		if err != nil {
			continue
		}

		// Parse categories JSON
		var categories []string
		if book.Categories != "" {
			json.Unmarshal([]byte(book.Categories), &categories)
		}

		books = append(books, BookSearchResult{
			ID:            book.ID,
			Title:         book.Title,
			Authors:       []string{book.Author},
			Description:   book.Description,
			CoverURL:      book.CoverURL,
			PublishedDate: book.PublishedDate,
		})
	}

	return books, nil
}

// cacheBooks stores Google Books results in our local database
func (h *BookHandler) cacheBooks(ctx context.Context, books []BookSearchResult) {
	for _, book := range books {
		// Check if book already exists
		var exists bool
		err := h.DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM books WHERE id = $1)", book.ID).Scan(&exists)
		if err != nil || exists {
			continue
		}

		// Insert new book
		_, err = h.DB.Exec(ctx, `
			INSERT INTO books (id, title, author, description, cover_url, isbn, published_date, pages, categories, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		`, 
			book.ID,
			book.Title,
			book.Authors[0], // Take first author
			book.Description,
			book.CoverURL,
			"", // ISBN not available from Google Books API in this format
			book.PublishedDate,
			0, // Pages not available
			"[]", // Empty categories for now
		)
		
		if err != nil {
			// Log error but continue with other books
			fmt.Printf("Failed to cache book %s: %v\n", book.Title, err)
		}
	}
}

// cacheGoogleBook stores a single Google Book in our local database
func (h *BookHandler) cacheGoogleBook(ctx context.Context, bookData interface{}) {
	// This would need to be implemented based on the structure of bookData
	// For now, we'll implement a basic version
	fmt.Printf("Caching Google Book: %v\n", bookData)
}

func (h *BookHandler) searchGoogleBooks(query string) ([]BookSearchResult, error) {
	apiKey := getEnv("GOOGLE_BOOKS_API_KEY", "")
	baseURL := "https://www.googleapis.com/books/v1/volumes"
	
	params := url.Values{}
	params.Add("q", query)
	params.Add("maxResults", "20")
	if apiKey != "" {
		params.Add("key", apiKey)
	}

	resp, err := http.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Items []struct {
			ID         string `json:"id"`
			VolumeInfo struct {
				Title         string   `json:"title"`
				Authors       []string `json:"authors"`
				Description   string   `json:"description"`
				ImageLinks    struct {
					Thumbnail string `json:"thumbnail"`
				} `json:"imageLinks"`
				PublishedDate string `json:"publishedDate"`
			} `json:"volumeInfo"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	books := make([]BookSearchResult, 0, len(result.Items))
	for _, item := range result.Items {
		books = append(books, BookSearchResult{
			ID:            item.ID,
			Title:         item.VolumeInfo.Title,
			Authors:       item.VolumeInfo.Authors,
			Description:   item.VolumeInfo.Description,
			CoverURL:      item.VolumeInfo.ImageLinks.Thumbnail,
			PublishedDate: item.VolumeInfo.PublishedDate,
		})
	}

	return books, nil
}

func (h *BookHandler) fetchGoogleBook(bookID string) (map[string]interface{}, error) {
	apiKey := getEnv("GOOGLE_BOOKS_API_KEY", "")
	url := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes/%s", bookID)
	
	if apiKey != "" {
		url += "?key=" + apiKey
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (h *BookHandler) cacheBook(ctx context.Context, bookData map[string]interface{}) error {
	// Simplified caching - in production, parse and store properly
	volumeInfo, ok := bookData["volumeInfo"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid book data")
	}

	query := `
		INSERT INTO books (id, title, authors, description, cover_url, raw_data, api_source, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET
			updated_at = NOW()
	`

	_, err := h.DB.Exec(ctx, query,
		bookData["id"],
		volumeInfo["title"],
		[]string{}, // Simplified
		volumeInfo["description"],
		"",
		bookData,
		"google",
	)

	return err
}

// GetBookReviews gets public reviews for a book
func (h *BookHandler) GetBookReviews(c echo.Context) error {
	bookID := c.Param("id")
	if bookID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "book ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	query := `
		SELECT l.id, l.user_id, l.status, l.rating, l.review, l.notes, l.created_at, l.updated_at,
		       u.username, u.name, u.picture
		FROM logs l
		JOIN users u ON l.user_id = u.id
		WHERE l.book_id = $1 AND l.is_public = true
		ORDER BY l.created_at DESC
		LIMIT 50
	`

	rows, err := h.DB.Query(ctx, query, bookID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch reviews",
		})
	}
	defer rows.Close()

	reviews := []map[string]interface{}{}
	for rows.Next() {
		var review struct {
			ID        string
			UserID    string
			Status    string
			Rating    *int
			Review    *string
			Notes     *string
			CreatedAt time.Time
			UpdatedAt time.Time
			Username  string
			Name      string
			Picture   *string
		}

		err := rows.Scan(
			&review.ID, &review.UserID, &review.Status, &review.Rating,
			&review.Review, &review.Notes, &review.CreatedAt, &review.UpdatedAt,
			&review.Username, &review.Name, &review.Picture,
		)
		if err != nil {
			continue
		}

		reviews = append(reviews, map[string]interface{}{
			"id":         review.ID,
			"status":     review.Status,
			"rating":     review.Rating,
			"review":     review.Review,
			"notes":      review.Notes,
			"created_at": review.CreatedAt,
			"updated_at": review.UpdatedAt,
			"user": map[string]interface{}{
				"id":       review.UserID,
				"username": review.Username,
				"name":     review.Name,
				"picture":  review.Picture,
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"reviews": reviews,
		"count":   len(reviews),
	})
}

// GetBookStats gets community statistics for a book
func (h *BookHandler) GetBookStats(c echo.Context) error {
	bookID := c.Param("id")
	if bookID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "book ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	query := `
		SELECT 
			COUNT(CASE WHEN status = 'want_to_read' THEN 1 END) as want_to_read,
			COUNT(CASE WHEN status = 'reading' THEN 1 END) as reading,
			COUNT(CASE WHEN status = 'read' THEN 1 END) as read,
			COUNT(CASE WHEN status = 'dnf' THEN 1 END) as dnf,
			AVG(CASE WHEN rating IS NOT NULL THEN rating END) as avg_rating,
			COUNT(CASE WHEN rating IS NOT NULL THEN 1 END) as rating_count
		FROM logs
		WHERE book_id = $1 AND is_public = true
	`

	var stats struct {
		WantToRead  int
		Reading     int
		Read        int
		DNF         int
		AvgRating   *float64
		RatingCount int
	}

	err := h.DB.QueryRow(ctx, query, bookID).Scan(
		&stats.WantToRead, &stats.Reading, &stats.Read, &stats.DNF,
		&stats.AvgRating, &stats.RatingCount,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch book stats",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"want_to_read": stats.WantToRead,
		"reading":      stats.Reading,
		"read":         stats.Read,
		"dnf":          stats.DNF,
		"avg_rating":   stats.AvgRating,
		"rating_count": stats.RatingCount,
	})
}

// GetBookLists gets lists that contain this book
func (h *BookHandler) GetBookLists(c echo.Context) error {
	bookID := c.Param("id")
	if bookID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "book ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	query := `
		SELECT l.id, l.name, l.description, l.is_public, l.theme_color, l.items_count,
		       u.username, u.name as creator_name
		FROM lists l
		JOIN list_items li ON l.id = li.list_id
		JOIN users u ON l.user_id = u.id
		WHERE li.book_id = $1 AND l.is_public = true
		ORDER BY l.created_at DESC
		LIMIT 10
	`

	rows, err := h.DB.Query(ctx, query, bookID)
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
			Name        string
			Description *string
			IsPublic    bool
			ThemeColor  string
			ItemsCount  int
			Username    string
			CreatorName string
		}

		err := rows.Scan(
			&list.ID, &list.Name, &list.Description, &list.IsPublic,
			&list.ThemeColor, &list.ItemsCount, &list.Username, &list.CreatorName,
		)
		if err != nil {
			continue
		}

		lists = append(lists, map[string]interface{}{
			"id":          list.ID,
			"name":        list.Name,
			"description": list.Description,
			"is_public":   list.IsPublic,
			"theme_color": list.ThemeColor,
			"items_count": list.ItemsCount,
			"creator": map[string]interface{}{
				"username": list.Username,
				"name":     list.CreatorName,
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"lists": lists,
		"count": len(lists),
	})
}

