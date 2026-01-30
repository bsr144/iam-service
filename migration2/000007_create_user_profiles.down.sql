-- ============================================================================
-- Migration: 000007_create_user_profiles (DOWN)
-- ============================================================================

DROP TRIGGER IF EXISTS trg_user_profiles_updated_at ON user_profiles;
DROP TABLE IF EXISTS user_profiles;
