#!/bin/sh
set -e

echo "ENTRYPOINT: Starting entrypoint..."

echo "Environment Variables:"
echo "PGHOST=${PGHOST}"
echo "PGPORT=${PGPORT}"
echo "PGUSER=${PGUSER}"
echo "PGPASSWORD=${PGPASSWORD}"
echo "POSTGRES_DB=${POSTGRES_DB}"
echo "GATEWAY_SERVICE_PORT=${GATEWAY_SERVICE_PORT}"
echo "DATABASE_PUBLIC_URL=${DATABASE_PUBLIC_URL}"

DATABASE_URL="${DATABASE_PUBLIC_URL}"

echo "FUCK: $DATABASE_URL"

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL at ${PGHOST}:${PGPORT}..."
until nc -z "${PGHOST}" "${PGPORT}"; do
  echo "PostgreSQL is unavailable - retrying in 2 seconds"
  sleep 2
done
echo "PostgreSQL is ready!"

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
exec /app/core-service/core-service || { echo "Core service failed to start"; exit 1; }
