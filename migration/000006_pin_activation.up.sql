CREATE TABLE pin_verification_logs (
    pin_verification_log_id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    result BOOLEAN NOT NULL,
    failure_reason VARCHAR(100),
    ip_address INET,
    user_agent TEXT,
    operation VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pin_verification_user ON pin_verification_logs(user_id, created_at);
CREATE INDEX idx_pin_verification_result ON pin_verification_logs(result, created_at);
CREATE INDEX idx_pin_verification_tenant ON pin_verification_logs(tenant_id, created_at);

COMMENT ON TABLE pin_verification_logs IS 'Audit log of all PIN verification attempts';
COMMENT ON COLUMN pin_verification_logs.operation IS 'The sensitive operation that required PIN verification';

CREATE TABLE user_activation_tracking (
    user_activation_tracking_id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id) ON DELETE CASCADE,


    admin_created BOOLEAN DEFAULT FALSE,
    admin_created_at TIMESTAMPTZ,
    admin_created_by UUID REFERENCES users(user_id) ON DELETE SET NULL,


    user_completed BOOLEAN DEFAULT FALSE,
    user_completed_at TIMESTAMPTZ,
    otp_verified_at TIMESTAMPTZ,
    profile_completed_at TIMESTAMPTZ,
    pin_set_at TIMESTAMPTZ,


    activated_at TIMESTAMPTZ,
    activation_method VARCHAR(50),


    status_history JSONB DEFAULT '[]'::JSONB,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_activation_user ON user_activation_tracking(user_id);
CREATE INDEX idx_activation_tenant ON user_activation_tracking(tenant_id);
CREATE INDEX idx_activation_status ON user_activation_tracking(activated_at);

COMMENT ON TABLE user_activation_tracking IS 'Tracks the dual registration flow (admin + participant) for audit trail';
COMMENT ON COLUMN user_activation_tracking.activation_method IS 'Which registration path was taken: admin_first, user_first, or simultaneous';
COMMENT ON COLUMN user_activation_tracking.status_history IS 'Array of status transition objects with timestamp and triggered_by';

CREATE TABLE admin_api_keys (
    admin_api_key_id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id) ON DELETE CASCADE,

    key_name VARCHAR(100) NOT NULL,
    key_hash VARCHAR(255) NOT NULL UNIQUE,
    key_prefix VARCHAR(20) NOT NULL,

    created_by UUID REFERENCES users(user_id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    expires_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    revoked_by UUID REFERENCES users(user_id) ON DELETE SET NULL,
    revoked_reason VARCHAR(255),

    last_used_at TIMESTAMPTZ,
    last_used_ip INET,

    ip_whitelist JSONB DEFAULT '[]'::JSONB,

    is_active BOOLEAN DEFAULT TRUE,

    CONSTRAINT unique_api_key_name UNIQUE (tenant_id, key_name)
);

CREATE INDEX idx_admin_api_keys_tenant ON admin_api_keys(tenant_id) WHERE is_active = TRUE;
CREATE INDEX idx_admin_api_keys_hash ON admin_api_keys(key_hash) WHERE is_active = TRUE;

COMMENT ON TABLE admin_api_keys IS 'API keys for secret endpoints (admin/approver creation via API)';
COMMENT ON COLUMN admin_api_keys.key_prefix IS 'First 8 characters of key for identification in logs (BGS_acme_ADMIN_)';
COMMENT ON COLUMN admin_api_keys.ip_whitelist IS 'JSON array of allowed IP addresses/CIDR ranges (e.g., ["192.168.1.0/24", "10.0.0.1"])';
