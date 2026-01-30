-- ============================================================================
-- Migration: 000024_create_pin_verification_logs
-- Description: Create pin_verification_logs table for PIN audit trail
-- ============================================================================

CREATE TABLE pin_verification_logs (
    -- Primary Key
    pin_verification_log_id UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    user_id             UUID NOT NULL,
    tenant_id           UUID NOT NULL,

    -- Verification Result
    result              BOOLEAN NOT NULL,
    failure_reason      VARCHAR(50),

    -- Operation Context
    operation           VARCHAR(100),

    -- Context Tracking
    ip_address          INET,
    user_agent          TEXT,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_pin_verification_logs_user FOREIGN KEY (user_id)
        REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_pin_verification_logs_tenant FOREIGN KEY (tenant_id)
        REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    CONSTRAINT chk_failure_reason CHECK (failure_reason IN (
        'invalid_pin', 'rate_limited', 'account_locked', 'pin_expired'
    ))
);

-- Indexes
CREATE INDEX idx_pin_verification_logs_user_id ON pin_verification_logs(user_id);
CREATE INDEX idx_pin_verification_logs_tenant_id ON pin_verification_logs(tenant_id);
CREATE INDEX idx_pin_verification_logs_recent ON pin_verification_logs(user_id, created_at DESC);
CREATE INDEX idx_pin_verification_logs_failed ON pin_verification_logs(user_id, created_at)
    WHERE result = FALSE;

-- Comments
COMMENT ON TABLE pin_verification_logs IS 'Audit log of PIN verification attempts';
COMMENT ON COLUMN pin_verification_logs.result IS 'Whether verification succeeded';
COMMENT ON COLUMN pin_verification_logs.failure_reason IS 'Reason for failure if result is FALSE';
COMMENT ON COLUMN pin_verification_logs.operation IS 'What operation required PIN verification';
