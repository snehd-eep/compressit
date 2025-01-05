#!/bin/bash

# Update and install Ghostscript
echo "Installing Ghostscript..."
apt-get update && apt-get install -y ghostscript

# Install Go dependencies
echo "Installing Go dependencies..."
go mod tidy

# Compile the Go application
echo "Building the Go application..."
go build -o compressit .

# You can also run tests before building
# echo "Running tests..."
# go test ./...

# Any other build-related tasks can go here
echo "Build complete!"
