#!/usr/bin/env bash

gomplate -f ./build/nginx/nginx.template.conf -o ./build/nginx/nginx.conf
#docker compose -f ./build/docker-compose.yaml up --build --remove-orphans --force-recreate --detach \
docker compose -f ./build/docker-compose.yaml up --remove-orphans --detach \
postgres nginx certbot-renew
docker compose -f ./build/docker-compose.yaml run --rm postgres-migrator
