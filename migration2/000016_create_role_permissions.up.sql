-- ============================================================================
-- Migration: 000016_create_role_permissions
-- Description: Create role_permissions junction table (M:N between roles and permissions)
-- ============================================================================

CREATE TABLE role_permissions (
    -- Primary Key
    role_permission_id  UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    role_id             UUID NOT NULL,
    permission_id       UUID NOT NULL,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID,
    deleted_at          TIMESTAMPTZ,

    -- Constraints
    CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id)
        REFERENCES roles(role_id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id)
        REFERENCES permissions(permission_id) ON DELETE CASCADE,
    CONSTRAINT uq_role_permissions UNIQUE (role_id, permission_id)
);

-- Indexes
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id) WHERE deleted_at IS NULL;

-- Comments
COMMENT ON TABLE role_permissions IS 'Junction table: which permissions belong to which roles';
