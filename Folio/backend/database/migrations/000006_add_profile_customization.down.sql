-- Drop trigger
DROP TRIGGER IF EXISTS trigger_update_user_reading_stats ON logs;

-- Drop function
DROP FUNCTION IF EXISTS update_user_reading_stats();

-- Drop indexes
DROP INDEX IF EXISTS idx_user_reading_stats_year_month;
DROP INDEX IF EXISTS idx_user_reading_stats_user_id;

-- Drop table
DROP TABLE IF EXISTS user_reading_stats;

-- Remove columns from users
ALTER TABLE users DROP COLUMN IF EXISTS reading_goal_year;
ALTER TABLE users DROP COLUMN IF EXISTS reading_goal;
ALTER TABLE users DROP COLUMN IF EXISTS bio;
ALTER TABLE users DROP COLUMN IF EXISTS banner_url;
ALTER TABLE users DROP COLUMN IF EXISTS favorite_book_ids;
