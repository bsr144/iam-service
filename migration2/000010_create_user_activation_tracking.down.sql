-- ============================================================================
-- Migration: 000010_create_user_activation_tracking (DOWN)
-- ============================================================================

DROP TRIGGER IF EXISTS trg_user_activation_tracking_updated_at ON user_activation_tracking;
DROP TABLE IF EXISTS user_activation_tracking;
