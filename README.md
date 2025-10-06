# Microservices Project
Deployment: https://cards-webapp-220185652609.europe-west1.run.app/

-This project consists of 5 microservices: 4 Go-based backend services and 1 Flutter web application. All services are containerized with Docker and ready for deployment.

## Services Description

### 1. Cards Service (Go)
- **Port**: 8080 (default)
- **Purpose**: Main API service handling card operations
- **Dependencies**: PostgreSQL, Redis
- **Endpoints**:
  - `POST /v1/register` - User registration
  - `POST /v1/issue` - Card issuance
  - `POST /v1/webhook` - Webhook handling
  - `GET /v1/:citizen_id/cards` - Get user cards
  - `GET /health` - Health check

### 2. Issuer Service (Go)
- **Port**: 8080 (default)
- **Purpose**: Handles card issuance logic
- **Dependencies**: Webhook URL for notifications
- **Endpoints**:
  - `POST /v1/cards` - Issue new card
  - `GET /health` - Health check

### 3. Notifications Service (Go)
- **Port**: 8080 (default)
- **Purpose**: Real-time notifications using Server-Sent Events (SSE)
- **Endpoints**:
  - `GET /notifications/stream` - SSE connection for real-time notifications
  - `POST /notify` - Send notification
  - `GET /health` - Health check

### 4. Webhook Service (Go)
- **Port**: 8080 (default)
- **Purpose**: Event forwarding and webhook management
- **Dependencies**: Redis
- **Endpoints**:
  - `POST /suscribe` - Subscribe to events
  - `POST /request` - Forward requests
  - `POST /response` - Forward responses

### 5. Webapp (Flutter)
- **Port**: 80 (nginx)
- **Purpose**: Frontend web application
- **Dependencies**: Cards Service, Notifications Service URLs

## Prerequisites

- Docker and Docker Compose
- Docker Registry access (for pushing images)
- PostgreSQL database
- Redis instance

## Environment Variables

### Cards Service
Create a `.env` file in the `cards/` directory with:
```env
PORT=8080
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=your_redis_password
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
```

### Issuer Service
Create a `.env` file in the `issuer/` directory with:
```env
PORT=8080
WEBHOOK_URL=https://your-webhook-service-url.com
```

### Notifications Service
Create a `.env` file in the `notifications/` directory with:
```env
PORT=8080
```

### Webhook Service
Create a `.env` file in the `webhook/` directory with:
```env
PORT=8080
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=your_redis_password
```

## Deployment Instructions

### 1. Build and Push Docker Images

For each Go service (cards, issuer, notifications, webhook):

```bash
# Navigate to service directory
cd cards  # or issuer, notifications, webhook

# Build the Docker image
docker build -t your-registry/cards-service:latest .

# Push to registry
docker push your-registry/cards-service:latest
```

### 2. Build and Push Flutter Webapp

```bash
# Navigate to webapp directory
cd webapp

# Build the Docker image
docker build -t your-registry/webapp:latest .

# Push to registry
docker push your-registry/webapp:latest
```

### 3. Configure Flutter App URLs

Before building the Flutter app, update the service URLs in `webapp/lib/env/env.dart`:

```dart
class Env {
  static String get issueServiceUrl {
    final url = 'https://your-issuer-service-url.com';
    print('Issue Service URL: $url');
    return url;
  }
  
  static String get notificationServiceUrl {
    final url = 'https://your-notifications-service-url.com';
    print('Notification Service URL: $url');
    return url;
  }
}
```

### 4. Deploy Services

You can deploy the services using any container orchestration platform:

#### Using Docker Compose (Local Development)
```yaml
version: '3.8'
services:
  cards-service:
    image: your-registry/cards-service:latest
    ports:
      - "8080:8080"
    environment:
      - REDIS_ADDR=redis:6379
      - DB_HOST=postgres:5432
    depends_on:
      - redis
      - postgres

  issuer-service:
    image: your-registry/issuer-service:latest
    ports:
      - "8081:8080"
    environment:
      - WEBHOOK_URL=http://webhook-service:8080

  notifications-service:
    image: your-registry/notifications-service:latest
    ports:
      - "8082:8080"

  webhook-service:
    image: your-registry/webhook-service:latest
    ports:
      - "8083:8080"
    environment:
      - REDIS_ADDR=redis:6379
    depends_on:
      - redis

  webapp:
    image: your-registry/webapp:latest
    ports:
      - "80:80"

  redis:
    image: redis:alpine
    ports:
      - "6379:6379"

  postgres:
    image: postgres:13
    environment:
      - POSTGRES_DB=your_db_name
      - POSTGRES_USER=your_db_user
      - POSTGRES_PASSWORD=your_db_password
    ports:
      - "5432:5432"
```

#### Using Kubernetes
Create deployment manifests for each service with appropriate environment variables and service configurations.

#### Using Cloud Platforms
- **AWS**: Use ECS, EKS, or App Runner
- **Google Cloud**: Use Cloud Run or GKE
- **Azure**: Use Container Instances or AKS

## Service Communication

1. **Flutter Webapp** → **Cards Service**: API calls for card operations
2. **Flutter Webapp** → **Notifications Service**: SSE connection for real-time updates
3. **Cards Service** → **Issuer Service**: Card issuance requests
4. **Cards Service** → **Webhook Service**: Event forwarding
5. **Issuer Service** → **Notifications Service**: Send notifications via webhook

## Health Checks

All services include health check endpoints at `/health`:
- Cards Service: `GET /health`
- Issuer Service: `GET /health`
- Notifications Service: `GET /health`
- Webhook Service: No explicit health endpoint (check port connectivity)

## Monitoring and Logs

- All services log to stdout/stderr for container log collection
- Use your container orchestration platform's logging solution
- Consider implementing structured logging with JSON format

## Security Considerations

- Configure CORS appropriately for production
- Use HTTPS for all service communication
- Implement proper authentication and authorization
- Secure database and Redis connections
- Use secrets management for sensitive environment variables

## Scaling

- **Cards Service**: Can be scaled horizontally (stateless)
- **Notifications Service**: Requires sticky sessions for SSE connections
- **Webhook Service**: Can be scaled horizontally (stateless)
- **Issuer Service**: Can be scaled horizontally (stateless)
- **Webapp**: Can be scaled horizontally (static files)

## Troubleshooting

1. **Service not starting**: Check environment variables and dependencies
2. **Database connection issues**: Verify PostgreSQL connection and credentials
3. **Redis connection issues**: Verify Redis service and credentials
4. **Flutter app not connecting**: Check service URLs in `env.dart`
5. **SSE not working**: Verify user tokens and connection management

## Development

For local development, you can run services individually:

```bash
# Run a Go service locally
cd cards
go mod download
go run main.go

# Run Flutter app locally
cd webapp
flutter pub get
flutter run -d web-server --web-port 8080
```

## Contributing

1. Make changes to the respective service
2. Test locally with Docker
3. Update environment variables as needed
4. Rebuild and push Docker images
5. Update deployment configurations
