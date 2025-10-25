# 🔐 Credenciales de Demo - Vendix

## ✅ Tenant y Usuario de Prueba Creados

Se ha creado exitosamente un tenant demo con un usuario administrador para pruebas.

---

## 📋 Información del Tenant

```
Tenant ID:   b6272a90-643a-409d-8b66-416756cc01f8
Tenant Name: Demo Company
Tenant Slug: demo
Schema Name: tenant_demo
Status:      active
```

---

## 👤 Credenciales de Usuario Admin

```
Email:    admin@demo.com
Password: demo123
Role:     admin
```

**✅ ACTUALIZADO:** El password ahora es `demo123` (antes era `password123`)

---

## 🚀 Cómo Usar en el Frontend

### 1. Acceder a la Aplicación

```
http://localhost:3000
```

### 2. Login

1. Ir a la página de login
2. Ingresar credenciales:
   - **Email**: `admin@demo.com`
   - **Password**: `demo123`
   - **Tenant ID**: `demo` (si es requerido)
3. Click en "Iniciar Sesión"

### 3. Navegación

Una vez logueado, puedes acceder a:
- **Dashboard** - Vista principal con KPIs
- **Ventas** - (en desarrollo)
- **Clientes** - Gestión de clientes
- **Inventario** - Gestión de productos
- **Gastos** - (en desarrollo)
- **Facturación** - Gestión de facturas

---

## 🔧 Endpoints de API Funcionando

### Autenticación

```bash
# Login
curl -X POST http://localhost:8080/api/v1/tenant/auth/login \
  -H 'Content-Type: application/json' \
  -H 'X-Tenant-ID: demo' \
  -d '{
    "email": "admin@demo.com",
    "password": "demo123"
  }'
```

### Productos

```bash
# Obtener token del login primero, luego:

# Listar productos
curl -X GET http://localhost:8080/api/v1/tenant/products \
  -H 'Authorization: Bearer <TOKEN>' \
  -H 'X-Tenant-ID: demo'

# Crear producto
curl -X POST http://localhost:8080/api/v1/tenant/products \
  -H 'Authorization: Bearer <TOKEN>' \
  -H 'X-Tenant-ID: demo' \
  -H 'Content-Type: application/json' \
  -d '{
    "code": "PROD-001",
    "name": "Laptop Dell XPS 15",
    "description": "Laptop de alto rendimiento",
    "product_type": "product",
    "unit": "unidad",
    "price": 75000.00,
    "cost": 60000.00,
    "tax_rate": 0.18
  }'
```

---

## 🎯 Estado Actual

### ✅ Componentes Funcionando

- ✅ Backend API corriendo en puerto 8080
- ✅ Frontend corriendo en puerto 3000
- ✅ Tenant demo creado e inicializado
- ✅ Usuario admin creado
- ✅ Rutas de API corregidas
- ✅ Autenticación funcionando
- ✅ Productos endpoint funcionando

### 🗄️ Base de Datos

- ✅ Schema público con tabla de tenants
- ✅ Schema `tenant_demo` creado
- ✅ Tablas de tenant inicializadas:
  - users
  - customers
  - products
  - invoices
  - payments
  - y más...

---

## 🔑 Token de Ejemplo

```
Access Token (válido por 15 minutos):
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNGEzYzVhMmYtNjNkOS00NDA2LTgzZmUtMmY5NmIyYTg3ZjI1IiwiZW1haWwiOiJhZG1pbkBkZW1vLmNvbSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc2MDY3NDQ4NywibmJmIjoxNzYwNjczNTg3LCJpYXQiOjE3NjA2NzM1ODcsImp0aSI6IjdkYTcyMGUxLTcxYTUtNDU2ZC1hOGVmLTdlOTNhOThlNWRhYiJ9.L4x1JmTrQiU291Y_GwT1_7Qj-DMYON1DwA7INWv6fXE

Refresh Token (válido por 7 días):
96886bb5-8efd-450c-b511-3355c80aa484
```

**Nota**: Este token es de ejemplo. Cada login genera un nuevo token.

---

## 📊 Datos de Prueba

### Crear Producto de Prueba

