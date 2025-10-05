-- Create annotations table for capturing thoughts, highlights, and notes
CREATE TABLE IF NOT EXISTS annotations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id VARCHAR(255) REFERENCES books(id) ON DELETE CASCADE, -- Nullable for unassociated notes
    log_id UUID REFERENCES logs(id) ON DELETE SET NULL,

    type VARCHAR(20) NOT NULL CHECK (type IN ('highlight', 'note')),
    content TEXT NOT NULL,
    context TEXT, -- For storing surrounding text of a highlight
    page_number INTEGER,
    tags TEXT[] DEFAULT '{}',
    source VARCHAR(50) DEFAULT 'quick-capture', -- 'quick-capture', 'list-curation', 'import'

    -- For intelligent association
    is_associated BOOLEAN DEFAULT false, -- True if linked to a book
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Critical Indexes for Performance
CREATE INDEX idx_annotations_user_book ON annotations(user_id, book_id);
CREATE INDEX idx_annotations_user_unassociated ON annotations(user_id) WHERE is_associated = false;
CREATE INDEX idx_annotations_tags ON annotations USING GIN(tags);
CREATE INDEX idx_annotations_content_fts ON annotations USING GIN (to_tsvector('english', content));
CREATE INDEX idx_annotations_created_at ON annotations(created_at DESC);
CREATE INDEX idx_annotations_user_created ON annotations(user_id, created_at DESC);

-- Trigger for updating updated_at
CREATE TRIGGER update_annotations_updated_at BEFORE UPDATE ON annotations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
