-- Remove list metadata fields
ALTER TABLE lists DROP COLUMN IF EXISTS header_image_url;
ALTER TABLE lists DROP COLUMN IF EXISTS theme_color;

-- Drop indexes
DROP INDEX IF EXISTS idx_lists_theme_color;
DROP INDEX IF EXISTS idx_lists_header_image;
