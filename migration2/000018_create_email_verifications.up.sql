-- ============================================================================
-- Migration: 000018_create_email_verifications
-- Description: Create email_verifications table for OTP-based email verification
-- ============================================================================

CREATE TABLE email_verifications (
    -- Primary Key
    email_verification_id UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    tenant_id           UUID NOT NULL,
    user_id             UUID NOT NULL,

    -- Verification Details
    email               VARCHAR(255) NOT NULL,
    otp_code            VARCHAR(10) NOT NULL,
    otp_hash            VARCHAR(255) NOT NULL,
    otp_type            VARCHAR(50) NOT NULL DEFAULT 'email_verification',

    -- Lifecycle
    expires_at          TIMESTAMPTZ NOT NULL,
    verified_at         TIMESTAMPTZ,

    -- Context Tracking
    ip_address          INET,
    user_agent          TEXT,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_email_verifications_tenant FOREIGN KEY (tenant_id)
        REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    CONSTRAINT fk_email_verifications_user FOREIGN KEY (user_id)
        REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT chk_otp_type CHECK (otp_type IN (
        'email_verification', 'registration', 'pin_reset',
        'password_reset', 'admin_invitation'
    ))
);

-- Indexes
CREATE INDEX idx_email_verifications_user_id ON email_verifications(user_id);
CREATE INDEX idx_email_verifications_email ON email_verifications(LOWER(email));
CREATE INDEX idx_email_verifications_pending ON email_verifications(user_id, otp_type)
    WHERE verified_at IS NULL AND expires_at > NOW();

-- Comments
COMMENT ON TABLE email_verifications IS 'OTP-based email verification records';
COMMENT ON COLUMN email_verifications.otp_hash IS 'Bcrypt hash of OTP code';
COMMENT ON COLUMN email_verifications.otp_type IS 'Purpose of the OTP';
