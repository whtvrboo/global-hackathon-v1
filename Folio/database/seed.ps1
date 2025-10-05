# Seed script for Folio database (Windows PowerShell)
# This script drops the database, runs migrations, and seeds with demo data

Write-Host "🌱 Seeding Folio database with demo data..." -ForegroundColor Green

# Check if we're in the right directory
if (-not (Test-Path "docker-compose.yml")) {
    Write-Host "❌ Please run this script from the Folio root directory" -ForegroundColor Red
    exit 1
}

# Stop any running containers
Write-Host "🛑 Stopping existing containers..." -ForegroundColor Yellow
docker-compose down

# Remove the database volume to start fresh
Write-Host "🗑️  Removing existing database volume..." -ForegroundColor Yellow
docker volume rm folio_postgres_data 2>$null

# Start the database
Write-Host "🚀 Starting database..." -ForegroundColor Yellow
docker-compose up -d postgres

# Wait for database to be ready
Write-Host "⏳ Waiting for database to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

# Run migrations
Write-Host "📊 Running database migrations..." -ForegroundColor Yellow
docker-compose exec -T postgres psql -U folio -d folio -f /docker-entrypoint-initdb.d/000001_create_initial_tables.up.sql
docker-compose exec -T postgres psql -U folio -d folio -f /docker-entrypoint-initdb.d/000002_add_guest_support.up.sql
docker-compose exec -T postgres psql -U folio -d folio -f /docker-entrypoint-initdb.d/000003_add_likes_comments.up.sql
docker-compose exec -T postgres psql -U folio -d folio -f /docker-entrypoint-initdb.d/000004_add_custom_lists.up.sql

# Seed the database
Write-Host "🌱 Seeding database with demo data..." -ForegroundColor Yellow
docker-compose exec -T postgres psql -U folio -d folio -f /docker-entrypoint-initdb.d/books_seed.sql
docker-compose exec -T postgres psql -U folio -d folio -f /docker-entrypoint-initdb.d/seed.sql

Write-Host "✅ Database seeded successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "📊 Demo data includes:" -ForegroundColor Cyan
Write-Host "   • 5 users with realistic profiles" -ForegroundColor White
Write-Host "   • 65+ real books across all genres" -ForegroundColor White
Write-Host "   • 25+ reading logs with ratings and reviews" -ForegroundColor White
Write-Host "   • Social interactions (likes and comments)" -ForegroundColor White
Write-Host "   • 12 custom lists with curated books" -ForegroundColor White
Write-Host "   • Follower relationships" -ForegroundColor White
Write-Host "   • Comprehensive book database for search" -ForegroundColor White
Write-Host ""
Write-Host "🚀 You can now start the full application with:" -ForegroundColor Cyan
Write-Host "   docker-compose up" -ForegroundColor White
Write-Host ""
Write-Host "👤 Demo users:" -ForegroundColor Cyan
Write-Host "   • alex_reader (Alex Johnson)" -ForegroundColor White
Write-Host "   • bookworm_sarah (Sarah Chen)" -ForegroundColor White
Write-Host "   • mike_literature (Mike Rodriguez)" -ForegroundColor White
Write-Host "   • emma_books (Emma Thompson)" -ForegroundColor White
Write-Host "   • david_reader (David Kim)" -ForegroundColor White
