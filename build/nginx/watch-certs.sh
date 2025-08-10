#!/usr/bin/env sh

echo "Waiting certificate changes in: $CERTS_DIR"

while true; do
    inotifywait -e modify,create,delete,move "$CERTS_DIR"
    echo "Detected change in cert files. Reloading nginx..."
    nginx -s reload
done
