DROP TRIGGER IF EXISTS set_users_updated_at ON users;
DROP FUNCTION IF EXISTS trigger_set_updated_at();
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;