#!/usr/bin/env bash

USE_NGINX_INIT_CONF=true gomplate -f ./build/nginx/nginx.template.conf -o ./build/nginx/nginx.conf
gomplate -f ./build/nginx/conf.d/default.template.conf -o ./build/nginx/conf.d/default.conf
docker compose -f ./build/docker-compose.yaml up --build --remove-orphans --force-recreate --detach certbot
gomplate -f ./build/nginx/nginx.template.conf -o ./build/nginx/nginx.conf
