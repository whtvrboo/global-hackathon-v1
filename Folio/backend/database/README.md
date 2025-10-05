# Backend Database Package

This package provides an elegant, production-ready database connection and migration management system for Folio.

## Features

✨ **Automatic Migration on Startup** - No manual steps required  
✨ **Embedded Migrations** - SQL files are compiled into the Go binary  
✨ **Connection Pooling** - Optimized for performance with pgxpool  
✨ **Health Checks** - Built-in database connectivity verification  
✨ **Retry Logic** - Automatically waits for database to be ready  
✨ **Zero External Dependencies** - Everything runs in the Go binary

## Architecture

### Why Embedded Migrations?

Traditional migration tools (like `golang-migrate` CLI or Makefile scripts) require:

- External tooling installation
- Manual migration commands
- Separate deployment steps
- More moving parts to manage

Our solution embeds migrations directly in the Go binary using `//go:embed`, which means:

- ✅ Migrations run automatically on startup
- ✅ No external tools needed
- ✅ Single binary deployment
- ✅ Guaranteed version consistency
- ✅ Idempotent (safe to run multiple times)

### How It Works

```go
//go:embed migrations/*.sql
var migrationsFS embed.FS
```

This directive tells Go to embed all `.sql` files from the `migrations/` directory at compile time. The files become part of the binary and are available at runtime.

## Usage

### In Your Application

```go
import "folio/api/database"

func main() {
    // Load configuration from environment
    dbConfig := database.LoadConfigFromEnv()

    // Wait for database to be ready
    ctx := context.Background()
    if err := database.WaitForDatabase(ctx, dbConfig, 15); err != nil {
        log.Fatal(err)
    }

    // Run migrations automatically
    if err := database.RunMigrations(dbConfig); err != nil {
        log.Fatal(err)
    }

    // Create connection pool
    pool, err := database.NewPool(ctx, dbConfig)
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()

    // Use pool for queries
    var result int
    pool.QueryRow(ctx, "SELECT 1").Scan(&result)
}
```

### Configuration

The package reads from environment variables:

```bash
DB_HOST=localhost       # Default: localhost
DB_PORT=5432           # Default: 5432
DB_USER=folio_user     # Default: folio_user
DB_PASSWORD=***        # Default: folio_password
DB_NAME=folio_db       # Default: folio_db
```

### Connection Pool Settings

Configured in `NewPool()`:

- **MaxConns**: 10 (maximum connections)
- **MinConns**: 2 (minimum idle connections)
- **MaxConnLifetime**: 1 hour
- **MaxConnIdleTime**: 30 minutes

These values are optimized for a small-to-medium application. Adjust based on your load.

## Migration Files

### Naming Convention

```
migrations/
├── 000001_create_initial_tables.up.sql
├── 000001_create_initial_tables.down.sql
├── 000002_add_user_preferences.up.sql
└── 000002_add_user_preferences.down.sql
```

- **Format**: `{version}_{description}.{direction}.sql`
- **Version**: 6-digit number (000001, 000002, etc.)
- **Direction**: `up` (apply) or `down` (rollback)

### Creating New Migrations

1. Create new migration files:

```bash
# Up migration (changes)
touch migrations/000002_add_feature.up.sql

# Down migration (rollback)
touch migrations/000002_add_feature.down.sql
```

2. Write SQL:

```sql
-- 000002_add_feature.up.sql
ALTER TABLE users ADD COLUMN bio TEXT;

-- 000002_add_feature.down.sql
ALTER TABLE users DROP COLUMN bio;
```

3. Rebuild and restart:

```bash
docker-compose up --build
```

The new migrations will be embedded and applied automatically!

## API Reference

### `LoadConfigFromEnv() *Config`

Loads database configuration from environment variables with sensible defaults.

### `NewPool(ctx context.Context, config *Config) (*pgxpool.Pool, error)`

Creates a new connection pool with optimized settings.

### `RunMigrations(config *Config) error`

Applies all pending migrations from embedded files. Idempotent - safe to run multiple times.

### `WaitForDatabase(ctx context.Context, config *Config, maxRetries int) error`

Waits for database to be ready, retrying with exponential backoff.

### `Config.GetConnectionString() string`

Returns PostgreSQL connection string for the configuration.

## Advanced Usage

### Manual Migration Control

If you need to run migrations manually (e.g., in tests):

```go
// Skip migrations in main.go
// Then run manually:

import "folio/api/database"

func setupTestDB(t *testing.T) {
    config := &database.Config{
        Host:     "localhost",
        Port:     "5432",
        User:     "test_user",
        Password: "test_pass",
        Database: "test_db",
    }

    if err := database.RunMigrations(config); err != nil {
        t.Fatal(err)
    }
}
```

### Custom Pool Configuration

```go
// Get base config
poolConfig, _ := pgxpool.ParseConfig(config.GetConnectionString())

// Customize
poolConfig.MaxConns = 20
poolConfig.MinConns = 5

// Create pool
pool, _ := pgxpool.NewWithConfig(ctx, poolConfig)
```

## Troubleshooting

### Migrations Won't Apply

```bash
# Check if migrations are embedded
go build -o main .
strings main | grep "CREATE TABLE"  # Should show SQL

# Check migration version
docker exec -it folio-db psql -U folio_user -d folio_db \
  -c "SELECT version, dirty FROM schema_migrations;"
```

### Database Connection Issues

```bash
# Test connection manually
docker exec -it folio-db psql -U folio_user -d folio_db -c "SELECT 1;"

# Check backend logs
docker-compose logs api
```

### Dirty Migration State

If a migration fails mid-way, the database may be in a "dirty" state:

```sql
-- Connect to database
docker exec -it folio-db psql -U folio_user -d folio_db

-- Check current version
SELECT * FROM schema_migrations;

-- Force version (use carefully!)
UPDATE schema_migrations SET dirty = false WHERE version = 1;
```

## Performance Considerations

### Pool Sizing

The default configuration (10 max, 2 min) is suitable for:

- Development
- Small production deployments (<1000 concurrent users)

For high-traffic applications:

```go
poolConfig.MaxConns = 50
poolConfig.MinConns = 10
```

### Query Optimization

Always use `context.Context` with timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := pool.Query(ctx, "SELECT * FROM users")
```

## Security Best Practices

1. **Never hardcode credentials** - Always use environment variables
2. **Use connection pooling** - Prevents connection exhaustion
3. **Set query timeouts** - Prevents slow queries from hanging
4. **Validate inputs** - Use parameterized queries to prevent SQL injection

Example of parameterized query:

```go
// ✅ Safe
pool.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", userID)

// ❌ Unsafe
pool.QueryRow(ctx, fmt.Sprintf("SELECT * FROM users WHERE id = %s", userID))
```

## Comparison: Before vs After

### Before (Manual Migrations)

```bash
# Developer workflow
cd database
make migrate-up           # Manual step
cd ../backend
go run main.go

# Docker deployment
docker-compose up -d db
docker-compose run --rm migrate  # Extra container
docker-compose up -d api
```

### After (Embedded Migrations)

```bash
# Developer workflow
go run main.go            # That's it!

# Docker deployment
docker-compose up -d      # That's it!
```

## Credits

Built with:

- [pgx/v5](https://github.com/jackc/pgx) - PostgreSQL driver
- [golang-migrate/migrate](https://github.com/golang-migrate/migrate) - Migration library
- [embed](https://pkg.go.dev/embed) - Standard library

This elegant solution demonstrates the power of Go's `embed` package and composable libraries! 🚀
