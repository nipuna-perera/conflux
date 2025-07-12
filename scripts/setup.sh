#!/bin/bash
# Development environment setup script
# Installs dependencies and prepares the development environment
# Sets up database connections and runs initial migrations

set -e

echo "Setting up full-stack development environment..."

# Check for required tools
command -v go >/dev/null 2>&1 || { echo "Go is required but not installed. Aborting." >&2; exit 1; }
command -v node >/dev/null 2>&1 || { echo "Node.js is required but not installed. Aborting." >&2; exit 1; }
command -v docker >/dev/null 2>&1 || { echo "Docker is required but not installed. Aborting." >&2; exit 1; }

# Setup backend
echo "Setting up Go backend..."
cd backend
go mod download
go mod tidy

# Setup frontend
echo "Setting up SvelteKit frontend..."
cd ../frontend
npm install

# Copy environment files
echo "Setting up environment configuration..."
cd ..
cp .env.example .env

echo "Setup complete! Run 'make dev' to start the development environment."
