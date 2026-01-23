CREATE TABLE tenants (
    tenant_id UUID PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    database_name VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT valid_slug CHECK (slug ~ '^[a-z0-9-]+$'),
    CONSTRAINT valid_database_name CHECK (database_name ~ '^tenant_[a-z0-9_]+$')
);

CREATE INDEX idx_tenants_slug ON tenants(slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_tenants_status ON tenants(status) WHERE deleted_at IS NULL;

COMMENT ON TABLE tenants IS 'Core tenant identity (organizations using the IAM system)';
COMMENT ON COLUMN tenants.slug IS 'URL-friendly unique identifier (lowercase, hyphens only)';
COMMENT ON COLUMN tenants.database_name IS 'PostgreSQL database name for tenant data (format: tenant_{slug})';

CREATE TABLE tenant_settings (
    tenant_setting_id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL UNIQUE REFERENCES tenants(tenant_id) ON DELETE CASCADE,

    subscription_tier VARCHAR(50) DEFAULT 'standard',
    max_branches INTEGER DEFAULT 10,
    max_employees INTEGER DEFAULT 10000,

    contact_email VARCHAR(255),
    contact_phone VARCHAR(50),
    contact_address TEXT,

    default_language VARCHAR(10) DEFAULT 'en',
    timezone VARCHAR(50) DEFAULT 'Asia/Jakarta',

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tenant_settings_tenant ON tenant_settings(tenant_id);

COMMENT ON TABLE tenant_settings IS 'Tenant configuration, subscription limits, and contact information';
COMMENT ON COLUMN tenant_settings.subscription_tier IS 'Subscription plan (standard, premium, enterprise)';
COMMENT ON COLUMN tenant_settings.max_branches IS 'Maximum number of branches allowed';
COMMENT ON COLUMN tenant_settings.max_employees IS 'Maximum number of employees allowed';

CREATE TABLE branches (
    branch_id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    is_headquarters BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT unique_branch_code_per_tenant UNIQUE (tenant_id, code)
);

CREATE INDEX idx_branches_tenant ON branches(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_branches_code ON branches(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_branches_active ON branches(is_active) WHERE deleted_at IS NULL;

COMMENT ON TABLE branches IS 'Core branch identity within a tenant';
COMMENT ON COLUMN branches.is_headquarters IS 'Indicates if this is the main/headquarters branch';

CREATE TABLE branch_contacts (
    branch_contact_id UUID PRIMARY KEY DEFAULT uuidv7(),
    branch_id UUID NOT NULL UNIQUE REFERENCES branches(branch_id) ON DELETE CASCADE,

    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(20),

    phone VARCHAR(50),
    email VARCHAR(255),

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_branch_contacts_branch ON branch_contacts(branch_id);

COMMENT ON TABLE branch_contacts IS 'Branch contact information and address details';
