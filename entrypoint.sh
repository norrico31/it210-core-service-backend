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
until nc -z postgres 5432; do
  echo "Waiting for PostgreSQL..."
  sleep 2
done

echo "PostgreSQL is up!"

DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Run migrations
if [ -x /app/auth-service/migrate ]; then
  echo "ENTRYPOINT: Running migrations..."
  /app/auth-service/migrate -path=/app/auth-service/cmd/migrate/migrations -database "${DATABASE_URL}" up || { echo "Migration failed"; exit 1; }
else
  echo "Migration binary not found."
fi

# Run seeding
if [ -x /app/auth-service/seed ]; then
  echo "ENTRYPOINT: Running seeding..."
  /app/auth-service/seed || { echo "Seeding failed"; }
else
  echo "Seeder binary not found."
fi

# Start the auth service
echo "ENTRYPOINT: Starting auth service..."
exec /app/auth-service/auth-service || { echo "Auth service failed to start"; }
