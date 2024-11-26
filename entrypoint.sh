#!/bin/sh
set -e

echo "ENTRYPOINT: Starting entrypoint..."

echo "Environment Variables:"
echo "JWT_SECRET=${JWT_SECRET}"
echo "DB_HOST=${DB_HOST}"
echo "PORT=${PORT}"
echo "DB_PORT=${DB_PORT}"
echo "DB_USER=${DB_USER}"
echo "DB_PASSWORD=${DB_PASSWORD}"
echo "DB_NAME=${DB_NAME}"
echo "DB_ADDRESS=${DB_ADDRESS}"
echo "GATEWAY_SERVICE_PORT=${GATEWAY_SERVICE_PORT}"

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL at ${DB_HOST}:${DB_PORT}..."
until nc -z "${DB_HOST}" "${DB_PORT}"; do
  echo "PostgreSQL is unavailable - retrying in 2 seconds"
  sleep 2
done
echo "PostgreSQL is ready!"


echo "PostgreSQL is up!"

DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}"

# Run migrations
if [ -x /app/core-service/migrate ]; then
  echo "ENTRYPOINT: Running migrations..."
  /app/core-service/migrate -path=/app/core-service/cmd/migrate/migrations -database "${DATABASE_URL}" up || { echo "Migration failed"; exit 1; }
else
  echo "Migration binary not found."
fi

# Run seeding
if [ -x /app/core-service/seed ]; then
  echo "ENTRYPOINT: Running seeding..."
  /app/core-service/seed || { echo "Seeding failed"; }
else
  echo "Seeder binary not found."
fi

# Start the core service
echo "ENTRYPOINT: Starting core service..."
exec /app/core-service/core-service || { echo "Core service failed to start"; }
