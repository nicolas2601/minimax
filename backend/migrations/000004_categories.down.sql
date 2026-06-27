DROP TRIGGER IF EXISTS set_categories_updated_at ON categories;
DROP INDEX IF EXISTS idx_categories_parent_id;
DROP INDEX IF EXISTS idx_categories_user_id;
DROP TABLE IF EXISTS categories;