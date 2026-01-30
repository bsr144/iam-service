-- ============================================================================
-- Migration: 000013_create_applications
-- Description: Create applications table for registered apps that use IAM
-- ============================================================================

CREATE TABLE applications (
    -- Primary Key
    application_id      UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    tenant_id           UUID NOT NULL,

    -- Business Identifiers
    code                VARCHAR(50) NOT NULL,
    name                VARCHAR(255) NOT NULL,
    description         TEXT,

    -- Application Settings
    settings            JSONB NOT NULL DEFAULT '{}',

    -- Status Management
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    -- Optimistic Locking
    version             INTEGER NOT NULL DEFAULT 1,

    -- Constraints
    CONSTRAINT fk_applications_tenant FOREIGN KEY (tenant_id)
        REFERENCES tenants(tenant_id) ON DELETE RESTRICT,
    CONSTRAINT uq_applications_tenant_code UNIQUE (tenant_id, code)
);

-- Trigger for updated_at
CREATE TRIGGER trg_applications_updated_at
    BEFORE UPDATE ON applications
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Indexes
CREATE INDEX idx_applications_tenant_id ON applications(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_applications_tenant_code_active ON applications(tenant_id, code)
    WHERE deleted_at IS NULL AND is_active = TRUE;

-- Comments
COMMENT ON TABLE applications IS 'Applications that use IAM - each defines its own roles and permissions';
COMMENT ON COLUMN applications.code IS 'Unique identifier within tenant (e.g., pension-fund, backoffice)';
