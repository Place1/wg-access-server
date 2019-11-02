#!/bin/bash
set -eou pipefail

docker login

read -p "Version: " version

IMAGE="place1/wg-access-server:$version"

docker build -t "$IMAGE" .
docker push "$IMAGE"
