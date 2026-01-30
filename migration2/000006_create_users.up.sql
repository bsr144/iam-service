-- ============================================================================
-- Migration: 000006_create_users
-- Description: Create users table - core identity (normalized design)
-- ============================================================================

CREATE TABLE users (
    -- Primary Key
    user_id             UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    tenant_id           UUID,
    branch_id           UUID,

    -- Core Identity
    email               VARCHAR(255) NOT NULL,
    email_verified      BOOLEAN NOT NULL DEFAULT FALSE,

    -- User Lifecycle Status
    is_active           BOOLEAN NOT NULL DEFAULT FALSE,
    is_service_account  BOOLEAN NOT NULL DEFAULT FALSE,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    -- Optimistic Locking
    version             INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT fk_users_tenant FOREIGN KEY (tenant_id)
        REFERENCES tenants(tenant_id) ON DELETE RESTRICT,
    CONSTRAINT fk_users_branch FOREIGN KEY (branch_id)
        REFERENCES branches(branch_id) ON DELETE SET NULL,
    CONSTRAINT uq_users_tenant_email UNIQUE (tenant_id, email),
    CONSTRAINT chk_users_email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

-- Trigger for updated_at
CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Indexes
CREATE INDEX idx_users_tenant_id ON users(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_branch_id ON users(branch_id) WHERE deleted_at IS NULL AND branch_id IS NOT NULL;
CREATE INDEX idx_users_tenant_email_active ON users(tenant_id, LOWER(email))
    WHERE deleted_at IS NULL AND is_active = TRUE;
CREATE INDEX idx_users_email ON users(LOWER(email)) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_service_accounts ON users(tenant_id)
    WHERE is_service_account = TRUE AND deleted_at IS NULL;

-- Comments
COMMENT ON TABLE users IS 'Core user identity - minimal fields, authentication/profile/security in related tables';
COMMENT ON COLUMN users.email IS 'Login identifier, unique within tenant, stored lowercase';
COMMENT ON COLUMN users.is_active IS 'Whether user can log in';
COMMENT ON COLUMN users.is_service_account IS 'Service accounts for API access (non-human)';
