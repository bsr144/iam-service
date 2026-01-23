CREATE TABLE auth_logs (
    auth_log_id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(user_id) ON DELETE SET NULL,
    event_type VARCHAR(50) NOT NULL,
    email VARCHAR(255),
    ip_address INET,
    user_agent TEXT,
    mfa_method VARCHAR(20),
    failure_reason VARCHAR(255),
    metadata JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_auth_logs_tenant ON auth_logs(tenant_id, created_at);
CREATE INDEX idx_auth_logs_user ON auth_logs(user_id, created_at);
CREATE INDEX idx_auth_logs_event ON auth_logs(event_type, created_at);

COMMENT ON TABLE auth_logs IS 'Immutable log of all authentication events';

CREATE TABLE permission_checks (
    permission_check_id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(user_id) ON DELETE SET NULL,
    permission_code VARCHAR(100) NOT NULL,
    resource_id UUID,
    resource_type VARCHAR(50),
    branch_id UUID,
    result BOOLEAN NOT NULL,
    metadata JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_permission_checks_tenant ON permission_checks(tenant_id, created_at);
CREATE INDEX idx_permission_checks_user ON permission_checks(user_id, created_at);
CREATE INDEX idx_permission_checks_result ON permission_checks(result, created_at);

COMMENT ON TABLE permission_checks IS 'Sampled log of permission checks (all denials + 10% of approvals)';

CREATE TABLE admin_audit_logs (
    admin_audit_log_id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID,
    before_state JSONB,
    after_state JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_admin_audit_tenant ON admin_audit_logs(tenant_id, created_at);
CREATE INDEX idx_admin_audit_user ON admin_audit_logs(user_id, created_at);
CREATE INDEX idx_admin_audit_entity ON admin_audit_logs(entity_type, entity_id, created_at);

COMMENT ON TABLE admin_audit_logs IS 'Immutable log of administrative actions with before/after state';
