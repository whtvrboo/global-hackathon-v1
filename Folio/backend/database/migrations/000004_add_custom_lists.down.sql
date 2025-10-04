-- Drop trigger
DROP TRIGGER IF EXISTS trigger_update_list_items_count ON list_items;

-- Drop function
DROP FUNCTION IF EXISTS update_list_items_count();

-- Drop tables
DROP TABLE IF EXISTS list_items;
DROP TABLE IF EXISTS lists;

