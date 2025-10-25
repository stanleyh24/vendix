# Architecture Documentation

## System Overview

Vendix is a multi-tenant electronic invoicing system built for the Dominican Republic market. The system follows a modular monolith architecture, designed to be easily split into microservices if needed.

## Core Architectural Principles

1. **Multi-Tenancy**: Isolated data per tenant using PostgreSQL schemas
2. **Modularity**: Each business domain is self-contained
3. **Separation of Concerns**: Clear layers (API, Service, Repository)
4. **Testability**: Interfaces and dependency injection
5. **Scalability**: Stateless design, horizontal scaling capable

## System Architecture Diagram

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Browser   │────▶│   Frontend   │────▶│  API Server │
│   (React)   │     │  (Vite/React)│     │   (Fiber)   │
└─────────────┘     └──────────────┘     └──────┬──────┘
                                                 │
                    ┌────────────────────────────┼────────────────┐
                    │                            │                │
             ┌──────▼──────┐            ┌───────▼────────┐  ┌───▼───┐
             │ PostgreSQL  │            │  Background    │  │ Redis │
             │ (Multi-     │            │    Worker      │  │       │
             │  Tenant)    │            │   (Asynq)      │  └───────┘
             └─────────────┘            └────────────────┘
                    │
          ┌─────────┼─────────┐
          │         │         │
    ┌─────▼───┐ ┌──▼────┐ ┌──▼────┐
    │ public  │ │tenant1│ │tenant2│
    │ schema  │ │schema │ │schema │
    └─────────┘ └───────┘ └───────┘
