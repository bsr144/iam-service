CREATE TABLE saml_configurations (
    saml_configuration_id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    idp_entity_id VARCHAR(500) NOT NULL,
    idp_sso_url VARCHAR(500) NOT NULL,
    idp_slo_url VARCHAR(500),
    idp_certificate TEXT NOT NULL,
    sp_entity_id VARCHAR(500) NOT NULL,
    sp_acs_url VARCHAR(500) NOT NULL,
    sp_slo_url VARCHAR(500),
    attribute_mapping JSONB NOT NULL DEFAULT '{
        "email": "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress",
        "name": "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name",
        "roles": "http://schemas.microsoft.com/ws/2008/06/identity/claims/role"
    }'::JSONB,
    role_mapping JSONB DEFAULT '{}'::JSONB,
    auto_provision_users BOOLEAN DEFAULT TRUE,
    default_branch_id UUID REFERENCES branches(branch_id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT unique_saml_per_tenant UNIQUE (tenant_id)
);

CREATE INDEX idx_saml_tenant ON saml_configurations(tenant_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE saml_configurations IS 'SAML 2.0 SSO configuration per tenant (Future Stage)';
COMMENT ON COLUMN saml_configurations.attribute_mapping IS 'Maps SAML assertion attributes to user fields';
COMMENT ON COLUMN saml_configurations.role_mapping IS 'Maps SAML groups to internal role codes';
COMMENT ON COLUMN saml_configurations.auto_provision_users IS 'Automatically create users on first SSO login';

CREATE TABLE mfa_devices (
    mfa_device_id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    device_type VARCHAR(20) NOT NULL CHECK (device_type IN ('totp', 'backup_codes')),
    device_name VARCHAR(100),
    secret_encrypted TEXT NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_mfa_devices_user ON mfa_devices(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_mfa_devices_type ON mfa_devices(device_type) WHERE deleted_at IS NULL;

COMMENT ON TABLE mfa_devices IS 'TOTP devices and backup codes for multi-factor authentication (Future Stage)';
COMMENT ON COLUMN mfa_devices.secret_encrypted IS 'AES-256 encrypted TOTP secret (for totp) or JSON array of backup codes (for backup_codes)';
