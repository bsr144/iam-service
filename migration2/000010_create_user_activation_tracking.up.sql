-- ============================================================================
-- Migration: 000010_create_user_activation_tracking
-- Description: Create user_activation_tracking table for registration workflow
-- ============================================================================

CREATE TABLE user_activation_tracking (
    -- Primary Key
    user_activation_tracking_id UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    user_id             UUID NOT NULL,
    tenant_id           UUID,

    -- Admin Registration (Step 1 - Admin creates user)
    admin_created       BOOLEAN NOT NULL DEFAULT FALSE,
    admin_created_at    TIMESTAMPTZ,
    admin_created_by    UUID,

    -- User Completion (Step 2 - User completes registration)
    user_completed      BOOLEAN NOT NULL DEFAULT FALSE,
    user_completed_at   TIMESTAMPTZ,
    otp_verified_at     TIMESTAMPTZ,
    profile_completed_at TIMESTAMPTZ,
    pin_set_at          TIMESTAMPTZ,

    -- Final Activation
    activated_at        TIMESTAMPTZ,
    activation_method   VARCHAR(50),

    -- Status History (JSONB array of status transitions)
    status_history      JSONB NOT NULL DEFAULT '[]'::jsonb,

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_user_activation_tracking_user FOREIGN KEY (user_id)
        REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_user_activation_tracking_tenant FOREIGN KEY (tenant_id)
        REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    CONSTRAINT uq_user_activation_tracking_user_id UNIQUE (user_id),
    CONSTRAINT chk_activation_method CHECK (activation_method IN ('admin_first', 'user_first', 'simultaneous'))
);

-- Trigger for updated_at
CREATE TRIGGER trg_user_activation_tracking_updated_at
    BEFORE UPDATE ON user_activation_tracking
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Indexes
CREATE INDEX idx_user_activation_pending_user ON user_activation_tracking(tenant_id)
    WHERE admin_created = TRUE AND user_completed = FALSE AND activated_at IS NULL;
CREATE INDEX idx_user_activation_pending_admin ON user_activation_tracking(tenant_id)
    WHERE user_completed = TRUE AND admin_created = FALSE AND activated_at IS NULL;

-- Comments
COMMENT ON TABLE user_activation_tracking IS 'Tracks user registration/activation workflow progress';
COMMENT ON COLUMN user_activation_tracking.activation_method IS 'How the user was activated: admin_first, user_first, or simultaneous';
COMMENT ON COLUMN user_activation_tracking.status_history IS 'JSONB array of {status, timestamp, triggered_by} transitions';
