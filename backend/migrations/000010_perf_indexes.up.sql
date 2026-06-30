-- Migration 000010: performance indexes for reports aggregations.
--
-- All existing transactions queries used the (user_id, date) composite
-- index which already filtered efficiently. The aggregations
-- (SumByCategory, SumByAccount, MonthlyTrend) need to group by
-- category_id or account_id AFTER filtering by user+date — adding
-- composite indexes that cover (user_id, group_key, date) lets
-- Postgres do an index-only scan when the rowset is large.

CREATE INDEX IF NOT EXISTS idx_transactions_user_cat_date
  ON transactions (user_id, category_id, date DESC)
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_transactions_user_acct_date
  ON transactions (user_id, account_id, date DESC)
  WHERE deleted_at IS NULL;

-- Travel settlement greedy algorithm — covers the CTE in
-- (group_id, paid_by) and (group_id, user_id) joins.
CREATE INDEX IF NOT EXISTS idx_travel_expenses_group_paid
  ON travel_expenses (group_id, paid_by);

CREATE INDEX IF NOT EXISTS idx_travel_shares_expense_user
  ON travel_expense_shares (expense_id, user_id);

-- Recurring: lookup by user + next_run_at for /recurring/generate-today.
-- is_active is the boolean column (not "enabled").
CREATE INDEX IF NOT EXISTS idx_recurring_rules_user_next
  ON recurring_rules (user_id, next_run_date)
  WHERE is_active = true;
