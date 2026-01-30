-- ============================================================================
-- Migration: 000001_create_extensions_and_functions
-- Description: Create PostgreSQL extensions and common trigger functions
-- Requires: PostgreSQL 18+ for native uuidv7() support
-- ============================================================================

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================================
-- FUNCTION: update_updated_at_column
-- Description: Trigger function to automatically update updated_at timestamp
-- ============================================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION update_updated_at_column() IS
    'Trigger function to automatically set updated_at to current timestamp on UPDATE';

-- ============================================================================
-- NOTE: PostgreSQL 18+ provides native uuidv7() function
-- UUIDv7 is time-ordered (timestamp in first 48 bits) for:
--   - Better B-tree index performance (sequential inserts)
--   - Natural time-based sorting
--   - Global uniqueness without coordination
-- ============================================================================
