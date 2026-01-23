ALTER TABLE users ALTER COLUMN tenant_id DROP NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS unique_platform_role_code 
ON roles (code) WHERE tenant_id IS NULL;

INSERT INTO roles (role_id, tenant_id, code, name, description, scope_level, is_system, is_active, created_at, updated_at)
VALUES (
    uuidv7(),
    NULL, 
    'PLATFORM_ADMIN',
    'Platform Administrator',
    'Super administrator with cross-tenant access and platform-level permissions. Can create roles and special accounts for any tenant.',
    'system',
    TRUE,
    TRUE,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT (code) WHERE tenant_id IS NULL DO NOTHING;

INSERT INTO permissions (permission_id, code, name, description, module, resource, action, scope_level, is_system, created_at, updated_at) VALUES
(gen_random_uuid(), 'platform:role:create', 'Create Tenant Role', 'Create custom roles for any tenant', 'platform', 'role', 'create', 'system', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(gen_random_uuid(), 'platform:role:read', 'Read Tenant Roles', 'View roles across all tenants', 'platform', 'role', 'read', 'system', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(gen_random_uuid(), 'platform:role:update', 'Update Tenant Role', 'Modify roles for any tenant', 'platform', 'role', 'update', 'system', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(gen_random_uuid(), 'platform:role:delete', 'Delete Tenant Role', 'Delete roles for any tenant', 'platform', 'role', 'delete', 'system', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'platform:account:create', 'Create Platform Account', 'Create special accounts with system roles for any tenant', 'platform', 'account', 'create', 'system', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(gen_random_uuid(), 'platform:account:read', 'Read Platform Accounts', 'View special accounts across all tenants', 'platform', 'account', 'read', 'system', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'platform:tenant:read', 'Read All Tenants', 'View all tenant information', 'platform', 'tenant', 'read', 'system', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(gen_random_uuid(), 'platform:tenant:create', 'Create Tenant', 'Create new tenant', 'platform', 'tenant', 'create', 'system', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'platform:audit:read', 'Read Platform Audit', 'View audit logs across all tenants', 'platform', 'audit', 'read', 'system', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (code) DO NOTHING;

INSERT INTO role_permissions (role_permission_id, role_id, permission_id, created_at)
SELECT
    uuidv7(),
    r.role_id,
    p.permission_id,
    CURRENT_TIMESTAMP
FROM roles r
CROSS JOIN permissions p
WHERE r.code = 'PLATFORM_ADMIN'
    AND r.is_system = TRUE
    AND r.tenant_id IS NULL
    AND p.code LIKE 'platform:%'
ON CONFLICT ON CONSTRAINT unique_role_permission DO NOTHING;

CREATE INDEX IF NOT EXISTS idx_users_platform_admin ON users(tenant_id) WHERE tenant_id IS NULL;

COMMENT ON INDEX idx_users_platform_admin IS 'Index for identifying platform administrators (users with NULL tenant_id)';

COMMENT ON COLUMN users.tenant_id IS 'Tenant ID for tenant-scoped users. NULL indicates platform administrator.';
