#!/bin/bash

# Test Docker Hub Setup

echo "🧪 Testing Docker Hub Setup"
echo "==========================="

DOCKER_USERNAME="tusharrishav"

# Check Docker login status
echo "1. Checking Docker login status..."
if docker info | grep -q "Username"; then
    echo "✅ Logged in as: $(docker info | grep 'Username' | cut -d' ' -f2)"
else
    echo "❌ Not logged in to Docker Hub"
    echo "Run: ./docker-login.sh"
    exit 1
fi

# Test Docker Hub connectivity
echo ""
echo "2. Testing Docker Hub connectivity..."
if docker pull hello-world:latest > /dev/null 2>&1; then
    echo "✅ Successfully pulled test image from Docker Hub"
else
    echo "❌ Failed to pull from Docker Hub"
    exit 1
fi

# Test building and pushing a simple image
echo ""
echo "3. Testing image build and push..."
echo "Building test image..."

# Create a simple test Dockerfile
cat > test.Dockerfile << EOF
FROM alpine:latest
RUN echo "Hello from Docker Hub test" > /test.txt
CMD ["cat", "/test.txt"]
EOF

# Build test image
docker build -f test.Dockerfile -t $DOCKER_USERNAME/test-image:latest .

if [ $? -eq 0 ]; then
    echo "✅ Test image built successfully"
    
    # Push test image
    echo "Pushing test image to Docker Hub..."
    docker push $DOCKER_USERNAME/test-image:latest
    
    if [ $? -eq 0 ]; then
        echo "✅ Test image pushed successfully"
        echo "You can view it at: https://hub.docker.com/r/$DOCKER_USERNAME/test-image"
    else
        echo "❌ Failed to push test image"
        exit 1
    fi
else
    echo "❌ Failed to build test image"
    exit 1
fi

# Cleanup
echo ""
echo "4. Cleaning up test files..."
rm -f test.Dockerfile
docker rmi $DOCKER_USERNAME/test-image:latest hello-world:latest > /dev/null 2>&1

echo ""
echo "🎉 Docker Hub setup test completed successfully!"
echo "You're ready to deploy the text-to-video application." 