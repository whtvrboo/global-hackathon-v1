-- Drop triggers
DROP TRIGGER IF EXISTS update_logs_updated_at ON logs;
DROP TRIGGER IF EXISTS update_books_updated_at ON books;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order (respecting foreign key dependencies)
DROP TABLE IF EXISTS followers;
DROP TABLE IF EXISTS watchlists;
DROP TABLE IF EXISTS logs;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";

