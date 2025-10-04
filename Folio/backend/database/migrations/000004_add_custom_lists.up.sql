-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create lists table
CREATE TABLE lists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_public BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for lists
CREATE INDEX idx_lists_user_id ON lists(user_id);
CREATE INDEX idx_lists_created_at ON lists(created_at DESC);
CREATE INDEX idx_lists_is_public ON lists(is_public);

-- Create list_items junction table
CREATE TABLE list_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    list_id UUID NOT NULL REFERENCES lists(id) ON DELETE CASCADE,
    book_id VARCHAR(255) NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    notes TEXT,
    item_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(list_id, book_id)
);

-- Create indexes for list_items
CREATE INDEX idx_list_items_list_id ON list_items(list_id);
CREATE INDEX idx_list_items_book_id ON list_items(book_id);
CREATE INDEX idx_list_items_order ON list_items(list_id, item_order);

-- Add items_count to lists table
ALTER TABLE lists ADD COLUMN items_count INTEGER DEFAULT 0;

-- Create function to update items_count
CREATE OR REPLACE FUNCTION update_list_items_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE lists SET items_count = items_count + 1 WHERE id = NEW.list_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE lists SET items_count = GREATEST(items_count - 1, 0) WHERE id = OLD.list_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for items_count
CREATE TRIGGER trigger_update_list_items_count
AFTER INSERT OR DELETE ON list_items
FOR EACH ROW EXECUTE FUNCTION update_list_items_count();

