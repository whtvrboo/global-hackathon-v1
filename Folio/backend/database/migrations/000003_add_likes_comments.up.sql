-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create log_likes table
CREATE TABLE log_likes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    log_id UUID NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, log_id)
);

-- Create indexes for likes
CREATE INDEX idx_log_likes_user_id ON log_likes(user_id);
CREATE INDEX idx_log_likes_log_id ON log_likes(log_id);
CREATE INDEX idx_log_likes_created_at ON log_likes(created_at DESC);

-- Create log_comments table
CREATE TABLE log_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    log_id UUID NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for comments
CREATE INDEX idx_log_comments_user_id ON log_comments(user_id);
CREATE INDEX idx_log_comments_log_id ON log_comments(log_id);
CREATE INDEX idx_log_comments_created_at ON log_comments(created_at DESC);

-- Add counts to logs table
ALTER TABLE logs ADD COLUMN likes_count INTEGER DEFAULT 0;
ALTER TABLE logs ADD COLUMN comments_count INTEGER DEFAULT 0;

-- Create function to update likes_count
CREATE OR REPLACE FUNCTION update_log_likes_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE logs SET likes_count = likes_count + 1 WHERE id = NEW.log_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE logs SET likes_count = GREATEST(likes_count - 1, 0) WHERE id = OLD.log_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for likes_count
CREATE TRIGGER trigger_update_log_likes_count
AFTER INSERT OR DELETE ON log_likes
FOR EACH ROW EXECUTE FUNCTION update_log_likes_count();

-- Create function to update comments_count
CREATE OR REPLACE FUNCTION update_log_comments_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE logs SET comments_count = comments_count + 1 WHERE id = NEW.log_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE logs SET comments_count = GREATEST(comments_count - 1, 0) WHERE id = OLD.log_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for comments_count
CREATE TRIGGER trigger_update_log_comments_count
AFTER INSERT OR DELETE ON log_comments
FOR EACH ROW EXECUTE FUNCTION update_log_comments_count();

