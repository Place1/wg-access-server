#!/bin/sh
set -e

SCRIPT_DIR="$(dirname $0)"
OUT_DIR="$SCRIPT_DIR/proto/proto"

mkdir -p "$OUT_DIR" || true

protoc \
  -I proto/ \
  --go_out="$OUT_DIR" --go_opt=paths=source_relative \
  --go-grpc_out="$OUT_DIR" --go-grpc_opt=paths=source_relative \
  proto/*.proto
