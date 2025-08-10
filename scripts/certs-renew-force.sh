#!/usr/bin/env bash

docker compose -f ./build/docker-compose.yaml run --rm --name certbot-renew-dry certbot-renew certbot renew --force-renewal
sudo openssl x509 -enddate -noout -in ./build/certbot/letsencrypt/live/$PUBLIC_DOMAIN/fullchain.pem
