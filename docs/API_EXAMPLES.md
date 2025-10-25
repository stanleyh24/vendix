# API Examples

This document contains practical examples for using the Vendix API.

## Table of Contents
1. [Tenant Management](#tenant-management)
2. [Authentication](#authentication)
3. [Customer Management](#customer-management)
4. [Product Management](#product-management)
5. [Invoice Management](#invoice-management)
6. [Payment Management](#payment-management)

## Tenant Management

### Create a New Tenant

```bash
curl -X POST http://localhost:8080/api/v1/public/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Demo Company SRL",
    "slug": "demo",
    "domain": "demo.billing.example.com",
    "plan_id": null
  }'
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Demo Company SRL",
  "slug": "demo",
  "schema_name": "tenant_demo",
  "domain": "demo.billing.example.com",
  "status": "active",
  "created_at": "2025-10-16T10:00:00Z",
  "updated_at": "2025-10-16T10:00:00Z"
}
```

### Initialize Tenant Schema

```bash
curl -X POST http://localhost:8080/api/v1/public/tenants/550e8400-e29b-41d4-a716-446655440000/init
```

Response:
```json
{
  "message": "Tenant schema initialized successfully"
}
```

### List All Tenants

```bash
curl http://localhost:8080/api/v1/public/tenants
```

## Authentication

### Register a New User

```bash
curl -X POST http://localhost:8080/api/v1/tenant/auth/register \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: demo" \
  -d '{
    "email": "admin@demo.com",
    "password": "SecurePassword123!",
    "first_name": "Juan",
    "last_name": "Pérez",
    "role": "admin"
  }'
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "550e8400-e29b-41d4-a716-446655440001",
  "token_type": "Bearer",
  "expires_in": 900,
  "user": {
    "id": "660e8400-e29b-41d4-a716-446655440002",
    "email": "admin@demo.com",
    "first_name": "Juan",
    "last_name": "Pérez",
    "role": "admin",
    "is_active": true,
    "created_at": "2025-10-16T10:05:00Z"
  }
}
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/tenant/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: demo" \
  -d '{
    "email": "admin@demo.com",
    "password": "SecurePassword123!"
  }'
```

### Refresh Access Token

```bash
curl -X POST http://localhost:8080/api/v1/tenant/auth/refresh \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: demo" \
  -d '{
    "refresh_token": "550e8400-e29b-41d4-a716-446655440001"
  }'
```

### Get Current User

```bash
curl http://localhost:8080/api/v1/tenant/auth/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-Tenant-ID: demo"
```

### Logout

```bash
curl -X POST http://localhost:8080/api/v1/tenant/auth/logout \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-Tenant-ID: demo" \
  -d '{
    "refresh_token": "550e8400-e29b-41d4-a716-446655440001"
  }'
```

## Customer Management

### Create a Customer

```bash
curl -X POST http://localhost:8080/api/v1/tenant/customers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-Tenant-ID: demo" \
  -d '{
    "customer_type": "business",
    "tax_id": "131234567",
    "name": "Comercial Ejemplo SRL",
    "email": "contacto@ejemplo.com.do",
    "phone": "+1 809-555-1234",
    "address": "Av. 27 de Febrero #123",
    "city": "Santo Domingo",
    "state": "Distrito Nacional",
    "postal_code": "10101",
    "country": "DO"
  }'
```

Response:
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440003",
  "customer_type": "business",
  "tax_id": "131234567",
  "name": "Comercial Ejemplo SRL",
  "email": "contacto@ejemplo.com.do",
  "phone": "+1 809-555-1234",
  "address": "Av. 27 de Febrero #123",
  "city": "Santo Domingo",
  "state": "Distrito Nacional",
  "postal_code": "10101",
  "country": "DO",
  "is_active": true,
  "created_at": "2025-10-16T10:10:00Z",
  "updated_at": "2025-10-16T10:10:00Z"
}
```

### List Customers

```bash
curl http://localhost:8080/api/v1/tenant/customers \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-Tenant-ID: demo"
```

### Get Customer by ID

```bash
curl http://localhost:8080/api/v1/tenant/customers/770e8400-e29b-41d4-a716-446655440003 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-Tenant-ID: demo"
```

### Update Customer

```bash
curl -X PUT http://localhost:8080/api/v1/tenant/customers/770e8400-e29b-41d4-a716-446655440003 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-Tenant-ID: demo" \
  -d '{
    "phone": "+1 809-555-5678",
    "email": "nuevo@ejemplo.com.do"
  }'
```

### Delete Customer

```bash
curl -X DELETE http://localhost:8080/api/v1/tenant/customers/770e8400-e29b-41d4-a716-446655440003 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-Tenant-ID: demo"
```

## Product Management

### Create a Product

```bash
curl -X POST http://localhost:8080/api/v1/tenant/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-Tenant-ID: demo" \
  -d '{
    "code": "PROD-001",
    "name": "Laptop Dell XPS 15",
    "description": "Laptop de alta gama con procesador Intel i7",
    "product_type": "product",
    "unit": "unidad",
    "price": 75000.00,
    "cost": 60000.00,
    "tax_rate": 0.18,
    "stock_quantity": 50.00
  }'
```

**Nota sobre inventario:**
- El campo `stock_quantity` indica la cantidad disponible en stock
- Solo aplica para productos (`product_type: "product"`), no para servicios
- Al crear una factura/venta, el stock se reduce automáticamente
- Si el stock es insuficiente, la venta será rechazada

Response:
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440004",
  "code": "PROD-001",
  "name": "Laptop Dell XPS 15",
  "description": "Laptop de alta gama con procesador Intel i7",
  "product_type": "product",
  "unit": "unidad",
  "price": 75000.00,
  "cost": 60000.00,
  "tax_rate": 0.18,
  "stock_quantity": 50.00,
  "is_active": true,
  "created_at": "2025-10-16T10:15:00Z",
  "updated_at": "2025-10-16T10:15:00Z"
}
```

### List Products

```bash
curl http://localhost:8080/api/v1/tenant/products \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-Tenant-ID: demo"
```

## Invoice Management

### Create an Invoice (Crédito Fiscal - NCF 01)

```bash
curl -X POST http://localhost:8080/api/v1/tenant/invoices \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-Tenant-ID: demo" \
  -d '{
    "customer_id": "770e8400-e29b-41d4-a716-446655440003",
    "ncf_type": "01",
    "issue_date": "2025-10-16",
    "due_date": "2025-11-16",
    "notes": "Factura con crédito fiscal para empresa",
    "lines": [
      {
        "product_id": "880e8400-e29b-41d4-a716-446655440004",
        "description": "Laptop Dell XPS 15",
        "quantity": 2,
        "unit_price": 75000.00,
        "tax_rate": 0.18
      }
    ]
  }'
```

### Create an Invoice (Consumidor Final - NCF 02)

```bash
curl -X POST http://localhost:8080/api/v1/tenant/invoices \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-Tenant-ID: demo" \
  -d '{
    "ncf_type": "02",
    "issue_date": "2025-10-16",
    "due_date": "2025-11-16",
    "notes": "Venta rápida - Consumidor final",
    "lines": [
      {
        "product_id": "880e8400-e29b-41d4-a716-446655440004",
        "description": "Laptop Dell XPS 15",
        "quantity": 1,
        "unit_price": 75000.00,
        "tax_rate": 0.18
      }
    ]
  }'
```

**Nota sobre customer_id:**
- El campo `customer_id` es **opcional**
- Si se omite o es `null`, el sistema automáticamente usa el "Cliente Genérico" (ID: 00000000-0000-0000-0000-000000000001)
- El cliente genérico se crea automáticamente en cada tenant para ventas rápidas sin identificación específica
- Útil para ventas al por menor donde no se requiere identificar al cliente

**Nota sobre ncf_type (Tipo de Comprobante Fiscal):**
- El campo `ncf_type` especifica el tipo de NCF a emitir
- Tipos disponibles:
  - **"01"**: Crédito Fiscal - Para empresas con RNC (permite deducir ITBIS)
  - **"02"**: Consumidor Final - Para personas sin RNC o ventas generales (default)
- Si se omite, el sistema usa "02" (Consumidor Final) por defecto
- Solo son válidos los valores "01" y "02" en esta versión

**Nota sobre reducción de stock:**
- Al crear una factura con productos (no servicios), el sistema valida que hay suficiente stock
- Si hay stock suficiente, se reduce automáticamente la cantidad vendida
- Se registra un movimiento de inventario en la tabla `inventory_movements`
- Si el stock es insuficiente, la API retornará un error 400 indicando el problema

### Send Invoice to DGII

```bash
curl -X POST http://localhost:8080/api/v1/tenant/invoices/990e8400-e29b-41d4-a716-446655440005/send \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-Tenant-ID: demo"
```

### Cancel Invoice

```bash
curl -X POST http://localhost:8080/api/v1/tenant/invoices/990e8400-e29b-41d4-a716-446655440005/cancel \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-Tenant-ID: demo" \
  -d '{
    "reason": "Error en la facturación"
  }'
```

## Payment Management

### Record a Payment

```bash
curl -X POST http://localhost:8080/api/v1/tenant/payments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "X-Tenant-ID: demo" \
  -d '{
    "customer_id": "770e8400-e29b-41d4-a716-446655440003",
    "payment_method": "bank_transfer",
    "payment_date": "2025-10-20",
    "amount": 70800.00,
    "currency": "DOP",
    "reference": "TRANS-20251020-001",
    "notes": "Pago de factura octubre",
    "allocations": [
      {
        "invoice_id": "990e8400-e29b-41d4-a716-446655440005",
        "amount": 70800.00
      }
    ]
  }'
```

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "error": "Error message description",
  "code": 400
}
```

Common HTTP status codes:
- `200 OK` - Success
- `201 Created` - Resource created
- `204 No Content` - Success with no response body
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

