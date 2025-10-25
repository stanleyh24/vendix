# Deployment Guide

This guide covers deploying the Vendix application to production.

## Prerequisites

- Docker and Docker Compose installed
- PostgreSQL 15+ database
- Redis 7+ instance
- S3-compatible storage (AWS S3, MinIO, etc.)
- Domain name configured
- SSL certificates (Let's Encrypt recommended)

## Environment Configuration

### Production Environment Variables

Create a `.env` file with production values:

```env
# Application
APP_ENV=production
APP_PORT=8080
APP_NAME=Vendix
APP_URL=https://billing.yourcompany.com

# Database
DB_HOST=your-db-host.rds.amazonaws.com
DB_PORT=5432
DB_USER=vendix
DB_PASSWORD=your-secure-password
DB_NAME=vendix
DB_SSL_MODE=require
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5

# JWT - Use strong random secrets
JWT_SECRET=generate-a-strong-random-secret-key
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# Redis
REDIS_ADDR=your-redis-host:6379
REDIS_PASSWORD=your-redis-password
REDIS_DB=0

# Stripe
STRIPE_SECRET_KEY=sk_live_your_stripe_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret

# Email (SMTP)
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASSWORD=your-sendgrid-api-key
SMTP_FROM=noreply@yourcompany.com

# S3 Storage
S3_ENDPOINT=s3.amazonaws.com
S3_ACCESS_KEY=your-aws-access-key
S3_SECRET_KEY=your-aws-secret-key
S3_BUCKET=billing-documents
S3_REGION=us-east-1
S3_USE_SSL=true

# DGII (Configure when ready for production)
DGII_ENABLED=false
DGII_API_URL=https://api.dgii.gov.do
DGII_CERT_PATH=/certs/dgii.pem
DGII_KEY_PATH=/certs/dgii.key

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_MAX_REQUESTS=100
RATE_LIMIT_WINDOW=1m

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Metrics
METRICS_ENABLED=true
METRICS_PORT=9090
```

## Docker Deployment

### 1. Build Production Images

```bash
# Build API image
docker build -t your-registry/vendix-api:latest -f Dockerfile .

# Build frontend image
docker build -t your-registry/vendix-frontend:latest -f frontend/Dockerfile ./frontend

# Push to registry
docker push your-registry/vendix-api:latest
docker push your-registry/vendix-frontend:latest
```

### 2. Production Docker Compose

Create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  api:
    image: your-registry/vendix-api:latest
    restart: unless-stopped
    ports:
      - "8080:8080"
    env_file:
      - .env
    depends_on:
      - redis
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  worker:
    image: your-registry/vendix-api:latest
    restart: unless-stopped
    command: ["./worker"]
    env_file:
      - .env
    depends_on:
      - redis

  frontend:
    image: your-registry/vendix-frontend:latest
    restart: unless-stopped
    ports:
      - "3000:3000"
    environment:
      VITE_API_URL: https://api.yourcompany.com

  redis:
    image: redis:7-alpine
    restart: unless-stopped
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data

volumes:
  redis_data:
```

### 3. Deploy

```bash
docker-compose -f docker-compose.prod.yml up -d
```

## Kubernetes Deployment

### 1. Create Namespace

```bash
kubectl create namespace vendix
```

### 2. Create Secrets

```bash
kubectl create secret generic vendix-secrets \
  --from-env-file=.env \
  --namespace=vendix
```

### 3. Apply Kubernetes Manifests

Create `k8s/deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vendix-api
  namespace: vendix
spec:
  replicas: 3
  selector:
    matchLabels:
      app: vendix-api
  template:
    metadata:
      labels:
        app: vendix-api
    spec:
      containers:
      - name: api
        image: your-registry/vendix-api:latest
        ports:
        - containerPort: 8080
        envFrom:
        - secretRef:
            name: vendix-secrets
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: vendix-api
  namespace: vendix
spec:
  selector:
    app: vendix-api
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

Apply:

```bash
kubectl apply -f k8s/deployment.yaml
```

## Database Setup

### 1. Create Production Database

```sql
CREATE DATABASE vendix;
CREATE USER vendix WITH ENCRYPTED PASSWORD 'your-secure-password';
GRANT ALL PRIVILEGES ON DATABASE vendix TO vendix;
```

### 2. Run Migrations

Migrations run automatically on API startup. To run manually:

```bash
# Connect to a running API container
docker exec -it vendix-api-1 ./api
```

## SSL/TLS Configuration

### Using Let's Encrypt with Nginx

Create `nginx.conf`:

```nginx
upstream api {
    server api:8080;
}

server {
    listen 80;
    server_name billing.yourcompany.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name billing.yourcompany.com;

    ssl_certificate /etc/letsencrypt/live/billing.yourcompany.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/billing.yourcompany.com/privkey.pem;

    location / {
        proxy_pass http://api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Monitoring & Logging

### 1. Prometheus Metrics

Metrics are exposed on port 9090:

```yaml
scrape_configs:
  - job_name: 'vendix'
    static_configs:
      - targets: ['api:9090']
```

### 2. Structured Logging

Logs are output in JSON format for easy parsing:

```bash
# View logs
docker logs -f vendix-api-1

# With jq for pretty printing
docker logs -f vendix-api-1 | jq
```

## Backup Strategy

### Database Backups

```bash
# Daily automated backup
0 2 * * * pg_dump -h $DB_HOST -U $DB_USER vendix | gzip > /backups/vendix_$(date +%Y%m%d).sql.gz
```

### S3 Document Backups

Enable versioning on your S3 bucket for document protection.

## Health Checks

### API Health Endpoint

```bash
curl https://billing.yourcompany.com/health
```

Expected response:
```json
{
  "status": "ok",
  "app": "Vendix"
}
```

## Scaling Considerations

### Horizontal Scaling

The application is stateless and can be scaled horizontally:

```bash
# Scale API to 5 replicas
kubectl scale deployment/vendix-api --replicas=5 -n vendix
```

### Database Connection Pooling

Configure appropriate connection pool sizes in your `.env`:

```env
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
```

## Security Checklist

- [ ] Change all default passwords
- [ ] Use strong JWT secret
- [ ] Enable SSL/TLS
- [ ] Configure firewall rules
- [ ] Enable rate limiting
- [ ] Set up monitoring and alerts
- [ ] Regular security updates
- [ ] Database encryption at rest
- [ ] Regular backups
- [ ] Audit logging enabled

## Troubleshooting

### Database Connection Issues

```bash
# Test database connection
docker exec -it vendix-api-1 psql -h $DB_HOST -U $DB_USER -d vendix
```

### Redis Connection Issues

```bash
# Test Redis connection
docker exec -it vendix-redis-1 redis-cli ping
```

### View Application Logs

```bash
# API logs
docker logs -f vendix-api-1

# Worker logs
docker logs -f vendix-worker-1
```

## Support

For production issues, contact: support@yourcompany.com

