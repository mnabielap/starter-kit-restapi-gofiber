#!/bin/sh

# Stop execution if any command fails
set -e

echo "Starting GoFiber Application..."

# Execute the command passed to the container (CMD)
exec "$@"