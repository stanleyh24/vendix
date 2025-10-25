# 🔧 Solución: Error de Número de Factura Duplicada

## Problema

Al crear una factura aparece el siguiente error:
```
Error al procesar venta
failed to create invoice: ERROR: duplicate key value violates unique constraint "invoices_invoice_number_key" (SQLSTATE 23505)
```

## Causa

Este error ocurre cuando el sistema intenta crear una factura con un número que ya existe en la base de datos. Esto puede suceder por:

1. **Facturas de prueba existentes**: Hay facturas creadas anteriormente que ya usan ese número
2. **Error en la generación del número**: La función que genera números secuenciales no estaba manejando correctamente el caso de "no hay registros"
3. **Condición de carrera**: Dos ventas procesándose al mismo tiempo pueden generar el mismo número

## Solución Implementada

### 1. Código Corregido ✅

He corregido la función `GetNextInvoiceNumber` en el archivo:
```
/internal/modules/invoices/repository.go
```

**Cambios realizados**:
- ✅ Uso de `sql.NullInt64` para manejar valores NULL
- ✅ Consulta más robusta con `MAX()` y `COALESCE()`
- ✅ Validación del formato del número de factura
- ✅ Extracción correcta de la parte numérica
- ✅ Manejo adecuado de errores

**Nueva implementación**:
```go
func (r *Repository) GetNextInvoiceNumber(ctx context.Context, schema string, prefix string) (string, error) {
    var maxNum sql.NullInt64
    query := fmt.Sprintf(`
        SELECT COALESCE(
            MAX(
                CAST(
                    SUBSTRING(invoice_number FROM LENGTH($1) + 1) AS INTEGER
                )
            ), 
            0
        ) as max_num
        FROM %s.invoices 
        WHERE invoice_number ~ ('^' || $1 || '[0-9]+$')
          AND LENGTH(invoice_number) = LENGTH($1) + 8
    `, schema)

    err := r.db.GetContext(ctx, &maxNum, query, prefix)
    if err != nil {
        return "", fmt.Errorf("failed to get next invoice number: %w", err)
    }

    nextNum := int(maxNum.Int64) + 1
    return fmt.Sprintf("%s%08d", prefix, nextNum), nil
}
```

### 2. Limpiar Facturas Duplicadas

Si ya tienes facturas duplicadas en tu base de datos, necesitas limpiarlas.

#### Opción A: Script Automático

```bash
cd /home/scorpion/desa/Projects/vendix
./scripts/fix-invoice-numbers.sh
```

Este script te mostrará:
- Facturas existentes
- Números duplicados
- Instrucciones para limpiar

#### Opción B: Manualmente con PostgreSQL

```bash
# Conectar a la base de datos
psql vendix -U vendix

# Cambiar al schema del tenant (ejemplo: demo)
SET search_path TO demo;

# Ver facturas existentes
SELECT id, invoice_number, total, created_at 
FROM invoices 
ORDER BY created_at DESC;

# Si hay duplicados, eliminar TODAS las facturas (SOLO EN DESARROLLO)
DELETE FROM invoice_lines;
DELETE FROM invoices;

# Salir
\q
```

⚠️ **ADVERTENCIA**: Esto eliminará todas las facturas. Solo hazlo en ambiente de desarrollo.

#### Opción C: Eliminar Solo las Duplicadas

Si solo quieres eliminar las facturas duplicadas manteniendo la primera:

```sql
SET search_path TO demo;

-- Ver duplicados
SELECT invoice_number, COUNT(*) as count 
FROM invoices 
GROUP BY invoice_number 
HAVING COUNT(*) > 1;

-- Eliminar duplicados manteniendo el más antiguo
DELETE FROM invoice_lines 
WHERE invoice_id IN (
    SELECT id FROM (
        SELECT id, 
               ROW_NUMBER() OVER (PARTITION BY invoice_number ORDER BY created_at) as rn
        FROM invoices
    ) t
    WHERE t.rn > 1
);

DELETE FROM invoices 
WHERE id IN (
    SELECT id FROM (
        SELECT id, 
               ROW_NUMBER() OVER (PARTITION BY invoice_number ORDER BY created_at) as rn
        FROM invoices
    ) t
    WHERE t.rn > 1
);
```

