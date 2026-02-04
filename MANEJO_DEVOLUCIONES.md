# Manejo de Devoluciones en Vendix

## Visión General

El sistema de devoluciones permite procesar devoluciones de productos vendidos a través de ventas (POS) o facturas, restaurando el inventario, manejando reembolsos y generando los asientos contables correspondientes.

## Arquitectura Propuesta

### 1. Nuevo Módulo: `internal/modules/returns`

**Estructura:**
```
internal/modules/returns/
├── models.go       # Estructuras de datos (Return, ReturnLine, ReturnStatus, etc.)
├── repository.go   # Acceso a base de datos
├── service.go      # Lógica de negocio
└── handler.go      # Endpoints HTTP
```

### 2. Estructura de Datos

#### Tabla: `returns`
```sql
CREATE TABLE IF NOT EXISTS returns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    return_number VARCHAR(50) UNIQUE NOT NULL,
    sale_id UUID REFERENCES sales(id),           -- NULL si es devolución de factura
    invoice_id UUID REFERENCES invoices(id),     -- NULL si es devolución de venta
    customer_id UUID REFERENCES customers(id),
    
    -- Información de la devolución
    return_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    return_type VARCHAR(20) NOT NULL,            -- 'sale', 'invoice'
    return_reason VARCHAR(100),                  -- 'defect', 'wrong_item', 'customer_request', etc.
    notes TEXT,
    
    -- Totales
    subtotal DECIMAL(15,2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    total DECIMAL(15,2) NOT NULL DEFAULT 0,
    
    -- Estado y tracking
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'approved', 'rejected', 'completed', 'refunded'
    refund_status VARCHAR(20),                    -- 'pending', 'partial', 'completed', 'failed'
    refund_method VARCHAR(20),                    -- 'cash', 'card', 'transfer', 'credit_note'
    refund_amount DECIMAL(15,2) DEFAULT 0,
    
    -- Auditoría
    created_by UUID REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_returns_sale_id ON returns(sale_id);
CREATE INDEX idx_returns_invoice_id ON returns(invoice_id);
CREATE INDEX idx_returns_customer_id ON returns(customer_id);
CREATE INDEX idx_returns_status ON returns(status);
CREATE INDEX idx_returns_return_date ON returns(return_date);
CREATE INDEX idx_returns_return_number ON returns(return_number);
```

#### Tabla: `return_lines`
```sql
CREATE TABLE IF NOT EXISTS return_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    return_id UUID NOT NULL REFERENCES returns(id) ON DELETE CASCADE,
    
    -- Referencia a la línea original
    sale_line_id UUID REFERENCES sale_lines(id),     -- NULL si es factura
    invoice_line_id UUID REFERENCES invoice_lines(id), -- NULL si es venta
    
    -- Información del producto devuelto
    product_id UUID REFERENCES products(id),
    line_number INT NOT NULL,
    description VARCHAR(500) NOT NULL,
    quantity DECIMAL(10,2) NOT NULL,
    unit_price DECIMAL(15,2) NOT NULL,
    tax_rate DECIMAL(5,4) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    line_total DECIMAL(15,2) NOT NULL,
    
    -- Estado del producto devuelto
    item_condition VARCHAR(20),                  -- 'new', 'used', 'defective', 'damaged'
    restock_quantity DECIMAL(10,2),             -- Cantidad que se puede revender
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_return_lines_return_id ON return_lines(return_id);
CREATE INDEX idx_return_lines_product_id ON return_lines(product_id);
```

### 3. Modelos Go

