# 🌱 Datos de Ejemplo - Vendix

## 📋 Resumen

Se han generado datos de ejemplo completos para probar todas las funcionalidades del sistema Vendix. Los datos incluyen clientes, productos, facturas y gastos con diferentes estados y escenarios realistas.

---

## 🚀 Cómo Cargar los Datos

### Ejecución Simple

```bash
./scripts/seed-demo-data.sh
```

### Ejecución Manual

```bash
# Copiar el archivo SQL al contenedor
docker cp scripts/seed-demo-data.sql vendix-postgres-dev:/tmp/seed-demo-data.sql

# Ejecutar el script
docker exec -i vendix-postgres-dev psql -U postgres -d vendix -f /tmp/seed-demo-data.sql

# Limpiar
docker exec vendix-postgres-dev rm /tmp/seed-demo-data.sql
```

---

## 📊 Datos Generados

### 👥 Clientes (9 en total)

#### Individuales (5):
1. **Juan Pérez**
   - RNC: 001-1234567-8
   - Email: juan.perez@email.com
   - Teléfono: 809-555-0101

2. **María García**
   - RNC: 001-9876543-2
   - Email: maria.garcia@email.com
   - Teléfono: 809-555-0102

3. **Carlos Rodríguez**
   - RNC: 001-5555666-7
   - Email: carlos.rod@email.com
   - Teléfono: 809-555-0103

4. **Ana Martínez**
   - RNC: 001-3332221-1
   - Email: ana.martinez@email.com
   - Teléfono: 809-555-0104

#### Empresas (4):
5. **Supermercado La Económica SRL**
   - RNC: 130-56789-0
   - Email: contacto@laeconomica.do
   - Teléfono: 809-555-0201

6. **Ferretería El Constructor SA**
   - RNC: 130-11223-3
   - Email: ventas@elconstructor.do
   - Teléfono: 809-555-0202

7. **Restaurante Don Pepe**
   - RNC: 131-99887-7
   - Email: admin@donpepe.do
   - Teléfono: 809-555-0301

8. **Farmacia Moderna**
   - RNC: 130-77665-5
   - Email: info@farmaciamoderna.do
   - Teléfono: 809-555-0401

9. **Gastos Generales** (cliente interno para gastos)
   - RNC: 999-9999999-9

---

### 🛍️ Productos (13 en total)

#### Electrónica (5):
- **LAPTOP-001**: Laptop Dell XPS 15 - RD$ 75,000.00
- **LAPTOP-002**: Laptop HP Pavilion - RD$ 45,000.00
- **MONITOR-001**: Monitor LG 27" - RD$ 18,000.00
- **MOUSE-001**: Mouse Logitech MX Master - RD$ 4,500.00
- **KEYBOARD-001**: Teclado Mecánico Corsair - RD$ 8,500.00

#### Papelería (3):
- **PAPER-001**: Resma de Papel A4 - RD$ 350.00
- **PEN-001**: Bolígrafos BIC (Caja 50) - RD$ 450.00
- **FOLDER-001**: Folder Manila (Paquete 25) - RD$ 280.00

#### Servicios (3):
- **SERV-001**: Consultoría IT (Hora) - RD$ 3,500.00
- **SERV-002**: Desarrollo Web - RD$ 85,000.00
- **SERV-003**: Soporte Técnico Mensual - RD$ 15,000.00

#### Alimentos (2):
- **FOOD-001**: Café Premium (Libra) - RD$ 450.00
- **FOOD-002**: Galletas Surtidas - RD$ 185.00

---

### 📝 Facturas (7 en total)

#### 1. INV-00000001 - **PAGADA** ✅
- **Cliente**: Juan Pérez
- **Monto**: RD$ 182,310.00
- **Fecha**: Hace 30 días
- **Items**:
  - 2x Laptop Dell XPS 15
  - 1x Mouse Logitech MX Master

#### 2. INV-00000002 - **ENVIADA** 📤
- **Cliente**: Supermercado La Económica SRL
- **Monto**: RD$ 41,300.00
- **Fecha**: Hace 15 días
- **Items**:
  - 10x Consultoría IT (horas)

#### 3. INV-00000003 - **PENDIENTE** ⏳
- **Cliente**: Restaurante Don Pepe
- **Monto**: RD$ 100,300.00
- **Fecha**: Hace 10 días
- **Items**:
  - 1x Desarrollo Web Corporativo

#### 4. INV-00000004 - **PAGADA** ✅
- **Cliente**: María García
- **Monto**: RD$ 21,240.00
- **Fecha**: Hace 5 días
- **Items**:
  - 1x Monitor LG 27"

#### 5. INV-00000005 - **PENDIENTE** ⏳
- **Cliente**: Ferretería El Constructor SA
- **Monto**: RD$ 66,080.00
- **Fecha**: Hace 3 días
- **Items**:
  - 5x Teclado Mecánico Corsair
  - 1x Mouse Logitech MX Master

#### 6. INV-00000006 - **BORRADOR** 📄
- **Cliente**: Supermercado La Económica SRL
- **Monto**: RD$ 17,700.00
- **Fecha**: Hoy
- **Items**:
  - 1x Plan de Soporte Técnico Mensual

#### 7. INV-00000007 - **ANULADA** ❌
- **Cliente**: Juan Pérez
- **Monto**: RD$ 88,500.00
- **Fecha**: Hace 20 días
- **Motivo**: Error en especificaciones
- **Items**:
  - 1x Laptop Dell XPS 15

---

### 💸 Gastos (10 en total)

