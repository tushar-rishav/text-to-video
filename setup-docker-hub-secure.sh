#!/bin/bash

# Secure Docker Hub Setup Script for Kubernetes

echo "🔐 Secure Docker Hub Setup for Kubernetes"
echo "========================================="

# Configuration
DOCKER_USERNAME="tusharrishav"
NAMESPACE="text-to-video"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed. Please install kubectl first."
    exit 1
fi

echo "📝 Setting up Docker Hub authentication..."

# Login to Docker Hub
if ! docker info | grep -q "Username"; then
    echo "Please login to Docker Hub:"
    echo "Username: $DOCKER_USERNAME"
    echo "Password: Use your Docker Hub password or access token"
    echo ""
    echo "Note: For security, it's recommended to use a Docker Hub access token"
    echo "Create one at: https://hub.docker.com/settings/security"
    echo ""
    
    docker login
    
    if [ $? -ne 0 ]; then
        echo "❌ Failed to login to Docker Hub"
        exit 1
    fi
else
    echo "✅ Already logged in to Docker Hub"
    docker info | grep "Username"
fi

# Create namespace if it doesn't exist
echo "🏗️  Creating namespace..."
kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

# Create Docker registry secret from existing Docker config
echo "🔑 Creating Docker registry secret from existing config..."
if [ -f /home/ubuntu/.docker/config.json ]; then
    echo "Using existing Docker config file..."
    kubectl create secret generic docker-registry-secret \
        --from-file=.dockerconfigjson=/home/ubuntu/.docker/config.json \
        --type=kubernetes.io/dockerconfigjson \
        --namespace=$NAMESPACE \
        --dry-run=client -o yaml | kubectl apply -f -
    
    if [ $? -eq 0 ]; then
        echo "✅ Docker registry secret created successfully"
    else
        echo "❌ Failed to create Docker registry secret"
        exit 1
    fi
else
    echo "❌ Docker config file not found at /home/ubuntu/.docker/config.json"
    echo "Please run 'docker login' first"
    exit 1
fi

echo ""
echo "✅ Secure Docker Hub setup completed!"
echo ""
echo "📋 Summary:"
echo "- Docker Hub Username: $DOCKER_USERNAME"
echo "- Registry: docker.io/$DOCKER_USERNAME"
echo "- Kubernetes Secret: docker-registry-secret"
echo "- Namespace: $NAMESPACE"
echo "- Method: Using existing Docker config (more secure)"
echo ""
echo "🚀 You can now run: ./deploy.sh" 