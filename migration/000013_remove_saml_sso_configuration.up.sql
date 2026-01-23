DROP TRIGGER IF EXISTS update_saml_configurations_updated_at ON public.saml_configurations;

DROP INDEX IF EXISTS public.idx_saml_tenant;

DROP TABLE IF EXISTS public.saml_configurations;
