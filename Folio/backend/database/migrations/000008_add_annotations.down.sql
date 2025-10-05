-- Drop annotations table and related indexes
DROP TRIGGER IF EXISTS update_annotations_updated_at ON annotations;
DROP INDEX IF EXISTS idx_annotations_user_created;
DROP INDEX IF EXISTS idx_annotations_created_at;
DROP INDEX IF EXISTS idx_annotations_content_fts;
DROP INDEX IF EXISTS idx_annotations_tags;
DROP INDEX IF EXISTS idx_annotations_user_unassociated;
DROP INDEX IF EXISTS idx_annotations_user_book;
DROP TABLE IF EXISTS annotations;