### 3. Reiniciar el Backend

Después de hacer los cambios en el código, **DEBES REINICIAR** el backend:

```bash
# Detener el backend actual (Ctrl+C en la terminal donde está corriendo)

# Recompilar y ejecutar
cd /home/scorpion/desa/Projects/vendix
make dev

# O si usas go run directamente
go run cmd/api/main.go
```

## Pasos para Resolver (Resumen)

### Paso 1: Verificar el Estado Actual
```bash
cd /home/scorpion/desa/Projects/vendix
./scripts/fix-invoice-numbers.sh
```

### Paso 2: Limpiar Facturas (si es necesario)
```bash
# Solo en desarrollo - esto borra todas las facturas
psql vendix -U vendix -c "SET search_path TO demo; DELETE FROM invoice_lines; DELETE FROM invoices;"
```

### Paso 3: Reiniciar el Backend
```bash
# Detener el backend actual (Ctrl+C)
cd /home/scorpion/desa/Projects/vendix
make dev
```

### Paso 4: Probar de Nuevo
1. Abre el frontend: http://localhost:5173
2. Ve a "Ventas"
3. Agrega productos y completa una venta
4. ✅ Debería funcionar correctamente ahora

## Verificación

Después de aplicar la solución, verifica que:

- [ ] El backend se reinició correctamente
- [ ] No hay facturas duplicadas en la base de datos
- [ ] Puedes crear una nueva venta sin errores
- [ ] Los números de factura son secuenciales (INV-00000001, INV-00000002, etc.)

## Prueba de Concepto

```bash
# Crear 3 ventas seguidas y verificar los números
# Deberías ver:
# - Primera venta: INV-00000001
# - Segunda venta: INV-00000002
# - Tercera venta: INV-00000003
```

## Para Producción

⚠️ **Importante**: En producción, NO elimines facturas. En su lugar:

1. **Identifica la factura problemática**
2. **Márcala como cancelada** (no la borres)
3. **Crea una nueva factura** con el siguiente número disponible
4. **Registra el incidente** para auditoría

## Mejoras Futuras

Para evitar este problema completamente en el futuro, considera:

### Opción 1: Usar Secuencias de PostgreSQL
```sql
-- Crear una secuencia por tenant
CREATE SEQUENCE demo.invoice_number_seq START 1;

-- Usar en la aplicación
SELECT nextval('demo.invoice_number_seq');
```

### Opción 2: Advisory Locks
```go
// Usar un lock de PostgreSQL para evitar condiciones de carrera
_, err := tx.Exec("SELECT pg_advisory_xact_lock($1)", hashString("invoice_number"))
```

### Opción 3: Tabla de Contadores
```sql
-- Tabla para mantener contadores
CREATE TABLE demo.counters (
    name VARCHAR(50) PRIMARY KEY,
    value INTEGER NOT NULL DEFAULT 1
);

-- Incrementar atómicamente
UPDATE demo.counters 
SET value = value + 1 
WHERE name = 'invoice_number' 
RETURNING value;
```

## Soporte

Si el problema persiste:

1. **Verifica los logs del backend**:
   ```bash
   tail -f /home/scorpion/desa/Projects/vendix/tmp/app.log
   ```

2. **Verifica el estado de la base de datos**:
   ```bash
   psql vendix -U vendix -c "SET search_path TO demo; SELECT COUNT(*) FROM invoices;"
   ```

3. **Verifica que el código se recompiló**:
   ```bash
   # El backend debe mostrar algo como:
   # INFO Starting server on port 8080
   ```

## Resumen

✅ **Código corregido** - La función ahora usa una consulta más robusta  
✅ **Script de limpieza** - Puedes verificar y limpiar duplicados  
🔄 **Reinicia el backend** - IMPORTANTE para que los cambios tomen efecto  
✅ **Prueba de nuevo** - Crea una venta y verifica que funcione  

---

**Fecha**: 19 de octubre de 2025  
**Estado**: ✅ Solucionado  
**Versión**: 1.1

