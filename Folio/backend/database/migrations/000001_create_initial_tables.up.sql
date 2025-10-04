-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    google_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    picture TEXT,
    bio TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_google_id ON users(google_id);
CREATE INDEX idx_users_email ON users(email);

-- Books table (cached book data from external APIs)
CREATE TABLE IF NOT EXISTS books (
    id VARCHAR(255) PRIMARY KEY, -- External API book ID
    title VARCHAR(500) NOT NULL,
    authors TEXT[], -- Array of author names
    description TEXT,
    cover_url TEXT,
    published_date VARCHAR(50),
    page_count INTEGER,
    isbn_10 VARCHAR(20),
    isbn_13 VARCHAR(20),
    categories TEXT[],
    language VARCHAR(10),
    publisher VARCHAR(255),
    rating DECIMAL(3, 2),
    ratings_count INTEGER,
    api_source VARCHAR(50), -- 'google' or 'openlibrary'
    raw_data JSONB, -- Store full API response
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_books_title ON books(title);
CREATE INDEX idx_books_authors ON books USING GIN(authors);

-- Logs table (user book entries)
CREATE TABLE IF NOT EXISTS logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id VARCHAR(255) NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('want_to_read', 'reading', 'read', 'dnf')),
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    review TEXT,
    notes TEXT,
    start_date DATE,
    finish_date DATE,
    is_public BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_logs_user_id ON logs(user_id);
CREATE INDEX idx_logs_book_id ON logs(book_id);
CREATE INDEX idx_logs_status ON logs(status);
CREATE INDEX idx_logs_created_at ON logs(created_at DESC);
CREATE INDEX idx_logs_user_created ON logs(user_id, created_at DESC);

-- Watchlists table (books users want to track)
CREATE TABLE IF NOT EXISTS watchlists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id VARCHAR(255) NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    notes TEXT,
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, book_id)
);

CREATE INDEX idx_watchlists_user_id ON watchlists(user_id);
CREATE INDEX idx_watchlists_book_id ON watchlists(book_id);
CREATE INDEX idx_watchlists_created_at ON watchlists(created_at DESC);

-- Followers table (social following relationships)
CREATE TABLE IF NOT EXISTS followers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(follower_id, following_id),
    CHECK (follower_id != following_id)
);

CREATE INDEX idx_followers_follower_id ON followers(follower_id);
CREATE INDEX idx_followers_following_id ON followers(following_id);
CREATE INDEX idx_followers_created_at ON followers(created_at DESC);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for updating updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_books_updated_at BEFORE UPDATE ON books
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_logs_updated_at BEFORE UPDATE ON logs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

