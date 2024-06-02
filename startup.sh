#!/bin/sh
set -e

# Ensure Go binaries are in PATH
export PATH=$PATH:/go/bin

# Run go mod tidy
make tidy

# Run database migrations
make migrate-up

# Seed the database
make seed

# Start the Go application
./main
