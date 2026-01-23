CREATE TABLE verifications (
    verification_id UUID PRIMARY KEY DEFAULT uuidv7(),

    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,

    purpose VARCHAR(100) NOT NULL,

    verification_method VARCHAR(50) NOT NULL,

    delivery_target VARCHAR(255) NOT NULL,
    delivery_channel VARCHAR(50),
    delivery_status VARCHAR(50) DEFAULT 'pending',
    delivery_attempts INT DEFAULT 0,
    last_delivery_attempt_at TIMESTAMPTZ,
    delivery_error TEXT,

    max_attempts INT NOT NULL DEFAULT 3,
    attempts_used INT NOT NULL DEFAULT 0,
    locked_until TIMESTAMPTZ,

    status VARCHAR(50) NOT NULL DEFAULT 'pending',

    ip_address INET,
    user_agent TEXT,
    metadata JSONB DEFAULT '{}',

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,
    verified_at TIMESTAMPTZ,

    verification_result JSONB,
    failure_reason TEXT,

    CONSTRAINT verifications_status_check CHECK (status IN ('pending', 'sent', 'verified', 'failed', 'expired', 'locked')),
    CONSTRAINT verifications_attempts_check CHECK (attempts_used <= max_attempts),
    CONSTRAINT verifications_delivery_status_check CHECK (delivery_status IN ('pending', 'sent', 'failed', 'delivered', 'bounced'))
);
CREATE INDEX idx_verifications_entity ON verifications(entity_type, entity_id)
    WHERE status IN ('pending', 'sent');

CREATE INDEX idx_verifications_purpose ON verifications(purpose, status);

CREATE INDEX idx_verifications_expires ON verifications(expires_at)
    WHERE status IN ('pending', 'sent');

CREATE INDEX idx_verifications_status ON verifications(status)
    WHERE status IN ('pending', 'sent');

CREATE INDEX idx_verifications_created_at ON verifications(created_at);
COMMENT ON TABLE verifications IS 'Generic verification framework - purpose and method agnostic';
COMMENT ON COLUMN verifications.entity_type IS 'What resource needs verification (registration, user, session, etc.)';
COMMENT ON COLUMN verifications.entity_id IS 'Polymorphic reference to the entity being verified';
COMMENT ON COLUMN verifications.purpose IS 'Why verification is needed (register, reset_password, step_up_auth, etc.)';
COMMENT ON COLUMN verifications.verification_method IS 'How to verify (otp_email, pin, biometric, liveness, etc.)';
COMMENT ON COLUMN verifications.delivery_target IS 'Where to deliver the challenge (email, phone, device)';
COMMENT ON COLUMN verifications.max_attempts IS 'Maximum verification attempts before locking';

CREATE TABLE verification_challenges (
    challenge_id UUID PRIMARY KEY DEFAULT uuidv7(),
    verification_id UUID NOT NULL REFERENCES verifications(verification_id) ON DELETE CASCADE,

    challenge_type VARCHAR(50) NOT NULL,

    challenge_hash TEXT,
    challenge_data JSONB,

    otp_code_prefix VARCHAR(2),
    otp_delivery_id VARCHAR(255),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,

    is_used BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT verification_challenges_type_check CHECK (
        challenge_type IN ('otp', 'totp', 'pin', 'biometric_hash', 'liveness_token', 'webauthn_challenge')
    )
);
CREATE UNIQUE INDEX idx_verification_challenge_active ON verification_challenges(verification_id)
    WHERE is_used = FALSE AND expires_at > CURRENT_TIMESTAMP;
CREATE INDEX idx_verification_challenge_expires ON verification_challenges(expires_at)
    WHERE is_used = FALSE;

CREATE INDEX idx_verification_challenge_verification ON verification_challenges(verification_id);

CREATE INDEX idx_verification_challenge_created_at ON verification_challenges(created_at);
COMMENT ON TABLE verification_challenges IS 'Concrete verification proof - stores OTP, PIN, biometric data, etc.';
COMMENT ON COLUMN verification_challenges.challenge_hash IS 'bcrypt/SHA256 hash of OTP/PIN/token for secure comparison';
COMMENT ON COLUMN verification_challenges.challenge_data IS 'JSONB for complex challenges (biometric templates, webauthn data)';
COMMENT ON COLUMN verification_challenges.otp_code_prefix IS 'First 2 digits for logging/debugging (NOT security-critical)';
COMMENT ON COLUMN verification_challenges.is_used IS 'Prevents challenge reuse (replay attack protection)';
