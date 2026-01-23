CREATE TABLE registrations (
    registration_id UUID PRIMARY KEY DEFAULT uuidv7(),

    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(branch_id) ON DELETE SET NULL,
    email VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,

    user_agent TEXT,
    ip_address INET,
    referrer TEXT,
    metadata JSONB DEFAULT '{}',

    status VARCHAR(50) NOT NULL DEFAULT 'pending_verification',

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL '24 hours'),
    verified_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,

    user_id UUID REFERENCES users(user_id) ON DELETE SET NULL,
    cancellation_reason TEXT,

    CONSTRAINT registrations_status_check CHECK (status IN ('pending_verification', 'verified', 'completed', 'expired', 'cancelled')),

    CONSTRAINT registrations_unique_email_per_tenant_pending UNIQUE (tenant_id, email)
);
CREATE UNIQUE INDEX idx_registrations_unique_pending ON registrations(tenant_id, email)
    WHERE status IN ('pending_verification', 'verified');
CREATE INDEX idx_registrations_status ON registrations(status)
    WHERE status IN ('pending_verification', 'verified');

CREATE INDEX idx_registrations_expires ON registrations(expires_at)
    WHERE status IN ('pending_verification', 'verified');

CREATE INDEX idx_registrations_email ON registrations(email)
    WHERE status = 'pending_verification';

CREATE INDEX idx_registrations_tenant ON registrations(tenant_id);

CREATE INDEX idx_registrations_created_at ON registrations(created_at);
COMMENT ON TABLE registrations IS 'Registration process resource - exists before user creation to prevent email hijacking';
COMMENT ON COLUMN registrations.status IS 'Lifecycle: pending_verification → verified → completed';
COMMENT ON COLUMN registrations.expires_at IS 'Registration expires after 24h if not completed';
COMMENT ON COLUMN registrations.password_hash IS 'Temporarily stored until registration completes, then moved to user_credentials';
COMMENT ON COLUMN registrations.user_id IS 'Set after registration completes and user is created';
COMMENT ON COLUMN registrations.metadata IS 'Flexible JSONB for campaign tracking, device info, etc.';
