# 📚 Folio - Discover Books Through the People Who Read Them

A modern social reading platform where **curated lists and authentic reviews** take center stage. Built with the "Sovereign Stack" philosophy - fully self-hosted, beautifully designed, and focused on meaningful book discovery.

## ✨ What Makes Folio Different

**Lists as Art.** Every list is a beautifully designed, shareable page that showcases a curator's taste. No spreadsheets, no boring grids - just magazine-quality presentation.

**Instant Reviews.** Log a book, give it a rating, and we'll prompt you for one thought you'd share with a friend. That's it. Your review is live, helping others discover their next read.

**Social Discovery.** Follow curators whose taste you trust. See their lists, reviews, and reading activity in a clean, TikTok-style feed.

## 🎯 Philosophy: The Sovereign Stack

Folio is built to be **fully self-contained and deployable anywhere**. No vendor lock-in, no proprietary services - just Docker, PostgreSQL, Go, and Vue. Your data, your infrastructure, your control.

## 🏗️ Architecture

### Tech Stack

- **Frontend**: Vue 3 + Vite + TailwindCSS + DaisyUI
- **Backend**: Go 1.22 + Echo Framework
- **Database**: PostgreSQL 16
- **Orchestration**: Docker Compose
- **Migrations**: golang-migrate (embedded in Go binary)

### Project Structure

```
Folio/
├── frontend/              # Vue 3 frontend
│   ├── src/
│   │   ├── components/   # Reusable Vue components
│   │   ├── views/        # Page components
│   │   ├── router/       # Vue Router configuration
│   │   └── stores/       # Pinia state management
│   ├── Dockerfile        # Multi-stage Node + Nginx build
│   └── nginx.conf        # Nginx configuration
├── backend/              # Go API server
│   ├── database/         # Database connection & migrations
│   │   ├── migrations/   # SQL migration files (embedded)
│   │   └── database.go   # Connection pool & auto-migration
│   ├── main.go           # Application entry point
│   ├── Dockerfile        # Multi-stage Go build
│   └── go.mod            # Go dependencies
├── database/             # Database utilities (legacy)
│   ├── Makefile          # Migration commands (Unix/Mac)
│   └── run.ps1           # Migration commands (Windows)
└── docker-compose.yml    # Full stack orchestration
```

## 🚀 Quick Start (5 Minutes to Launch)

### Prerequisites

