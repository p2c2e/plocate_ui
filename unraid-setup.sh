#!/bin/bash

# Plocate UI - Unraid Setup Script
# This script helps you set up Plocate UI on your Unraid server

set -e

echo "=================================="
echo "  Plocate UI - Unraid Setup"
echo "=================================="
echo ""

# Check if running on Unraid
if [ ! -d "/mnt/user" ]; then
    echo "Warning: /mnt/user not found. Are you running on Unraid?"
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Create necessary directories
echo "Creating directories..."
mkdir -p /mnt/cache/appdata/plocate-ui/db
mkdir -p /mnt/cache/appdata/plocate-ui

# Copy config if it doesn't exist
if [ ! -f "config.yml" ]; then
    echo "Creating default config.yml..."
    cp config.example.yml config.yml
fi

# Prompt for configuration
echo ""
echo "Configuration:"
echo "-------------"
echo "Current config.yml will be used. Edit it to:"
echo "  1. Specify paths to index (e.g., /mnt/user/media)"
echo "  2. Set indexing schedule"
echo "  3. Configure database location"
echo ""
read -p "Would you like to edit config.yml now? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    ${EDITOR:-nano} config.yml
fi

# Build Docker image
echo ""
echo "Building Docker image..."
docker build -t plocate-ui:latest .

# Stop existing container if running
if [ "$(docker ps -q -f name=plocate-ui)" ]; then
    echo "Stopping existing container..."
    docker stop plocate-ui
    docker rm plocate-ui
fi

# Run container
echo ""
echo "Starting container..."
docker-compose up -d

# Wait for container to start
echo "Waiting for container to start..."
sleep 3

# Check if container is running
if [ "$(docker ps -q -f name=plocate-ui)" ]; then
    echo ""
    echo "=================================="
    echo "  Setup Complete!"
    echo "=================================="
    echo ""
    echo "Plocate UI is now running!"
    echo ""
    echo "Access it at: http://$(hostname -I | awk '{print $1}'):8080"
    echo ""
    echo "Useful commands:"
    echo "  - View logs: docker logs -f plocate-ui"
    echo "  - Stop: docker-compose down"
    echo "  - Restart: docker-compose restart"
    echo ""
else
    echo ""
    echo "Error: Container failed to start"
    echo "Check logs with: docker logs plocate-ui"
    exit 1
fi
