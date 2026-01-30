-- ============================================================================
-- Migration: 000015_create_permissions
-- Description: Create permissions table for fine-grained permissions
-- ============================================================================

CREATE TABLE permissions (
    -- Primary Key
    permission_id       UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    application_id      UUID NOT NULL,

    -- Business Identifiers
    code                VARCHAR(100) NOT NULL,
    name                VARCHAR(255) NOT NULL,
    description         TEXT,

    -- Permission Classification (for UI grouping)
    resource_type       VARCHAR(50),
    action              VARCHAR(50),

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
    CONSTRAINT fk_permissions_application FOREIGN KEY (application_id)
        REFERENCES applications(application_id) ON DELETE CASCADE,
    CONSTRAINT uq_permissions_application_code UNIQUE (application_id, code),
    CONSTRAINT chk_permissions_code_format CHECK (code ~ '^[a-z_]+:[a-z_]+$')
);

-- Trigger for updated_at
CREATE TRIGGER trg_permissions_updated_at
    BEFORE UPDATE ON permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Indexes
CREATE INDEX idx_permissions_application_id ON permissions(application_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_permissions_app_code_active ON permissions(application_id, code)
    WHERE deleted_at IS NULL AND is_active = TRUE;
CREATE INDEX idx_permissions_resource_type ON permissions(application_id, resource_type)
    WHERE deleted_at IS NULL;

-- Comments
COMMENT ON TABLE permissions IS 'Fine-grained permissions in format resource:action';
COMMENT ON COLUMN permissions.code IS 'Permission code in format resource:action (e.g., user:create, loan:approve)';
COMMENT ON COLUMN permissions.resource_type IS 'Resource type extracted from code for grouping';
COMMENT ON COLUMN permissions.action IS 'Action extracted from code for grouping';
