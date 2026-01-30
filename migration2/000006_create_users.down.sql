-- ============================================================================
-- Migration: 000006_create_users (DOWN)
-- ============================================================================

DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP TABLE IF EXISTS users;
