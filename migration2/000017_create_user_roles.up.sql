-- ============================================================================
-- Migration: 000017_create_user_roles
-- Description: Create user_roles table for user-role assignments
-- ============================================================================

CREATE TABLE user_roles (
    -- Primary Key
    user_role_id        UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Core Assignment
    user_id             UUID NOT NULL,
    role_id             UUID NOT NULL,

    -- Optional Scoping
    branch_id           UUID,

    -- Assignment Metadata
    assigned_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assigned_by         UUID,
    effective_from      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    effective_to        TIMESTAMPTZ,

    -- Status Management
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    -- Constraints
    CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id)
        REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id)
        REFERENCES roles(role_id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_branch FOREIGN KEY (branch_id)
        REFERENCES branches(branch_id) ON DELETE CASCADE,
    CONSTRAINT uq_user_roles_user_role_branch UNIQUE (user_id, role_id, COALESCE(branch_id, '00000000-0000-0000-0000-000000000000'::uuid))
);

-- Trigger for updated_at
CREATE TRIGGER trg_user_roles_updated_at
    BEFORE UPDATE ON user_roles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Indexes
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id)
    WHERE deleted_at IS NULL AND is_active = TRUE;
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id)
    WHERE deleted_at IS NULL AND is_active = TRUE;
CREATE INDEX idx_user_roles_branch_id ON user_roles(branch_id)
    WHERE deleted_at IS NULL AND branch_id IS NOT NULL;
CREATE INDEX idx_user_roles_user_active ON user_roles(user_id, is_active)
    WHERE deleted_at IS NULL AND is_active = TRUE;
CREATE INDEX idx_user_roles_effective ON user_roles(user_id)
    WHERE deleted_at IS NULL AND is_active = TRUE
    AND (effective_to IS NULL OR effective_to > NOW());

-- Comments
COMMENT ON TABLE user_roles IS 'Assigns roles to users, with optional branch-level scoping';
COMMENT ON COLUMN user_roles.branch_id IS 'NULL means tenant-wide; specific ID means branch-scoped';
COMMENT ON COLUMN user_roles.effective_to IS 'Optional expiration for time-bound roles';
