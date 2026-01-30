-- ============================================================================
-- Migration: 000009_create_user_security_states
-- Description: Create user_security_states table - login tracking and security state (1:1 with users)
-- ============================================================================

CREATE TABLE user_security_states (
    -- Primary Key (same as user_id for 1:1 relationship)
    user_id                 UUID PRIMARY KEY,

    -- Failed Attempt Tracking
    failed_login_attempts   INTEGER NOT NULL DEFAULT 0,
    failed_pin_attempts     INTEGER NOT NULL DEFAULT 0,
    locked_until            TIMESTAMPTZ,

    -- Login History
    last_login_at           TIMESTAMPTZ,
    last_login_ip           INET,
    last_login_user_agent   TEXT,

    -- Verification Status
    email_verified_at       TIMESTAMPTZ,
    phone_verified_at       TIMESTAMPTZ,

    -- Security Events
    password_reset_at       TIMESTAMPTZ,
    last_activity_at        TIMESTAMPTZ,

    -- Audit Fields
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_user_security_states_user FOREIGN KEY (user_id)
        REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT chk_user_security_failed_login CHECK (failed_login_attempts >= 0),
    CONSTRAINT chk_user_security_failed_pin CHECK (failed_pin_attempts >= 0)
);

-- Trigger for updated_at
CREATE TRIGGER trg_user_security_states_updated_at
    BEFORE UPDATE ON user_security_states
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Indexes
CREATE INDEX idx_user_security_locked ON user_security_states(locked_until)
    WHERE locked_until IS NOT NULL AND locked_until > NOW();
CREATE INDEX idx_user_security_last_login ON user_security_states(last_login_at DESC)
    WHERE last_login_at IS NOT NULL;

-- Comments
COMMENT ON TABLE user_security_states IS 'Security state that changes frequently during authentication';
COMMENT ON COLUMN user_security_states.failed_login_attempts IS 'Counter for failed password attempts, resets on success';
COMMENT ON COLUMN user_security_states.failed_pin_attempts IS 'Counter for failed PIN attempts, resets on success';
COMMENT ON COLUMN user_security_states.locked_until IS 'Account is locked until this timestamp (NULL = not locked)';
