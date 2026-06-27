DROP TRIGGER IF EXISTS set_accounts_updated_at ON accounts;
DROP INDEX IF EXISTS idx_accounts_user_id;
DROP TABLE IF EXISTS accounts;