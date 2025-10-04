package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"folio/api/auth"
	"folio/api/database"
	"folio/api/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	DB *pgxpool.Pool
}

func main() {
	// Load configuration
	dbConfig := database.LoadConfigFromEnv()
	port := getEnv("SERVER_PORT", "8080")

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Wait for database to be ready
	if err := database.WaitForDatabase(ctx, dbConfig, 15); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations(dbConfig); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create connection pool
	pool, err := database.NewPool(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}
	defer pool.Close()

	// Create app instance
	app := &App{DB: pool}

	// Setup Echo
	e := echo.New()
	e.HidePort = true
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	setupRoutes(e, app)

	// Start server
	log.Printf("🚀 Server starting on port %s", port)
	
	// Graceful shutdown
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

func setupRoutes(e *echo.Echo, app *App) {
	// Initialize OAuth
	auth.InitOAuth()

	// Create handlers
	authHandler := &handlers.AuthHandler{DB: app.DB}
	bookHandler := &handlers.BookHandler{DB: app.DB}
	logHandler := &handlers.LogHandler{DB: app.DB}
	socialHandler := &handlers.SocialHandler{DB: app.DB}
	discoverHandler := &handlers.DiscoverHandler{DB: app.DB}
	guestHandler := &handlers.GuestHandler{DB: app.DB}
	listHandler := &handlers.ListHandler{DB: app.DB}

	// API routes
	api := e.Group("/api")

	// Health check endpoint
	api.GET("/health", app.healthCheck)

	// Public endpoints
	api.GET("/auth/google", authHandler.GoogleLogin)
	api.GET("/auth/google/callback", authHandler.GoogleCallback)
	api.GET("/auth/google/convert", authHandler.ConvertGuestToUser)
	api.POST("/auth/guest", guestHandler.CreateGuestUser)
	api.GET("/search", bookHandler.SearchBooks)
	api.GET("/books/:id", bookHandler.GetBook)
	api.GET("/discover", discoverHandler.GetRecommendations)
	api.GET("/discover/lists", discoverHandler.GetTrendingLists)

	// Protected endpoints
	protected := api.Group("", auth.JWTMiddleware)
	protected.GET("/me", authHandler.GetMe)
	protected.PUT("/me/profile", authHandler.UpdateProfile)
	protected.GET("/guest/me", guestHandler.GetGuestUser)
	protected.POST("/logs", logHandler.CreateLog)
	protected.GET("/logs/:id", logHandler.GetSingleLog)
	protected.GET("/users/:username/logs", logHandler.GetUserLogs)
	protected.GET("/feed", logHandler.GetFeed)
	protected.POST("/users/:username/follow", socialHandler.FollowUser)
	protected.DELETE("/users/:username/follow", socialHandler.UnfollowUser)
	protected.POST("/discover/swipe", discoverHandler.RecordSwipe)
	
	// Like and comment endpoints
	protected.POST("/logs/:id/like", socialHandler.ToggleLike)
	protected.GET("/logs/:id/comments", socialHandler.GetLogComments)
	protected.POST("/logs/:id/comments", socialHandler.CreateComment)
	protected.DELETE("/comments/:commentId", socialHandler.DeleteComment)
	
	// List endpoints
	protected.POST("/lists", listHandler.CreateList)
	protected.GET("/users/:username/lists", listHandler.GetUserLists)
	protected.GET("/lists/:id", listHandler.GetList)
	protected.PUT("/lists/:id", listHandler.UpdateList)
	protected.DELETE("/lists/:id", listHandler.DeleteList)
	protected.POST("/lists/:id/items", listHandler.AddBookToList)
	protected.PUT("/lists/:id/items/:itemId/order", listHandler.UpdateListItemOrder)
	protected.DELETE("/lists/:id/items/:itemId", listHandler.RemoveBookFromList)
	
	// List social features
	protected.POST("/lists/:id/like", listHandler.LikeList)
	protected.DELETE("/lists/:id/like", listHandler.UnlikeList)
	protected.GET("/lists/:id/comments", listHandler.GetListComments)
	protected.POST("/lists/:id/comments", listHandler.AddListComment)
}

// healthCheck performs a database query and returns system status
func (app *App) healthCheck(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	// Test database connection
	var result int
	err := app.DB.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		log.Printf("Health check failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Database connection failed",
			"error":   err.Error(),
		})
	}

	// Get database stats
	stats := app.DB.Stat()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"message": "All systems operational",
		"database": map[string]interface{}{
			"connected":      true,
			"total_conns":    stats.TotalConns(),
			"acquired_conns": stats.AcquiredConns(),
			"idle_conns":     stats.IdleConns(),
		},
		"timestamp": time.Now().UTC(),
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

