CREATE TABLE permissions (
    permission_id UUID PRIMARY KEY DEFAULT uuidv7(),
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    module VARCHAR(50) NOT NULL,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    scope_level VARCHAR(20) DEFAULT 'branch' CHECK (scope_level IN ('system', 'tenant', 'branch', 'self')),
    is_system BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_permissions_code ON permissions(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_permissions_module ON permissions(module) WHERE deleted_at IS NULL;

COMMENT ON TABLE permissions IS 'Granular permissions for all system operations';
COMMENT ON COLUMN permissions.code IS 'Permission code (format: {resource}:{action}, e.g., employee:create)';
COMMENT ON COLUMN permissions.scope_level IS 'Defines the scope at which this permission operates';
CREATE TABLE roles (
    role_id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    parent_role_id UUID REFERENCES roles(role_id) ON DELETE SET NULL,
    scope_level VARCHAR(20) DEFAULT 'branch' CHECK (scope_level IN ('system', 'tenant', 'branch')),
    is_system BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_roles_tenant ON roles(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_roles_code ON roles(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_roles_parent ON roles(parent_role_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE roles IS 'Roles with hierarchical support (child roles inherit parent permissions)';
COMMENT ON COLUMN roles.parent_role_id IS 'Parent role for inheritance (NULL = root role)';
CREATE TABLE role_permissions (
    role_permission_id UUID PRIMARY KEY DEFAULT uuidv7(),
    role_id UUID NOT NULL REFERENCES roles(role_id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(permission_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_role_permission UNIQUE (role_id, permission_id)
);

CREATE INDEX idx_role_permissions_role ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission ON role_permissions(permission_id);

COMMENT ON TABLE role_permissions IS 'Many-to-many mapping of roles to permissions';
CREATE TABLE user_roles (
    user_role_id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(role_id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(branch_id) ON DELETE CASCADE,
    effective_from TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    effective_to TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT unique_user_role_branch UNIQUE (user_id, role_id, branch_id)
);

CREATE INDEX idx_user_roles_user ON user_roles(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_roles_role ON user_roles(role_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_roles_branch ON user_roles(branch_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_roles_effective ON user_roles(effective_from, effective_to) WHERE deleted_at IS NULL;

COMMENT ON TABLE user_roles IS 'Assigns roles to users with optional branch scoping and time-bound validity';
COMMENT ON COLUMN user_roles.branch_id IS 'If NULL, role applies tenant-wide; otherwise scoped to specific branch';
COMMENT ON COLUMN user_roles.effective_from IS 'Role becomes active from this timestamp';
COMMENT ON COLUMN user_roles.effective_to IS 'Role expires after this timestamp (NULL = permanent)';
