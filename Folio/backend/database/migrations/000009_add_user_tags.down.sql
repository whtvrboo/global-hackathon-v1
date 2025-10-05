-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_update_user_tags ON annotations;
DROP FUNCTION IF EXISTS update_user_tags();

-- Drop user_tags table
DROP TABLE IF EXISTS user_tags;
