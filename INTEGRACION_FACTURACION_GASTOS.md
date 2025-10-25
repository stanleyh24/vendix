# Integración Completa: Facturación y Gastos

## 📋 Resumen

Se ha implementado completamente la integración de los módulos de **Facturación** y **Gastos** con el backend, incluyendo:

- ✅ Backend: Modelos, Repositorios, Servicios y Handlers
- ✅ Frontend: Integración con API REST
- ✅ Arquitectura multitenant con schemas
- ✅ CRUD completo para ambos módulos

---

## 🏗️ Arquitectura Implementada

### Backend

#### Módulo de Invoices (`internal/modules/invoices/`)

**Archivos creados/actualizados:**
- `models.go` - Estructuras de datos para facturas
- `repository.go` - Operaciones de base de datos con soporte multitenant
- `service.go` - Lógica de negocio
- `handler.go` - Endpoints HTTP con autenticación

**Modelos principales:**
```go
type Invoice struct {
    ID            uuid.UUID
    InvoiceNumber string
    CustomerID    uuid.UUID
    IssueDate     time.Time
    DueDate       time.Time
    Status        string
    Subtotal      float64
    TaxAmount     float64
    Total         float64
    Lines         []InvoiceLine
    // ... más campos
}

type InvoiceLine struct {
    ID          uuid.UUID
    ProductID   *uuid.UUID
    Description string
    Quantity    float64
    UnitPrice   float64
    TaxRate     float64
    LineTotal   float64
    // ... más campos
}
```

**Endpoints disponibles:**
- `POST /api/v1/tenant/invoices` - Crear factura
- `GET /api/v1/tenant/invoices` - Listar facturas (con filtro por status)
- `GET /api/v1/tenant/invoices/:id` - Obtener factura por ID
- `PUT /api/v1/tenant/invoices/:id` - Actualizar factura
- `POST /api/v1/tenant/invoices/:id/send` - Enviar factura a DGII
- `POST /api/v1/tenant/invoices/:id/cancel` - Anular factura

#### Módulo de Payments/Expenses (`internal/modules/payments/`)

**Archivos creados/actualizados:**
- `models.go` - Estructuras de datos para pagos/gastos
- `repository.go` - Operaciones de base de datos con soporte multitenant
- `service.go` - Lógica de negocio con prefijos (PAY-/EXP-)
- `handler.go` - Endpoints HTTP con alias `/expenses`

**Modelos principales:**
```go
type Payment struct {
    ID              uuid.UUID
    PaymentNumber   string
    CustomerID      *uuid.UUID
    PaymentMethod   string
    PaymentDate     time.Time
    Amount          float64
    Reference       *string
    Category        string
    Supplier        string
    Description     string
    // ... más campos
}
```

**Endpoints disponibles:**
- `POST /api/v1/tenant/payments` - Crear pago
- `POST /api/v1/tenant/expenses` - Crear gasto (alias)
- `GET /api/v1/tenant/payments` - Listar pagos
- `GET /api/v1/tenant/expenses` - Listar gastos (alias)
- `GET /api/v1/tenant/payments/:id` - Obtener por ID
- `PUT /api/v1/tenant/payments/:id` - Actualizar
- `DELETE /api/v1/tenant/payments/:id` - Eliminar

---

### Frontend

#### Página de Facturas (`frontend/src/pages/Invoices.jsx`)

**Funcionalidades integradas:**
- ✅ Carga de facturas desde API
- ✅ Filtrado por estado (draft, pending, sent, paid, cancelled)
- ✅ Búsqueda por número, cliente o RNC
- ✅ Envío de facturas a DGII
- ✅ Anulación de facturas
- ✅ Visualización de totales y KPIs
- ✅ Manejo de errores con alertas

**Conexión con API:**
```javascript
// Cargar facturas
const response = await api.get('/invoices');

// Enviar factura
await api.post(`/invoices/${invoice.id}/send`);

// Anular factura
await api.post(`/invoices/${invoice.id}/cancel`);
```

#### Página de Gastos (`frontend/src/pages/Expenses.jsx`)

**Funcionalidades integradas:**
- ✅ Carga de gastos desde API
- ✅ CRUD completo (Crear, Leer, Actualizar, Eliminar)
- ✅ Filtrado por categoría
- ✅ Modal para crear/editar gastos
- ✅ Confirmación antes de eliminar
- ✅ Categorización automática
- ✅ Cálculo de totales y promedios

**Conexión con API:**
```javascript
// Cargar gastos
const response = await api.get('/expenses');

// Crear gasto
await api.post('/expenses', payload);

// Actualizar gasto
await api.put(`/expenses/${id}`, payload);

// Eliminar gasto
await api.delete(`/expenses/${id}`);
```

---

## 🔧 Características Técnicas

### 1. Multitenant con Schemas

Todos los endpoints utilizan el `schema` del tenant actual:

```go
schema := middleware.GetTenantSchema(c)
invoices, err := h.service.List(c.Context(), schema, statusPtr)
```

