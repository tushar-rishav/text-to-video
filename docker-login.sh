#!/bin/bash

# Docker Hub Authentication Script

echo "🔐 Docker Hub Authentication Setup"
echo "=================================="

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if already logged in
if docker info | grep -q "Username"; then
    echo "✅ Already logged in to Docker Hub"
    docker info | grep "Username"
else
    echo "📝 Please login to Docker Hub..."
    echo "You will be prompted for your Docker Hub username and password/token"
    echo ""
    echo "Note: For security, it's recommended to use a Docker Hub access token instead of your password"
    echo "You can create one at: https://hub.docker.com/settings/security"
    echo ""
    
    docker login
    
    if [ $? -eq 0 ]; then
        echo "✅ Successfully logged in to Docker Hub"
    else
        echo "❌ Failed to login to Docker Hub"
        exit 1
    fi
fi

echo ""
echo "📋 Docker Hub Configuration:"
echo "Username: tusharrishav"
echo "Registry: docker.io/tusharrishav"
echo ""
echo "Your images will be pushed to:"
echo "- docker.io/tusharrishav/text-to-video-api:latest"
echo "- docker.io/tusharrishav/text-to-video-service:latest"
echo "- docker.io/tusharrishav/text-to-video-frontend:latest" 