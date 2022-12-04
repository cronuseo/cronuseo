CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
SELECT uuid_generate_v4();
CREATE TABLE if not exists ORG(
   org_id uuid PRIMARY KEY,
   org_key VARCHAR(40) NOT NULL,
   name VARCHAR(40) NOT NULL);
   
 CREATE TABLE if not exists ORG_USER(
   user_id uuid PRIMARY KEY,
   username VARCHAR(40) NOT NULL,
   firstname VARCHAR(40),
   lastname VARCHAR(40),
   org_id uuid,
   CONSTRAINT FK_ORG_ORG_USER FOREIGN KEY(org_id) REFERENCES ORG(org_id)
   );
 
 CREATE TABLE if not exists ORG_ROLE(
   role_id uuid PRIMARY KEY,
   role_key VARCHAR(40) NOT NULL,
   name VARCHAR(40) NOT NULL,
   org_id uuid,
   CONSTRAINT FK_ORG_ORG_ROLE FOREIGN KEY(org_id) REFERENCES ORG(org_id)
   );

 CREATE TABLE if not exists ORG_RESOURCE(
   resource_id uuid PRIMARY KEY,
   resource_key VARCHAR(40) NOT NULL,
   name VARCHAR(40) NOT NULL,
   org_id uuid,
   CONSTRAINT FK_ORG_ORG_RESOURCE FOREIGN KEY(org_id) REFERENCES ORG(org_id)
   );
   
 CREATE TABLE if not exists USER_ROLE(
   role_id uuid,
   user_id uuid,
   CONSTRAINT FK_ORG_ROLE_USER_ROLE FOREIGN KEY(role_id) REFERENCES ORG_ROLE(role_id),
   CONSTRAINT FK_ORG_USER_USER_ROLE FOREIGN KEY(user_id) REFERENCES ORG_USER(user_id)
   );

 CREATE TABLE if not exists RES_ACTION(
   action_id uuid PRIMARY KEY,
   action_key VARCHAR(40) NOT NULL,
   name VARCHAR(40) NOT NULL,
   resource_id uuid,
   CONSTRAINT FK_ORG_RESOURCE_PERMISSION FOREIGN KEY(resource_id) REFERENCES ORG_RESOURCE(resource_id)
   );