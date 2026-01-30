-- ============================================================================
-- Migration: 000013_create_applications (DOWN)
-- ============================================================================

DROP TRIGGER IF EXISTS trg_applications_updated_at ON applications;
DROP TABLE IF EXISTS applications;
