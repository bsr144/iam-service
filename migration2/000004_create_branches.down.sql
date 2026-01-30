-- ============================================================================
-- Migration: 000004_create_branches (DOWN)
-- ============================================================================

DROP TRIGGER IF EXISTS trg_branches_updated_at ON branches;
DROP TABLE IF EXISTS branches;
