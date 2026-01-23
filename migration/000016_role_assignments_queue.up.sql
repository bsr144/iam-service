CREATE TABLE role_assignments_queue (
    queue_id UUID PRIMARY KEY DEFAULT uuidv7(),

    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id) ON DELETE CASCADE,

    role_id UUID NOT NULL REFERENCES roles(role_id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(product_id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(branch_id) ON DELETE SET NULL,

    effective_from TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    effective_to TIMESTAMPTZ,

    status VARCHAR(50) NOT NULL DEFAULT 'pending',

    assigned_by UUID NOT NULL REFERENCES users(user_id),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    batch_id UUID,
    batch_total INT,
    batch_sequence INT,

    processed_at TIMESTAMPTZ,
    processing_started_at TIMESTAMPTZ,
    failure_reason TEXT,
    retry_count INT DEFAULT 0,

    user_role_id UUID REFERENCES user_roles(user_role_id) ON DELETE SET NULL,

    metadata JSONB DEFAULT '{}',

    CONSTRAINT role_assignments_queue_status_check CHECK (
        status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')
    )
);
CREATE INDEX idx_role_queue_status ON role_assignments_queue(status)
    WHERE status IN ('pending', 'processing');

CREATE INDEX idx_role_queue_user ON role_assignments_queue(user_id, status);

CREATE INDEX idx_role_queue_batch ON role_assignments_queue(batch_id)
    WHERE batch_id IS NOT NULL;

CREATE INDEX idx_role_queue_assigned_by ON role_assignments_queue(assigned_by, assigned_at);

CREATE INDEX idx_role_queue_tenant ON role_assignments_queue(tenant_id);

CREATE INDEX idx_role_queue_processing ON role_assignments_queue(processing_started_at)
    WHERE status = 'processing';
CREATE UNIQUE INDEX idx_role_queue_unique_pending ON role_assignments_queue(user_id, role_id, COALESCE(product_id, '00000000-0000-0000-0000-000000000000'::UUID))
    WHERE status IN ('pending', 'processing');
COMMENT ON TABLE role_assignments_queue IS 'Queue for admin role assignments before user activation - supports bulk operations';
COMMENT ON COLUMN role_assignments_queue.batch_id IS 'Groups bulk role assignments together for progress tracking';
COMMENT ON COLUMN role_assignments_queue.batch_total IS 'Total assignments in batch (for progress percentage)';
COMMENT ON COLUMN role_assignments_queue.batch_sequence IS 'Order in batch (1, 2, 3, ..., batch_total)';
COMMENT ON COLUMN role_assignments_queue.assigned_by IS 'Admin who initiated the role assignment';
COMMENT ON COLUMN role_assignments_queue.user_role_id IS 'Set after processing completes, links to created user_role';
COMMENT ON COLUMN role_assignments_queue.retry_count IS 'Number of retry attempts for failed assignments';
