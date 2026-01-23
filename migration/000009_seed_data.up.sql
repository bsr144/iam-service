INSERT INTO permissions (code, name, description, module, resource, action, scope_level, is_system) VALUES
('auth:login', 'Login', 'Authenticate and obtain access token', 'iam', 'auth', 'login', 'system', TRUE),
('auth:logout', 'Logout', 'Revoke access token', 'iam', 'auth', 'logout', 'self', TRUE),
('auth:refresh', 'Refresh Token', 'Refresh access token', 'iam', 'auth', 'refresh', 'self', TRUE),

('auth:setup-pin', 'Setup PIN', 'Set up 6-digit PIN during registration', 'iam', 'auth', 'setup_pin', 'self', TRUE),
('auth:verify-pin', 'Verify PIN', 'Verify PIN for sensitive operations', 'iam', 'auth', 'verify_pin', 'self', TRUE),
('auth:reset-pin', 'Reset PIN', 'Reset forgotten PIN (self-service)', 'iam', 'auth', 'reset_pin', 'self', TRUE),

('auth:mfa:setup', 'Setup MFA', 'Configure TOTP multi-factor authentication', 'iam', 'auth', 'mfa_setup', 'self', TRUE),
('auth:mfa:verify', 'Verify MFA', 'Verify MFA code during login', 'iam', 'auth', 'mfa_verify', 'self', TRUE),

('user:create', 'Create User', 'Create new user account', 'iam', 'user', 'create', 'tenant', TRUE),
('user:read', 'Read User', 'View user details', 'iam', 'user', 'read', 'branch', TRUE),
('user:read:all', 'Read All Users', 'View all users (tenant-wide)', 'iam', 'user', 'read', 'tenant', TRUE),
('user:update', 'Update User', 'Modify user details', 'iam', 'user', 'update', 'branch', TRUE),
('user:update:self', 'Update Own Profile', 'Modify own user profile', 'iam', 'user', 'update', 'self', TRUE),
('user:delete', 'Delete User', 'Soft-delete user account', 'iam', 'user', 'delete', 'tenant', TRUE),
('user:reset-password', 'Reset User Password', 'Initiate password reset for user', 'iam', 'user', 'reset_password', 'branch', TRUE),
('user:reset-pin-admin', 'Admin Reset User PIN', 'Reset user PIN as admin', 'iam', 'user', 'reset_pin', 'branch', TRUE),
('user:assign-role', 'Assign Role to User', 'Assign roles to user', 'iam', 'user', 'assign_role', 'tenant', TRUE),
('user:approve-registration', 'Approve Registration', 'Approve pending participant registrations', 'iam', 'user', 'approve', 'branch', TRUE),
('user:create-via-api', 'Create User via Secret API', 'Create admin/approver via secret API', 'iam', 'user', 'create_api', 'tenant', TRUE),

('role:create', 'Create Role', 'Create new role', 'iam', 'role', 'create', 'tenant', TRUE),
('role:read', 'Read Role', 'View role details', 'iam', 'role', 'read', 'tenant', TRUE),
('role:update', 'Update Role', 'Modify role details', 'iam', 'role', 'update', 'tenant', TRUE),
('role:delete', 'Delete Role', 'Delete role', 'iam', 'role', 'delete', 'tenant', TRUE),
('role:assign-permission', 'Assign Permission to Role', 'Assign permissions to role', 'iam', 'role', 'assign_permission', 'tenant', TRUE),

('permission:read', 'Read Permission', 'View available permissions', 'iam', 'permission', 'read', 'system', TRUE),

('audit:read', 'Read Audit Logs', 'View audit logs', 'iam', 'audit', 'read', 'tenant', TRUE),
('audit:read:auth', 'Read Auth Logs', 'View authentication logs', 'iam', 'audit', 'read', 'tenant', TRUE),

('employee:create', 'Create Employee', 'Register new employee', 'employee', 'employee', 'create', 'branch', TRUE),
('employee:read', 'Read Employee', 'View employee details', 'employee', 'employee', 'read', 'branch', TRUE),
('employee:read:all', 'Read All Employees', 'View all employees (tenant-wide)', 'employee', 'employee', 'read', 'tenant', TRUE),
('employee:update', 'Update Employee', 'Modify employee details', 'employee', 'employee', 'update', 'branch', TRUE),
('employee:delete', 'Delete Employee', 'Soft-delete employee', 'employee', 'employee', 'delete', 'branch', TRUE),
('employee:import', 'Import Employees', 'Bulk import employees via Excel', 'employee', 'employee', 'import', 'branch', TRUE),
('employee:export', 'Export Employees', 'Export employee data to Excel', 'employee', 'employee', 'export', 'branch', TRUE),

