-- ============================================================================
-- Migration: 000020_create_refresh_tokens
-- Description: Create refresh_tokens table for JWT refresh token management
-- ============================================================================

CREATE TABLE refresh_tokens (
    -- Primary Key
    refresh_token_id    UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    tenant_id           UUID NOT NULL,
    user_id             UUID NOT NULL,

    -- Token Data
    token_hash          VARCHAR(255) NOT NULL,
    token_family        UUID NOT NULL,

    -- Lifecycle
    expires_at          TIMESTAMPTZ NOT NULL,
    revoked_at          TIMESTAMPTZ,
    revoked_reason      VARCHAR(100),
    replaced_by_token_id UUID,

    -- Context Tracking
    ip_address          INET,
    user_agent          TEXT,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_refresh_tokens_tenant FOREIGN KEY (tenant_id)
        REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id)
        REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT uq_refresh_tokens_hash UNIQUE (token_hash)
);

-- Indexes
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token_family ON refresh_tokens(token_family);
CREATE INDEX idx_refresh_tokens_active ON refresh_tokens(user_id)
    WHERE revoked_at IS NULL AND expires_at > NOW();
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens(expires_at)
    WHERE revoked_at IS NULL;

-- Comments
COMMENT ON TABLE refresh_tokens IS 'JWT refresh tokens for session management';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'SHA-256 hash of the refresh token';
COMMENT ON COLUMN refresh_tokens.token_family IS 'Family ID for token rotation detection';
COMMENT ON COLUMN refresh_tokens.replaced_by_token_id IS 'Reference to the new token that replaced this one';
