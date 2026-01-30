-- ============================================================================
-- Migration: 000019_create_password_reset_tokens
-- Description: Create password_reset_tokens table for password recovery
-- ============================================================================

CREATE TABLE password_reset_tokens (
    -- Primary Key
    password_reset_token_id UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    tenant_id           UUID NOT NULL,
    user_id             UUID NOT NULL,

    -- Token Data
    token_hash          VARCHAR(255) NOT NULL,

    -- Lifecycle
    expires_at          TIMESTAMPTZ NOT NULL,
    used_at             TIMESTAMPTZ,

    -- Context Tracking
    ip_address          INET,
    user_agent          TEXT,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_password_reset_tokens_tenant FOREIGN KEY (tenant_id)
        REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    CONSTRAINT fk_password_reset_tokens_user FOREIGN KEY (user_id)
        REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT uq_password_reset_tokens_hash UNIQUE (token_hash)
);

-- Indexes
CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);
CREATE INDEX idx_password_reset_tokens_pending ON password_reset_tokens(user_id)
    WHERE used_at IS NULL AND expires_at > NOW();

-- Comments
COMMENT ON TABLE password_reset_tokens IS 'Tokens for password reset requests';
COMMENT ON COLUMN password_reset_tokens.token_hash IS 'SHA-256 hash of the reset token';
