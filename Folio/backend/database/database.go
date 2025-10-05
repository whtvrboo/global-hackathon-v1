package database

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

//go:embed seeds/*.sql
var seedsFS embed.FS

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// LoadConfigFromEnv loads database configuration from environment variables
func LoadConfigFromEnv() *Config {
	return &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "folio_user"),
		Password: getEnv("DB_PASSWORD", "folio_password"),
		Database: getEnv("DB_NAME", "folio_db"),
	}
}

// GetConnectionString returns the PostgreSQL connection string
func (c *Config) GetConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.Database,
	)
}

// NewPool creates a new connection pool
func NewPool(ctx context.Context, config *Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(config.GetConnectionString())
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	// Configure pool settings
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("✓ Database connection pool established")
	return pool, nil
}

// RunMigrations applies database migrations automatically
func RunMigrations(config *Config) error {
	log.Println("Running database migrations...")

	// Create source driver from embedded files
	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithSourceInstance(
		"iofs",
		sourceDriver,
		config.GetConnectionString(),
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Check if database is in dirty state before running migrations
	version, dirty, err := m.Version()
	if dirty {
		log.Printf("⚠ Database is in dirty state at version %d, attempting to fix...", version)
		// Force the version to clean state - this allows the migration to retry
		if err := m.Force(int(version)); err != nil {
			return fmt.Errorf("failed to force version: %w", err)
		}
		log.Printf("✓ Forced database to clean state at version %d", version)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	version, dirty, err = m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if dirty {
		log.Printf("⚠ Database is in dirty state at version %d", version)
	} else if err == migrate.ErrNilVersion {
		log.Println("✓ Database migrations completed (no version yet)")
	} else {
		log.Printf("✓ Database migrations completed (version: %d)", version)
	}

	// Optionally run seed data after successful migrations
	if isTruthy(getEnv("DB_AUTO_SEED", "false")) {
        if err := runSeeds(config); err != nil {
            return err
        }
    }

	return nil
}

// WaitForDatabase waits for the database to be ready
func WaitForDatabase(ctx context.Context, config *Config, maxRetries int) error {
	log.Println("Waiting for database to be ready...")

	for i := 0; i < maxRetries; i++ {
		pool, err := pgxpool.New(ctx, config.GetConnectionString())
		if err == nil {
			if err := pool.Ping(ctx); err == nil {
				pool.Close()
				log.Println("✓ Database is ready")
				return nil
			}
			pool.Close()
		}

		log.Printf("Database not ready, retrying in 2s... (%d/%d)", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("database not ready after %d attempts", maxRetries)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// isTruthy returns true if the string looks like a true value
func isTruthy(v string) bool {
    switch strings.ToLower(strings.TrimSpace(v)) {
    case "1", "t", "true", "y", "yes", "on":
        return true
    default:
        return false
    }
}


// runSeeds loads and executes seed SQL based on environment variables.
// Controls:
// - DB_AUTO_SEED=true to enable seeding (checked in RunMigrations)
// - DB_SEED_DATA one of: "books", "full", "demo" (default: "full")
func runSeeds(config *Config) error {
    mode := strings.ToLower(strings.TrimSpace(getEnv("DB_SEED_DATA", "full")))

    var seedFile string
    switch mode {
    case "books":
        seedFile = "books_seed.sql"
    case "demo":
        seedFile = "demo_seed.sql"
    case "full":
        fallthrough
    default:
        seedFile = "seed.sql"
    }

    // Read the chosen seed file from embedded FS
    seedSQLBytes, err := seedsFS.ReadFile("seeds/" + seedFile)
    if err != nil {
        return fmt.Errorf("failed to read seed file %s: %w", seedFile, err)
    }
    seedSQL := string(seedSQLBytes)

    // Expand simple psql include used in seed.sql: \i books_seed.sql
    // This allows executing via pgx without psql.
    if strings.Contains(seedSQL, "\\i ") || strings.Contains(seedSQL, "\\i\t") {
        // Only handle includes of books_seed.sql which is the only one present
        includeTarget := "books_seed.sql"
        inclBytes, inclErr := seedsFS.ReadFile("seeds/" + includeTarget)
        if inclErr != nil {
            return fmt.Errorf("failed to read included seed file %s: %w", includeTarget, inclErr)
        }
        // Replace any line starting with \i books_seed.sql with file contents
        lines := strings.Split(seedSQL, "\n")
        for i, line := range lines {
            trimmed := strings.TrimSpace(line)
            if strings.HasPrefix(trimmed, "\\i ") || strings.HasPrefix(trimmed, "\\i\t") {
                if strings.Contains(trimmed, includeTarget) {
                    lines[i] = string(inclBytes)
                }
            }
        }
        seedSQL = strings.Join(lines, "\n")
    }

    ctx := context.Background()
    pool, err := pgxpool.New(ctx, config.GetConnectionString())
    if err != nil {
        return fmt.Errorf("failed to connect for seeding: %w", err)
    }
    defer pool.Close()

    log.Printf("Running database seed: %s (mode=%s)", seedFile, mode)
    // Execute as a single multi-statement script
    if _, err := pool.Exec(ctx, seedSQL); err != nil {
        return fmt.Errorf("failed to execute seed script %s: %w", seedFile, err)
    }
    log.Println("✓ Database seed completed")
    return nil
}

