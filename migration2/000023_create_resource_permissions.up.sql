-- ============================================================================
-- Migration: 000023_create_resource_permissions
-- Description: Create resource_permissions table for fine-grained resource-level ACL
-- ============================================================================

CREATE TABLE resource_permissions (
    -- Primary Key
    resource_permission_id UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Who has the permission
    user_id             UUID NOT NULL,
    application_id      UUID NOT NULL,

    -- What resource
    resource_type       VARCHAR(50) NOT NULL,
    resource_id         VARCHAR(255) NOT NULL,

    -- What permission
    permission_code     VARCHAR(100) NOT NULL,

    -- Grant Metadata
    granted_by          UUID,
    granted_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at          TIMESTAMPTZ,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    -- Constraints
    CONSTRAINT fk_resource_perms_user FOREIGN KEY (user_id)
        REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_resource_perms_application FOREIGN KEY (application_id)
        REFERENCES applications(application_id) ON DELETE CASCADE,
    CONSTRAINT uq_resource_perms UNIQUE (user_id, application_id, resource_type, resource_id, permission_code)
);

-- Indexes
CREATE INDEX idx_resource_perms_user_id ON resource_permissions(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_resource_perms_application_id ON resource_permissions(application_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_resource_perms_lookup ON resource_permissions(user_id, application_id, resource_type, resource_id)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_resource_perms_expires ON resource_permissions(expires_at)
    WHERE deleted_at IS NULL AND expires_at IS NOT NULL;

-- Comments
COMMENT ON TABLE resource_permissions IS 'Grants permissions on specific resource instances (object-level ACL)';
COMMENT ON COLUMN resource_permissions.resource_type IS 'Type of resource (e.g., document, report, loan)';
COMMENT ON COLUMN resource_permissions.resource_id IS 'Specific resource identifier';
COMMENT ON COLUMN resource_permissions.permission_code IS 'Permission granted (e.g., document:read)';
