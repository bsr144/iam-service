-- ============================================================================
-- Migration: 000002_create_tenants
-- Description: Create tenants table - root table for multi-tenant organization
-- Requires: PostgreSQL 18+ for native uuidv7()
-- ============================================================================

CREATE TABLE tenants (
    -- Primary Key (UUIDv7 for time-ordered, sortable IDs)
    tenant_id           UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Business Identifiers
    name                VARCHAR(255) NOT NULL,
    slug                VARCHAR(100) NOT NULL,
    database_name       VARCHAR(100),

    -- Status Management
    status              VARCHAR(20) NOT NULL DEFAULT 'active',

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    -- Optimistic Locking
    version             INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT uq_tenants_slug UNIQUE (slug),
    CONSTRAINT uq_tenants_database_name UNIQUE (database_name),
    CONSTRAINT chk_tenants_status CHECK (status IN ('active', 'inactive', 'suspended'))
);

-- Trigger for updated_at
CREATE TRIGGER trg_tenants_updated_at
    BEFORE UPDATE ON tenants
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Indexes
CREATE INDEX idx_tenants_slug_active ON tenants(slug) WHERE deleted_at IS NULL AND status = 'active';
CREATE INDEX idx_tenants_status ON tenants(status) WHERE deleted_at IS NULL;

-- Comments
COMMENT ON TABLE tenants IS 'Root table for multi-tenant organization management';
COMMENT ON COLUMN tenants.tenant_id IS 'UUIDv7 primary key - time-ordered for better index performance';
COMMENT ON COLUMN tenants.slug IS 'URL-safe unique identifier for the tenant';
COMMENT ON COLUMN tenants.status IS 'Tenant lifecycle status: active, inactive, suspended';
COMMENT ON COLUMN tenants.deleted_at IS 'Soft delete timestamp - NULL means active';
COMMENT ON COLUMN tenants.version IS 'Optimistic lock version - managed by application';
