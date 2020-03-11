#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
OUT_DIR="$SCRIPT_DIR/proto/proto"

mkdir -p "$OUT_DIR" || true

protoc \
  -I proto/ \
  proto/*.proto \
  --go_out="plugins=grpc:$OUT_DIR"

cd website || exit
npm run codegen
