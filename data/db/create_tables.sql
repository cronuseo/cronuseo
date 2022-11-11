CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
SELECT uuid_generate_v4();
CREATE TABLE if not exists organization(
   org_id uuid DEFAULT uuid_generate_v4 (),
   org_key VARCHAR(40) NOT NULL,
   name VARCHAR(40) NOT NULL,
   PRIMARY KEY ( org_id ));
   
CREATE TABLE if not exists tenant(
   tenant_id uuid DEFAULT uuid_generate_v4 (),
   tenant_key VARCHAR(40) NOT NULL,
   name VARCHAR(40) NOT NULL,
   org_id uuid,
   PRIMARY KEY ( tenant_id ),
   CONSTRAINT FK_organization_tenant FOREIGN KEY(org_id) REFERENCES organization(org_id)
   );
 
  CREATE TABLE if not exists project(
   project_id uuid DEFAULT uuid_generate_v4 (),
   project_key VARCHAR(40) NOT NULL,
   name VARCHAR(40),
   tenant_id uuid,
   PRIMARY KEY ( project_id ),
   CONSTRAINT FK_tenant_project FOREIGN KEY(tenant_id) REFERENCES tenant(tenant_id)
   ); 
  
 CREATE TABLE if not exists tenant_user(
   user_id uuid DEFAULT uuid_generate_v4 (),
   username VARCHAR(40) NOT NULL,
   first_name VARCHAR(100),
   last_name VARCHAR(100),
   tenant_id uuid,
   PRIMARY KEY ( user_id ),
   CONSTRAINT FK_tenant_org_user FOREIGN KEY(tenant_id) REFERENCES tenant(tenant_id)
   );
 
 CREATE TABLE if not exists tenant_group(
   group_id uuid DEFAULT uuid_generate_v4 (),
   group_key VARCHAR(40) NOT NULL,
   name VARCHAR(40) NOT NULL,
   tenant_id uuid,
   PRIMARY KEY ( group_id ),
   CONSTRAINT FK_tenant_org_group FOREIGN KEY(tenant_id) REFERENCES tenant(tenant_id)
   );
   
 CREATE TABLE if not exists group_user(
   group_id uuid,
   user_id uuid,
   CONSTRAINT FK_group_user_tenant_user FOREIGN KEY(user_id) references tenant_user(user_id),
   CONSTRAINT FK_group_user_tenant_group FOREIGN KEY(group_id) REFERENCES tenant_group(group_id)
   );
   
 CREATE TABLE if not exists resource (
   resource_id uuid DEFAULT uuid_generate_v4 (),
   resource_key VARCHAR(40) NOT NULL,
   name VARCHAR(40),
   tenant_id uuid,
   PRIMARY KEY ( resource_id ),
   CONSTRAINT FK_tenant_resource FOREIGN KEY(tenant_id) REFERENCES tenant(tenant_id)
   ); 
   
CREATE TABLE if not exists resource_action (
   resource_action_id uuid DEFAULT uuid_generate_v4 (),
   resource_action_key VARCHAR(40) NOT NULL,
   name VARCHAR(40),
   resource_id uuid,
   PRIMARY KEY ( resource_action_id ),
   CONSTRAINT FK_resource_resource_action FOREIGN KEY(resource_id) REFERENCES resource(resource_id)
   ); 
   
CREATE TABLE if not exists resource_role (
   resource_role_id uuid DEFAULT uuid_generate_v4 (),
   resource_role_key VARCHAR(40) NOT NULL,
   name VARCHAR(40),
   resource_id uuid,
   PRIMARY KEY ( resource_role_id ),
   CONSTRAINT FK_resource_resource_role FOREIGN KEY(resource_id) REFERENCES resource(resource_id)
   ); 
   
 CREATE TABLE if not exists action_role(
   resource_role_id uuid,
   resource_action_id uuid,
   CONSTRAINT FK_action_role_resource_role FOREIGN KEY(resource_role_id) REFERENCES resource_role(resource_role_id),
   CONSTRAINT FK_action_role_resource_action FOREIGN KEY(resource_action_id) REFERENCES resource_action(resource_action_id)
   );
   
CREATE TABLE if not exists user_resource_role(
   user_id uuid,
   resource_role_id uuid,
   CONSTRAINT FK_resource_role_user_resource_role FOREIGN KEY(resource_role_id) REFERENCES resource_role(resource_role_id),
   CONSTRAINT FK_tenant_user_user_resource_role FOREIGN KEY(user_id) REFERENCES tenant_user(user_id)
   );
   
CREATE TABLE if not exists group_resource_role(
   group_id uuid,
   resource_role_id uuid,
   CONSTRAINT FK_resource_role_group_resource_role FOREIGN KEY(resource_role_id) REFERENCES resource_role(resource_role_id),
   CONSTRAINT FK_tenant_user_group_resource_role FOREIGN KEY(group_id) REFERENCES tenant_group(group_id)
   );
   
CREATE TABLE if not exists resource_action_role(
   resource_action_id uuid,
   resource_role_id uuid,
   resource_id uuid,
   CONSTRAINT FK_resource_action_role_role FOREIGN KEY(resource_role_id) REFERENCES resource_role(resource_role_id),
   CONSTRAINT FK_resource_action_role_action FOREIGN KEY(resource_action_id) REFERENCES resource_action(resource_action_id),
   CONSTRAINT FK_resource_action_role_resource FOREIGN KEY(resource_id) REFERENCES resource(resource_id)
   );
  
CREATE TABLE if not exists resource_action_role_key(
   resource_action_key VARCHAR(40) NOT NULL,
   resource_role_key VARCHAR(40) NOT NULL,
   resource_key VARCHAR(160) NOT NULL
   );
   