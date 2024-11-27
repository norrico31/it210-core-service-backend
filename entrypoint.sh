#!/bin/sh
set -e

echo "ENTRYPOINT: Starting entrypoint..."

echo "Environment Variables:"
echo "DATABASE_PUBLIC_URL=${DATABASE_PUBLIC_URL}"

# Extract host and port for health check
DB_HOST=$(echo "${DATABASE_PUBLIC_URL}" | sed -n 's/.*@\(.*\):\([0-9]*\).*/\1/p')
DB_PORT=$(echo "${DATABASE_PUBLIC_URL}" | sed -n 's/.*:\([0-9]*\)\/.*/\1/p')

echo "Waiting for PostgreSQL at ${DB_HOST}:${DB_PORT}..."
RETRY_COUNT=0
MAX_RETRIES=15
until nc -z "${DB_HOST}" "${DB_PORT}"; do
  RETRY_COUNT=$((RETRY_COUNT + 1))
  if [ "${RETRY_COUNT}" -ge "${MAX_RETRIES}" ]; then
    echo "PostgreSQL is still unavailable after ${MAX_RETRIES} retries. Exiting."
    exit 1
  fi
  echo "PostgreSQL is unavailable - retrying in 2 seconds"
  sleep 2
done
echo "PostgreSQL is ready!"

# Run migrations
if [ -x /app/core-service/migrate ]; then
  echo "ENTRYPOINT: Running migrations..."
  if /app/core-service/migrate -path=/app/core-service/cmd/migrate/migrations -database "${DATABASE_PUBLIC_URL}"; then
    echo "Migrations applied successfully."
  else
    echo "Migration failed. Exiting."
    exit 1
  fi
else
  echo "Migration executable not found. Skipping migrations."
fi

# Run seeding
if [ -x /app/core-service/seed ]; then
  echo "ENTRYPOINT: Running seeding..."
  if /app/core-service/seed; then
    echo "Seeding completed successfully."
  else
    echo "Seeding failed. Exiting."
    exit 1
  fi
else
  echo "Seeder executable not found. Skipping seeding."
fi

# Start the service
echo "ENTRYPOINT: Starting core service..."
exec /app/core-service/core-service || { echo "Core service failed to start"; exit 1; }
