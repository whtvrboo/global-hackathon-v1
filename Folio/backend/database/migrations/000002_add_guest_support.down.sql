-- Rollback guest user support
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_user_identity;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_guest_session_id;
DROP INDEX IF EXISTS idx_users_is_guest;

ALTER TABLE users DROP COLUMN IF EXISTS is_guest;
ALTER TABLE users DROP COLUMN IF EXISTS guest_session_id;
ALTER TABLE users DROP COLUMN IF EXISTS converted_at;

-- Restore NOT NULL constraints
ALTER TABLE users ALTER COLUMN google_id SET NOT NULL;
ALTER TABLE users ALTER COLUMN email SET NOT NULL;

-- Restore unique index on email
CREATE UNIQUE INDEX idx_users_email ON users(email);