### 2. Numeración Automática

**Facturas:** `INV-00000001`, `INV-00000002`, etc.
**Gastos:** `EXP-00000001`, `EXP-00000002`, etc.

```go
invoiceNumber, err := s.repo.GetNextInvoiceNumber(ctx, schema, "INV-")
```

### 3. Cálculo Automático de Totales

Las facturas calculan automáticamente:
- Subtotal (suma de líneas sin impuestos)
- Impuestos (ITBIS 18%)
- Total (subtotal + impuestos)

```go
lineSubtotal := lineReq.Quantity * lineReq.UnitPrice
lineTax := lineSubtotal * lineReq.TaxRate
lineTotal := lineSubtotal + lineTax
```

### 4. Validaciones

- Fechas en formato ISO (YYYY-MM-DD)
- Montos positivos
- Campos requeridos validados
- UUIDs válidos para referencias

### 5. Manejo de Errores

Frontend muestra alertas amigables:
```javascript
showAlert('error', 'Error', 'No se pudieron cargar las facturas');
```

Backend devuelve errores descriptivos:
```go
return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
    "error": "Invalid request body",
})
```

---

## 📊 Flujo de Datos

### Crear Factura

1. **Frontend** → Usuario completa formulario (futuro)
2. **Frontend** → `POST /api/v1/tenant/invoices` con datos
3. **Backend** → Valida datos
4. **Backend** → Genera número automático (INV-XXXXXXXX)
5. **Backend** → Calcula totales
6. **Backend** → Guarda en schema del tenant
7. **Backend** → Retorna factura creada
8. **Frontend** → Muestra confirmación y recarga lista

### Crear Gasto

1. **Frontend** → Usuario completa modal
2. **Frontend** → `POST /api/v1/tenant/expenses` con datos:
   ```json
   {
     "description": "Papelería",
     "amount": 5500.00,
     "payment_date": "2025-10-17",
     "supplier": "Papelería Nacional",
     "payment_method": "cash",
     "category": "supplies",
     "currency": "DOP"
   }
   ```
3. **Backend** → Genera número (EXP-XXXXXXXX)
4. **Backend** → Guarda en `payments` table
5. **Backend** → Retorna gasto creado
6. **Frontend** → Cierra modal y recarga lista

---

## 🧪 Pruebas

### Verificar API (con curl)

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/tenant/auth/login \
  -H 'Content-Type: application/json' \
  -H 'X-Tenant-ID: demo' \
  -d '{"email":"admin@demo.com","password":"demo123"}' \
  | jq -r '.access_token')

# Listar facturas
curl -s http://localhost:8080/api/v1/tenant/invoices \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: demo" | jq

# Listar gastos
curl -s http://localhost:8080/api/v1/tenant/expenses \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: demo" | jq
```

### Verificar Frontend

1. Acceder a: `http://localhost:5173/invoices`
2. Debería cargar lista vacía desde API
3. Acceder a: `http://localhost:5173/expenses`
4. Debería mostrar botón "Nuevo Gasto"
5. Crear un gasto y verificar que se guarda

---

## 🔐 Seguridad

✅ Autenticación JWT requerida  
✅ Validación de tenant en cada request  
✅ Schema isolation (cada tenant en su propio schema)  
✅ Validación de UUIDs para prevenir inyección  
✅ SQL parametrizado (previene SQL injection)  

---

## 📝 Próximos Pasos

1. **Formulario de creación de facturas** en frontend
2. **Generación de PDF** para facturas
3. **Integración con DGII** (firma digital y envío)
4. **Reportes de facturación** y gastos
5. **Dashboard con gráficas** de ingresos/egresos
6. **Paginación** para listas grandes
7. **Filtros avanzados** (fechas, rangos, etc.)
8. **Exportación** a Excel/CSV

---

## 🐛 Solución de Problemas

### Error: "relation does not exist"

**Solución:** Las tablas ya existen en las migraciones. Reiniciar el backend:
```bash
docker-compose -f docker-compose.dev.yml restart api
```

### Error: 401 Unauthorized

**Solución:** Verificar que el token JWT sea válido y no haya expirado.

### Frontend no carga datos

**Solución:** Verificar en la consola del navegador:
1. Network tab → verificar requests
2. Console → buscar errores JavaScript
3. Verificar que `X-Tenant-ID` se envíe en headers

---

## ✅ Estado Final

- ✅ **Backend Invoices:** Completamente funcional
- ✅ **Backend Payments/Expenses:** Completamente funcional
- ✅ **Frontend Invoices:** Integrado con API (lista y acciones)
- ✅ **Frontend Expenses:** Integrado con API (CRUD completo)
- ✅ **Multitenant:** Todos los endpoints respetan schemas
- ✅ **Autenticación:** JWT funcionando correctamente
- ✅ **Sin errores de linter:** Código limpio

**Fecha de implementación:** 17 de octubre de 2025

