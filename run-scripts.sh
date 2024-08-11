#!/bin/bash

# Get the container ID or name of the PostgreSQL container
CONTAINER_NAME=$(docker-compose ps -q db)

# Wait for the PostgreSQL container to be ready
until docker exec -it $CONTAINER_NAME psql -U admin -d postgres -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"

# Run DDL and DML scripts
docker exec -it $CONTAINER_NAME psql -U admin -d postgres -f /docker-entrypoint-initdb.d/01-create-schema.sql
docker exec -it $CONTAINER_NAME psql -U admin -d postgres -f /docker-entrypoint-initdb.d/02-populate-data.sql
