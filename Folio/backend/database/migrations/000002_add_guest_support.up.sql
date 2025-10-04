-- Add guest user support
ALTER TABLE users ADD COLUMN is_guest BOOLEAN DEFAULT false;
ALTER TABLE users ADD COLUMN guest_session_id VARCHAR(255) UNIQUE;
ALTER TABLE users ADD COLUMN converted_at TIMESTAMPTZ;

-- Create index for guest session lookups
CREATE INDEX idx_users_guest_session_id ON users(guest_session_id);
CREATE INDEX idx_users_is_guest ON users(is_guest);

-- Update the users table to allow NULL google_id for guest users
ALTER TABLE users ALTER COLUMN google_id DROP NOT NULL;
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;

-- Add constraint to ensure either google_id or guest_session_id is present
ALTER TABLE users ADD CONSTRAINT check_user_identity 
  CHECK (
    (google_id IS NOT NULL AND is_guest = false) OR 
    (guest_session_id IS NOT NULL AND is_guest = true)
  );

-- Update the unique constraint on email to allow NULL for guests
DROP INDEX IF EXISTS idx_users_email;
CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE email IS NOT NULL;
