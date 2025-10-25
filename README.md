# Vendix - Electronic Invoicing for Dominican Republic

A complete SaaS billing and electronic invoicing system designed specifically for the Dominican Republic, featuring DGII (Dirección General de Impuestos Internos) compliance, multi-tenant architecture, and modern cloud-native design.

## 🌟 Features

### Core Functionality
- **Multi-Tenant Architecture**: Each tenant has isolated data in separate PostgreSQL schemas
- **Electronic Invoicing**: Full DGII compliance with NCF (Número de Comprobante Fiscal) support
- **Customer Management**: Manage individual and business customers with RNC/Cédula support
- **Product/Service Catalog**: Flexible product and service management with ITBIS tax rates
- **Invoice Generation**: Create, send, and track invoices with automatic DGII submission
- **Payment Processing**: Record payments with Stripe integration and automatic reconciliation
- **Reports & Analytics**: Sales reports, tax reports (Libro de Ventas/Compras), and custom dashboards
- **Accounting Integration**: Basic journal entries and month-end closing

### Technical Features
- **JWT Authentication**: Secure authentication with access and refresh tokens
- **Role-Based Access Control (RBAC)**: Admin, Accountant, Billing Clerk, and Viewer roles
- **Plan-Based Features**: Module access control based on subscription plans
- **Background Jobs**: Asynchronous processing for invoicing, emails, and DGII submission
- **RESTful API**: Comprehensive API with Swagger documentation
- **Rate Limiting**: Per-tenant and per-user rate limiting
- **Audit Logging**: Immutable audit trails for all operations

## 🏗️ Architecture

### Tech Stack
- **Backend**: Go 1.21+ with Fiber framework
- **Database**: PostgreSQL 15+ (multi-tenant with schema-per-tenant)
- **Cache/Queue**: Redis 7+
- **Job Processing**: Asynq (Redis-based)
- **Frontend**: React 18 with Vite, TailwindCSS, and React Router
- **Storage**: S3-compatible object storage (MinIO for development)
- **Payments**: Stripe integration
- **Deployment**: Docker & Docker Compose

### Project Structure

```
vendix/
├── cmd/
│   ├── api/          # API server entry point
│   └── worker/       # Background worker entry point
├── internal/
│   ├── auth/         # Authentication (JWT, password hashing)
│   ├── config/       # Configuration management
│   ├── database/     # Database connection and migrations
│   ├── dgii/         # DGII compliance module
│   ├── jobs/         # Background job handlers
│   ├── logger/       # Structured logging
│   ├── middleware/   # HTTP middleware (auth, tenant, rate limit)
│   ├── modules/      # Business modules
│   │   ├── auth/     # User authentication endpoints
│   │   ├── tenants/  # Tenant management
│   │   ├── customers/# Customer management
│   │   ├── products/ # Product catalog
│   │   ├── invoices/ # Invoice management
│   │   ├── payments/ # Payment processing
│   │   ├── reports/  # Reports and analytics
│   │   ├── accounting/# Accounting operations
│   │   └── webhooks/ # Webhook handlers (Stripe)
│   └── server/       # HTTP server setup
├── frontend/         # React frontend application
├── scripts/          # Utility scripts
├── docs/            # Documentation and Swagger specs
├── .github/         # CI/CD workflows
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Docker and Docker Compose
- Node.js 20+ (for frontend development)
- Make (optional, for using Makefile commands)

### Development Mode (with Hot-Reload) ⚡

**Recommended for development** - changes are reflected instantly without rebuilding images.

1. **Clone the repository**
```bash
git clone <repository-url>
cd vendix
```

2. **Create environment file**
```bash
cp env.example .env
# Edit .env with your configuration
```

3. **Start all services in development mode**
```bash
make dev-up
```

4. **View logs**
```bash
make dev-logs
```

**What you get:**
- ✅ Hot-reload for Go backend (using Air)
- ✅ Hot-reload for React frontend (using Vite HMR)
- ✅ Instant code changes without rebuilding
- ✅ All services included (PostgreSQL, Redis, MinIO)

**Access the application:**
- Frontend: http://localhost:3000
- API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- MinIO Console: http://localhost:9001

**Stop development environment:**
```bash
make dev-down
```

📖 **See [DEV_SETUP.md](./DEV_SETUP.md) for detailed development guide**

### Production Mode (Optimized Build)

For production deployments or testing production builds:

1. **Clone the repository**
```bash
git clone <repository-url>
cd vendix
```

2. **Create environment file**
```bash
cp env.example .env
# Edit .env with your configuration
```

3. **Start all services**
```bash
make docker-up
# or: docker-compose up -d
```

This will start:
- PostgreSQL (port 5432)
- Redis (port 6379)
- MinIO (ports 9000, 9001)
- API Server (port 8080)
- Worker (background jobs)
- Frontend (port 3000)

4. **Create a test tenant**
```bash
bash scripts/create-test-tenant.sh
```

5. **Access the application**
- Frontend: http://localhost:3000
- API: http://localhost:8080
- API Documentation: http://localhost:8080/swagger/index.html
- MinIO Console: http://localhost:9001

### Running Locally Without Docker (Advanced)

If you prefer to run services directly on your machine for maximum performance:

1. **Start infrastructure services**
```bash
docker-compose up -d postgres redis minio
```

2. **Run the API server**
```bash
make run-api
```

3. **Run the worker** (in another terminal)
```bash
make run-worker
```

4. **Run the frontend** (in another terminal)
```bash
cd frontend
npm install
npm run dev
```

5. **Create a test tenant**
```bash
make create-tenant SLUG=demo
# Note the returned tenant ID
make init-tenant ID=<tenant-id>
```

## 📖 Usage Guide

### Creating Your First Tenant

```bash
curl -X POST http://localhost:8080/api/v1/public/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Company",
    "slug": "mycompany",
    "plan_id": null
  }'
