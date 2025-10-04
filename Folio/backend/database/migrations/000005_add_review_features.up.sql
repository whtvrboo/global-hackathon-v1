-- Add spoiler flag and review formatting to logs
ALTER TABLE logs ADD COLUMN spoiler_flag BOOLEAN DEFAULT false;
ALTER TABLE logs ADD COLUMN review_format VARCHAR(20) DEFAULT 'text';

-- Add social features to lists
ALTER TABLE lists ADD COLUMN likes_count INTEGER DEFAULT 0;
ALTER TABLE lists ADD COLUMN comments_count INTEGER DEFAULT 0;

-- Create list_likes table
CREATE TABLE list_likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_id UUID NOT NULL REFERENCES lists(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(list_id, user_id)
);

-- Create list_comments table
CREATE TABLE list_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_id UUID NOT NULL REFERENCES lists(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_list_likes_list_id ON list_likes(list_id);
CREATE INDEX idx_list_likes_user_id ON list_likes(user_id);
CREATE INDEX idx_list_comments_list_id ON list_comments(list_id);
CREATE INDEX idx_list_comments_user_id ON list_comments(user_id);

-- Create triggers to update list like counts
CREATE OR REPLACE FUNCTION update_list_likes_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE lists SET likes_count = likes_count + 1 WHERE id = NEW.list_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE lists SET likes_count = likes_count - 1 WHERE id = OLD.list_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_list_likes_count
    AFTER INSERT OR DELETE ON list_likes
    FOR EACH ROW EXECUTE FUNCTION update_list_likes_count();

-- Create triggers to update list comment counts
CREATE OR REPLACE FUNCTION update_list_comments_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE lists SET comments_count = comments_count + 1 WHERE id = NEW.list_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE lists SET comments_count = comments_count - 1 WHERE id = OLD.list_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_list_comments_count
    AFTER INSERT OR DELETE ON list_comments
    FOR EACH ROW EXECUTE FUNCTION update_list_comments_count();
