#!/usr/bin/env bash

set -e

GOMPLATE_URL="github.com/hairyhenderson/gomplate/v4/cmd/gomplate@latest"
GOMPLATE_BIN="$(command -v gomplate || true)"

if [ -z "$GOMPLATE_BIN" ]; then
  echo "gomplate not found, installing..."
  go install "$GOMPLATE_URL"
else
  echo "gomplate found at : $GOMPLATE_BIN"
fi
