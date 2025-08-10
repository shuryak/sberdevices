#!/usr/bin/env bash

docker compose -f ./build/docker-compose.yaml up --build --remove-orphans --force-recreate backend-delve
