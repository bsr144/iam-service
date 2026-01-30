-- ============================================================================
-- Migration: 000001_create_extensions_and_functions (DOWN)
-- Description: Drop common trigger functions
-- ============================================================================

DROP FUNCTION IF EXISTS update_updated_at_column();

-- Note: Extensions are not dropped to avoid affecting other schemas

