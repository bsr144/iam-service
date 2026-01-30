-- ============================================================================
-- Migration: 000007_create_user_profiles
-- Description: Create user_profiles table - personal profile information (1:1 with users)
-- ============================================================================

CREATE TABLE user_profiles (
    -- Primary Key (same as user_id for 1:1 relationship)
    user_id             UUID PRIMARY KEY,

    -- Required Profile Fields
    first_name          VARCHAR(100) NOT NULL,
    last_name           VARCHAR(100) NOT NULL,

    -- Optional Profile Fields
    phone               VARCHAR(50),
    date_of_birth       DATE,
    gender              VARCHAR(20),
    marital_status      VARCHAR(20),
    address             TEXT,
    id_number           VARCHAR(100),
    avatar_url          VARCHAR(500),

    -- Localization
    preferred_language  VARCHAR(10) NOT NULL DEFAULT 'en',
    timezone            VARCHAR(50) NOT NULL DEFAULT 'Asia/Jakarta',

    -- Flexible Metadata (tenant-specific fields)
    metadata            JSONB NOT NULL DEFAULT '{}',

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_user_profiles_user FOREIGN KEY (user_id)
        REFERENCES users(user_id) ON DELETE CASCADE
);

-- Trigger for updated_at
CREATE TRIGGER trg_user_profiles_updated_at
    BEFORE UPDATE ON user_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Indexes for search
CREATE INDEX idx_user_profiles_name ON user_profiles(LOWER(first_name), LOWER(last_name));
CREATE INDEX idx_user_profiles_phone ON user_profiles(phone) WHERE phone IS NOT NULL;

-- Comments
COMMENT ON TABLE user_profiles IS 'Personal profile information separated from user identity';
COMMENT ON COLUMN user_profiles.metadata IS 'Flexible JSON storage for tenant-specific profile fields';
COMMENT ON COLUMN user_profiles.id_number IS 'National ID - consider application-level encryption';
