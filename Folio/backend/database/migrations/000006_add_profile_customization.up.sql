-- Add profile customization fields to users table
ALTER TABLE users ADD COLUMN favorite_book_ids TEXT[] DEFAULT '{}';
ALTER TABLE users ADD COLUMN banner_url TEXT;
ALTER TABLE users ADD COLUMN bio TEXT;
ALTER TABLE users ADD COLUMN reading_goal INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN reading_goal_year INTEGER DEFAULT EXTRACT(YEAR FROM NOW());

-- Create user_reading_stats table for advanced analytics
CREATE TABLE user_reading_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    year INTEGER NOT NULL,
    month INTEGER NOT NULL,
    books_read INTEGER DEFAULT 0,
    pages_read INTEGER DEFAULT 0,
    avg_rating DECIMAL(3,2) DEFAULT 0.0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, year, month)
);

-- Create indexes
CREATE INDEX idx_user_reading_stats_user_id ON user_reading_stats(user_id);
CREATE INDEX idx_user_reading_stats_year_month ON user_reading_stats(year, month);

-- Create function to update reading stats when a log is created/updated
CREATE OR REPLACE FUNCTION update_user_reading_stats()
RETURNS TRIGGER AS $$
DECLARE
    finish_year INTEGER;
    finish_month INTEGER;
    book_pages INTEGER;
BEGIN
    -- Only update stats for completed reads
    IF NEW.status = 'read' AND NEW.finish_date IS NOT NULL THEN
        finish_year := EXTRACT(YEAR FROM NEW.finish_date::DATE);
        finish_month := EXTRACT(MONTH FROM NEW.finish_date::DATE);
        
        -- Get book pages
        SELECT pages INTO book_pages FROM books WHERE id = NEW.book_id;
        IF book_pages IS NULL THEN
            book_pages := 0;
        END IF;
        
        -- Insert or update stats
        INSERT INTO user_reading_stats (user_id, year, month, books_read, pages_read, avg_rating)
        VALUES (NEW.user_id, finish_year, finish_month, 1, book_pages, COALESCE(NEW.rating, 0))
        ON CONFLICT (user_id, year, month)
        DO UPDATE SET
            books_read = user_reading_stats.books_read + 1,
            pages_read = user_reading_stats.pages_read + book_pages,
            avg_rating = (user_reading_stats.avg_rating * user_reading_stats.books_read + COALESCE(NEW.rating, 0)) / (user_reading_stats.books_read + 1),
            updated_at = NOW();
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for reading stats
CREATE TRIGGER trigger_update_user_reading_stats
    AFTER INSERT OR UPDATE ON logs
    FOR EACH ROW EXECUTE FUNCTION update_user_reading_stats();