```go
// ReturnStatus representa el estado de una devolución
type ReturnStatus string

const (
    ReturnStatusPending   ReturnStatus = "pending"    // Pendiente de aprobación
    ReturnStatusApproved  ReturnStatus = "approved"   // Aprobada, pendiente procesamiento
    ReturnStatusRejected  ReturnStatus = "rejected"   // Rechazada
    ReturnStatusCompleted ReturnStatus = "completed"  // Completada (inventario restaurado)
    ReturnStatusRefunded  ReturnStatus = "refunded"   // Reembolsada completamente
)

// ReturnType representa el tipo de devolución
type ReturnType string

const (
    ReturnTypeSale    ReturnType = "sale"    // Devolución de venta POS
    ReturnTypeInvoice ReturnType = "invoice" // Devolución de factura
)

// RefundMethod representa el método de reembolso
type RefundMethod string

const (
    RefundMethodCash       RefundMethod = "cash"        // Efectivo
    RefundMethodCard       RefundMethod = "card"        // Tarjeta (reversar transacción)
    RefundMethodTransfer   RefundMethod = "transfer"    // Transferencia bancaria
    RefundMethodCreditNote RefundMethod = "credit_note" // Nota de crédito (facturación)
)

// ReturnReason representa la razón de la devolución
type ReturnReason string

const (
    ReturnReasonDefect         ReturnReason = "defect"          // Producto defectuoso
    ReturnReasonWrongItem      ReturnReason = "wrong_item"      // Producto incorrecto
    ReturnReasonCustomerRequest ReturnReason = "customer_request" // Solicitud del cliente
    ReturnReasonDamaged        ReturnReason = "damaged"         // Producto dañado en tránsito
    ReturnReasonNotAsDescribed ReturnReason = "not_as_described" // No coincide con descripción
)

// Return representa una devolución
type Return struct {
    ID            uuid.UUID
    ReturnNumber  string
    SaleID        *uuid.UUID
    InvoiceID     *uuid.UUID
    CustomerID    *uuid.UUID
    ReturnDate    time.Time
    ReturnType    ReturnType
    ReturnReason  *ReturnReason
    Notes         *string
    Subtotal      float64
    TaxAmount     float64
    Total         float64
    Status        ReturnStatus
    RefundStatus  *string
    RefundMethod  *RefundMethod
    RefundAmount  float64
    CreatedBy     *uuid.UUID
    ApprovedBy    *uuid.UUID
    ApprovedAt    *time.Time
    CreatedAt     time.Time
    UpdatedAt     time.Time
    
    // Relations
    Lines         []ReturnLine
    CustomerName  string
    SaleNumber    *string
    InvoiceNumber *string
}

// ReturnLine representa una línea de devolución
type ReturnLine struct {
    ID             uuid.UUID
    ReturnID       uuid.UUID
    SaleLineID     *uuid.UUID
    InvoiceLineID  *uuid.UUID
    ProductID      *uuid.UUID
    LineNumber     int
    Description    string
    Quantity       float64
    UnitPrice      float64
    TaxRate        float64
    TaxAmount      float64
    LineTotal      float64
    ItemCondition  *string
    RestockQuantity *float64
    CreatedAt      time.Time
}
```

## Flujo de Proceso

### 1. Creación de Devolución

**Endpoint:** `POST /api/v1/tenant/returns`

**Proceso:**
1. Usuario selecciona la venta/factura original
2. Selecciona las líneas y cantidades a devolver
3. Especifica razón y condición de los productos
4. Sistema valida:
   - Que la venta/factura existe y está completa
   - Que las cantidades no excedan lo vendido
   - Que no se haya devuelto previamente todo
5. Crea el registro de devolución con status `pending`

### 2. Aprobación de Devolución

**Endpoint:** `PATCH /api/v1/tenant/returns/:id/approve`

**Proceso:**
1. Supervisor revisa y aprueba la devolución
2. Cambia status a `approved`
3. Sistema prepara:
   - Restauración de inventario (si aplica)
   - Proceso de reembolso
   - Asientos contables

### 3. Procesamiento de Devolución

**Endpoint:** `POST /api/v1/tenant/returns/:id/process`

**Proceso:**
1. **Restauración de Inventario:**
   - Por cada línea de devolución:
     - Si `item_condition` es 'new' o 'used': restaura stock completo
     - Si es 'defective' o 'damaged': solo restaura según `restock_quantity`
     - Registra movimiento en `inventory_movements` con `movement_type = 'return'`

2. **Actualización de Venta/Factura Original:**
   - Marca líneas como devueltas (nueva columna `returned_quantity` en sale_lines/invoice_lines)
   - Si todas las líneas están devueltas: cambia status a `refunded`

