-- ============================================================================
-- Migration: 000021_create_verifications
-- Description: Create verifications table for generic verification framework
-- ============================================================================

CREATE TABLE verifications (
    -- Primary Key
    verification_id     UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Polymorphic Reference (what is being verified)
    entity_type         VARCHAR(50) NOT NULL,
    entity_id           UUID NOT NULL,

    -- Verification Context
    purpose             VARCHAR(50) NOT NULL,
    verification_method VARCHAR(50) NOT NULL,

    -- Delivery Information
    delivery_target     VARCHAR(255) NOT NULL,
    delivery_channel    VARCHAR(20),
    delivery_status     VARCHAR(20) NOT NULL DEFAULT 'pending',
    delivery_attempts   INTEGER NOT NULL DEFAULT 0,
    last_delivery_attempt_at TIMESTAMPTZ,
    delivery_error      TEXT,

    -- Security Constraints
    max_attempts        INTEGER NOT NULL DEFAULT 3,
    attempts_used       INTEGER NOT NULL DEFAULT 0,
    locked_until        TIMESTAMPTZ,

    -- State
    status              VARCHAR(20) NOT NULL DEFAULT 'pending',

    -- Context Tracking
    ip_address          INET,
    user_agent          TEXT,
    metadata            JSONB NOT NULL DEFAULT '{}',

    -- Lifecycle
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at          TIMESTAMPTZ NOT NULL,
    verified_at         TIMESTAMPTZ,

    -- Result
    verification_result JSONB,
    failure_reason      TEXT,

    -- Constraints
    CONSTRAINT chk_entity_type CHECK (entity_type IN (
        'registration', 'user', 'password_reset', 'email_change', 'step_up_auth', 'session'
    )),
    CONSTRAINT chk_purpose CHECK (purpose IN (
        'register', 'reset_password', 'change_email', 'sensitive_operation', 'login_mfa', 'step_up'
    )),
    CONSTRAINT chk_verification_method CHECK (verification_method IN (
        'otp_email', 'otp_sms', 'pin', 'totp', 'biometric', 'liveness', 'webauthn'
    )),
    CONSTRAINT chk_delivery_channel CHECK (delivery_channel IN (
        'email', 'sms', 'push', 'app', 'device'
    )),
    CONSTRAINT chk_delivery_status CHECK (delivery_status IN (
        'pending', 'sent', 'failed', 'delivered', 'bounced'
    )),
    CONSTRAINT chk_status CHECK (status IN (
        'pending', 'sent', 'verified', 'failed', 'expired', 'locked'
    )),
    CONSTRAINT chk_attempts CHECK (attempts_used >= 0 AND attempts_used <= max_attempts)
);

-- Indexes
CREATE INDEX idx_verifications_entity ON verifications(entity_type, entity_id);
CREATE INDEX idx_verifications_pending ON verifications(entity_id, purpose)
    WHERE status IN ('pending', 'sent') AND expires_at > NOW();
CREATE INDEX idx_verifications_expires ON verifications(expires_at)
    WHERE status IN ('pending', 'sent');

-- Comments
COMMENT ON TABLE verifications IS 'Generic verification framework supporting multiple methods';
COMMENT ON COLUMN verifications.entity_type IS 'Type of entity being verified';
COMMENT ON COLUMN verifications.entity_id IS 'ID of the entity being verified';
COMMENT ON COLUMN verifications.purpose IS 'Why verification is needed';
COMMENT ON COLUMN verifications.verification_method IS 'How to verify (otp, pin, biometric, etc.)';
