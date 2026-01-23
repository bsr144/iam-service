DROP INDEX IF EXISTS idx_user_roles_user_product;
DROP INDEX IF EXISTS idx_user_roles_product;
DROP INDEX IF EXISTS idx_roles_product;
DROP INDEX IF EXISTS idx_roles_code_product;
DROP INDEX IF EXISTS idx_roles_code_tenant;
DROP INDEX IF EXISTS idx_roles_code_system;
DROP INDEX IF EXISTS idx_products_active;
DROP INDEX IF EXISTS idx_products_code;
DROP INDEX IF EXISTS idx_products_tenant;

DROP FUNCTION IF EXISTS get_active_products(UUID);
DROP FUNCTION IF EXISTS get_user_products(UUID);

DROP FUNCTION IF EXISTS user_has_permission(UUID, VARCHAR, UUID, UUID);

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

DROP VIEW IF EXISTS v_user_roles;

CREATE VIEW v_user_roles AS
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

DROP VIEW IF EXISTS v_user_permissions;

CREATE VIEW v_user_permissions AS
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
DROP INDEX IF EXISTS idx_user_roles_user_product;
DROP INDEX IF EXISTS idx_user_roles_product;
ALTER TABLE user_roles DROP CONSTRAINT IF EXISTS unique_user_role_product_branch;
ALTER TABLE user_roles DROP COLUMN IF EXISTS product_id;
ALTER TABLE user_roles ADD CONSTRAINT unique_user_role_branch
    UNIQUE (user_id, role_id, branch_id);
ALTER TABLE roles DROP CONSTRAINT IF EXISTS check_product_requires_tenant;
ALTER TABLE roles DROP COLUMN IF EXISTS product_id;
CREATE INDEX idx_roles_code ON roles(code) WHERE deleted_at IS NULL;
DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT permission_id FROM permissions WHERE code LIKE 'product:%'
);
DELETE FROM permissions WHERE code LIKE 'product:%';

DROP TRIGGER IF EXISTS update_products_updated_at ON products;
DROP TABLE IF EXISTS products;

DELETE FROM schema_version WHERE version = '2.1';
