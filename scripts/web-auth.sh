#!/usr/bin/env bash

docker compose -f ./build/docker-compose.yaml up --build --remove-orphans --force-recreate --detach web-auth
#docker compose -f ./build/docker-compose.yaml up --remove-orphans --detach web-auth
