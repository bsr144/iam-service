DELETE FROM user_roles
WHERE role_id IN (
    SELECT role_id FROM roles WHERE code = 'PLATFORM_ADMIN' AND tenant_id IS NULL
);

DELETE FROM user_security
WHERE user_id IN (SELECT user_id FROM users WHERE tenant_id IS NULL);

DELETE FROM user_profiles
WHERE user_id IN (SELECT user_id FROM users WHERE tenant_id IS NULL);

DELETE FROM user_credentials
WHERE user_id IN (SELECT user_id FROM users WHERE tenant_id IS NULL);

DELETE FROM user_activation_tracking
WHERE user_id IN (SELECT user_id FROM users WHERE tenant_id IS NULL);

DELETE FROM users WHERE tenant_id IS NULL;

DROP INDEX IF EXISTS idx_users_platform_admin;

DELETE FROM role_permissions
WHERE role_id IN (
    SELECT role_id FROM roles WHERE code = 'PLATFORM_ADMIN' AND tenant_id IS NULL
);

DELETE FROM roles WHERE code = 'PLATFORM_ADMIN' AND tenant_id IS NULL;

DELETE FROM permissions WHERE code LIKE 'platform:%';


COMMENT ON COLUMN users.tenant_id IS NULL;
