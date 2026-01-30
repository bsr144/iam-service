-- ============================================================================
-- Migration: 000012_create_user_branches
-- Description: Create user_branches junction table (M:N between users and branches)
-- ============================================================================

CREATE TABLE user_branches (
    -- Primary Key
    user_branch_id      UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys (Composite relationship)
    user_id             UUID NOT NULL,
    branch_id           UUID NOT NULL,

    -- Assignment Details
    is_primary          BOOLEAN NOT NULL DEFAULT FALSE,
    assigned_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assigned_by         UUID,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ,

    -- Constraints
    CONSTRAINT fk_user_branches_user FOREIGN KEY (user_id)
        REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_user_branches_branch FOREIGN KEY (branch_id)
        REFERENCES branches(branch_id) ON DELETE CASCADE,
    CONSTRAINT uq_user_branches_user_branch UNIQUE (user_id, branch_id)
);

-- Trigger for updated_at
CREATE TRIGGER trg_user_branches_updated_at
    BEFORE UPDATE ON user_branches
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Ensure only one primary branch per user
CREATE UNIQUE INDEX idx_user_branches_primary
    ON user_branches(user_id)
    WHERE is_primary = TRUE AND deleted_at IS NULL;

-- Indexes
CREATE INDEX idx_user_branches_user_id ON user_branches(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_branches_branch_id ON user_branches(branch_id) WHERE deleted_at IS NULL;

-- Comments
COMMENT ON TABLE user_branches IS 'Junction table linking users to their assigned branches';
COMMENT ON COLUMN user_branches.is_primary IS 'Only one branch can be primary per user (enforced by partial unique index)';
