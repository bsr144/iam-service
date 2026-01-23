CREATE TABLE public.saml_configurations (
    saml_configuration_id uuid DEFAULT uuidv7() NOT NULL,
    tenant_id uuid NOT NULL,
    idp_entity_id varchar(500) NOT NULL,
    idp_sso_url varchar(500) NOT NULL,
    idp_slo_url varchar(500) NULL,
    idp_certificate text NOT NULL,
    sp_entity_id varchar(500) NOT NULL,
    sp_acs_url varchar(500) NOT NULL,
    sp_slo_url varchar(500) NULL,
    attribute_mapping jsonb DEFAULT '{"name": "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name", "email": "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress", "roles": "http://schemas.microsoft.com/ws/2008/06/identity/claims/role"}'::jsonb NOT NULL,
    role_mapping jsonb DEFAULT '{}'::jsonb NULL,
    auto_provision_users bool DEFAULT true NULL,
    default_branch_id uuid NULL,
    is_active bool DEFAULT true NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
    deleted_at timestamptz NULL,
    CONSTRAINT saml_configurations_pkey PRIMARY KEY (saml_configuration_id),
    CONSTRAINT unique_saml_per_tenant UNIQUE (tenant_id),
    CONSTRAINT saml_configurations_default_branch_id_fkey
        FOREIGN KEY (default_branch_id)
        REFERENCES public.branches(branch_id),
    CONSTRAINT saml_configurations_tenant_id_fkey
        FOREIGN KEY (tenant_id)
        REFERENCES public.tenants(tenant_id)
        ON DELETE CASCADE
);

CREATE INDEX idx_saml_tenant
    ON public.saml_configurations USING btree (tenant_id)
    WHERE (deleted_at IS NULL);

CREATE TRIGGER update_saml_configurations_updated_at
BEFORE UPDATE ON public.saml_configurations
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
