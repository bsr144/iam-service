DROP TABLE IF EXISTS schema_version;


DROP FUNCTION IF EXISTS get_user_branches(UUID);
DROP FUNCTION IF EXISTS user_has_permission(UUID, VARCHAR, UUID);


DROP VIEW IF EXISTS v_user_roles;
DROP VIEW IF EXISTS v_user_permissions;
DROP VIEW IF EXISTS v_users_complete;


DROP TRIGGER IF EXISTS update_user_activation_tracking_updated_at ON user_activation_tracking;
DROP TRIGGER IF EXISTS update_roles_updated_at ON roles;
DROP TRIGGER IF EXISTS update_permissions_updated_at ON permissions;
DROP TRIGGER IF EXISTS update_saml_configurations_updated_at ON saml_configurations;
DROP TRIGGER IF EXISTS update_user_security_updated_at ON user_security;
DROP TRIGGER IF EXISTS update_user_profiles_updated_at ON user_profiles;
DROP TRIGGER IF EXISTS update_user_credentials_updated_at ON user_credentials;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_branch_contacts_updated_at ON branch_contacts;
DROP TRIGGER IF EXISTS update_branches_updated_at ON branches;
DROP TRIGGER IF EXISTS update_tenant_settings_updated_at ON tenant_settings;
DROP TRIGGER IF EXISTS update_tenants_updated_at ON tenants;


DROP FUNCTION IF EXISTS update_updated_at_column();
