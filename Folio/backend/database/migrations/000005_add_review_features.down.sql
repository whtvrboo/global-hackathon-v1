-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_list_comments_count ON list_comments;
DROP TRIGGER IF EXISTS trigger_update_list_likes_count ON list_likes;

-- Drop functions
DROP FUNCTION IF EXISTS update_list_comments_count();
DROP FUNCTION IF EXISTS update_list_likes_count();

-- Drop indexes
DROP INDEX IF EXISTS idx_list_comments_user_id;
DROP INDEX IF EXISTS idx_list_comments_list_id;
DROP INDEX IF EXISTS idx_list_likes_user_id;
DROP INDEX IF EXISTS idx_list_likes_list_id;

-- Drop tables
DROP TABLE IF EXISTS list_comments;
DROP TABLE IF EXISTS list_likes;

-- Remove columns from lists
ALTER TABLE lists DROP COLUMN IF EXISTS comments_count;
ALTER TABLE lists DROP COLUMN IF EXISTS likes_count;

-- Remove columns from logs
ALTER TABLE logs DROP COLUMN IF EXISTS review_format;
ALTER TABLE logs DROP COLUMN IF EXISTS spoiler_flag;
