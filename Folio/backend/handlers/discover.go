package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"folio/api/auth"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type DiscoverHandler struct {
	DB *pgxpool.Pool
}

type RecommendationRequest struct {
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
	Genres []string `json:"genres"`
}

// GetRecommendations returns personalized book recommendations
func (h *DiscoverHandler) GetRecommendations(c echo.Context) error {
	userID := auth.GetUserID(c)
	
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	limit := 20 // Default cards per request
	if l := c.QueryParam("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	var recommendations []map[string]interface{}

	if userID != "" {
		// Personalized recommendations for logged-in users
		recommendations = h.getPersonalizedRecommendations(ctx, userID, limit)
	} else {
		// Generic trending/popular books for guests
		recommendations = h.getTrendingBooks(limit)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

func (h *DiscoverHandler) getPersonalizedRecommendations(ctx context.Context, userID string, limit int) []map[string]interface{} {
	// Get user's reading history
	readCategories := h.getUserCategories(ctx, userID)
	readAuthors := h.getUserAuthors(ctx, userID)
	
	// Fetch recommendations from Google Books based on user preferences
	recommendations := []map[string]interface{}{}
	
	// Strategy 1: Books in similar categories (60%)
	if len(readCategories) > 0 {
		categoryBooks := h.fetchBooksByCategory(readCategories, int(float64(limit)*0.6))
		recommendations = append(recommendations, categoryBooks...)
	}
	
	// Strategy 2: Books by similar authors (20%)
	if len(readAuthors) > 0 {
		authorBooks := h.fetchBooksByAuthor(readAuthors, int(float64(limit)*0.2))
		recommendations = append(recommendations, authorBooks...)
	}
	
	// Strategy 3: Trending/Popular books (20%)
	trendingBooks := h.getTrendingBooks(int(float64(limit) * 0.2))
	recommendations = append(recommendations, trendingBooks...)
	
	// Filter out books user has already logged
	recommendations = h.filterAlreadyLogged(ctx, userID, recommendations)
	
	// Shuffle for variety
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(recommendations), func(i, j int) {
		recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
	})
	
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}
	
	return recommendations
}

func (h *DiscoverHandler) getUserCategories(ctx context.Context, userID string) []string {
	query := `
		SELECT DISTINCT unnest(b.categories) as category
		FROM logs l
		JOIN books b ON l.book_id = b.id
		WHERE l.user_id = $1 
		  AND l.rating >= 4
		  AND b.categories IS NOT NULL
		LIMIT 10
	`
	
	rows, err := h.DB.Query(ctx, query, userID)
	if err != nil {
		return []string{}
	}
	defer rows.Close()
	
	categories := []string{}
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err == nil {
			categories = append(categories, category)
		}
	}
	
	return categories
}

func (h *DiscoverHandler) getUserAuthors(ctx context.Context, userID string) []string {
	query := `
		SELECT DISTINCT unnest(b.authors) as author
		FROM logs l
		JOIN books b ON l.book_id = b.id
		WHERE l.user_id = $1 
		  AND l.rating >= 4
		  AND b.authors IS NOT NULL
		LIMIT 5
	`
	
	rows, err := h.DB.Query(ctx, query, userID)
	if err != nil {
		return []string{}
	}
	defer rows.Close()
	
	authors := []string{}
	for rows.Next() {
		var author string
		if err := rows.Scan(&author); err == nil {
			authors = append(authors, author)
		}
	}
	
	return authors
}

func (h *DiscoverHandler) fetchBooksByCategory(categories []string, limit int) []map[string]interface{} {
	if len(categories) == 0 {
		return []map[string]interface{}{}
	}
	
	// Pick a random category for variety
	category := categories[rand.Intn(len(categories))]
	return h.searchGoogleBooks(fmt.Sprintf("subject:%s", category), limit)
}

func (h *DiscoverHandler) fetchBooksByAuthor(authors []string, limit int) []map[string]interface{} {
	if len(authors) == 0 {
		return []map[string]interface{}{}
	}
	
	// Pick a random author
	author := authors[rand.Intn(len(authors))]
	return h.searchGoogleBooks(fmt.Sprintf("inauthor:%s", author), limit)
}

