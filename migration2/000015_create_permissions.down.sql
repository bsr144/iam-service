-- ============================================================================
-- Migration: 000015_create_permissions (DOWN)
-- ============================================================================

DROP TRIGGER IF EXISTS trg_permissions_updated_at ON permissions;
DROP TABLE IF EXISTS permissions;
