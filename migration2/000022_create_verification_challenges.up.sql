-- ============================================================================
-- Migration: 000022_create_verification_challenges
-- Description: Create verification_challenges table for storing verification proof
-- ============================================================================

CREATE TABLE verification_challenges (
    -- Primary Key
    challenge_id        UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Key
    verification_id     UUID NOT NULL,

    -- Challenge Data
    challenge_type      VARCHAR(50) NOT NULL,
    challenge_hash      VARCHAR(255),
    challenge_data      JSONB,

    -- OTP-specific (for logging/debugging)
    otp_code_prefix     VARCHAR(10),
    otp_delivery_id     VARCHAR(255),

    -- Lifecycle
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at          TIMESTAMPTZ NOT NULL,
    used_at             TIMESTAMPTZ,

    -- Security
    is_used             BOOLEAN NOT NULL DEFAULT FALSE,

    -- Constraints
    CONSTRAINT fk_verification_challenges_verification FOREIGN KEY (verification_id)
        REFERENCES verifications(verification_id) ON DELETE CASCADE,
    CONSTRAINT chk_challenge_type CHECK (challenge_type IN (
        'otp', 'totp', 'pin', 'biometric_hash', 'liveness_token', 'webauthn_challenge'
    ))
);

-- Indexes
CREATE INDEX idx_verification_challenges_verification_id ON verification_challenges(verification_id);
CREATE INDEX idx_verification_challenges_active ON verification_challenges(verification_id)
    WHERE is_used = FALSE AND expires_at > NOW();

-- Comments
COMMENT ON TABLE verification_challenges IS 'Stores verification proof data (OTP codes, biometric hashes, etc.)';
COMMENT ON COLUMN verification_challenges.challenge_hash IS 'Bcrypt hash for OTP/PIN challenges';
COMMENT ON COLUMN verification_challenges.challenge_data IS 'JSONB for complex challenge types (WebAuthn, etc.)';
COMMENT ON COLUMN verification_challenges.otp_code_prefix IS 'First 2 digits for support purposes';
