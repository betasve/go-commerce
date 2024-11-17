-- Drop the trigger
DROP TRIGGER set_updated_at ON users;

-- Drop the trigger function
DROP FUNCTION update_updated_at_column();

-- Drop the users table's email index
DROP INDEX IF EXISTS idx_users_email;

-- Drop the users table
DROP TABLE users;
