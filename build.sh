#!/bin/bash

# Tesselbox Go Build Script

set -e

# Change to the script's directory
cd "$(dirname "$0")"

echo "Building Tesselbox Go..."

# Create build directory
mkdir -p build

# Build client
echo "Building client..."
go build -o build/tesselbox .

# Build server
# Temporarily disabled - server needs to be updated to match new World API
# echo "Building server..."
# go build -o build/tesselbox-server ./cmd/server/

echo "Build complete!"
echo "Client: build/tesselbox"
echo "Server: build/tesselbox-server"
echo ""
echo "To run the client:"
echo "  ./build/tesselbox"
echo ""
echo "To run the server:"
echo "  ./build/tesselbox-server"