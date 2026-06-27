DROP TRIGGER IF EXISTS set_budgets_updated_at ON budgets;
DROP INDEX IF EXISTS idx_budgets_period;
DROP INDEX IF EXISTS idx_budgets_user_category;
DROP INDEX IF EXISTS idx_budgets_user_id;
DROP TABLE IF EXISTS budgets;