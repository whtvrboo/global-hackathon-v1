-- Drop triggers
DROP TRIGGER IF EXISTS trigger_update_log_comments_count ON log_comments;
DROP TRIGGER IF EXISTS trigger_update_log_likes_count ON log_likes;

-- Drop functions
DROP FUNCTION IF EXISTS update_log_comments_count();
DROP FUNCTION IF EXISTS update_log_likes_count();

-- Remove columns from logs table
ALTER TABLE logs DROP COLUMN IF EXISTS comments_count;
ALTER TABLE logs DROP COLUMN IF EXISTS likes_count;

-- Drop tables
DROP TABLE IF EXISTS log_comments;
DROP TABLE IF EXISTS log_likes;

