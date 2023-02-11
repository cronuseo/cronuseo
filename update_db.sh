#!/usr/bin/env bash

CONTAINER_ID='db'

username=$(yq '.organization.config.username' ./config/local.yml)
password=$(yq '.organization.config.password' ./config/local.yml)
organization=$(yq '.organization.config.organization' ./config/local.yml)

DB_HOST='localhost'
DB_PORT='5432'
DB_USERNAME='postgres'
DB_PASSWORD='postgres' 
DB_NAME='cronuseo' 


# Continuously check the status of the database
while true; do
  # Use docker exec to run psql inside the container
  result=$(docker exec "$CONTAINER_ID" psql -U postgres -c '\l')

  # Check the exit status of docker exec
  if [ $? -eq 0 ]; then
    # The database is up, print a message
    echo "The database is up and accepting connections."
    break
  else
    # The database is down, print a message
    echo "The database is down."
  fi

  # Sleep for a short time before checking again
  sleep 5
done

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
