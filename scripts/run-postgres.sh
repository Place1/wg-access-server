#!/bin/bash
set -eou pipefail

NAME="wg-access-server"

if [[ ! "$(docker ps -aqf name=$NAME)" ]]; then
  docker run \
    -e 'POSTGRES_USER=postgres' \
    -e 'POSTGRES_PASSWORD=example' \
    -e 'POSTGRES_DB=postgres' \
    -p 5432:5432 \
    -d \
    --name "$NAME" \
    postgres:11-alpine
else
  docker start "$NAME"
fi

echo "started postgres: $NAME"
