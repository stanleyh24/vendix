# Sistema de Inventario - Implementación Completa

## 📋 Resumen

Se ha implementado un sistema completo de control de inventario que permite:
- Gestionar cantidades en stock para productos
- Reducción automática de stock al realizar ventas
- Validación de stock disponible antes de procesar ventas
- Auditoría completa de movimientos de inventario
- Diferenciación entre productos físicos y servicios

---

## 🗄️ Cambios en Base de Datos

### Nueva Migración (Versión 17)

Se agregó una nueva migración que incluye:

#### 1. Campo `stock_quantity` en tabla `products`
```sql
ALTER TABLE products ADD COLUMN IF NOT EXISTS stock_quantity DECIMAL(10, 2) DEFAULT 0 NOT NULL;
CREATE INDEX IF NOT EXISTS idx_products_stock_quantity ON products(stock_quantity);
```

#### 2. Tabla de Auditoría `inventory_movements`
```sql
CREATE TABLE IF NOT EXISTS inventory_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    movement_type VARCHAR(50) NOT NULL, -- 'sale', 'purchase', 'adjustment', 'return'
    quantity DECIMAL(10, 2) NOT NULL,
    previous_stock DECIMAL(10, 2) NOT NULL,
    new_stock DECIMAL(10, 2) NOT NULL,
    reference_type VARCHAR(50), -- 'invoice', 'credit_note', 'manual'
    reference_id UUID,
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Ubicación**: `/internal/database/tenant_migrations.go`

---

## 🔧 Cambios en Backend (Go)

### 1. Modelo de Productos
**Archivo**: `/internal/modules/products/models.go`

Se agregó el campo `StockQuantity` a:
- `Product` struct
- `CreateProductRequest`
- `UpdateProductRequest`  
- `ProductResponse`

```go
type Product struct {
    // ... otros campos
    StockQuantity float64   `db:"stock_quantity" json:"stock_quantity"`
    // ...
}
```

### 2. Repositorio de Productos
**Archivo**: `/internal/modules/products/repository.go`

Se actualizaron todos los queries SQL para incluir `stock_quantity` y se agregaron nuevos métodos:

```go
// UpdateStock actualiza la cantidad en stock de un producto
func (r *Repository) UpdateStock(ctx context.Context, schema string, productID uuid.UUID, quantity float64) error

// GetStockQuantity obtiene la cantidad actual en stock de un producto
func (r *Repository) GetStockQuantity(ctx context.Context, schema string, productID uuid.UUID) (float64, error)

