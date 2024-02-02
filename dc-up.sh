#!/bin/bash

bash ./docker/docker-compose-check.sh
if [[ $? -eq 1 ]]; then exit 1; fi

echo "Starting docker compose in the foreground ..."

docker compose --env-file ./docker/environments/app.env up --no-deps