func (h *DiscoverHandler) getTrendingBooks(limit int) []map[string]interface{} {
	// Popular/trending searches
	queries := []string{
		"bestseller 2024",
		"popular fiction",
		"award winning books",
		"recommended reads",
		"top rated books",
	}
	
	query := queries[rand.Intn(len(queries))]
	return h.searchGoogleBooks(query, limit)
}

func (h *DiscoverHandler) searchGoogleBooks(query string, maxResults int) []map[string]interface{} {
	apiKey := getEnv("GOOGLE_BOOKS_API_KEY", "")
	url := fmt.Sprintf(
		"https://www.googleapis.com/books/v1/volumes?q=%s&maxResults=%d&orderBy=relevance",
		strings.ReplaceAll(query, " ", "+"),
		maxResults,
	)
	
	if apiKey != "" {
		url += "&key=" + apiKey
	}
	
	resp, err := http.Get(url)
	if err != nil {
		return []map[string]interface{}{}
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []map[string]interface{}{}
	}
	
	var result struct {
		Items []struct {
			ID         string `json:"id"`
			VolumeInfo struct {
				Title         string   `json:"title"`
				Authors       []string `json:"authors"`
				Description   string   `json:"description"`
				Categories    []string `json:"categories"`
				PageCount     int      `json:"pageCount"`
				PublishedDate string   `json:"publishedDate"`
				ImageLinks    struct {
					Thumbnail      string `json:"thumbnail"`
					SmallThumbnail string `json:"smallThumbnail"`
				} `json:"imageLinks"`
				AverageRating float64 `json:"averageRating"`
				RatingsCount  int     `json:"ratingsCount"`
			} `json:"volumeInfo"`
		} `json:"items"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return []map[string]interface{}{}
	}
	
	books := []map[string]interface{}{}
	for _, item := range result.Items {
		// Get high quality cover image
		coverURL := item.VolumeInfo.ImageLinks.Thumbnail
		if coverURL != "" {
			// Upgrade to larger image
			coverURL = strings.Replace(coverURL, "zoom=1", "zoom=2", 1)
		}
		
		book := map[string]interface{}{
			"id":             item.ID,
			"title":          item.VolumeInfo.Title,
			"authors":        item.VolumeInfo.Authors,
			"description":    item.VolumeInfo.Description,
			"categories":     item.VolumeInfo.Categories,
			"cover_url":      coverURL,
			"page_count":     item.VolumeInfo.PageCount,
			"published_date": item.VolumeInfo.PublishedDate,
			"rating":         item.VolumeInfo.AverageRating,
			"ratings_count":  item.VolumeInfo.RatingsCount,
		}
		books = append(books, book)
	}
	
	return books
}

func (h *DiscoverHandler) filterAlreadyLogged(ctx context.Context, userID string, books []map[string]interface{}) []map[string]interface{} {
	if len(books) == 0 {
		return books
	}
	
	// Get list of book IDs user has already logged
	query := `SELECT book_id FROM logs WHERE user_id = $1`
	rows, err := h.DB.Query(ctx, query, userID)
	if err != nil {
		return books
	}
	defer rows.Close()
	
	loggedBooks := make(map[string]bool)
	for rows.Next() {
		var bookID string
		if err := rows.Scan(&bookID); err == nil {
			loggedBooks[bookID] = true
		}
	}
	
	// Filter out logged books
	filtered := []map[string]interface{}{}
	for _, book := range books {
		bookID, ok := book["id"].(string)
		if !ok || !loggedBooks[bookID] {
			filtered = append(filtered, book)
		}
	}
	
	return filtered
}

// RecordSwipe records user's swipe action for future recommendations
func (h *DiscoverHandler) RecordSwipe(c echo.Context) error {
	userID := auth.GetUserID(c)
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}
	
	var req struct {
		BookID string `json:"book_id"`
		Action string `json:"action"` // "like" or "pass"
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}
	
	// Could store swipe history in a separate table for better recommendations
	// For now, just acknowledge
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"action":  req.Action,
	})
}

