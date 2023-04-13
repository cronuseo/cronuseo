#!/bin/bash -e



APP_ENV=${APP_ENV:-local}

echo "[`date`] Running entrypoint script in the '${APP_ENV}' environment..."

CONFIG_FILE=./config/${APP_ENV}.yml

if [[ -z ${APP_DSN} ]]; then
  export APP_DSN=`sed -n 's/^dsn:[[:space:]]*"\(.*\)"/\1/p' ${CONFIG_FILE}`
fi

echo "[`date`] Starting server..."
./server -config ${CONFIG_FILE} 