// RecordInventoryMovement registra un movimiento de inventario para auditoría
func (r *Repository) RecordInventoryMovement(ctx context.Context, schema string, productID uuid.UUID, movementType string, quantity, previousStock, newStock float64, referenceType *string, referenceID *uuid.UUID, notes *string, createdBy *uuid.UUID) error
```

### 3. Servicio de Productos
**Archivo**: `/internal/modules/products/service.go`

Se actualizó para:
- Incluir `stock_quantity` al crear productos
- Permitir actualización de stock manualmente

### 4. Servicio de Facturas (Invoices)
**Archivo**: `/internal/modules/invoices/service.go`

**Cambios principales**:
- Se importa el módulo de productos
- Se inyecta el repositorio de productos en el constructor
- Al crear una factura, se ejecuta la siguiente lógica:

```go
// Para cada línea de la factura:
1. Verificar si es un producto (no servicio)
2. Validar que hay suficiente stock
3. Reducir el stock automáticamente
4. Registrar el movimiento en inventory_movements
5. Si no hay stock suficiente, retornar error 400
```

**Flujo de validación**:
```go
if product.ProductType == "product" {
    if product.StockQuantity < line.Quantity {
        return error "stock insuficiente"
    }
    // Reducir stock
    // Registrar movimiento
}
```

---

## 🎨 Cambios en Frontend (React)

### Archivo: `/frontend/src/pages/Products.jsx`

#### 1. Estado del Formulario
Se agregó `stock_quantity` al estado inicial:
```javascript
const [formData, setFormData] = useState({
  // ... otros campos
  stock_quantity: '0',
});
```

#### 2. Visualización en Tarjetas de Productos
Se muestra el stock con código de colores:
- **Rojo**: Stock agotado (≤ 0)
- **Amarillo**: Stock bajo (≤ 10)
- **Verde**: Stock normal (> 10)

```jsx
{product.product_type === 'product' && (
  <div className="flex items-center justify-between text-sm">
    <span className="text-gray-600">Stock:</span>
    <span className={`font-bold ${
      product.stock_quantity <= 0 ? 'text-red-600' :
      product.stock_quantity <= 10 ? 'text-yellow-600' :
      'text-[#00C853]'
    }`}>
      {product.stock_quantity.toFixed(2)} {product.unit}
    </span>
  </div>
)}
```

#### 3. Formulario de Creación/Edición
Se agregó campo de stock (solo visible para productos físicos):
```jsx
{formData.product_type === 'product' && (
  <div>
    <label>Cantidad en Stock *</label>
    <input
      type="number"
      step="0.01"
      min="0"
      value={formData.stock_quantity}
      onChange={(e) => setFormData({ ...formData, stock_quantity: e.target.value })}
    />
    <p className="text-xs text-gray-500">
      La cantidad en stock se reducirá automáticamente al realizar ventas
    </p>
  </div>
)}
```

---

## 📚 Documentación API

### Archivo: `/docs/API_EXAMPLES.md`

Se actualizó con:

#### Ejemplo de Creación de Producto
```json
{
  "code": "PROD-001",
  "name": "Laptop Dell XPS 15",
  "product_type": "product",
  "price": 75000.00,
  "stock_quantity": 50.00  // ← Nuevo campo
}
```

#### Notas sobre Inventario
- El campo `stock_quantity` indica la cantidad disponible en stock
- Solo aplica para productos (`product_type: "product"`), no para servicios
- Al crear una factura/venta, el stock se reduce automáticamente
- Si el stock es insuficiente, la venta será rechazada

#### Ejemplo de Creación de Factura
```json
{
  "customer_id": "...",
  "lines": [
    {
      "product_id": "...",
      "quantity": 2,  // ← Se validará contra stock disponible
      "unit_price": 75000.00
    }
  ]
}
```

**Comportamiento**:
- Al crear la factura, el sistema valida que hay suficiente stock
- Si hay stock suficiente, se reduce automáticamente la cantidad vendida
- Se registra un movimiento de inventario
- Si el stock es insuficiente, retorna error 400

---

## 🔄 Flujo Completo de una Venta

### 1. Usuario crea un producto
```
POST /api/v1/tenant/products
{
  "code": "LAPTOP-001",
  "name": "Laptop Dell",
  "product_type": "product",
  "stock_quantity": 100.00
}
```

### 2. El sistema guarda el producto
- Se crea el registro en la tabla `products`
- Campo `stock_quantity` = 100.00

### 3. Usuario crea una factura
```
POST /api/v1/tenant/invoices
{
  "lines": [
    {
      "product_id": "...",
      "quantity": 5
    }
  ]
}
```

### 4. El sistema procesa la venta
```
1. Obtener producto
2. Validar: stock_quantity (100) >= quantity (5) ✓
3. Crear factura
4. Reducir stock: UPDATE products SET stock_quantity = 95
5. Registrar movimiento en inventory_movements:
   - movement_type: "sale"
   - quantity: 5
   - previous_stock: 100
   - new_stock: 95
   - reference_type: "invoice"
6. Retornar factura creada
```

### 5. Stock actualizado
- El producto ahora tiene `stock_quantity` = 95.00
- Hay un registro de auditoría del movimiento

---

## 🚀 Migraciones Necesarias

### Ejecutar Migraciones

Para aplicar los cambios en un tenant existente:

```bash
# Opción 1: Usando script de migración
./scripts/migrate-tenants.sh

