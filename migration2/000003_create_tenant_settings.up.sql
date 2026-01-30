-- ============================================================================
-- Migration: 000003_create_tenant_settings
-- Description: Create tenant_settings table for tenant-specific configuration
-- ============================================================================

CREATE TABLE tenant_settings (
    -- Primary Key
    tenant_setting_id   UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    tenant_id           UUID NOT NULL,

    -- Subscription & Limits
    subscription_tier   VARCHAR(50) NOT NULL DEFAULT 'standard',
    max_branches        INTEGER NOT NULL DEFAULT 10,
    max_employees       INTEGER NOT NULL DEFAULT 10000,

    -- Contact Information
    contact_email       VARCHAR(255),
    contact_phone       VARCHAR(50),
    contact_address     TEXT,

    -- Localization
    default_language    VARCHAR(10) NOT NULL DEFAULT 'en',
    timezone            VARCHAR(50) NOT NULL DEFAULT 'Asia/Jakarta',

    -- Security Policies (JSONB for flexibility)
    password_policy     JSONB NOT NULL DEFAULT '{
        "min_length": 8,
        "require_uppercase": true,
        "require_lowercase": true,
        "require_number": true,
        "require_special": true,
        "history_count": 5,
        "max_age_days": 90
    }'::jsonb,

    pin_policy          JSONB NOT NULL DEFAULT '{
        "length": 6,
        "max_attempts": 10,
        "lockout_duration_minutes": 30
    }'::jsonb,

    session_policy      JSONB NOT NULL DEFAULT '{
        "access_token_ttl_minutes": 30,
        "refresh_token_ttl_days": 7,
        "max_concurrent_sessions": 5
    }'::jsonb,

    -- Features
    approval_required   BOOLEAN NOT NULL DEFAULT true,
    mfa_required        BOOLEAN NOT NULL DEFAULT false,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_tenant_settings_tenant FOREIGN KEY (tenant_id)
        REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    CONSTRAINT uq_tenant_settings_tenant_id UNIQUE (tenant_id)
);

-- Trigger for updated_at
CREATE TRIGGER trg_tenant_settings_updated_at
    BEFORE UPDATE ON tenant_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE tenant_settings IS 'Tenant-specific configuration and policies';
COMMENT ON COLUMN tenant_settings.password_policy IS 'JSON object containing password policy configuration';
COMMENT ON COLUMN tenant_settings.pin_policy IS 'JSON object containing PIN policy configuration';
COMMENT ON COLUMN tenant_settings.session_policy IS 'JSON object containing session policy configuration';
