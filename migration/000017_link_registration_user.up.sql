ALTER TABLE users ADD COLUMN registration_id UUID REFERENCES registrations(registration_id) ON DELETE SET NULL;
ALTER TABLE users ADD COLUMN registration_completed_at TIMESTAMPTZ;

CREATE INDEX idx_users_registration ON users(registration_id)
    WHERE registration_id IS NOT NULL;

COMMENT ON COLUMN users.registration_id IS 'Link to original registration process (audit trail)';
COMMENT ON COLUMN users.registration_completed_at IS 'When registration flow completed and user was created';