3. **Reembolso:**
   - Según `refund_method`:
     - **Cash:** Registra salida de efectivo
     - **Card:** Requiere reversar transacción (integración con pasarela)
     - **Transfer:** Registra transferencia bancaria
     - **Credit Note:** Crea nota de crédito (nueva factura con total negativo)

4. **Contabilidad:**
   - Genera asiento contable automático:
     - **DEBE:**
       - Devoluciones y Descuentos en Ventas (cuenta de ingresos negativa)
       - ITBIS por Pagar (reversión de impuesto)
     - **HABER:**
       - Caja/Banco (si reembolso en efectivo)
       - Cuentas por Cobrar (si nota de crédito)
     - Restaura inventario (Débito a Inventario)

5. Cambia status a `completed` o `refunded`

### 4. Rechazo de Devolución

**Endpoint:** `PATCH /api/v1/tenant/returns/:id/reject`

**Proceso:**
1. Cambia status a `rejected`
2. Registra razón de rechazo en `notes`
3. No modifica inventario ni ventas

## Integración con Módulos Existentes

### Ventas (`sales`)
- Agregar campo `returned_quantity` a `sale_lines`
- Método `GetSaleForReturn()` para obtener venta con información completa
- Actualizar status cuando todas las líneas están devueltas

### Facturas (`invoices`)
- Similar a ventas: `returned_quantity` en `invoice_lines`
- Manejar notas de crédito como facturas especiales

### Inventario (`products`)
- Usar `UpdateStock()` para restaurar inventario
- `RecordInventoryMovement()` con `movement_type = 'return'`

### Contabilidad (`accounting`)
- Generar asientos contables automáticos en procesamiento
- Usar cuentas configuradas en `account_mappings`:
  - `sales_return` → Devoluciones y Descuentos
  - `tax_return` → ITBIS por Pagar (reversión)

### Caja (`cashregisters`)
- Si reembolso en efectivo: registrar salida
- Afectar cierre de caja

## Endpoints API

```go
// Listar devoluciones
GET /api/v1/tenant/returns
  Query params:
    - status: filter by status
    - sale_id: filter by sale
    - invoice_id: filter by invoice
    - customer_id: filter by customer
    - start_date, end_date: filter by date range

// Obtener devolución por ID
GET /api/v1/tenant/returns/:id

// Crear devolución
POST /api/v1/tenant/returns
  Body: CreateReturnRequest

// Aprobar devolución
PATCH /api/v1/tenant/returns/:id/approve

// Rechazar devolución
PATCH /api/v1/tenant/returns/:id/reject
  Body: { reason: string }

// Procesar devolución
POST /api/v1/tenant/returns/:id/process
  Body: ProcessReturnRequest
    - refund_method: 'cash' | 'card' | 'transfer' | 'credit_note'
    - refund_amount: decimal (opcional, por defecto total)

// Obtener venta/factura para devolución
GET /api/v1/tenant/sales/:id/for-return
GET /api/v1/tenant/invoices/:id/for-return
```

## Consideraciones Especiales

### 1. Devoluciones Parciales
- Se pueden hacer múltiples devoluciones de la misma venta
- Validar que la suma de devoluciones no exceda la cantidad original

### 2. Productos No Reembolsables
- Agregar campo `is_returnable` en `products`
- Validar antes de procesar devolución

### 3. Tiempo Límite para Devoluciones
- Configurar política de días para aceptar devoluciones (ej: 30 días)
- Validar en creación de devolución

### 4. Notas de Crédito (DGII)
- Si se requiere nota de crédito fiscal:
  - Generar NCF tipo 04 (Nota de Crédito)
  - Enviar a DGII
  - Vincular a factura original

### 5. Reportes
- Reporte de devoluciones por período
- Análisis de razones de devolución
- Impacto en ingresos

## Migración SQL