- **Docker & Docker Compose** installed ([Get Docker](https://docs.docker.com/get-docker/))
- **2GB RAM** available
- **Ports available**: 80 (frontend), 8080 (API), 5432 (database)

### Option 1: Quick Launch (Recommended)

```bash
# Clone the repository
git clone <your-repo-url>
cd Folio

# Copy environment file
cp env.example .env

# Start everything with one command
docker-compose up --build

# ✨ That's it! Visit http://localhost
```

The application will:

1. 🗄️ Initialize PostgreSQL database
2. 🔄 Run all migrations automatically
3. 🚀 Start the Go API server
4. 🎨 Build and serve the Vue frontend
5. ✅ Be ready at **http://localhost**

### Option 2: Development Setup

For active development with hot-reload:

```bash
# Terminal 1 - Start database
docker-compose up db

# Terminal 2 - Run backend
cd backend
go run main.go

# Terminal 3 - Run frontend with hot-reload
cd frontend
npm install
npm run dev

# Visit http://localhost:5173 (Vite dev server)
```

### First Time Setup

**No manual steps required!** The backend automatically:

- Waits for PostgreSQL to be ready
- Runs all database migrations
- Seeds demo data (optional)
- Starts serving requests

**Visit http://localhost** and you'll see:

- A beautiful discover feed with curated lists
- Search functionality powered by Google Books API
- Guest browsing or sign up with Google OAuth

## 🎨 Core Features

### 📋 Beautiful Lists

- **Magazine-quality presentation** with header images and themed colors
- **Drag-and-drop reordering** for perfect curation
- **Curator notes** for each book explaining why it's on the list
- **Public/private visibility** controls

### ⭐ One-Click Reviews

- Log a book with a status (Want to Read, Reading, Read, DNF)
- Give it a star rating
- **Instant prompt**: "What's one thought you'd share with a friend?"
- Your review goes live immediately on the book's page

### 🔍 Smart Discovery

- **Unified feed** showing lists, reviews, and reading activity
- **Search** powered by Google Books API with 10M+ titles
- **Follow curators** whose taste resonates with you
- **Community stats** showing what others are reading

### 👤 Profile Pages

- Your lists displayed as a **gallery of taste**
- Reading stats and analytics
- Followers/following social graph
- Custom bio and profile customization

### 📱 Progressive Web App (PWA)

- **Install on mobile** for app-like experience
- **Offline support** for cached content
- **Push notifications** for social activity (coming soon)

## 🔧 Development

### Backend Development

```bash
cd backend

# Run locally (requires Go 1.22+)
go run main.go

# Run tests
go test ./...

# Format code
go fmt ./...

# Add a new migration
# Create files: backend/database/migrations/000XXX_description.up.sql
#               backend/database/migrations/000XXX_description.down.sql
# Migrations run automatically on startup
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Run dev server (hot reload at localhost:5173)
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Key Technologies

- **Backend**: Go 1.22, Echo framework, GORM, JWT auth
- **Frontend**: Vue 3 (Composition API), Vite, TailwindCSS, Pinia
- **Database**: PostgreSQL 16 with UUID primary keys
- **Search**: Google Books API integration
- **Auth**: Google OAuth 2.0 + Guest mode

### Database Management

The backend handles migrations automatically, but you can also manage them manually:

#### Windows (PowerShell)

```powershell
cd database
.\run.ps1 migrate-up      # Apply migrations
.\run.ps1 migrate-down    # Rollback
.\run.ps1 migrate-version # Check version
.\run.ps1 psql            # Connect to DB
```

#### Unix/Mac

```bash
cd database
make migrate-up      # Apply migrations
make migrate-down    # Rollback
make migrate-version # Check version
make psql            # Connect to DB
```

## 📊 Database Schema

### Core Tables

- **users**: User accounts and profiles (Google OAuth)
- **books**: Cached book data from external APIs (Google Books, Open Library)
- **logs**: User reading logs (status, rating, reviews)
- **watchlists**: Books users want to track
- **followers**: Social following relationships

All tables use:

- UUID primary keys for distributed scalability
- TIMESTAMPTZ for timezone-aware timestamps
- Proper indexes for query performance
- Foreign key constraints for data integrity

## 🎨 Design System

Folio uses a Linear-inspired design language:

- Clean, minimal interface
- Subtle animations and transitions
- Focus on content and readability
- Consistent spacing and typography
- DaisyUI component library for rapid development

## 🔌 API Endpoints

### Public Endpoints (No Auth Required)

```
GET  /api/health                    # System health check
GET  /api/books/search?q=...        # Search books via Google Books API
GET  /api/books/:id                 # Get book details
GET  /api/books/:id/reviews         # Get public reviews for a book
GET  /api/books/:id/lists           # Get lists containing a book
GET  /api/lists/:id                 # Get public list details
GET  /api/discover                  # Get discovery feed (lists, books)
GET  /api/users/:username/profile   # Get public profile
```

### Protected Endpoints (Requires JWT)

```
POST /api/logs                      # Create reading log
PUT  /api/logs/:id                  # Update reading log (for quick reviews)
GET  /api/me/logs                   # Get current user's logs
GET  /api/me/lists                  # Get current user's lists
POST /api/lists                     # Create a new list
PUT  /api/lists/:id                 # Update list metadata
POST /api/lists/:id/items           # Add book to list
PUT  /api/lists/:id/items/:item_id  # Update list item notes
DELETE /api/lists/:id/items/:item_id # Remove book from list
POST /api/users/:username/follow    # Follow a user
DELETE /api/users/:username/follow  # Unfollow a user
GET  /api/feed                      # Get personalized feed
```

### Authentication

```
GET  /api/auth/google               # Initiate Google OAuth flow
GET  /api/auth/google/callback      # OAuth callback handler
POST /api/auth/guest                # Create guest account
POST /api/auth/guest/convert        # Convert guest to full account
```

## 🔐 Authentication

Folio supports two authentication modes:

### Google OAuth 2.0 (Primary)

1. User clicks "Sign in with Google"
2. Redirected to Google's consent screen
3. Backend exchanges code for access token
4. User profile is fetched and stored in database
5. JWT is generated and returned to frontend
6. JWT is stored in localStorage and used for API requests

**Setup Required:**

- Create a Google Cloud project
- Enable Google+ API
- Create OAuth 2.0 credentials
- Add authorized redirect URIs
- Set `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, and `GOOGLE_REDIRECT_URL` in `.env`

### Guest Mode (Frictionless Onboarding)

1. User clicks "Continue as Guest"
2. Backend creates a temporary guest account
3. Guest can browse, create lists, and log books
4. **Conversion prompt** appears after meaningful activity
5. Guest can convert to full account via Google OAuth
6. All guest data is preserved after conversion

**Benefits:**

- Zero-friction onboarding
- Users can try before committing
- Increases conversion rates
- No data loss on conversion

## 🚢 Production Deployment

### Docker Compose (Recommended)

The application is designed to run anywhere Docker is available:

```bash
# Use production compose file
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose logs -f backend frontend

# Stop services
docker-compose down

# Update to latest version
git pull
docker-compose -f docker-compose.prod.yml up -d --build
```

### Platform-Specific Guides

#### Railway (Easiest)

1. Fork this repository
2. Create new project on Railway
3. Connect your GitHub repository
4. Railway auto-detects Docker Compose
5. Set environment variables in dashboard
6. Deploy! Railway provides HTTPS domain automatically

**Required Environment Variables:**

```bash
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=https://your-app.railway.app/api/auth/google/callback
JWT_SECRET=your_random_secret_key_min_32_chars
```

#### Fly.io

```bash
# Install flyctl
curl -L https://fly.io/install.sh | sh

# Login
fly auth login

# Launch app (follow prompts)
fly launch

# Set secrets
fly secrets set GOOGLE_CLIENT_ID=xxx
fly secrets set GOOGLE_CLIENT_SECRET=xxx
fly secrets set JWT_SECRET=xxx

# Deploy
fly deploy
```

#### DigitalOcean / AWS / GCP

1. Create a VM with Docker installed
2. Clone repository
3. Copy `env.prod.example` to `.env` and configure
4. Run `docker-compose -f docker-compose.prod.yml up -d`
5. Set up reverse proxy (Nginx/Caddy) for HTTPS
6. Point domain to server IP

### Environment Variables Reference

```bash
# Database Configuration
DB_HOST=db                          # Docker service name or host
DB_PORT=5432
DB_USER=folio_user
DB_PASSWORD=your_secure_password    # Change in production!
DB_NAME=folio_db

# Server Configuration
SERVER_PORT=8080
FRONTEND_URL=http://localhost       # Your frontend URL

# Google OAuth (Required for production)
GOOGLE_CLIENT_ID=xxx.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=xxx
GOOGLE_REDIRECT_URL=https://yourdomain.com/api/auth/google/callback

# Security
JWT_SECRET=your_random_secret_minimum_32_characters
ALLOWED_ORIGINS=https://yourdomain.com

# Optional: External Services
GOOGLE_BOOKS_API_KEY=xxx            # For higher rate limits
SENTRY_DSN=xxx                      # For error tracking
```

### SSL/HTTPS Setup

**Option 1: Caddy (Automatic HTTPS)**

```bash
# Install Caddy
sudo apt install caddy

# Create Caddyfile
cat > /etc/caddy/Caddyfile <<EOF
yourdomain.com {
    reverse_proxy localhost:80
}
EOF

# Restart Caddy
sudo systemctl restart caddy
```

**Option 2: Let's Encrypt + Nginx**

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d yourdomain.com

# Auto-renewal is configured automatically
```

### Health Checks

Monitor your deployment:

```bash
# API health
curl https://yourdomain.com/api/health

# Expected response:
{
  "status": "healthy",
  "database": "connected",
  "version": "1.0.0"
}
```

### Backup Strategy

```bash
# Backup database
docker exec folio-db pg_dump -U folio_user folio_db > backup.sql

# Restore database
docker exec -i folio-db psql -U folio_user folio_db < backup.sql

# Automate with cron
0 2 * * * /path/to/backup-script.sh
```

## 🧪 Testing

```bash
# Backend tests
cd backend && go test ./...

# Frontend tests
cd frontend && npm test

# Integration tests
docker-compose up -d
curl http://localhost:8080/api/health
```

## 📝 Feature Roadmap

### ✅ Completed (v1.0)

- [x] Docker orchestration with auto-migrations
- [x] Google OAuth + Guest mode authentication
- [x] Book search via Google Books API
- [x] Reading logs with status tracking
- [x] Beautiful list creation and management
- [x] One-click review prompts
- [x] Social following system
- [x] Unified discovery feed
- [x] Profile pages with stats
- [x] PWA support with offline mode
- [x] Responsive mobile-first design

### 🚧 In Progress (v1.1)

- [ ] Comments on lists and reviews
- [ ] Advanced search filters
- [ ] Reading challenges and goals
- [ ] Book recommendations engine
- [ ] Export data (JSON, CSV)

### 🔮 Future (v2.0+)

- [ ] Reading clubs and group discussions
- [ ] Book clubs with scheduled reads
- [ ] Activity notifications
- [ ] Email digests
- [ ] Import from Goodreads
- [ ] Mobile apps (React Native)
- [ ] API rate limiting and caching
- [ ] Full-text search with Elasticsearch

## 🤝 Contributing

This is a hackathon project, but contributions are welcome!

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 📄 License

See LICENSE file for details.

## 🎯 Demo Flow

Want to see Folio in action? Here's the perfect demo path:

1. **Land on Discover** → See curated lists and trending books in a beautiful feed
2. **Click a List** → Experience magazine-quality list presentation with book covers
3. **View a Profile** → See a curator's gallery of taste with their best lists
4. **Search a Book** → Find any book from 10M+ titles via Google Books
5. **Log the Book** → Mark as "Read" and give it 5 stars
6. **The Magic** → Modal transforms: "What's one thought you'd share with a friend?"
7. **Publish Review** → Your review is instantly live on the book's page

**This is the wedge:** Lists and reviews as the hero, not features buried in menus.

## 🤝 Contributing

This project welcomes contributions! Here's how:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests and linting
5. Commit with clear messages (`git commit -m 'Add amazing feature'`)
6. Push to your branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

**Areas we'd love help with:**

- 🎨 UI/UX improvements
- 🔍 Advanced search features
- 📱 Mobile app development
- 🌍 Internationalization
- 📚 Book data sources beyond Google Books
- ⚡ Performance optimizations

## 📄 License

See LICENSE file for details.

## 🙏 Acknowledgments

- Built for the **Global Hackathon v1**
- Inspired by **Goodreads** (functionality) and **Linear** (design philosophy)
- Book data powered by **Google Books API**
- Hosted on the **Sovereign Stack** (Docker + PostgreSQL + Go + Vue)
- Design system inspired by **TailwindCSS** and **modern web aesthetics**

## 🐛 Troubleshooting

### Database connection failed

```bash
# Check if database is running
docker-compose ps

# View database logs
docker-compose logs db

# Restart database
docker-compose restart db
```

### Frontend not loading

```bash
# Check if all services are running
docker-compose ps

# Rebuild frontend
docker-compose up --build frontend

# Check Nginx logs
docker-compose logs frontend
```

### OAuth not working

1. Verify `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET` are set
2. Check `GOOGLE_REDIRECT_URL` matches your Google Cloud Console
3. Ensure redirect URI is added to authorized URIs in Google Console
4. Try guest mode to test functionality without OAuth

### Port already in use

```bash
# Find process using port 80
sudo lsof -i :80

# Kill process or change port in docker-compose.yml
```

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/whtvrboo/global-hackathon-v1/issues)
- **Discussions**: [GitHub Discussions](https://github.com/whtvrboo/global-hackathon-v1/discussions)
- **Email**: your-email@example.com

---

**Made with ❤️ and lots of ☕ by readers, for readers.**

_Folio: Where lists are art, and reviews are conversations._
