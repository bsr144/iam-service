-- ============================================================================
-- Migration: 000004_create_branches
-- Description: Create branches table for organizational units within tenant
-- ============================================================================

CREATE TABLE branches (
    -- Primary Key
    branch_id           UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    tenant_id           UUID NOT NULL,

    -- Business Identifiers
    code                VARCHAR(50) NOT NULL,
    name                VARCHAR(255) NOT NULL,

    -- Branch Type
    is_headquarters     BOOLEAN NOT NULL DEFAULT FALSE,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    -- Optimistic Locking
    version             INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT fk_branches_tenant FOREIGN KEY (tenant_id)
        REFERENCES tenants(tenant_id) ON DELETE RESTRICT,
    CONSTRAINT uq_branches_tenant_code UNIQUE (tenant_id, code)
);

-- Trigger for updated_at
CREATE TRIGGER trg_branches_updated_at
    BEFORE UPDATE ON branches
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Indexes
CREATE INDEX idx_branches_tenant_id ON branches(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_branches_tenant_code_active ON branches(tenant_id, code) WHERE deleted_at IS NULL AND is_active = TRUE;
CREATE INDEX idx_branches_is_headquarters ON branches(tenant_id) WHERE is_headquarters = TRUE AND deleted_at IS NULL;

-- Comments
COMMENT ON TABLE branches IS 'Organizational units within a tenant for geographic or departmental scoping';
COMMENT ON COLUMN branches.code IS 'Unique identifier within tenant (e.g., HQ, BRANCH-001)';
COMMENT ON COLUMN branches.is_headquarters IS 'Indicates if this is the main/headquarters branch';
