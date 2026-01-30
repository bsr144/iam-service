-- ============================================================================
-- Migration: 000008_create_user_credentials
-- Description: Create user_credentials table - authentication credentials (1:1 with users)
-- ============================================================================

CREATE TABLE user_credentials (
    -- Primary Key (same as user_id for 1:1 relationship)
    user_id                 UUID PRIMARY KEY,

    -- Password Authentication
    password_hash           VARCHAR(255),
    password_set_at         TIMESTAMPTZ,
    password_changed_at     TIMESTAMPTZ,
    password_expires_at     TIMESTAMPTZ,

    -- PIN Authentication (6-digit second factor)
    pin_hash                VARCHAR(255),
    pin_set_at              TIMESTAMPTZ,
    pin_changed_at          TIMESTAMPTZ,
    pin_expires_at          TIMESTAMPTZ,

    -- OAuth (Google)
    google_id               VARCHAR(255),

    -- Audit Fields
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_user_credentials_user FOREIGN KEY (user_id)
        REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT uq_user_credentials_google_id UNIQUE (google_id)
);

-- Trigger for updated_at
CREATE TRIGGER trg_user_credentials_updated_at
    BEFORE UPDATE ON user_credentials
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Indexes
CREATE INDEX idx_user_credentials_google_id ON user_credentials(google_id)
    WHERE google_id IS NOT NULL;
CREATE INDEX idx_user_credentials_password_expires ON user_credentials(password_expires_at)
    WHERE password_expires_at IS NOT NULL;

-- Comments
COMMENT ON TABLE user_credentials IS 'Authentication credentials separated from user identity';
COMMENT ON COLUMN user_credentials.password_hash IS 'Bcrypt hash with cost factor 12';
COMMENT ON COLUMN user_credentials.pin_hash IS 'Bcrypt hash of 6-digit PIN for 2FA';
COMMENT ON COLUMN user_credentials.google_id IS 'Google OAuth subject ID for social login';
