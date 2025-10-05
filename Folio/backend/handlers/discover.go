package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"folio/api/auth"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"sort"
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

type RecommendationReason struct {
	Type        string `json:"type"`        // "category", "author", "friend", "trending", "similar"
	Value       string `json:"value"`       // The specific category, author, etc.
	Confidence  int    `json:"confidence"`  // 1-100 confidence score
	Description string `json:"description"` // Human-readable explanation
}

type PersonalizedRecommendation struct {
	Book   map[string]interface{} `json:"book"`
	Reason RecommendationReason   `json:"reason"`
	Score  float64                `json:"score"`
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

	var recommendations []PersonalizedRecommendation

	if userID != "" {
		// Personalized recommendations for logged-in users
		recommendations = h.getPersonalizedRecommendations(ctx, userID, limit)
	} else {
		// Generic trending/popular books for guests
		recommendations = h.getTrendingRecommendations(limit)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

// GetTrendingLists returns popular public lists
func (h *DiscoverHandler) GetTrendingLists(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	limit := 10
	if l := c.QueryParam("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	query := `
		SELECT l.id, l.user_id, l.name, l.description, l.is_public, l.items_count, l.created_at, l.updated_at,
		       u.username, u.name as user_name, u.picture
		FROM lists l
		JOIN users u ON l.user_id = u.id
		WHERE l.is_public = true AND l.items_count > 0
		ORDER BY l.items_count DESC, l.created_at DESC
		LIMIT $1
	`

	rows, err := h.DB.Query(ctx, query, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch trending lists",
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
			Username    string
			UserName    string
			Picture     *string
		}

		err := rows.Scan(
			&list.ID, &list.UserID, &list.Name, &list.Description, &list.IsPublic,
			&list.ItemsCount, &list.CreatedAt, &list.UpdatedAt,
			&list.Username, &list.UserName, &list.Picture,
		)
		if err != nil {
			continue
		}

		lists = append(lists, map[string]interface{}{
			"id":          list.ID,
			"name":        list.Name,
			"description": list.Description,
			"is_public":   list.IsPublic,
			"items_count": list.ItemsCount,
			"created_at":  list.CreatedAt,
			"updated_at":  list.UpdatedAt,
			"user": map[string]interface{}{
				"id":       list.UserID,
				"username": list.Username,
				"name":     list.UserName,
				"picture":  list.Picture,
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"lists": lists,
		"count": len(lists),
	})
}

func (h *DiscoverHandler) getPersonalizedRecommendations(ctx context.Context, userID string, limit int) []PersonalizedRecommendation {
	// Get user's reading profile
	userProfile := h.getUserReadingProfile(ctx, userID)
	
	// Get all possible recommendations with reasons
	allRecommendations := []PersonalizedRecommendation{}
	
	// Strategy 1: Books similar to highly-rated books (40%)
	similarBooks := h.getSimilarBooksRecommendations(ctx, userID, userProfile, int(float64(limit)*0.4))
	allRecommendations = append(allRecommendations, similarBooks...)
	
	// Strategy 2: Books from friends' favorites (25%)
	friendBooks := h.getFriendRecommendations(ctx, userID, int(float64(limit)*0.25))
	allRecommendations = append(allRecommendations, friendBooks...)
	
	// Strategy 3: Books in favorite categories (20%)
	categoryBooks := h.getCategoryRecommendations(ctx, userID, userProfile, int(float64(limit)*0.2))
	allRecommendations = append(allRecommendations, categoryBooks...)
	
	// Strategy 4: Trending books with high ratings (10%)
	trendingBooks := h.getTrendingRecommendations(int(float64(limit) * 0.1))
	allRecommendations = append(allRecommendations, trendingBooks...)
	
	// Strategy 5: Serendipity - random high-quality books (5%)
	serendipityBooks := h.getSerendipityRecommendations(int(float64(limit) * 0.05))
	allRecommendations = append(allRecommendations, serendipityBooks...)
	
	// Remove duplicates and filter already logged books
	allRecommendations = h.deduplicateAndFilterRecommendations(ctx, userID, allRecommendations)
	
	// Sort by score (personalization strength)
	sort.Slice(allRecommendations, func(i, j int) bool {
		return allRecommendations[i].Score > allRecommendations[j].Score
	})
	
	// Take top recommendations
	if len(allRecommendations) > limit {
		allRecommendations = allRecommendations[:limit]
	}
	
	// Add some randomness to the top results to avoid always showing the same order
	rand.Seed(time.Now().UnixNano())
	if len(allRecommendations) > 3 {
		// Shuffle the top 3 positions slightly
		topThree := allRecommendations[:3]
		rand.Shuffle(len(topThree), func(i, j int) {
			topThree[i], topThree[j] = topThree[j], topThree[i]
		})
		allRecommendations = append(topThree, allRecommendations[3:]...)
	}
	
	return allRecommendations
}

// getUserReadingProfile creates a comprehensive profile of user's reading preferences
func (h *DiscoverHandler) getUserReadingProfile(ctx context.Context, userID string) map[string]interface{} {
	profile := map[string]interface{}{
		"favorite_categories": []string{},
		"favorite_authors":    []string{},
		"favorite_books":      []string{},
		"reading_patterns":    map[string]interface{}{},
		"review_keywords":     []string{},
	}
	
	// Get favorite categories (from highly-rated books)
	query := `
		SELECT unnest(b.categories) as category, COUNT(*) as count, AVG(l.rating) as avg_rating
		FROM logs l
		JOIN books b ON l.book_id = b.id
		WHERE l.user_id = $1 AND l.rating >= 4 AND b.categories IS NOT NULL
		GROUP BY category
		ORDER BY count DESC, avg_rating DESC
		LIMIT 10
	`
	
	rows, err := h.DB.Query(ctx, query, userID)
	if err == nil {
		defer rows.Close()
		categories := []string{}
		for rows.Next() {
			var category string
			var count int
			var avgRating float64
			if err := rows.Scan(&category, &count, &avgRating); err == nil {
				categories = append(categories, category)
			}
		}
		profile["favorite_categories"] = categories
	}
	
	// Get favorite authors
	query = `
		SELECT unnest(b.authors) as author, COUNT(*) as count, AVG(l.rating) as avg_rating
		FROM logs l
		JOIN books b ON l.book_id = b.id
		WHERE l.user_id = $1 AND l.rating >= 4 AND b.authors IS NOT NULL
		GROUP BY author
		ORDER BY count DESC, avg_rating DESC
		LIMIT 5
	`
	
	rows, err = h.DB.Query(ctx, query, userID)
	if err == nil {
		defer rows.Close()
		authors := []string{}
		for rows.Next() {
			var author string
			var count int
			var avgRating float64
			if err := rows.Scan(&author, &count, &avgRating); err == nil {
				authors = append(authors, author)
			}
		}
		profile["favorite_authors"] = authors
	}
	
	// Get favorite books (highest rated)
	query = `
		SELECT l.book_id, l.rating
		FROM logs l
		WHERE l.user_id = $1 AND l.rating >= 4
		ORDER BY l.rating DESC, l.created_at DESC
		LIMIT 5
	`
	
	rows, err = h.DB.Query(ctx, query, userID)
	if err == nil {
		defer rows.Close()
		books := []string{}
		for rows.Next() {
			var bookID string
			var rating int
			if err := rows.Scan(&bookID, &rating); err == nil {
				books = append(books, bookID)
			}
		}
		profile["favorite_books"] = books
	}
	
	// Extract keywords from reviews
	query = `
		SELECT l.review
		FROM logs l
		WHERE l.user_id = $1 AND l.review IS NOT NULL AND l.review != ''
		LIMIT 20
	`
	
	rows, err = h.DB.Query(ctx, query, userID)
	if err == nil {
		defer rows.Close()
		keywords := []string{}
		for rows.Next() {
			var review string
			if err := rows.Scan(&review); err == nil {
				// Extract meaningful keywords from review
				reviewKeywords := h.extractKeywordsFromReview(review)
				keywords = append(keywords, reviewKeywords...)
			}
		}
		profile["review_keywords"] = keywords
	}
	
	return profile
}

// extractKeywordsFromReview extracts meaningful keywords from a review
func (h *DiscoverHandler) extractKeywordsFromReview(review string) []string {
	// Simple keyword extraction - in a real system, you'd use NLP
	keywords := []string{}
	
	// Remove common words and extract meaningful terms
	commonWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "is": true, "are": true, "was": true, "were": true,
		"be": true, "been": true, "have": true, "has": true, "had": true, "do": true,
		"does": true, "did": true, "will": true, "would": true, "could": true, "should": true,
		"this": true, "that": true, "these": true, "those": true, "i": true, "you": true,
		"he": true, "she": true, "it": true, "we": true, "they": true, "me": true,
		"him": true, "her": true, "us": true, "them": true, "my": true, "your": true,
		"his": true, "its": true, "our": true, "their": true,
	}
	
	// Split into words and filter
	words := regexp.MustCompile(`\b\w+\b`).FindAllString(strings.ToLower(review), -1)
	for _, word := range words {
		if len(word) > 3 && !commonWords[word] {
			keywords = append(keywords, word)
		}
	}
	
	// Return unique keywords, limited to 10
	seen := make(map[string]bool)
	result := []string{}
	for _, keyword := range keywords {
		if !seen[keyword] && len(result) < 10 {
			seen[keyword] = true
			result = append(result, keyword)
		}
	}
	
	return result
}

// getSimilarBooksRecommendations finds books similar to user's favorites
func (h *DiscoverHandler) getSimilarBooksRecommendations(ctx context.Context, userID string, profile map[string]interface{}, limit int) []PersonalizedRecommendation {
	recommendations := []PersonalizedRecommendation{}
	
	favoriteBooks, ok := profile["favorite_books"].([]string)
	if !ok || len(favoriteBooks) == 0 {
		return recommendations
	}
	
	// For each favorite book, find similar books
	for _, bookID := range favoriteBooks {
		similarBooks := h.findSimilarBooks(ctx, bookID, limit/len(favoriteBooks)+1)
		for _, book := range similarBooks {
			recommendations = append(recommendations, PersonalizedRecommendation{
				Book: book,
				Reason: RecommendationReason{
					Type:        "similar",
					Value:       bookID,
					Confidence:  85,
					Description: "Similar to a book you loved",
				},
				Score: 0.9,
			})
		}
	}
	
	return recommendations
}

// getFriendRecommendations gets books from friends' favorites
func (h *DiscoverHandler) getFriendRecommendations(ctx context.Context, userID string, limit int) []PersonalizedRecommendation {
	recommendations := []PersonalizedRecommendation{}
	
	query := `
		SELECT DISTINCT l.book_id, COUNT(*) as friend_count
		FROM logs l
		JOIN followers f ON l.user_id = f.following_id
		WHERE f.follower_id = $1 
		  AND l.rating >= 4 
		  AND l.book_id NOT IN (SELECT book_id FROM logs WHERE user_id = $1)
		GROUP BY l.book_id
		ORDER BY friend_count DESC
		LIMIT $2
	`
	
	rows, err := h.DB.Query(ctx, query, userID, limit)
	if err != nil {
		return recommendations
	}
	defer rows.Close()
	
	for rows.Next() {
		var bookID string
		var friendCount int
		if err := rows.Scan(&bookID, &friendCount); err != nil {
			continue
		}
		
		// Get book details
		book := h.getBookDetails(ctx, bookID)
		if book != nil {
			recommendations = append(recommendations, PersonalizedRecommendation{
				Book: book,
				Reason: RecommendationReason{
					Type:        "friend",
					Value:       fmt.Sprintf("%d friends", friendCount),
					Confidence:  75,
					Description: fmt.Sprintf("Loved by %d of your friends", friendCount),
				},
				Score: 0.8,
			})
		}
	}
	
	return recommendations
}

// getCategoryRecommendations gets books in user's favorite categories
func (h *DiscoverHandler) getCategoryRecommendations(ctx context.Context, userID string, profile map[string]interface{}, limit int) []PersonalizedRecommendation {
	recommendations := []PersonalizedRecommendation{}
	
	favoriteCategories, ok := profile["favorite_categories"].([]string)
	if !ok || len(favoriteCategories) == 0 {
		return recommendations
	}
	
	// Pick a random favorite category for variety
	category := favoriteCategories[rand.Intn(len(favoriteCategories))]
	books := h.searchGoogleBooks(fmt.Sprintf("subject:%s", category), limit)
	
	for _, book := range books {
		recommendations = append(recommendations, PersonalizedRecommendation{
			Book: book,
			Reason: RecommendationReason{
				Type:        "category",
				Value:       category,
				Confidence:  70,
				Description: fmt.Sprintf("In %s, a genre you love", category),
			},
			Score: 0.7,
		})
	}
	
	return recommendations
}

// getTrendingRecommendations gets trending books
func (h *DiscoverHandler) getTrendingRecommendations(limit int) []PersonalizedRecommendation {
	recommendations := []PersonalizedRecommendation{}
	
	books := h.getTrendingBooks(limit)
	for _, book := range books {
		recommendations = append(recommendations, PersonalizedRecommendation{
			Book: book,
			Reason: RecommendationReason{
				Type:        "trending",
				Value:       "popular",
				Confidence:  60,
				Description: "Trending on Folio",
			},
			Score: 0.6,
		})
	}
	
	return recommendations
}

// getSerendipityRecommendations gets random high-quality books
func (h *DiscoverHandler) getSerendipityRecommendations(limit int) []PersonalizedRecommendation {
	recommendations := []PersonalizedRecommendation{}
	
	// Get random books with good ratings
	books := h.getRandomHighQualityBooks(limit)
	for _, book := range books {
		recommendations = append(recommendations, PersonalizedRecommendation{
			Book: book,
			Reason: RecommendationReason{
				Type:        "serendipity",
				Value:       "discovery",
				Confidence:  50,
				Description: "Something different you might love",
			},
			Score: 0.5,
		})
	}
	
	return recommendations
}

// deduplicateAndFilterRecommendations removes duplicates and already logged books
func (h *DiscoverHandler) deduplicateAndFilterRecommendations(ctx context.Context, userID string, recommendations []PersonalizedRecommendation) []PersonalizedRecommendation {
	seen := make(map[string]bool)
	filtered := []PersonalizedRecommendation{}
	
	// Get already logged books
	loggedBooks := h.getLoggedBookIDs(ctx, userID)
	
	for _, rec := range recommendations {
		bookID, ok := rec.Book["id"].(string)
		if !ok || seen[bookID] || loggedBooks[bookID] {
			continue
		}
		seen[bookID] = true
		filtered = append(filtered, rec)
	}
	
	return filtered
}

// Helper functions
func (h *DiscoverHandler) findSimilarBooks(ctx context.Context, bookID string, limit int) []map[string]interface{} {
	// Get book details first
	book := h.getBookDetails(ctx, bookID)
	if book == nil {
		return []map[string]interface{}{}
	}
	
	// Find books with similar categories or authors
	categories, _ := book["categories"].([]string)
	authors, _ := book["authors"].([]string)
	
	var searchQuery string
	if len(categories) > 0 {
		searchQuery = fmt.Sprintf("subject:%s", categories[0])
	} else if len(authors) > 0 {
		searchQuery = fmt.Sprintf("inauthor:%s", authors[0])
	} else {
		return []map[string]interface{}{}
	}
	
	return h.searchGoogleBooks(searchQuery, limit)
}

func (h *DiscoverHandler) getBookDetails(ctx context.Context, bookID string) map[string]interface{} {
	query := `
		SELECT id, title, author, description, cover_url, published_date, pages, categories
		FROM books
		WHERE id = $1
	`
	
	var book struct {
		ID            string
		Title         string
		Author        string
		Description   string
		CoverURL      string
		PublishedDate string
		Pages         int
		Categories    string
	}
	
	err := h.DB.QueryRow(ctx, query, bookID).Scan(
		&book.ID, &book.Title, &book.Author, &book.Description,
		&book.CoverURL, &book.PublishedDate, &book.Pages, &book.Categories,
	)
	
	if err != nil {
		return nil
	}
	
	var categories []string
	if book.Categories != "" {
		json.Unmarshal([]byte(book.Categories), &categories)
	}
	
	return map[string]interface{}{
		"id":             book.ID,
		"title":          book.Title,
		"authors":        []string{book.Author},
		"description":    book.Description,
		"cover_url":      book.CoverURL,
		"published_date": book.PublishedDate,
		"pages":          book.Pages,
		"categories":     categories,
	}
}

func (h *DiscoverHandler) getLoggedBookIDs(ctx context.Context, userID string) map[string]bool {
	loggedBooks := make(map[string]bool)
	
	query := `SELECT book_id FROM logs WHERE user_id = $1`
	rows, err := h.DB.Query(ctx, query, userID)
	if err != nil {
		return loggedBooks
	}
	defer rows.Close()
	
	for rows.Next() {
		var bookID string
		if err := rows.Scan(&bookID); err == nil {
			loggedBooks[bookID] = true
		}
	}
	
	return loggedBooks
}

func (h *DiscoverHandler) getRandomHighQualityBooks(limit int) []map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	query := `
		SELECT b.id, b.title, b.author, b.description, b.cover_url, b.published_date, b.pages, b.categories,
		       AVG(l.rating) as avg_rating, COUNT(l.id) as log_count
		FROM books b
		LEFT JOIN logs l ON b.id = l.book_id AND l.status = 'read'
		GROUP BY b.id, b.title, b.author, b.description, b.cover_url, b.published_date, b.pages, b.categories
		HAVING COUNT(l.id) > 0 AND AVG(l.rating) >= 3.5
		ORDER BY RANDOM()
		LIMIT $1
	`
	
	rows, err := h.DB.Query(ctx, query, limit)
	if err != nil {
		return []map[string]interface{}{}
	}
	defer rows.Close()
	
	var books []map[string]interface{}
	for rows.Next() {
		var book struct {
			ID            string
			Title         string
			Author        string
			Description   string
			CoverURL      string
			PublishedDate string
			Pages         int
			Categories    string
			AvgRating     *float64
			LogCount      int
		}
		
		err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.Description,
			&book.CoverURL, &book.PublishedDate, &book.Pages, &book.Categories,
			&book.AvgRating, &book.LogCount,
		)
		if err != nil {
			continue
		}
		
		var categories []string
		if book.Categories != "" {
			json.Unmarshal([]byte(book.Categories), &categories)
		}
		
		rating := 0.0
		if book.AvgRating != nil {
			rating = *book.AvgRating
		}
		
		books = append(books, map[string]interface{}{
			"id":             book.ID,
			"title":          book.Title,
			"authors":        []string{book.Author},
			"description":    book.Description,
			"cover_url":      book.CoverURL,
			"published_date": book.PublishedDate,
			"pages":          book.Pages,
			"categories":     categories,
			"rating":         rating,
			"log_count":      book.LogCount,
		})
	}
	
	return books
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get popular books from our local database
	query := `
		SELECT b.id, b.title, b.authors, b.description, b.cover_url, b.published_date, b.page_count, b.categories,
		       COUNT(l.id) as log_count,
		       AVG(l.rating) as avg_rating
		FROM books b
		LEFT JOIN logs l ON b.id = l.book_id AND l.status = 'read'
		GROUP BY b.id, b.title, b.authors, b.description, b.cover_url, b.published_date, b.page_count, b.categories
		HAVING COUNT(l.id) > 0
		ORDER BY log_count DESC, avg_rating DESC
		LIMIT $1
	`

	rows, err := h.DB.Query(ctx, query, limit)
	if err != nil {
		// Fallback to any books if no logs exist
		return h.getRandomBooks(ctx, limit)
	}
	defer rows.Close()

	var books []map[string]interface{}
	for rows.Next() {
		var book struct {
			ID            string
			Title         string
			Authors       []string
			Description   string
			CoverURL      string
			PublishedDate string
			PageCount     int
			Categories    []string
			LogCount      int
			AvgRating     *float64
		}

		err := rows.Scan(
			&book.ID, &book.Title, &book.Authors, &book.Description,
			&book.CoverURL, &book.PublishedDate, &book.PageCount, &book.Categories,
			&book.LogCount, &book.AvgRating,
		)
		if err != nil {
			continue
		}

		rating := 0.0
		if book.AvgRating != nil {
			rating = *book.AvgRating
		}

		books = append(books, map[string]interface{}{
			"id":            book.ID,
			"title":         book.Title,
			"authors":       book.Authors,
			"description":   book.Description,
			"cover_url":     book.CoverURL,
			"published_date": book.PublishedDate,
			"page_count":    book.PageCount,
			"categories":    book.Categories,
			"rating":        rating,
			"log_count":     book.LogCount,
		})
	}

	// If no books with logs, return random books
	if len(books) == 0 {
		return h.getRandomBooks(ctx, limit)
	}

	return books
}

// getRandomBooks returns random books from our database
func (h *DiscoverHandler) getRandomBooks(ctx context.Context, limit int) []map[string]interface{} {
	query := `
		SELECT id, title, authors, description, cover_url, published_date, page_count, categories
		FROM books
		ORDER BY RANDOM()
		LIMIT $1
	`

	rows, err := h.DB.Query(ctx, query, limit)
	if err != nil {
		return []map[string]interface{}{}
	}
	defer rows.Close()

	var books []map[string]interface{}
	for rows.Next() {
		var book struct {
			ID            string
			Title         string
			Authors       []string
			Description   string
			CoverURL      string
			PublishedDate string
			PageCount     int
			Categories    []string
		}

		err := rows.Scan(
			&book.ID, &book.Title, &book.Authors, &book.Description,
			&book.CoverURL, &book.PublishedDate, &book.PageCount, &book.Categories,
		)
		if err != nil {
			continue
		}

		books = append(books, map[string]interface{}{
			"id":            book.ID,
			"title":         book.Title,
			"authors":       book.Authors,
			"description":   book.Description,
			"cover_url":     book.CoverURL,
			"published_date": book.PublishedDate,
			"page_count":    book.PageCount,
			"categories":    book.Categories,
			"rating":        0.0,
		})
	}

	return books
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

