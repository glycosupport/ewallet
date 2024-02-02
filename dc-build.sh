#!/bin/bash

bash ./docker/docker-compose-check.sh
if [[ $? -eq 1 ]]; then exit 1; fi

if [ $# -eq 0 ]
then
    echo "Building docker compose"
else
    echo "Building docker compose with additional parameter $1 ..."
fi

docker compose --env-file ./docker/environments/app.env build $1