| # | Referencia | Descripción | Método | Monto | Fecha |
|---|------------|-------------|--------|-------|-------|
| 1 | EDESUR-OCT2025 | Electricidad - Octubre | Efectivo | RD$ 8,900.00 | Hace 25 días |
| 2 | INV-PAP-001 | Papelería y útiles | Transferencia | RD$ 5,500.00 | Hace 20 días |
| 3 | META-ADS-789 | Campaña redes sociales | Tarjeta | RD$ 15,000.00 | Hace 18 días |
| 4 | ALQ-OCT-2025 | Alquiler de oficina | Transferencia | RD$ 35,000.00 | Hace 15 días |
| 5 | GASOLINA-001 | Combustible vehículo | Efectivo | RD$ 2,800.00 | Hace 12 días |
| 6 | AGUA-OCT | Servicio de agua | Transferencia | RD$ 4,500.00 | Hace 10 días |
| 7 | INTERNET-FIB | Internet fibra óptica | Tarjeta | RD$ 12,000.00 | Hace 8 días |
| 8 | LIMPIEZA-OCT | Servicio de limpieza | Cheque | RD$ 8,500.00 | Hace 5 días |
| 9 | TELEFONIA | Servicios telefónicos | Transferencia | RD$ 3,200.00 | Hace 3 días |
| 10 | CAFE-SUP | Café y snacks oficina | Efectivo | RD$ 1,800.00 | Hace 1 día |

**Total de Gastos**: RD$ 97,200.00

---

## 📈 Resumen Financiero

### Facturas

| Estado | Cantidad | Total |
|--------|----------|-------|
| Pagadas | 2 | RD$ 203,550.00 |
| Pendientes | 2 | RD$ 166,380.00 |
| Enviadas | 1 | RD$ 41,300.00 |
| Borradores | 1 | RD$ 17,700.00 |
| Anuladas | 1 | RD$ 88,500.00 |
| **TOTAL** | **7** | **RD$ 517,430.00** |

### Ingresos vs Gastos

- **Ingresos Cobrados**: RD$ 203,550.00 (facturas pagadas)
- **Gastos**: RD$ 97,200.00
- **Balance**: RD$ 106,350.00 ✅

- **Por Cobrar**: RD$ 207,680.00 (pendientes + enviadas)
- **Potencial Total**: RD$ 314,030.00

---

## 🧪 Verificación de Datos

### Via API

```bash
# Obtener token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/tenant/auth/login \
  -H 'Content-Type: application/json' \
  -H 'X-Tenant-ID: demo' \
  -d '{"email":"admin@demo.com","password":"demo123"}' \
  | jq -r '.access_token')

# Verificar clientes
curl -s http://localhost:8080/api/v1/tenant/customers \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: demo" | jq 'length'
# → 9

# Verificar productos
curl -s http://localhost:8080/api/v1/tenant/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: demo" | jq 'length'
# → 13

# Verificar facturas
curl -s http://localhost:8080/api/v1/tenant/invoices \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: demo" | jq 'length'
# → 7

# Verificar gastos
curl -s http://localhost:8080/api/v1/tenant/expenses \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: demo" | jq 'length'
# → 10
```

### Via Frontend

Visita estas URLs para ver los datos:

- **Clientes**: http://localhost:5173/customers
- **Productos**: http://localhost:5173/products
- **Facturas**: http://localhost:5173/invoices
- **Gastos**: http://localhost:5173/expenses
- **Punto de Venta**: http://localhost:5173/sales

---

## 🔄 Regenerar Datos

Si necesitas regenerar los datos (esto borrará los datos existentes):

```bash
./scripts/seed-demo-data.sh
```

**⚠️ Advertencia**: Este script elimina todos los datos existentes del tenant demo antes de insertar los datos de ejemplo.

---

## 🎯 Casos de Uso para Pruebas

### 1. Flujo de Facturación Completa

1. Ir a **/sales** (Punto de Venta)
2. Seleccionar cliente: **Juan Pérez**
3. Agregar productos al carrito
4. Procesar venta
5. Ver factura generada en **/invoices**
6. Enviar factura a DGII
7. Marcar como pagada

### 2. Gestión de Gastos

1. Ir a **/expenses**
2. Clic en "Nuevo Gasto"
3. Completar formulario
4. Guardar
5. Verificar en lista de gastos

### 3. Reportes Financieros

1. Ver dashboard en **/** (home)
2. Comparar ingresos vs gastos
3. Ver gráficas de ventas
4. Analizar tendencias

### 4. Gestión de Inventario

1. Ir a **/products**
2. Verificar stock disponible
3. Procesar ventas
4. Ver actualización de stock en tiempo real

### 5. Contabilidad

1. Ir a **/accounting**
2. Ver plan de cuentas
3. Ver asientos contables generados por ventas
4. Generar balance de comprobación
5. Ver estados financieros

---

## 📝 Notas Adicionales

- Todos los productos tienen ITBIS del 18%
- Las fechas son relativas al día actual
- Los montos están en pesos dominicanos (DOP)
- Los clientes tienen RNC válidos según formato RD
- Los números de teléfono usan el formato dominicano (809)

---

## 🐛 Solución de Problemas

### Los datos no aparecen en el frontend

1. Verificar que el backend está corriendo:
   ```bash
   curl http://localhost:8080/health
   ```

2. Verificar autenticación:
   ```bash
   # Login
   curl -X POST http://localhost:8080/api/v1/tenant/auth/login \
     -H 'Content-Type: application/json' \
     -H 'X-Tenant-ID: demo' \
     -d '{"email":"admin@demo.com","password":"demo123"}'
   ```

3. Revisar console del navegador (F12)

### Error al ejecutar script

```bash
# Verificar que el contenedor está corriendo
docker ps | grep vendix-postgres-dev

# Verificar conexión a la base de datos
docker exec vendix-postgres-dev psql -U postgres -d vendix -c "\dn"
```

---

**Fecha de creación**: 17 de octubre de 2025  
**Versión**: 1.0

