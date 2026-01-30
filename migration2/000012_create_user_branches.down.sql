-- ============================================================================
-- Migration: 000012_create_user_branches (DOWN)
-- ============================================================================

DROP TRIGGER IF EXISTS trg_user_branches_updated_at ON user_branches;
DROP TABLE IF EXISTS user_branches;
