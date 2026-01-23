DELETE FROM role_permissions WHERE role_id IN (
    SELECT id FROM roles WHERE is_system = TRUE
);

DELETE FROM roles WHERE is_system = TRUE;

DELETE FROM permissions WHERE is_system = TRUE;