# Opción 2: Manualmente para tenant específico
psql -U vendix_user -d vendix_db -c "SET search_path TO tenant_demo; \
  ALTER TABLE products ADD COLUMN IF NOT EXISTS stock_quantity DECIMAL(10, 2) DEFAULT 0 NOT NULL; \
  -- (ejecutar resto del SQL de la migración)"
```

### Datos Existentes
Los productos existentes tendrán `stock_quantity` = 0 por defecto. Se debe actualizar manualmente según corresponda.

---

## ✅ Características Implementadas

- ✅ Campo de stock en productos
- ✅ Reducción automática de stock en ventas
- ✅ Validación de stock disponible
- ✅ Diferenciación productos vs servicios
- ✅ Tabla de auditoría de movimientos
- ✅ UI con indicadores visuales de stock
- ✅ Documentación API actualizada
- ✅ Mensajes de error informativos

---

## 🎯 Casos de Uso

### ✅ Caso 1: Venta Exitosa
```
Producto: Stock = 50
Venta: Cantidad = 10
Resultado: Stock = 40 ✓
```

### ❌ Caso 2: Stock Insuficiente
```
Producto: Stock = 5
Venta: Cantidad = 10
Resultado: Error 400 "stock insuficiente para X: disponible 5.00, solicitado 10.00"
```

### ✅ Caso 3: Servicio (sin stock)
```
Producto: Tipo = "service"
Venta: Cantidad = 10
Resultado: No se valida ni reduce stock ✓
```

---

## 🔮 Mejoras Futuras Sugeridas

1. **Alertas de Stock Bajo**
   - Notificaciones cuando stock < umbral mínimo
   - Dashboard con productos en stock bajo

2. **Reporte de Movimientos**
   - Endpoint para consultar `inventory_movements`
   - Exportación a Excel/PDF

3. **Gestión de Compras**
   - Endpoint para registrar compras de inventario
   - Movimiento tipo "purchase" para aumentar stock

4. **Stock por Ubicación**
   - Múltiples almacenes/ubicaciones
   - Transferencias entre ubicaciones

5. **Stock de Seguridad**
   - Campo `minimum_stock` en productos
   - Alertas automáticas

6. **Devoluciones**
   - Aumentar stock en notas de crédito
   - Movimiento tipo "return"

---

## 📝 Notas Técnicas

### Precisión Decimal
Se usa `DECIMAL(10, 2)` que permite:
- Hasta 99,999,999.99 unidades
- 2 decimales de precisión

### Tipos de Movimientos Soportados
- `sale`: Venta (reduce stock)
- `purchase`: Compra (aumenta stock)
- `adjustment`: Ajuste manual
- `return`: Devolución (aumenta stock)

### Validación de Stock
- Se valida ANTES de crear la factura
- Si falla, la factura NO se crea
- Transacción atómica garantizada

---

## 🐛 Troubleshooting

### Error: "column stock_quantity does not exist"
**Solución**: Ejecutar las migraciones pendientes
```bash
./scripts/migrate-tenants.sh
```

### Stock no se reduce al crear factura
**Verificar**:
1. El producto es tipo "product" (no "service")
2. El `product_id` está presente en la línea de factura
3. Revisar logs del backend para errores

### Frontend no muestra campo de stock
**Verificar**:
1. El tipo de producto es "product"
2. La API retorna el campo `stock_quantity`
3. Refrescar la página del navegador

---

## 📞 Soporte

Para reportar problemas o sugerir mejoras relacionadas con el sistema de inventario, contactar al equipo de desarrollo.

---

**Fecha de Implementación**: 19 de Octubre, 2025  
**Versión de Migración**: 17  
**Estado**: ✅ Completado y Probado

