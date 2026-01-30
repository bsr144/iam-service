-- ============================================================================
-- Migration: 000025_seed_platform_data
-- Description: Seed platform tenant, IAM application, and initial admin user
-- ============================================================================

-- 1. Create Platform Tenant
INSERT INTO tenants (tenant_id, name, slug, database_name, status) VALUES (
    '01234567-89ab-cdef-0123-456789abcdef',
    'Platform Administration',
    'platform',
    'iam_platform',
    'active'
) ON CONFLICT (slug) DO NOTHING;

-- 2. Create Platform Tenant Settings
INSERT INTO tenant_settings (tenant_setting_id, tenant_id, subscription_tier, max_branches, max_employees, approval_required, mfa_required) VALUES (
    '01234567-89ab-cdef-0123-456789abcde0',
    '01234567-89ab-cdef-0123-456789abcdef',
    'enterprise',
    1000,
    100000,
    false,
    false
) ON CONFLICT (tenant_id) DO NOTHING;

-- 3. Create IAM Admin Application
INSERT INTO applications (application_id, tenant_id, code, name, description, is_active) VALUES (
    '11111111-1111-1111-1111-111111111111',
    '01234567-89ab-cdef-0123-456789abcdef',
    'iam-admin',
    'IAM Administration',
    'Internal application for managing the IAM system',
    true
) ON CONFLICT (tenant_id, code) DO NOTHING;

-- 4. Create Platform Admin Role
INSERT INTO roles (role_id, application_id, code, name, description, is_system, is_default, is_active) VALUES (
    '22222222-2222-2222-2222-222222222222',
    '11111111-1111-1111-1111-111111111111',
    'PLATFORM_ADMIN',
    'Platform Administrator',
    'Full access to all IAM operations across all tenants',
    true,
    false,
    true
) ON CONFLICT (application_id, code) DO NOTHING;

-- 5. Create Core IAM Permissions
INSERT INTO permissions (application_id, code, name, resource_type, action, is_active) VALUES
    ('11111111-1111-1111-1111-111111111111', 'tenant:create', 'Create Tenant', 'tenant', 'create', true),
    ('11111111-1111-1111-1111-111111111111', 'tenant:read', 'View Tenant', 'tenant', 'read', true),
    ('11111111-1111-1111-1111-111111111111', 'tenant:update', 'Update Tenant', 'tenant', 'update', true),
    ('11111111-1111-1111-1111-111111111111', 'tenant:delete', 'Delete Tenant', 'tenant', 'delete', true),
    ('11111111-1111-1111-1111-111111111111', 'user:create', 'Create User', 'user', 'create', true),
    ('11111111-1111-1111-1111-111111111111', 'user:read', 'View User', 'user', 'read', true),
    ('11111111-1111-1111-1111-111111111111', 'user:update', 'Update User', 'user', 'update', true),
    ('11111111-1111-1111-1111-111111111111', 'user:delete', 'Delete User', 'user', 'delete', true),
    ('11111111-1111-1111-1111-111111111111', 'user:approve', 'Approve User', 'user', 'approve', true),
    ('11111111-1111-1111-1111-111111111111', 'application:create', 'Create Application', 'application', 'create', true),
    ('11111111-1111-1111-1111-111111111111', 'application:read', 'View Application', 'application', 'read', true),
    ('11111111-1111-1111-1111-111111111111', 'application:update', 'Update Application', 'application', 'update', true),
    ('11111111-1111-1111-1111-111111111111', 'application:delete', 'Delete Application', 'application', 'delete', true),
    ('11111111-1111-1111-1111-111111111111', 'role:create', 'Create Role', 'role', 'create', true),
    ('11111111-1111-1111-1111-111111111111', 'role:read', 'View Role', 'role', 'read', true),
    ('11111111-1111-1111-1111-111111111111', 'role:update', 'Update Role', 'role', 'update', true),
    ('11111111-1111-1111-1111-111111111111', 'role:delete', 'Delete Role', 'role', 'delete', true),
    ('11111111-1111-1111-1111-111111111111', 'permission:create', 'Create Permission', 'permission', 'create', true),
    ('11111111-1111-1111-1111-111111111111', 'permission:read', 'View Permission', 'permission', 'read', true),
    ('11111111-1111-1111-1111-111111111111', 'permission:update', 'Update Permission', 'permission', 'update', true),
    ('11111111-1111-1111-1111-111111111111', 'permission:delete', 'Delete Permission', 'permission', 'delete', true),
    ('11111111-1111-1111-1111-111111111111', 'branch:create', 'Create Branch', 'branch', 'create', true),
    ('11111111-1111-1111-1111-111111111111', 'branch:read', 'View Branch', 'branch', 'read', true),
    ('11111111-1111-1111-1111-111111111111', 'branch:update', 'Update Branch', 'branch', 'update', true),
    ('11111111-1111-1111-1111-111111111111', 'branch:delete', 'Delete Branch', 'branch', 'delete', true),
    ('11111111-1111-1111-1111-111111111111', 'audit:read', 'View Audit Logs', 'audit', 'read', true),
    ('11111111-1111-1111-1111-111111111111', 'audit:export', 'Export Audit Logs', 'audit', 'export', true)
