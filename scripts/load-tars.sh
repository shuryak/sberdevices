#!/usr/bin/env bash

set -e

for file in ./.tmp/*; do
  if [[ -f "$file" ]]; then
    echo "Loading $file into Docker..."
    docker load -i "$file"
  fi
done
