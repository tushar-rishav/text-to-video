#!/bin/bash

# Text-to-Video API image builder script

set -e
REGISTRY="tusharrishav"  # Your Docker Hub username
echo "🚀 Building docker images for Text-to-Video API..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Build Docker images
echo "📦 Building Docker images..."

echo "Building backend image..."
docker build -t $REGISTRY/text-to-video-api:latest ./backend

echo "Building video service image..."
docker build -t $REGISTRY/text-to-video-service:latest ./video-service

echo "Building frontend image..."
docker build -t $REGISTRY/text-to-video-frontend:latest ./frontend

# Login to Docker Hub (if not already logged in)
echo "🔐 Checking Docker Hub login..."
if ! docker info | grep -q "Username"; then
    echo "Please login to Docker Hub:"
    docker login
fi

# Push images to registry
echo "📤 Pushing images to Docker Hub..."
docker push $REGISTRY/text-to-video-api:latest
docker push $REGISTRY/text-to-video-service:latest
docker push $REGISTRY/text-to-video-frontend:latest