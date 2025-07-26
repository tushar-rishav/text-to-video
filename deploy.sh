#!/bin/bash

# Text-to-Video API Deployment Script

set -e

# Configuration
REGISTRY="tusharrishav"  # Your Docker Hub username
NAMESPACE="text-to-video"

echo "🚀 Starting Text-to-Video API deployment..."

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed. Please install kubectl first."
    exit 1
fi

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

# Create namespace
echo "🏗️  Creating namespace..."
kubectl apply -f k8s/namespace.yaml

# Apply Kubernetes manifests
echo "📋 Applying Kubernetes manifests..."

# Apply manifests in order
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/mysql.yaml
kubectl apply -f k8s/redis.yaml

# Wait for database and Redis to be ready
echo "⏳ Waiting for database and Redis to be ready..."
kubectl wait --for=condition=ready pod -l app=mysql -n $NAMESPACE --timeout=300s
kubectl wait --for=condition=ready pod -l app=redis -n $NAMESPACE --timeout=300s

# Apply application manifests
kubectl apply -f k8s/backend.yaml
kubectl apply -f k8s/video-service.yaml
kubectl apply -f k8s/frontend.yaml

# Wait for all pods to be ready
echo "⏳ Waiting for all pods to be ready..."
kubectl wait --for=condition=ready pod -l app=backend -n $NAMESPACE --timeout=300s
kubectl wait --for=condition=ready pod -l app=video-service -n $NAMESPACE --timeout=300s
kubectl wait --for=condition=ready pod -l app=frontend -n $NAMESPACE --timeout=300s

echo "✅ Deployment completed successfully!"

# Show service information
echo "📊 Service Information:"
echo "Frontend: http://localhost:3000 (if port-forwarded)"
echo "Backend API: http://localhost:8080 (if port-forwarded)"
echo "Video Service: http://localhost:8000 (if port-forwarded)"

# Show pod status
echo "📋 Pod Status:"
kubectl get pods -n $NAMESPACE

# Show services
echo "🌐 Services:"
kubectl get services -n $NAMESPACE

echo ""
echo "🎉 Text-to-Video API is now deployed!"
echo "To access the frontend, you can:"
echo "1. Port forward: kubectl port-forward -n $NAMESPACE svc/frontend-service 3000:3000"
echo "2. Or use an ingress controller if configured"
echo ""
echo "To check logs:"
echo "kubectl logs -n $NAMESPACE -l app=backend"
echo "kubectl logs -n $NAMESPACE -l app=video-service"
echo "kubectl logs -n $NAMESPACE -l app=frontend" 