('contribution:create', 'Upload Contribution', 'Upload contribution file', 'contribution', 'contribution', 'create', 'branch', TRUE),
('contribution:read', 'Read Contribution', 'View contribution details', 'contribution', 'contribution', 'read', 'branch', TRUE),
('contribution:read:all', 'Read All Contributions', 'View all contributions (tenant-wide)', 'contribution', 'contribution', 'read', 'tenant', TRUE),
('contribution:approve', 'Approve Contribution', 'Approve uploaded contribution', 'contribution', 'contribution', 'approve', 'branch', TRUE),
('contribution:verify', 'Verify Contribution', 'Final verification of contribution', 'contribution', 'contribution', 'verify', 'branch', TRUE),
('contribution:reject', 'Reject Contribution', 'Reject contribution upload', 'contribution', 'contribution', 'reject', 'branch', TRUE),

('allocation:create', 'Create Allocation', 'Create new allocation', 'allocation', 'allocation', 'create', 'branch', TRUE),
('allocation:read', 'Read Allocation', 'View allocation details', 'allocation', 'allocation', 'read', 'branch', TRUE),
('allocation:read:all', 'Read All Allocations', 'View all allocations (tenant-wide)', 'allocation', 'allocation', 'read', 'tenant', TRUE),
('allocation:process', 'Process Allocation', 'Calculate allocation distribution', 'allocation', 'allocation', 'process', 'branch', TRUE),
('allocation:validate', 'Validate Allocation', 'Validate calculated allocation', 'allocation', 'allocation', 'validate', 'branch', TRUE),
('allocation:verify', 'Verify Allocation', 'Final verification of allocation', 'allocation', 'allocation', 'verify', 'branch', TRUE),
('allocation:reject', 'Reject Allocation', 'Reject allocation', 'allocation', 'allocation', 'reject', 'branch', TRUE),

('master-data:create', 'Create Master Data', 'Create master data entries', 'master-data', 'master-data', 'create', 'tenant', TRUE),
('master-data:read', 'Read Master Data', 'View master data entries', 'master-data', 'master-data', 'read', 'branch', TRUE),
('master-data:update', 'Update Master Data', 'Modify master data entries', 'master-data', 'master-data', 'update', 'tenant', TRUE),
('master-data:delete', 'Delete Master Data', 'Delete master data entries', 'master-data', 'master-data', 'delete', 'tenant', TRUE),

('report:view', 'View Reports', 'View branch reports', 'report', 'report', 'view', 'branch', TRUE),
('report:view:tenant', 'View Tenant Reports', 'View consolidated tenant reports', 'report', 'report', 'view', 'tenant', TRUE),
('report:export', 'Export Reports', 'Export reports to PDF/Excel', 'report', 'report', 'export', 'branch', TRUE),

('translation:create', 'Create Translation', 'Add new translation', 'translation', 'translation', 'create', 'tenant', TRUE),
('translation:read', 'Read Translation', 'View translations', 'translation', 'translation', 'read', 'system', TRUE),
('translation:update', 'Update Translation', 'Modify translation', 'translation', 'translation', 'update', 'tenant', TRUE),
('translation:delete', 'Delete Translation', 'Delete translation', 'translation', 'translation', 'delete', 'tenant', TRUE)
ON CONFLICT (code) DO NOTHING;

INSERT INTO roles (code, name, description, scope_level, is_system) VALUES
('ADMIN', 'Administrator', 'Full administrative access within tenant (user management, approvals, master data)', 'tenant', TRUE),
('APPROVER', 'Approver', 'Can approve participant registrations and review submissions', 'branch', TRUE),
('PARTICIPANT', 'Participant', 'End user participant with self-service access', 'branch', TRUE)
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r, permissions p
WHERE r.code = 'ADMIN' AND r.is_system = TRUE
  AND p.scope_level IN ('tenant', 'branch', 'self')
  AND p.code NOT LIKE 'auth:mfa%'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r, permissions p
WHERE r.code = 'APPROVER' AND r.is_system = TRUE
  AND p.code IN (
    'auth:login', 'auth:logout', 'auth:refresh',
    'auth:setup-pin', 'auth:verify-pin', 'auth:reset-pin',
    'user:read', 'user:update:self',
    'user:approve-registration', 'user:reset-pin-admin',
    'employee:read',
    'contribution:read', 'contribution:approve', 'contribution:verify', 'contribution:reject',
    'allocation:read', 'allocation:validate', 'allocation:verify', 'allocation:reject',
    'report:view', 'report:export',
    'master-data:read'
  )
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r, permissions p
WHERE r.code = 'PARTICIPANT' AND r.is_system = TRUE
  AND p.code IN (
    'auth:login', 'auth:logout', 'auth:refresh',
    'auth:setup-pin', 'auth:verify-pin', 'auth:reset-pin',
    'user:update:self',
    'report:view'
  )
ON CONFLICT (role_id, permission_id) DO NOTHING;