```sql
-- Versión 24: create_returns_table
CREATE TABLE IF NOT EXISTS returns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    return_number VARCHAR(50) UNIQUE NOT NULL,
    sale_id UUID REFERENCES sales(id),
    invoice_id UUID REFERENCES invoices(id),
    customer_id UUID REFERENCES customers(id),
    return_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    return_type VARCHAR(20) NOT NULL CHECK (return_type IN ('sale', 'invoice')),
    return_reason VARCHAR(100),
    notes TEXT,
    subtotal DECIMAL(15,2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    total DECIMAL(15,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'completed', 'refunded')),
    refund_status VARCHAR(20),
    refund_method VARCHAR(20),
    refund_amount DECIMAL(15,2) DEFAULT 0,
    created_by UUID REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_returns_sale_id ON returns(sale_id);
CREATE INDEX idx_returns_invoice_id ON returns(invoice_id);
CREATE INDEX idx_returns_customer_id ON returns(customer_id);
CREATE INDEX idx_returns_status ON returns(status);
CREATE INDEX idx_returns_return_date ON returns(return_date);
CREATE INDEX idx_returns_return_number ON returns(return_number);

-- Versión 25: create_return_lines_table
CREATE TABLE IF NOT EXISTS return_lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    return_id UUID NOT NULL REFERENCES returns(id) ON DELETE CASCADE,
    sale_line_id UUID REFERENCES sale_lines(id),
    invoice_line_id UUID REFERENCES invoice_lines(id),
    product_id UUID REFERENCES products(id),
    line_number INT NOT NULL,
    description VARCHAR(500) NOT NULL,
    quantity DECIMAL(10,2) NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(15,2) NOT NULL,
    tax_rate DECIMAL(5,4) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    line_total DECIMAL(15,2) NOT NULL,
    item_condition VARCHAR(20) CHECK (item_condition IN ('new', 'used', 'defective', 'damaged')),
    restock_quantity DECIMAL(10,2),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_return_lines_return_id ON return_lines(return_id);
CREATE INDEX idx_return_lines_product_id ON return_lines(product_id);
CREATE INDEX idx_return_lines_sale_line_id ON return_lines(sale_line_id);
CREATE INDEX idx_return_lines_invoice_line_id ON return_lines(invoice_line_id);

-- Versión 26: add_returned_quantity_to_sale_lines
ALTER TABLE sale_lines 
ADD COLUMN IF NOT EXISTS returned_quantity DECIMAL(10,2) DEFAULT 0;

ALTER TABLE sale_lines 
ADD CONSTRAINT check_returned_quantity 
CHECK (returned_quantity >= 0 AND returned_quantity <= quantity);

-- Versión 27: add_returned_quantity_to_invoice_lines
ALTER TABLE invoice_lines 
ADD COLUMN IF NOT EXISTS returned_quantity DECIMAL(10,2) DEFAULT 0;

ALTER TABLE invoice_lines 
ADD CONSTRAINT check_returned_quantity_invoice 
CHECK (returned_quantity >= 0 AND returned_quantity <= quantity);

-- Versión 28: add_is_returnable_to_products
ALTER TABLE products 
ADD COLUMN IF NOT EXISTS is_returnable BOOLEAN DEFAULT true;
```

## Frontend

### Página: `/returns` o `/devoluciones`

**Componentes necesarios:**
1. Lista de devoluciones (tabla)
2. Modal para crear devolución:
   - Buscar venta/factura original
   - Seleccionar líneas y cantidades
   - Razón y condición
   - Método de reembolso
3. Detalle de devolución
4. Procesar devolución (si está aprobada)
5. Historial de devoluciones por venta/factura

**Integraciones:**
- Desde página de ventas: botón "Devolver"
- Desde página de facturas: botón "Devolver"
- Página dedicada de devoluciones

## Seguridad y Permisos

- **Crear devolución:** Todos los usuarios autorizados
- **Aprobar/Rechazar:** Solo supervisores/administradores
- **Procesar devolución:** Solo usuarios autorizados (puede requerir 2FA para montos altos)
- **Ver reportes:** Según permisos de módulo de reportes

## Próximos Pasos

1. ✅ Definir estructura de datos
2. ⏳ Crear migraciones SQL
3. ⏳ Implementar módulo backend (models, repository, service, handler)
4. ⏳ Integrar con módulos existentes
5. ⏳ Crear endpoints API
6. ⏳ Implementar frontend
7. ⏳ Agregar reportes de devoluciones
9. ⏳ Tests unitarios e integración

