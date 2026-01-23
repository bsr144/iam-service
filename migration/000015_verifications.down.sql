DROP INDEX IF EXISTS idx_verification_challenge_created_at;
DROP INDEX IF EXISTS idx_verification_challenge_verification;
DROP INDEX IF EXISTS idx_verification_challenge_expires;
DROP INDEX IF EXISTS idx_verification_challenge_active;

DROP TABLE IF EXISTS verification_challenges;

DROP INDEX IF EXISTS idx_verifications_created_at;
DROP INDEX IF EXISTS idx_verifications_status;
DROP INDEX IF EXISTS idx_verifications_expires;
DROP INDEX IF EXISTS idx_verifications_purpose;
DROP INDEX IF EXISTS idx_verifications_entity;

DROP TABLE IF EXISTS verifications;
