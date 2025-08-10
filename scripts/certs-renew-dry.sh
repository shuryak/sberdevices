#!/usr/bin/env bash

docker compose -f ./build/docker-compose.yaml run --rm --name certbot-renew-dry certbot-renew certbot renew --dry-run
sudo openssl x509 -enddate -noout -in ./build/certbot/letsencrypt/live/$PUBLIC_DOMAIN/fullchain.pem
