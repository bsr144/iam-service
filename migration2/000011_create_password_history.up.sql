-- ============================================================================
-- Migration: 000011_create_password_history
-- Description: Create password_history table to prevent password reuse
-- ============================================================================

CREATE TABLE password_history (
    -- Primary Key
    password_history_id UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    user_id             UUID NOT NULL,

    -- Password Hash
    password_hash       VARCHAR(255) NOT NULL,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_password_history_user FOREIGN KEY (user_id)
        REFERENCES users(user_id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_password_history_user_id ON password_history(user_id);
CREATE INDEX idx_password_history_user_recent ON password_history(user_id, created_at DESC);

-- Comments
COMMENT ON TABLE password_history IS 'Stores last N password hashes to prevent reuse (N configured per tenant)';
COMMENT ON COLUMN password_history.password_hash IS 'Bcrypt hash - compare new password against these';
