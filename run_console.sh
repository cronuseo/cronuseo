#!/usr/bin/env bash

# Start the backend using Docker Compose

# Wait for the backend API to start
while true; do
  HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v1/health)
  if [ "$HTTP_STATUS" -eq 200 ]; then
    break
  fi
  sleep 1
done

CONTAINER_ID='cronuseo-cronuseo-db-1'

username=$(yq '.organization.config.username' ./config/local.yml)
password=$(yq '.organization.config.password' ./config/local.yml)
organization=$(yq '.organization.config.organization' ./config/local.yml)

eval "$(grep ^DB_HOST= .env)"
eval "$(grep ^DB_PORT= .env)"
eval "$(grep ^DB_USERNAME= .env)"
eval "$(grep ^DB_PASSWORD= .env)"
eval "$(grep ^DB_NAME= .env)"

# clear tables
docker exec -ti -e "PGPASSWORD=$DB_PASSWORD" $CONTAINER_ID psql -h $DB_HOST -U $DB_USERNAME -d $DB_NAME -c "TRUNCATE org, org_user, org_role, org_resource, user_role, res_action, org_admin_user RESTART IDENTITY CASCADE;"

# generate uuid for organization
org_id=$(uuidgen | tr '[:upper:]' '[:lower:]')

# generate api key
key=$(head -c 32 /dev/urandom | base64)

# create organization
docker exec -ti -e "PGPASSWORD=$DB_PASSWORD" $CONTAINER_ID psql -h $DB_HOST -U $DB_USERNAME -d $DB_NAME -c "INSERT INTO org(org_key,name,org_id, org_api_key) VALUES('$organization','$organization','$org_id'::uuid,'$key');"

# generate uuid for admin user
user_id=$(uuidgen | tr '[:upper:]' '[:lower:]')

# create organization
docker exec -ti -e "PGPASSWORD=$DB_PASSWORD" $CONTAINER_ID psql -h $DB_HOST -U $DB_USERNAME -d $DB_NAME -c "INSERT INTO org_admin_user(user_id,username,password,org_id,is_super) VALUES('$user_id','$username',crypt('$password', gen_salt('bf')),'$org_id'::uuid,'true');"

# Get the console repository
dir=cronuseo-console
 
if [ -d "$dir" -a ! -h "$dir" ]
then
  echo "$dir exists"
  cd cronuseo-console
  git pull
else
  echo "$dir not exists"
  git clone https://github.com/shashimalcse/cronuseo-console
  cd cronuseo-console
fi

echo "NEXTAUTH_SECRET=secret" >> .env
echo "BASE_API=http://localhost:8080/api/v1" >> .env

echo "NEXTAUTH_SECRET=secret" >> .env.local
echo "BASE_API=http://localhost:8080/api/v1" >> .env.local

# Change into the console repository directory


# Start the frontend using Next.js
npm install
npm run dev
