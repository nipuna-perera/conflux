#!/bin/bash
# Development startup script
# Starts both backend and frontend in development mode

set -e

echo "Starting development environment..."

# Function to kill background processes on exit
cleanup() {
    echo "Stopping development servers..."
    kill $(jobs -p) 2>/dev/null || true
}
trap cleanup EXIT

# Start backend in background
echo "Starting Go backend..."
cd backend
go run cmd/server/main.go &

# Start frontend in background
echo "Starting SvelteKit frontend..."
cd ../frontend
npm run dev &

# Wait for processes
wait