```

## Layer Architecture

### 1. API Layer (`internal/server`, `internal/modules/*/handler.go`)

**Responsibilities:**
- HTTP request/response handling
- Request validation
- Route registration
- Middleware application

**Example:**
```go
func (h *Handler) Create(c *fiber.Ctx) error {
    var req CreateCustomerRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    customer, err := h.service.Create(c.Context(), schema, &req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.Status(201).JSON(customer)
}
```

### 2. Service Layer (`internal/modules/*/service.go`)

**Responsibilities:**
- Business logic
- Transaction orchestration
- Cross-module coordination
- Authorization checks

**Example:**
```go
func (s *Service) Create(ctx context.Context, schema string, req *CreateCustomerRequest) (*CustomerResponse, error) {
    // Business logic
    customer := &Customer{
        ID: uuid.New(),
        Name: req.Name,
        // ... validation and transformation
    }
    
    if err := s.repo.Create(ctx, schema, customer); err != nil {
        return nil, err
    }
    
    return customer.ToResponse(), nil
}
```

### 3. Repository Layer (`internal/modules/*/repository.go`)

**Responsibilities:**
- Database operations
- Query construction
- Schema-aware queries
- Data mapping

**Example:**
```go
func (r *Repository) Create(ctx context.Context, schema string, customer *Customer) error {
    query := fmt.Sprintf(`
        INSERT INTO %s.customers (id, name, email, ...)
        VALUES ($1, $2, $3, ...)
    `, schema)
    
    _, err := r.db.ExecContext(ctx, query, customer.ID, customer.Name, ...)
    return err
}
```

## Multi-Tenancy Implementation

### Schema Isolation

Each tenant gets a dedicated PostgreSQL schema:

```
Database: vendix
├── public schema (shared)
│   ├── tenants
│   ├── plans
│   ├── subscriptions
│   └── plan_modules
├── tenant_demo schema
│   ├── users
│   ├── customers
│   ├── products
│   ├── invoices
│   └── ...
└── tenant_acme schema
    ├── users
    ├── customers
    └── ...
```

### Tenant Resolution

1. **Request arrives** with `X-Tenant-ID` header or subdomain
2. **TenantMiddleware** resolves tenant from database
3. **Schema is set** in request context
4. **All queries** use the tenant's schema

```go
// Middleware sets tenant schema
func TenantMiddleware(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        tenant := resolveTenant(c)
        c.Locals("tenant_schema", tenant.SchemaName)
        return c.Next()
    }
}

// Repository uses schema from context
query := fmt.Sprintf("SELECT * FROM %s.customers", schema)
```

## Authentication & Authorization

### JWT Flow

```
1. User logs in ──▶ POST /auth/login
                      │
2. Validate credentials ◀──┐
                      │    │
3. Generate tokens   │    │
   - Access Token    │    Database
   - Refresh Token   │    │
                      │    │
4. Store refresh token ──▶┘
                      │
5. Return tokens ◀────┘
```

### Token Types

**Access Token (JWT)**
- Short-lived (15 minutes)
- Contains: user_id, email, role
- Used for API authentication

**Refresh Token (UUID)**
- Long-lived (7 days)
- Stored in database
- Used to get new access tokens

### Role-Based Access Control

Roles:
- **Admin**: Full access
- **Accountant**: Financial operations
- **Billing Clerk**: Invoice management
- **Viewer**: Read-only access

## Module System

### Plan-Based Module Access

```go
// Check if tenant has module access
func ModuleMiddleware(db *database.DB, moduleName string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        hasAccess := checkModuleAccess(tenantID, moduleName)
        if !hasAccess {
            return c.Status(403).JSON(fiber.Map{
                "error": "Module not available in your plan"
            })
        }
        return c.Next()
    }
}
```

### Available Modules

- `invoicing` - Invoice management
- `customers` - Customer management
- `products` - Product catalog
- `payments` - Payment processing
- `reports` - Reports and analytics
- `accounting` - Journal entries, accounting
- `pos` - Point of sale
- `api_access` - API integration

## Background Jobs

### Job Queue Architecture

```
API Request ──▶ Enqueue Job ──▶ Redis Queue
                                    │
                                    ▼
                              Worker Process
                                    │
                              ┌─────┴─────┐
                              │   Handle  │
                              │    Job    │
                              └─────┬─────┘
                                    ▼
                              Update Database
```

### Job Types

1. **Email Jobs** - Send invoices, notifications
2. **DGII Jobs** - Sign and submit invoices
3. **Billing Jobs** - Process recurring billing
4. **Report Jobs** - Generate monthly reports

## DGII Integration

### Adapter Pattern

```go
type DGIIAdapter interface {
    SignDocument(ctx context.Context, doc *ElectronicDocument) (*SignedDocument, error)
    SendDocument(ctx context.Context, signed *SignedDocument) (*DGIIResponse, error)
    CancelDocument(ctx context.Context, docID string, reason string) error
    GetDocumentStatus(ctx context.Context, docID string) (*DocumentStatus, error)
}
```

### Mock vs Real Implementation

- **MockAdapter**: For testing, no real DGII connection
- **RealAdapter**: Production implementation (to be implemented with real credentials)

## Database Design

### Public Schema Tables

```sql
-- Tenant management
CREATE TABLE tenants (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    slug VARCHAR(100) UNIQUE,
    schema_name VARCHAR(63) UNIQUE,
    status VARCHAR(50),
    plan_id UUID
);

-- Subscription management
CREATE TABLE plans (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    price DECIMAL(10,2),
    features JSONB
);

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY,
    tenant_id UUID REFERENCES tenants(id),
    plan_id UUID REFERENCES plans(id),
    status VARCHAR(50),
    current_period_end TIMESTAMP
);
```

### Tenant Schema Tables

Each tenant schema has:
- `users` - User accounts
- `customers` - Customer records
- `products` - Product catalog
- `invoices` - Invoice records
- `invoice_lines` - Invoice line items
- `payments` - Payment records
- `journal_entries` - Accounting entries
- `audit_logs` - Audit trail

## Scalability Considerations

### Horizontal Scaling

The application is **stateless** and can be scaled horizontally:

```bash
# Scale to 5 instances
kubectl scale deployment/vendix-api --replicas=5
```

### Database Scaling

- **Connection pooling**: Configured per instance
- **Read replicas**: Can add read-only replicas
- **Partitioning**: Schema-per-tenant provides natural partitioning

### Caching Strategy

- **Redis**: Session data, rate limiting
- **Application cache**: Reference data (plans, tax rates)

## Security Architecture

### Defense in Depth

1. **Network**: Firewall, VPC isolation
2. **Application**: JWT, RBAC, rate limiting
3. **Database**: Encrypted at rest, SSL connections
4. **Audit**: Immutable audit logs

### Data Isolation

- Schema-level isolation per tenant
- Row-level security possible for future enhancement
- Separate encryption keys per tenant (optional)

## Future Microservices Migration

The modular design allows easy extraction:

```
Current: Modular Monolith
└── vendix
    ├── tenants module
    ├── auth module
    ├── customers module
    └── invoices module

Future: Microservices
├── tenant-service
├── auth-service
├── customer-service
├── invoice-service
└── payment-service
```

## Performance Considerations

### Query Optimization

- Indexes on foreign keys
- Composite indexes for common queries
- Prepared statements (via pgx)

### Caching

- Module access checks cached in Redis
- Tenant resolution cached per request

### Rate Limiting

- Per-tenant limits
- Per-user limits
- Configurable windows

## Monitoring & Observability

### Structured Logging

```go
logger.Info("Invoice created",
    "invoice_id", invoice.ID,
    "tenant_id", tenantID,
    "amount", invoice.Total,
)
```

### Metrics (Prometheus)

- Request duration
- Error rates
- Active tenants
- Job queue depth

### Health Checks

- `/health` - Overall health
- Database connectivity
- Redis connectivity

## Testing Strategy

### Unit Tests

- Service layer logic
- Business rules
- Utilities (JWT, password hashing)

### Integration Tests

- API endpoints
- Database operations
- Background jobs

### End-to-End Tests

- Complete workflows
- Multi-tenant scenarios
- DGII submission flow

