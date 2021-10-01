#!/bin/sh
set -e

SCRIPT_DIR="$(dirname $0)"
OUT_DIR="$SCRIPT_DIR/proto/proto"

mkdir -p "$OUT_DIR" || true

protoc \
  -I proto/ \
  proto/*.proto \
  --go_opt=paths=source_relative \
  --go_out="plugins=grpc:$OUT_DIR"
