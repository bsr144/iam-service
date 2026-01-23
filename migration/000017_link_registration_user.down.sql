DROP INDEX IF EXISTS idx_users_registration;

ALTER TABLE users DROP COLUMN IF EXISTS registration_completed_at;
ALTER TABLE users DROP COLUMN IF EXISTS registration_id;
