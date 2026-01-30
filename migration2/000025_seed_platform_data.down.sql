-- ============================================================================
-- Migration: 000025_seed_platform_data (DOWN)
-- Description: Remove seed data (WARNING: This will delete platform admin!)
-- ============================================================================

-- Delete in reverse order of creation due to foreign key constraints

-- Delete user role assignment
DELETE FROM user_roles WHERE user_role_id = '55555555-5555-5555-5555-555555555555';

-- Delete user activation tracking
DELETE FROM user_activation_tracking WHERE user_activation_tracking_id = '44444444-4444-4444-4444-444444444444';

-- Delete user security state
DELETE FROM user_security_states WHERE user_id = '33333333-3333-3333-3333-333333333333';

-- Delete user credentials
DELETE FROM user_credentials WHERE user_id = '33333333-3333-3333-3333-333333333333';

-- Delete user profile
DELETE FROM user_profiles WHERE user_id = '33333333-3333-3333-3333-333333333333';

-- Delete admin user
DELETE FROM users WHERE user_id = '33333333-3333-3333-3333-333333333333';

-- Delete role permissions
DELETE FROM role_permissions WHERE role_id = '22222222-2222-2222-2222-222222222222';

-- Delete permissions
DELETE FROM permissions WHERE application_id = '11111111-1111-1111-1111-111111111111';

-- Delete Platform Admin role
DELETE FROM roles WHERE role_id = '22222222-2222-2222-2222-222222222222';

-- Delete IAM Admin application
DELETE FROM applications WHERE application_id = '11111111-1111-1111-1111-111111111111';

-- Delete tenant settings
DELETE FROM tenant_settings WHERE tenant_id = '01234567-89ab-cdef-0123-456789abcdef';

-- Delete Platform tenant
DELETE FROM tenants WHERE tenant_id = '01234567-89ab-cdef-0123-456789abcdef';
