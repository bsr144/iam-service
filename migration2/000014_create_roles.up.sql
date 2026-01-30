-- ============================================================================
-- Migration: 000014_create_roles
-- Description: Create roles table for application-specific roles
-- ============================================================================

CREATE TABLE roles (
    -- Primary Key
    role_id             UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    application_id      UUID NOT NULL,

    -- Business Identifiers
    code                VARCHAR(50) NOT NULL,
    name                VARCHAR(255) NOT NULL,
    description         TEXT,

    -- Role Type
    is_system           BOOLEAN NOT NULL DEFAULT FALSE,
    is_default          BOOLEAN NOT NULL DEFAULT FALSE,

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
    CONSTRAINT fk_roles_application FOREIGN KEY (application_id)
        REFERENCES applications(application_id) ON DELETE CASCADE,
    CONSTRAINT uq_roles_application_code UNIQUE (application_id, code)
);

-- Trigger for updated_at
CREATE TRIGGER trg_roles_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Indexes
CREATE INDEX idx_roles_application_id ON roles(application_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_roles_app_code_active ON roles(application_id, code)
    WHERE deleted_at IS NULL AND is_active = TRUE;
CREATE INDEX idx_roles_system ON roles(application_id) WHERE is_system = TRUE AND deleted_at IS NULL;

-- Comments
COMMENT ON TABLE roles IS 'Application-specific roles that bundle permissions together';
COMMENT ON COLUMN roles.code IS 'Unique within application (e.g., ADMIN, HR_STAFF, VIEWER)';
COMMENT ON COLUMN roles.is_system IS 'System roles (e.g., SUPER_ADMIN) cannot be modified or deleted';
COMMENT ON COLUMN roles.is_default IS 'Default role assigned to new users in application';
