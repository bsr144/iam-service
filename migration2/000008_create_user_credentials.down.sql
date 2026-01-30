-- ============================================================================
-- Migration: 000008_create_user_credentials (DOWN)
-- ============================================================================

DROP TRIGGER IF EXISTS trg_user_credentials_updated_at ON user_credentials;
DROP TABLE IF EXISTS user_credentials;
