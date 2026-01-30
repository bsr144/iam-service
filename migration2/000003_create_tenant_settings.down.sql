-- ============================================================================
-- Migration: 000003_create_tenant_settings (DOWN)
-- ============================================================================

DROP TRIGGER IF EXISTS trg_tenant_settings_updated_at ON tenant_settings;
DROP TABLE IF EXISTS tenant_settings;
