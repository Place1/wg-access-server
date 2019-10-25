#!/bin/bash
set -eou pipefail

docker pull docker:stable

read -p "Version: " version

IMAGE="place1/wireguard-access-server:$version"

docker login
docker build -t "$IMAGE" .
docker push "$IMAGE"
