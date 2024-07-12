#!/bin/sh
set -e

# Ensure Go binaries are in PATH
export PATH=$PATH:/go/bin

make tidy

make migrate-up

make seed

make swag

# Start the Go application
./main