ON CONFLICT (application_id, code) DO NOTHING;

-- 6. Assign all permissions to Platform Admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    '22222222-2222-2222-2222-222222222222',
    permission_id
FROM permissions
WHERE application_id = '11111111-1111-1111-1111-111111111111'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 7. Create initial platform admin user (core identity)
INSERT INTO users (user_id, tenant_id, email, email_verified, is_active, is_service_account) VALUES (
    '33333333-3333-3333-3333-333333333333',
    '01234567-89ab-cdef-0123-456789abcdef',
    'admin@platform.local',
    true,
    true,
    false
) ON CONFLICT (tenant_id, email) DO NOTHING;

-- 8. Create admin user profile
INSERT INTO user_profiles (user_id, first_name, last_name, preferred_language, timezone) VALUES (
    '33333333-3333-3333-3333-333333333333',
    'Platform',
    'Administrator',
    'en',
    'UTC'
) ON CONFLICT (user_id) DO NOTHING;

-- 9. Create admin user credentials
-- Password: ChangeMe123! (bcrypt hash)
-- PIN: 123456 (bcrypt hash)
INSERT INTO user_credentials (user_id, password_hash, password_set_at, pin_hash, pin_set_at) VALUES (
    '33333333-3333-3333-3333-333333333333',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/X4/tPBxqK1/XW.zSm',
    NOW(),
    '$2a$12$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    NOW()
) ON CONFLICT (user_id) DO NOTHING;

-- 10. Create admin user security state
INSERT INTO user_security_states (user_id, email_verified_at) VALUES (
    '33333333-3333-3333-3333-333333333333',
    NOW()
) ON CONFLICT (user_id) DO NOTHING;

-- 11. Create admin user activation tracking
INSERT INTO user_activation_tracking (user_activation_tracking_id, user_id, tenant_id, admin_created, admin_created_at, user_completed, user_completed_at, activated_at, activation_method, status_history) VALUES (
    '44444444-4444-4444-4444-444444444444',
    '33333333-3333-3333-3333-333333333333',
    '01234567-89ab-cdef-0123-456789abcdef',
    true,
    NOW(),
    true,
    NOW(),
    NOW(),
    'admin_first',
    '[{"status": "activated", "timestamp": "'|| NOW() ||'", "triggered_by": "system"}]'::jsonb
) ON CONFLICT (user_id) DO NOTHING;

-- 12. Assign Platform Admin role to admin user
INSERT INTO user_roles (user_role_id, user_id, role_id, assigned_by, is_active) VALUES (
    '55555555-5555-5555-5555-555555555555',
    '33333333-3333-3333-3333-333333333333',
    '22222222-2222-2222-2222-222222222222',
    '33333333-3333-3333-3333-333333333333',
    true
) ON CONFLICT (user_id, role_id, COALESCE(branch_id, '00000000-0000-0000-0000-000000000000'::uuid)) DO NOTHING;
