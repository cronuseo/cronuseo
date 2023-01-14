BEGIN;
TRUNCATE org, org_user, org_role, org_resource, user_role, res_action, org_admin_user RESTART IDENTITY CASCADE;
COMMIT;