CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT uuidv7(),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(branch_id) ON DELETE SET NULL,

    
    email VARCHAR(255) NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    email_verified_at TIMESTAMPTZ,

    
    is_service_account BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,

    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT unique_email_per_tenant UNIQUE (tenant_id, email)
);

CREATE INDEX idx_users_tenant ON users(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_branch ON users(branch_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_active ON users(is_active) WHERE deleted_at IS NULL;

COMMENT ON TABLE users IS 'Core user identity - frequently accessed, kept minimal for performance';
COMMENT ON COLUMN users.is_service_account IS 'Service accounts for API-to-API communication';


CREATE TABLE user_credentials (
    user_credential_id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
    
    password_hash VARCHAR(255),
    password_changed_at TIMESTAMPTZ,
    password_expires_at TIMESTAMPTZ,
    password_history JSONB DEFAULT '[]'::JSONB,

    pin_hash VARCHAR(255),
    pin_set_at TIMESTAMPTZ,
    pin_changed_at TIMESTAMPTZ,
    pin_expires_at TIMESTAMPTZ,
    pin_history JSONB DEFAULT '[]'::JSONB,

    sso_provider VARCHAR(50),
    sso_provider_id VARCHAR(255),

    mfa_enabled BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_credentials_user ON user_credentials(user_id);
CREATE INDEX idx_user_credentials_sso ON user_credentials(sso_provider, sso_provider_id) WHERE sso_provider IS NOT NULL;
CREATE UNIQUE INDEX idx_user_credentials_sso_provider_id ON user_credentials(sso_provider, sso_provider_id)
    WHERE sso_provider IS NOT NULL;

COMMENT ON TABLE user_credentials IS 'User authentication credentials - isolated for security';
COMMENT ON COLUMN user_credentials.password_hash IS 'Bcrypt hash with cost factor 12';
COMMENT ON COLUMN user_credentials.pin_hash IS 'Bcrypt hash of 6-digit PIN (cost factor 10)';
COMMENT ON COLUMN user_credentials.password_history IS 'JSON array of last 5 password hashes to prevent reuse';
COMMENT ON COLUMN user_credentials.pin_history IS 'JSON array of last 3 PIN hashes to prevent reuse';


CREATE TABLE user_profiles (
    user_profile_id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,

    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
 
    address TEXT,
    phone VARCHAR(50),

    gender VARCHAR(20) CHECK (gender IN ('male', 'female', 'other')),
    marital_status VARCHAR(20) CHECK (marital_status IN ('single', 'married', 'divorced', 'widowed')),
    date_of_birth DATE,
    place_of_birth VARCHAR(100),

    avatar_url VARCHAR(500),
    preferred_language VARCHAR(10) DEFAULT 'en',
    timezone VARCHAR(50) DEFAULT 'Asia/Jakarta',

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT valid_age CHECK (
        date_of_birth IS NULL OR
        (date_of_birth <= CURRENT_DATE - INTERVAL '18 years' AND
         date_of_birth >= CURRENT_DATE - INTERVAL '75 years')
    )
);

CREATE INDEX idx_user_profiles_user ON user_profiles(user_id);
CREATE INDEX idx_user_profiles_name ON user_profiles(first_name, last_name);

COMMENT ON TABLE user_profiles IS 'User personal information and preferences - PII encrypted at rest';
COMMENT ON COLUMN user_profiles.first_name IS 'Legal first name (required for all users)';
COMMENT ON COLUMN user_profiles.last_name IS 'Legal last name (required for all users)';
COMMENT ON COLUMN user_profiles.address IS 'Full residential address (encrypted at rest with AES-256)';
COMMENT ON COLUMN user_profiles.date_of_birth IS 'Date of birth (must be 18-75 years old, encrypted at rest)';


CREATE TABLE user_security (
    user_security_id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,

    last_login_at TIMESTAMPTZ,
    last_login_ip INET,
    
    failed_login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMPTZ,
    
    admin_registered_at TIMESTAMPTZ,
    user_registered_at TIMESTAMPTZ,
    admin_registered_by UUID REFERENCES users(user_id) ON DELETE SET NULL,
    
    invitation_token_hash VARCHAR(255),
    invitation_expires_at TIMESTAMPTZ,
    
    metadata JSONB DEFAULT '{}'::JSONB,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_security_user ON user_security(user_id);
CREATE INDEX idx_user_security_invitation ON user_security(invitation_token_hash)
    WHERE invitation_token_hash IS NOT NULL;

COMMENT ON TABLE user_security IS 'User runtime security state - login tracking, locking, registration flow';
COMMENT ON COLUMN user_security.admin_registered_at IS 'Timestamp when admin pre-registered this user';
COMMENT ON COLUMN user_security.user_registered_at IS 'Timestamp when user completed self-registration';
COMMENT ON COLUMN user_security.invitation_token_hash IS 'SHA-256 hash of invitation token (for admin-initiated registration)';
