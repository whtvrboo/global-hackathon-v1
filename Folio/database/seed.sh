#!/bin/bash

# Seed script for Folio database
# This script drops the database, runs migrations, and seeds with demo data

set -e

echo "🌱 Seeding Folio database with demo data..."

# Check if we're in the right directory
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ Please run this script from the Folio root directory"
    exit 1
fi

# Stop any running containers
echo "🛑 Stopping existing containers..."
docker-compose down

# Remove the database volume to start fresh
echo "🗑️  Removing existing database volume..."
docker volume rm folio_postgres_data 2>/dev/null || true

# Start the database
echo "🚀 Starting database..."
docker-compose up -d postgres

# Wait for database to be ready
echo "⏳ Waiting for database to be ready..."
sleep 10

# Run migrations
echo "📊 Running database migrations..."
docker-compose exec -T postgres psql -U folio -d folio -f /docker-entrypoint-initdb.d/000001_create_initial_tables.up.sql
docker-compose exec -T postgres psql -U folio -d folio -f /docker-entrypoint-initdb.d/000002_add_guest_support.up.sql
docker-compose exec -T postgres psql -U folio -d folio -f /docker-entrypoint-initdb.d/000003_add_likes_comments.up.sql
docker-compose exec -T postgres psql -U folio -d folio -f /docker-entrypoint-initdb.d/000004_add_custom_lists.up.sql

# Seed the database
echo "🌱 Seeding database with demo data..."
docker-compose exec -T postgres psql -U folio -d folio -f /docker-entrypoint-initdb.d/books_seed.sql
docker-compose exec -T postgres psql -U folio -d folio -f /docker-entrypoint-initdb.d/seed.sql

echo "✅ Database seeded successfully!"
echo ""
echo "📊 Demo data includes:"
echo "   • 5 users with realistic profiles"
echo "   • 65+ real books across all genres"
echo "   • 25+ reading logs with ratings and reviews"
echo "   • Social interactions (likes and comments)"
echo "   • 12 custom lists with curated books"
echo "   • Follower relationships"
echo "   • Comprehensive book database for search"
echo ""
echo "🚀 You can now start the full application with:"
echo "   docker-compose up"
echo ""
echo "👤 Demo users:"
echo "   • alex_reader (Alex Johnson)"
echo "   • bookworm_sarah (Sarah Chen)"
echo "   • mike_literature (Mike Rodriguez)"
echo "   • emma_books (Emma Thompson)"
echo "   • david_reader (David Kim)"
