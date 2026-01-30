-- ============================================================================
-- Migration: 000005_create_branch_contacts
-- Description: Create branch_contacts table for branch contact information
-- ============================================================================

CREATE TABLE branch_contacts (
    -- Primary Key
    branch_contact_id   UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Foreign Keys
    branch_id           UUID NOT NULL,

    -- Contact Information
    address             TEXT,
    city                VARCHAR(100),
    province            VARCHAR(100),
    postal_code         VARCHAR(20),
    country             VARCHAR(100) DEFAULT 'Indonesia',
    phone               VARCHAR(50),
    email               VARCHAR(255),
    fax                 VARCHAR(50),

    -- Coordinates (for mapping)
    latitude            DECIMAL(10, 8),
    longitude           DECIMAL(11, 8),

    -- Audit Fields
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_branch_contacts_branch FOREIGN KEY (branch_id)
        REFERENCES branches(branch_id) ON DELETE CASCADE,
    CONSTRAINT uq_branch_contacts_branch_id UNIQUE (branch_id)
);

-- Trigger for updated_at
CREATE TRIGGER trg_branch_contacts_updated_at
    BEFORE UPDATE ON branch_contacts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE branch_contacts IS 'Contact information for branches (1:1 relationship)';
