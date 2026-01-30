-- ============================================================================
-- Migration: 000005_create_branch_contacts (DOWN)
-- ============================================================================

DROP TRIGGER IF EXISTS trg_branch_contacts_updated_at ON branch_contacts;
DROP TABLE IF EXISTS branch_contacts;
