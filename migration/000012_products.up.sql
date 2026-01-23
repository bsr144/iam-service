CREATE TABLE products (
    product_id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    product_type VARCHAR(50) DEFAULT 'application',
    is_active BOOLEAN DEFAULT FALSE,
    licensed_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT unique_product_code_per_tenant UNIQUE (tenant_id, code)
);

CREATE INDEX idx_products_tenant ON products(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_code ON products(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_active ON products(is_active) WHERE deleted_at IS NULL AND is_active = TRUE;

COMMENT ON TABLE products IS 'Products/applications within a tenant (e.g., crm, erp, hr-portal)';
COMMENT ON COLUMN products.code IS 'Unique product code within tenant (lowercase, hyphen-separated)';
COMMENT ON COLUMN products.product_type IS 'Type of product (defaults to application, can be customized if needed in future)';
COMMENT ON COLUMN products.is_active IS 'Product activation status (FALSE by default, set to TRUE when purchased/licensed)';
COMMENT ON COLUMN products.licensed_until IS 'License expiry date for trial/demo (NULL = perpetual license)';

INSERT INTO permissions (code, name, description, module, resource, action, scope_level, is_system) VALUES
('product:create', 'Create Product', 'Create new product/application within tenant', 'iam', 'product', 'create', 'tenant', TRUE),
('product:read', 'Read Product', 'View product details', 'iam', 'product', 'read', 'tenant', TRUE),
('product:update', 'Update Product', 'Modify product details', 'iam', 'product', 'update', 'tenant', TRUE),
('product:delete', 'Delete Product', 'Soft-delete product', 'iam', 'product', 'delete', 'tenant', TRUE),
('product:list', 'List Products', 'List all products in tenant', 'iam', 'product', 'list', 'tenant', TRUE)
ON CONFLICT (code) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r, permissions p
WHERE r.code = 'ADMIN' AND r.is_system = TRUE
  AND p.code IN ('product:create', 'product:read', 'product:update', 'product:delete', 'product:list')
ON CONFLICT (role_id, permission_id) DO NOTHING;
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r, permissions p
WHERE r.code = 'APPROVER' AND r.is_system = TRUE
  AND p.code IN ('product:read', 'product:list')
ON CONFLICT (role_id, permission_id) DO NOTHING;

ALTER TABLE roles ADD COLUMN product_id UUID REFERENCES products(product_id) ON DELETE CASCADE;
ALTER TABLE roles ADD CONSTRAINT check_product_requires_tenant
    CHECK ((product_id IS NULL) OR (tenant_id IS NOT NULL));
DROP INDEX IF EXISTS idx_roles_code;
CREATE UNIQUE INDEX idx_roles_code_system ON roles(code)
    WHERE tenant_id IS NULL AND product_id IS NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX idx_roles_code_tenant ON roles(tenant_id, code)
    WHERE tenant_id IS NOT NULL AND product_id IS NULL AND deleted_at IS NULL;
CREATE UNIQUE INDEX idx_roles_code_product ON roles(tenant_id, product_id, code)
    WHERE product_id IS NOT NULL AND deleted_at IS NULL;

CREATE INDEX idx_roles_product ON roles(product_id) WHERE deleted_at IS NULL;

COMMENT ON COLUMN roles.product_id IS 'Product scope for role (NULL = tenant-wide or system role, set = product-specific role)';

ALTER TABLE user_roles ADD COLUMN product_id UUID REFERENCES products(product_id) ON DELETE CASCADE;
ALTER TABLE user_roles DROP CONSTRAINT IF EXISTS unique_user_role_branch;
ALTER TABLE user_roles ADD CONSTRAINT unique_user_role_product_branch
    UNIQUE (user_id, role_id, product_id, branch_id);

CREATE INDEX idx_user_roles_product ON user_roles(product_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_roles_user_product ON user_roles(user_id, product_id) WHERE deleted_at IS NULL;

COMMENT ON COLUMN user_roles.product_id IS 'Product context for role assignment (NULL = all products, set = specific product)';

DROP VIEW IF EXISTS v_user_permissions;

CREATE VIEW v_user_permissions AS
SELECT DISTINCT
    ur.user_id,
    u.tenant_id,
    ur.product_id,
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

COMMENT ON VIEW v_user_permissions IS 'Resolved user permissions through role assignments (includes product context)';

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
    ur.product_id,
    prod.code as product_code,
    prod.name as product_name,
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
LEFT JOIN products prod ON prod.product_id = ur.product_id AND prod.deleted_at IS NULL
LEFT JOIN branches b ON b.branch_id = ur.branch_id AND b.deleted_at IS NULL
INNER JOIN roles r ON r.role_id = ur.role_id AND r.deleted_at IS NULL
WHERE u.deleted_at IS NULL
  AND (ur.effective_from IS NULL OR ur.effective_from <= CURRENT_TIMESTAMP)
  AND (ur.effective_to IS NULL OR ur.effective_to > CURRENT_TIMESTAMP);

COMMENT ON VIEW v_user_roles IS 'Summary of active user role assignments (includes product and branch context)';

DROP FUNCTION IF EXISTS user_has_permission(UUID, VARCHAR, UUID);

CREATE OR REPLACE FUNCTION user_has_permission(
    p_user_id UUID,
    p_permission_code VARCHAR,
    p_branch_id UUID DEFAULT NULL,
    p_product_id UUID DEFAULT NULL
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
    
          AND (p_product_id IS NULL OR vup.product_id = p_product_id OR vup.product_id IS NULL)
    
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

COMMENT ON FUNCTION user_has_permission IS 'Check if user has a specific permission (considering scope and product context)';

CREATE OR REPLACE FUNCTION get_user_products(p_user_id UUID)
RETURNS TABLE (product_id UUID, product_code VARCHAR, product_name VARCHAR) AS $$
BEGIN
    RETURN QUERY
    SELECT DISTINCT p.product_id, p.code, p.name
    FROM user_roles ur
    INNER JOIN products p ON p.product_id = ur.product_id
    WHERE ur.user_id = p_user_id
      AND ur.deleted_at IS NULL
      AND p.deleted_at IS NULL
      AND p.is_active = TRUE
      AND (ur.effective_from IS NULL OR ur.effective_from <= CURRENT_TIMESTAMP)
      AND (ur.effective_to IS NULL OR ur.effective_to > CURRENT_TIMESTAMP)
    ORDER BY p.code;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_user_products IS 'Get all active products accessible by a user (based on role assignments)';

CREATE OR REPLACE FUNCTION get_active_products(p_tenant_id UUID)
RETURNS TABLE (product_id UUID, product_code VARCHAR, product_name VARCHAR, licensed_until TIMESTAMPTZ) AS $$
BEGIN
    RETURN QUERY
    SELECT p.product_id, p.code, p.name, p.licensed_until
    FROM products p
    WHERE p.tenant_id = p_tenant_id
      AND p.deleted_at IS NULL
      AND p.is_active = TRUE
      AND (p.licensed_until IS NULL OR p.licensed_until > CURRENT_TIMESTAMP)
    ORDER BY p.code;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_active_products IS 'Get all active and licensed products for a tenant (checks license expiry)';

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

INSERT INTO schema_version (version, description) VALUES
('2.1', 'Product-based multi-application RBAC: products table, product-scoped roles and user_roles')
ON CONFLICT (version) DO NOTHING;
