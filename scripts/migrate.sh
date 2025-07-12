#!/bin/bash
# Database migration script
# Runs database migrations for the current environment

set -e

echo "Running database migrations..."

cd backend
go run cmd/server/main.go migrate

echo "Migrations completed successfully!"
