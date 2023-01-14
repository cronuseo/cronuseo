CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

SELECT uuid_generate_v4();

CREATE TABLE
    if not exists ORG(
        id SERIAL,
        org_id uuid PRIMARY KEY,
        org_key VARCHAR(40) NOT NULL,
        name VARCHAR(40) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE
    if not exists ORG_USER(
        id SERIAL,
        user_id uuid PRIMARY KEY,
        username VARCHAR(40) NOT NULL,
        firstname VARCHAR(40),
        lastname VARCHAR(40),
        org_id uuid,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        CONSTRAINT FK_ORG_ORG_USER FOREIGN KEY(org_id) REFERENCES ORG(org_id)
    );

CREATE TABLE
    if not exists ORG_ROLE(
        id SERIAL,
        role_id uuid PRIMARY KEY,
        role_key VARCHAR(40) NOT NULL,
        name VARCHAR(40) NOT NULL,
        org_id uuid,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        CONSTRAINT FK_ORG_ORG_ROLE FOREIGN KEY(org_id) REFERENCES ORG(org_id)
    );

CREATE TABLE
    if not exists ORG_RESOURCE(
        id SERIAL,
        resource_id uuid PRIMARY KEY,
        resource_key VARCHAR(40) NOT NULL,
        name VARCHAR(40) NOT NULL,
        org_id uuid,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        CONSTRAINT FK_ORG_ORG_RESOURCE FOREIGN KEY(org_id) REFERENCES ORG(org_id)
    );

CREATE TABLE
    if not exists USER_ROLE(
        id SERIAL,
        role_id uuid,
        user_id uuid,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        CONSTRAINT FK_ORG_ROLE_USER_ROLE FOREIGN KEY(role_id) REFERENCES ORG_ROLE(role_id),
        CONSTRAINT FK_ORG_USER_USER_ROLE FOREIGN KEY(user_id) REFERENCES ORG_USER(user_id)
    );

CREATE TABLE
    if not exists RES_ACTION(
        id SERIAL,
        action_id uuid PRIMARY KEY,
        action_key VARCHAR(40) NOT NULL,
        name VARCHAR(40) NOT NULL,
        resource_id uuid,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        CONSTRAINT FK_ORG_RESOURCE_PERMISSION FOREIGN KEY(resource_id) REFERENCES ORG_RESOURCE(resource_id)
    );

CREATE TABLE
    if not exists ORG_ADMIN_USER(
        id SERIAL,
        user_id uuid PRIMARY KEY,
        username VARCHAR(40) NOT NULL,
        password CHAR(60) NOT NULL,
        is_super BOOLEAN NOT NULL,
        org_id uuid,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE OR REPLACE FUNCTION TRIGGER_SET_TIMESTAMP() 
RETURNS TRIGGER AS 
	$$ BEGIN NEW.updated_at = NOW();
	RETURN NEW;
END; 
$$ LANGUAGE plpgsql;
DO $$
DECLARE
    t text;
BEGIN
    FOR t IN 
        SELECT  table_name FROM information_schema.columns
             WHERE column_name = 'updated_at'    
    LOOP 


        EXECUTE format('CREATE TRIGGER set_timestamp
                        BEFORE UPDATE ON %I
                        FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp()',
                        t);
    END loop;
    END;

$$ LANGUAGE plpgsql;