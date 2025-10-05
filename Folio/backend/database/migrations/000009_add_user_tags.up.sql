-- Create user_tags table for tracking tag usage counts
CREATE TABLE IF NOT EXISTS user_tags (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tag TEXT NOT NULL,
    usage_count INTEGER DEFAULT 1,
    last_used_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, tag)
);

-- Index for efficient querying
CREATE INDEX idx_user_tags_user_count ON user_tags(user_id, usage_count DESC);
CREATE INDEX idx_user_tags_last_used ON user_tags(user_id, last_used_at DESC);

-- Function to update user_tags when annotations are inserted/updated
CREATE OR REPLACE FUNCTION update_user_tags()
RETURNS TRIGGER AS $$
DECLARE
    tag_item TEXT;
BEGIN
    -- Handle INSERT
    IF TG_OP = 'INSERT' THEN
        -- Process each tag in the new annotation
        IF NEW.tags IS NOT NULL THEN
            FOREACH tag_item IN ARRAY NEW.tags
            LOOP
                -- Insert or update tag count
                INSERT INTO user_tags (user_id, tag, usage_count, last_used_at)
                VALUES (NEW.user_id, tag_item, 1, NEW.created_at)
                ON CONFLICT (user_id, tag) 
                DO UPDATE SET 
                    usage_count = user_tags.usage_count + 1,
                    last_used_at = NEW.created_at;
            END LOOP;
        END IF;
        RETURN NEW;
    END IF;
    
    -- Handle UPDATE
    IF TG_OP = 'UPDATE' THEN
        -- Remove old tags
        IF OLD.tags IS NOT NULL THEN
            FOREACH tag_item IN ARRAY OLD.tags
            LOOP
                UPDATE user_tags 
                SET usage_count = usage_count - 1
                WHERE user_id = OLD.user_id AND tag = tag_item;
                
                -- Remove tag if count reaches 0
                DELETE FROM user_tags 
                WHERE user_id = OLD.user_id AND tag = tag_item AND usage_count <= 0;
            END LOOP;
        END IF;
        
        -- Add new tags
        IF NEW.tags IS NOT NULL THEN
            FOREACH tag_item IN ARRAY NEW.tags
            LOOP
                INSERT INTO user_tags (user_id, tag, usage_count, last_used_at)
                VALUES (NEW.user_id, tag_item, 1, NEW.updated_at)
                ON CONFLICT (user_id, tag) 
                DO UPDATE SET 
                    usage_count = user_tags.usage_count + 1,
                    last_used_at = NEW.updated_at;
            END LOOP;
        END IF;
        RETURN NEW;
    END IF;
    
    -- Handle DELETE
    IF TG_OP = 'DELETE' THEN
        IF OLD.tags IS NOT NULL THEN
            FOREACH tag_item IN ARRAY OLD.tags
            LOOP
                UPDATE user_tags 
                SET usage_count = usage_count - 1
                WHERE user_id = OLD.user_id AND tag = tag_item;
                
                -- Remove tag if count reaches 0
                DELETE FROM user_tags 
                WHERE user_id = OLD.user_id AND tag = tag_item AND usage_count <= 0;
            END LOOP;
        END IF;
        RETURN OLD;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for annotations table
DROP TRIGGER IF EXISTS trigger_update_user_tags ON annotations;
CREATE TRIGGER trigger_update_user_tags
    AFTER INSERT OR UPDATE OR DELETE ON annotations
    FOR EACH ROW EXECUTE FUNCTION update_user_tags();

-- Backfill existing annotations to populate user_tags
INSERT INTO user_tags (user_id, tag, usage_count, last_used_at)
SELECT 
    user_id,
    unnest(tags) as tag,
    COUNT(*) as usage_count,
    MAX(created_at) as last_used_at
FROM annotations 
WHERE tags IS NOT NULL AND array_length(tags, 1) > 0
GROUP BY user_id, unnest(tags)
ON CONFLICT (user_id, tag) 
DO UPDATE SET 
    usage_count = EXCLUDED.usage_count,
    last_used_at = EXCLUDED.last_used_at;
