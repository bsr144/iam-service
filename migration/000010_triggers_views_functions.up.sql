CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_tenants_updated_at BEFORE UPDATE ON tenants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tenant_settings_updated_at BEFORE UPDATE ON tenant_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_branches_updated_at BEFORE UPDATE ON branches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_branch_contacts_updated_at BEFORE UPDATE ON branch_contacts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_credentials_updated_at BEFORE UPDATE ON user_credentials
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_profiles_updated_at BEFORE UPDATE ON user_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_security_updated_at BEFORE UPDATE ON user_security
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_saml_configurations_updated_at BEFORE UPDATE ON saml_configurations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_permissions_updated_at BEFORE UPDATE ON permissions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_activation_tracking_updated_at BEFORE UPDATE ON user_activation_tracking
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


CREATE OR REPLACE VIEW v_users_complete AS
SELECT
    u.user_id,
    u.tenant_id,
    u.branch_id,
    u.email,
    u.email_verified,
    u.email_verified_at,
    u.is_service_account,
    u.is_active,
    u.created_at,
    u.updated_at,
    u.deleted_at,


    up.first_name,
    up.last_name,
    CONCAT(up.first_name, ' ', up.last_name) as full_name,
    up.address,
    up.phone,
    up.gender,
    up.marital_status,
    up.date_of_birth,
    up.place_of_birth,
    up.avatar_url,
    up.preferred_language,
    up.timezone,


    uc.password_hash IS NOT NULL as has_password,
    uc.pin_hash IS NOT NULL as has_pin,
    uc.password_expires_at,
    uc.pin_expires_at,
    uc.sso_provider,
    uc.mfa_enabled,


    us.last_login_at,
    us.last_login_ip,
    us.failed_login_attempts,
    us.locked_until,
    us.admin_registered_at,
    us.user_registered_at,
    us.admin_registered_by,
    us.invitation_expires_at
FROM users u
LEFT JOIN user_profiles up ON u.user_id = up.user_id
LEFT JOIN user_credentials uc ON u.user_id = uc.user_id
LEFT JOIN user_security us ON u.user_id = us.user_id;

COMMENT ON VIEW v_users_complete IS 'Complete user information from all normalized tables (excludes password/PIN hashes)';

CREATE OR REPLACE VIEW v_user_permissions AS
SELECT DISTINCT
    ur.user_id,
    u.tenant_id,
    ur.branch_id,
    p.code as permission_code,
    p.name as permission_name,
    p.module,
    p.resource,
    p.action,
    p.scope_level
FROM user_roles ur
INNER JOIN users u ON u.user_id = ur.user_id AND u.deleted_at IS NULL
INNER JOIN roles r ON r.role_id = ur.role_id AND r.deleted_at IS NULL
INNER JOIN role_permissions rp ON rp.role_id = r.role_id
INNER JOIN permissions p ON p.permission_id = rp.permission_id AND p.deleted_at IS NULL
WHERE ur.deleted_at IS NULL
  AND u.is_active = TRUE
  AND r.is_active = TRUE
  AND (ur.effective_from IS NULL OR ur.effective_from <= CURRENT_TIMESTAMP)
  AND (ur.effective_to IS NULL OR ur.effective_to > CURRENT_TIMESTAMP);

COMMENT ON VIEW v_user_permissions IS 'Resolved user permissions through role assignments';

CREATE OR REPLACE VIEW v_user_roles AS
SELECT
    u.user_id,
    u.email,
    up.first_name,
    up.last_name,
    CONCAT(up.first_name, ' ', up.last_name) as full_name,
    u.tenant_id,
    t.name as tenant_name,
    ur.branch_id,
    b.name as branch_name,
    r.role_id,
    r.code as role_code,
    r.name as role_name,
    ur.effective_from,
    ur.effective_to
FROM users u
INNER JOIN tenants t ON t.tenant_id = u.tenant_id AND t.deleted_at IS NULL
LEFT JOIN user_profiles up ON u.user_id = up.user_id
INNER JOIN user_roles ur ON ur.user_id = u.user_id AND ur.deleted_at IS NULL
LEFT JOIN branches b ON b.branch_id = ur.branch_id AND b.deleted_at IS NULL
INNER JOIN roles r ON r.role_id = ur.role_id AND r.deleted_at IS NULL
WHERE u.deleted_at IS NULL
  AND (ur.effective_from IS NULL OR ur.effective_from <= CURRENT_TIMESTAMP)
  AND (ur.effective_to IS NULL OR ur.effective_to > CURRENT_TIMESTAMP);

COMMENT ON VIEW v_user_roles IS 'Summary of active user role assignments';


CREATE OR REPLACE FUNCTION user_has_permission(
    p_user_id UUID,
    p_permission_code VARCHAR,
    p_branch_id UUID DEFAULT NULL
)
RETURNS BOOLEAN AS $$
DECLARE
    v_has_permission BOOLEAN;
BEGIN
    SELECT EXISTS (
        SELECT 1
        FROM v_user_permissions vup
        WHERE vup.user_id = p_user_id
          AND vup.permission_code = p_permission_code
          AND (
              vup.scope_level = 'system'
              OR vup.scope_level = 'tenant'
              OR (vup.scope_level = 'branch' AND (vup.branch_id = p_branch_id OR vup.branch_id IS NULL))
              OR vup.scope_level = 'self'
          )
    ) INTO v_has_permission;

    RETURN v_has_permission;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION user_has_permission IS 'Check if user has a specific permission (considering scope)';

CREATE OR REPLACE FUNCTION get_user_branches(p_user_id UUID)
RETURNS TABLE (branch_id UUID) AS $$
BEGIN
    RETURN QUERY
    SELECT DISTINCT ur.branch_id
    FROM user_roles ur
    WHERE ur.user_id = p_user_id
      AND ur.deleted_at IS NULL
      AND (ur.effective_from IS NULL OR ur.effective_from <= CURRENT_TIMESTAMP)
      AND (ur.effective_to IS NULL OR ur.effective_to > CURRENT_TIMESTAMP)
      AND ur.branch_id IS NOT NULL;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_user_branches IS 'Get all branches where user has role assignments';

CREATE TABLE IF NOT EXISTS schema_version (
    version VARCHAR(20) PRIMARY KEY,
    applied_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);

INSERT INTO schema_version (version, description) VALUES
('2.0', 'Normalized IAM schema: users split into 4 tables, tenants/branches normalized, improved performance and security')
ON CONFLICT (version) DO NOTHING;
