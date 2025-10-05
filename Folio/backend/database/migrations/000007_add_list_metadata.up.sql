-- Add list metadata fields for better visual presentation
ALTER TABLE lists ADD COLUMN header_image_url TEXT;
ALTER TABLE lists ADD COLUMN theme_color VARCHAR(7) DEFAULT '#6366f1';

-- Create indexes for better performance
CREATE INDEX idx_lists_theme_color ON lists(theme_color);
CREATE INDEX idx_lists_header_image ON lists(header_image_url) WHERE header_image_url IS NOT NULL;
