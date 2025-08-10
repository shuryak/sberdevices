#!/usr/bin/env bash

EXCLUDE_VARS=("OS" "ARCH" "TAG" "DOCKERFILE" "TAR_NAME")

is_excluded() {
    for var in "${EXCLUDE_VARS[@]}"; do
        [[ "$1" == "$var" ]] && return 0
    done
    return 1
}

BUILD_ARGS=()
while IFS='=' read -r key value; do
    if ! is_excluded "$key"; then
        BUILD_ARGS+=(--build-arg "$key=$value")
    fi
done < <(env)

docker buildx build --platform "$OS"/"$ARCH" -t "$TAG":latest -f "$DOCKERFILE" "${BUILD_ARGS[@]}" --load ./
docker save -o ./.tmp/"$TAR_NAME" "$TAG":latest
docker rmi -f "$TAG":latest

##!/usr/bin/env bash
#
#docker buildx build --platform "$OS"/"$ARCH" -t "$TAG":latest -f "$DOCKERFILE" --load ./
#docker save -o ./.tmp/"$TAR_NAME" "$TAG":latest