```javascript
// Desde el frontend (Console de DevTools)
const response = await api.post('/products', {
  code: 'LAPTOP-001',
  name: 'Laptop Dell XPS 15',
  description: 'Procesador Intel i7, 16GB RAM, 512GB SSD',
  product_type: 'product',
  unit: 'unidad',
  price: 75000.00,
  cost: 60000.00,
  tax_rate: 0.18
});

console.log(response.data);
```

### Crear Cliente de Prueba

```javascript
const response = await api.post('/customers', {
  tax_id: '001-0000000-0',
  name: 'Juan Pérez',
  email: 'juan@example.com',
  phone: '809-555-1234',
  address: 'Calle Principal #123, Santo Domingo'
});

console.log(response.data);
```

---

## 🛠️ Comandos Útiles

### Recrear Tenant Demo

Si necesitas recrear el tenant desde cero:

```bash
# Eliminar tenant existente (si existe)
# Nota: Esto es manual en la DB por ahora

# Crear nuevo tenant
bash scripts/create-test-tenant.sh

# Crear usuario admin
curl -X POST http://localhost:8080/api/v1/tenant/auth/register \
  -H 'Content-Type: application/json' \
  -H 'X-Tenant-ID: demo' \
  -d '{
    "email": "admin@demo.com",
    "password": "password123",
    "first_name": "Admin",
    "last_name": "Demo",
    "role": "admin"
  }'
```

### Verificar Estado del Sistema

```bash
# Ver todos los servicios
docker-compose -f docker-compose.dev.yml ps

# Ver logs del backend
docker-compose -f docker-compose.dev.yml logs -f api

# Ver logs del frontend
docker-compose -f docker-compose.dev.yml logs -f frontend
```

---

## 🎓 Flujo de Autenticación

```
1. Usuario abre http://localhost:3000
   ↓
2. Es redirigido a /login
   ↓
3. Ingresa credenciales:
   - Email: admin@demo.com
   - Password: password123
   - Tenant: demo
   ↓
4. Frontend hace POST a /api/v1/tenant/auth/login
   ↓
5. Backend:
   - Valida tenant "demo" existe
   - Valida credenciales
   - Genera JWT token
   - Retorna access_token + refresh_token
   ↓
6. Frontend:
   - Guarda tokens en authStore
   - Redirige a /dashboard
   ↓
7. Todas las peticiones posteriores incluyen:
   - Header: Authorization: Bearer <token>
   - Header: X-Tenant-ID: demo
```

---

## 🚀 Prueba Rápida

### Desde la Terminal

```bash
# 1. Login y guardar token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/tenant/auth/login \
  -H 'Content-Type: application/json' \
  -H 'X-Tenant-ID: demo' \
  -d '{"email":"admin@demo.com","password":"password123"}' \
  | grep -o '"access_token":"[^"]*' | sed 's/"access_token":"//')

# 2. Listar productos
curl -X GET http://localhost:8080/api/v1/tenant/products \
  -H "Authorization: Bearer $TOKEN" \
  -H 'X-Tenant-ID: demo'

# 3. Crear producto
curl -X POST http://localhost:8080/api/v1/tenant/products \
  -H "Authorization: Bearer $TOKEN" \
  -H 'X-Tenant-ID: demo' \
  -H 'Content-Type: application/json' \
  -d '{
    "code": "TEST-001",
    "name": "Producto de Prueba",
    "product_type": "product",
    "unit": "unidad",
    "price": 1000.00,
    "tax_rate": 0.18
  }'

# 4. Listar productos nuevamente
curl -X GET http://localhost:8080/api/v1/tenant/products \
  -H "Authorization: Bearer $TOKEN" \
  -H 'X-Tenant-ID: demo'
```

---

## 🎉 Todo Listo

**✅ Sistema Completamente Funcional**

- Tenant demo creado
- Usuario admin disponible
- API funcionando correctamente
- Frontend conectado al backend
- Rutas corregidas
- Error 404 solucionado

**Puedes empezar a usar la aplicación inmediatamente con las credenciales proporcionadas.** 🚀

---

**Credenciales:**
- Email: `admin@demo.com`
- Password: `password123`
- Tenant: `demo`

**URL:** http://localhost:3000

---

**Fecha:** Octubre 17, 2025  
**Versión:** 1.0  
**Estado:** ✅ Listo para Usar

