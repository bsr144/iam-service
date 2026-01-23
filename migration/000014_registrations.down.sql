-- Drop registrations table
DROP INDEX IF EXISTS idx_registrations_created_at;
DROP INDEX IF EXISTS idx_registrations_tenant;
DROP INDEX IF EXISTS idx_registrations_email;
DROP INDEX IF EXISTS idx_registrations_expires;
DROP INDEX IF EXISTS idx_registrations_status;
DROP INDEX IF EXISTS idx_registrations_unique_pending;

DROP TABLE IF EXISTS registrations;
