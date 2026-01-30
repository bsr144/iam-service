-- ============================================================================
-- Migration: 000017_create_user_roles (DOWN)
-- ============================================================================

DROP TRIGGER IF EXISTS trg_user_roles_updated_at ON user_roles;
DROP TABLE IF EXISTS user_roles;
