#!/bin/bash
set -eou pipefail

echo "echo current tags:"
curl -L -s 'https://registry.hub.docker.com/v2/repositories/place1/wg-access-server/tags?page_size=10' \
  | jq '.results[].name'

docker login

read -p "Version: " version

IMAGE="place1/wg-access-server:$version"

docker build -t "$IMAGE" .
docker push "$IMAGE"
