#!/bin/sh

set -e

echo "Starting migrations"

until pg_isready -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER"; do
  echo "Waiting for postgres..."
  sleep 2
done

echo $(env)
echo $(ls /migrations)

# Execute migrations in order
find /migrations -maxdepth 1 -type f -name '*.up.sql' | sort -V | while read -r file; do
  echo "Applying migration: $file"
  PGPASSWORD="$POSTGRES_PASSWORD" psql \
    -h "$POSTGRES_HOST" \
    -p "$POSTGRES_PORT" \
    -U "$POSTGRES_USER" \
    -d "$POSTGRES_DB" \
    -f "$file"
done

echo "Migrations finished"
