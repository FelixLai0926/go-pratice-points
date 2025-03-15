#!/bin/sh
APP_ENV=${APP_ENV:-example}
CONFIG_FILE="configs/${APP_ENV}.yaml"

# 拼接 URL
export DATABASE_URL=$(yq eval '.postgres | "postgres://" + .POSTGRES_USER + ":" + .POSTGRES_PASSWORD + "@" + .POSTGRES_HOST + ":" + (.POSTGRES_PORT|tostring) + "/" + .POSTGRES_DB + "?sslmode=disable"' "$CONFIG_FILE")
set -e


echo "Running database migrations..."
if migrate -path=/root/migrations -database "$DATABASE_URL" up; then
    echo "Migrations complete. Exiting..."
    exit 0
else
    echo "Migration failed! Exiting..."
    exit 1
fi