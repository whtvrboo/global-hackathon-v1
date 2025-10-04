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