```

### Initializing Tenant Schema

```bash
curl -X POST http://localhost:8080/api/v1/public/tenants/{tenant-id}/init
```

### Registering a User

```bash
curl -X POST http://localhost:8080/api/v1/tenant/auth/register \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: mycompany" \
  -d '{
    "email": "admin@mycompany.com",
    "password": "securepassword123",
    "first_name": "John",
    "last_name": "Doe",
    "role": "admin"
  }'
```

### Logging In

```bash
curl -X POST http://localhost:8080/api/v1/tenant/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID": mycompany" \
  -d '{
    "email": "admin@mycompany.com",
    "password": "securepassword123"
  }'
```

### Creating a Customer

```bash
curl -X POST http://localhost:8080/api/v1/tenant/customers \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: mycompany" \
  -H "Authorization: Bearer <access-token>" \
  -d '{
    "customer_type": "business",
    "tax_id": "123456789",
    "name": "Customer Inc.",
    "email": "contact@customer.com",
    "phone": "+1 809-555-0100",
    "address": "Calle Principal 123",
    "city": "Santo Domingo",
    "country": "DO"
  }'
```

## 🧪 Testing

### Run all tests
```bash
make test
```

### Run tests with coverage
```bash
make test-coverage
```

### Run specific module tests
```bash
go test -v ./internal/auth/...
go test -v ./internal/modules/customers/...
```

## 🔧 Development

### Available Make Commands

**Development (Hot-Reload):**
```bash
make dev-up              # Start development environment
make dev-down            # Stop development environment
make dev-logs            # View development logs
make dev-rebuild         # Rebuild dev containers
make dev-restart-api     # Restart only API service
make dev-restart-frontend # Restart only frontend service
```

**Production:**
```bash
make docker-up           # Start production environment
make docker-down         # Stop production environment
make docker-logs         # View production logs
make docker-rebuild      # Rebuild production containers
```

**Local Development:**
```bash
make build               # Build binaries
make run-api             # Run API locally
make run-worker          # Run worker locally
make test                # Run tests
make test-coverage       # Tests with coverage
make lint                # Run linter
make fmt                 # Format code
make swagger             # Generate Swagger docs
```

**Database:**
```bash
make create-tenant SLUG=demo  # Create test tenant
make init-tenant ID=<uuid>    # Initialize tenant schema
```

For a complete list of commands:
```bash
make help
```

## 🌐 API Documentation

Once the API server is running, visit:
- Swagger UI: http://localhost:8080/swagger/index.html

Key API endpoints:
- `POST /api/v1/public/tenants` - Create tenant
- `POST /api/v1/tenant/auth/register` - Register user
- `POST /api/v1/tenant/auth/login` - Login
- `GET /api/v1/tenant/customers` - List customers
- `POST /api/v1/tenant/customers` - Create customer
- `GET /api/v1/tenant/products` - List products
- `POST /api/v1/tenant/invoices` - Create invoice

## 📊 Database Schemas

### Public Schema (Shared)
- `tenants` - Tenant information
- `plans` - Subscription plans
- `plan_modules` - Module access per plan
- `subscriptions` - Tenant subscriptions
- `tenant_usage` - Usage metrics

### Tenant Schema (Per Tenant)
- `users` - Tenant users
- `refresh_tokens` - JWT refresh tokens
- `customers` - Customer records
- `products` - Product catalog
- `tax_rates` - Tax configurations
- `invoices` - Invoice records
- `invoice_lines` - Invoice line items
- `payments` - Payment records
- `payment_allocations` - Payment to invoice mapping
- `credit_notes` - Credit notes
- `journal_entries` - Accounting entries
- `audit_logs` - Audit trail
- `settings` - Tenant settings

## 🇩🇴 DGII Compliance

The system includes a DGII compliance module with:
- Electronic document signing interface
- Mock adapter for testing (no real DGII credentials required)
- Support for NCF (Número de Comprobante Fiscal)
- ITBIS (tax) calculations
- Document status tracking

To enable DGII integration in production:
1. Obtain DGII API credentials
2. Implement real DGII adapter (interface is provided)
3. Configure certificates in `DGII_CERT_PATH` and `DGII_KEY_PATH`
4. Set `DGII_ENABLED=true` in environment

## 🔐 Security

- JWT-based authentication with short-lived access tokens
- Refresh tokens for extended sessions
- bcrypt password hashing
- Role-based access control
- Rate limiting per tenant and user
- SQL injection protection (parameterized queries)
- CORS configuration
- Audit logging for all operations

## 🚀 Deployment

### Docker Production Build
```bash
docker build -t vendix-api .
docker run -p 8080:8080 vendix-api
```

### Environment Variables
See `env.example` for all available configuration options.

Key variables:
- `APP_ENV` - Environment (development, production)
- `DB_*` - Database configuration
- `JWT_SECRET` - Secret for JWT signing
- `STRIPE_SECRET_KEY` - Stripe API key
- `DGII_ENABLED` - Enable DGII integration

## 📝 License

MIT License - see LICENSE file for details

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Para Desarrolladores

Consulta estos recursos antes de contribuir:
- 🚀 [Inicio Rápido](QUICK_START.md) - Configuración en 5 minutos
- 📖 [Guía de Desarrollo](.github/DEVELOPMENT.md) - Flujo de trabajo y convenciones
- 🔧 [Configuración de Desarrollo](DEV_SETUP.md) - Detalles técnicos de hot-reload

### Proceso de Contribución

1. Fork el repositorio
2. Crea una rama para tu feature (`git checkout -b feat/nueva-funcionalidad`)
3. Desarrolla usando el entorno de desarrollo: `make dev-up`
4. Ejecuta tests: `make test`
5. Ejecuta linter: `make lint`
6. Commit tus cambios (`git commit -m 'feat: agregar nueva funcionalidad'`)
7. Push a tu fork (`git push origin feat/nueva-funcionalidad`)
8. Abre un Pull Request

## 📧 Support

For issues and questions:
- Open an issue on GitHub
- Email: support@billing.example.com

## 🙏 Acknowledgments

Built for businesses in the Dominican Republic to comply with DGII electronic invoicing requirements.

