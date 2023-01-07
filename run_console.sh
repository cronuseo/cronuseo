#!/usr/bin/env bash

# Start the backend using Docker Compose
docker-compose up -d

# Wait for the backend API to start
while true; do
  HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v1/health)
  if [ "$HTTP_STATUS" -eq 200 ]; then
    break
  fi
  sleep 1
done

# Get the console repository
git clone https://github.com/shashimalcse/cronuseo-console

# Change into the console repository directory
cd cronuseo-console

# Start the frontend using Next.js
npm install
npm run dev
