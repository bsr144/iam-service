DROP INDEX IF EXISTS idx_role_queue_unique_pending;
DROP INDEX IF EXISTS idx_role_queue_processing;
DROP INDEX IF EXISTS idx_role_queue_tenant;
DROP INDEX IF EXISTS idx_role_queue_assigned_by;
DROP INDEX IF EXISTS idx_role_queue_batch;
DROP INDEX IF EXISTS idx_role_queue_user;
DROP INDEX IF EXISTS idx_role_queue_status;

DROP TABLE IF EXISTS role_assignments_queue;
