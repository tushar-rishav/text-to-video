# Text-to-Video API

A scalable, event-driven text-to-video generation service built with Go, Python, and Svelte, deployed on Kubernetes with GPU support.

## Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Frontend  │    │   Go API    │    │   Python    │
│   (Svelte)  │◄──►│   Backend   │◄──►│   Video     │
│             │    │             │    │   Service   │
└─────────────┘    └─────────────┘    └─────────────┘
                           │                   │
                    ┌─────────────┐    ┌─────────────┐
                    │    Redis    │    │   MySQL     │
                    │   (Pub/Sub) │    │  Database   │
                    └─────────────┘    └─────────────┘
```

## Services

### 1. Go Backend API
- **Port**: 8080
- **Features**: 
  - REST API endpoints for job management
  - MySQL integration for job persistence
  - Redis pub/sub for event communication
  - File serving for generated videos

### 2. Python Video Service
- **Port**: 8000
- **Features**:
  - Video generation using genmo mochi-1 model
  - GPU-accelerated processing
  - Redis pub/sub for job consumption
  - Concurrent video processing

### 3. Svelte Frontend
- **Port**: 3000
- **Features**:
  - Text prompt submission
  - Job status monitoring
  - Video download interface
  - Real-time updates

## API Endpoints

### Submit Job
```http
POST /submit
Content-Type: application/json

{
  "prompt": "A cat playing in the garden"
}
```

### Get Job Status
```http
GET /status?job_id=123
```

### List Jobs
```http
GET /list?status=completed&offset=0&limit=10
```

### Get Video
```http
GET /video?job_id=123
```

## Prerequisites

- Kubernetes cluster with GPU nodes
- NVIDIA Container Runtime
- MySQL database
- Redis instance
- Docker

## Quick Start

### Prerequisites
- Docker installed and running
- Kubernetes cluster with GPU nodes
- kubectl configured
- Docker Hub account (https://hub.docker.com/repositories/tusharrishav)

### Setup Steps

1. **Clone the repository**
2. **Set up Docker Hub authentication** (choose one option):
   
   **Option A: Interactive setup (prompts for password)**
   ```bash
   chmod +x *.sh
   ./setup-docker-hub.sh
   ```
   
   **Option B: Secure setup (uses existing Docker config)**
   ```bash
   chmod +x *.sh
   ./setup-docker-hub-secure.sh
   ```
3. **Test Docker Hub setup**
   ```bash
   ./test-docker-hub.sh
   ```
4. **Deploy to Kubernetes**
   ```bash
   ./deploy.sh
   ```
5. **Access the frontend**
   ```bash
   kubectl port-forward -n text-to-video svc/frontend-service 3000:3000
   ```

See deployment section for detailed instructions.

## Development

### Local Development
```bash
# Start dependencies
docker-compose up -d mysql redis

# Run Go backend
cd backend && go run main.go

# Run Python service
cd video-service && python app.py

# Run Svelte frontend
cd frontend && npm run dev
```

### Building Images
```bash
# Build all images
docker build -t text-to-video-api:latest ./backend
docker build -t text-to-video-service:latest ./video-service
docker build -t text-to-video-frontend:latest ./frontend
```

## Deployment

### Docker Hub Configuration

The application is configured to use Docker Hub as the container registry:

- **Registry**: `docker.io/tusharrishav`
- **Images**:
  - `tusharrishav/text-to-video-api:latest`
  - `tusharrishav/text-to-video-service:latest`
  - `tusharrishav/text-to-video-frontend:latest`

### Authentication Setup

1. **Login to Docker Hub**:
   ```bash
   docker login
   # Username: tusharrishav
   # Password: Your Docker Hub password or access token
   ```

2. **Create Kubernetes secret** (automated by setup script):
   ```bash
   kubectl create secret docker-registry docker-registry-secret \
     --docker-server=https://index.docker.io/v1/ \
     --docker-username=tusharrishav \
     --docker-password=<your-password> \
     --namespace=text-to-video
   ```

### Security Note

For production deployments, it's recommended to:
- Use Docker Hub access tokens instead of passwords
- Store credentials in Kubernetes secrets
- Use private repositories for sensitive applications

See `k8s/` directory for Kubernetes manifests.

## License

MIT 