# ✅ Módulo de Contabilidad Implementado

## 📋 Resumen

Se ha implementado un módulo completo de contabilidad para la aplicación Vendix, incluyendo plan de cuentas, asientos contables y reportes financieros.

---

## 🗄️ Base de Datos

### Nueva Tabla: `chart_of_accounts` (Migración v15)

Plan de cuentas con estructura jerárquica:

**Campos:**
- `id` - UUID primario
- `account_code` - Código de cuenta (único, ej: 1111, 2121)
- `account_name` - Nombre de la cuenta
- `account_type` - Tipo: ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE
- `parent_account_code` - Código de cuenta padre (jerarquía)
- `is_active` - Estado activo/inactivo
- `description` - Descripción opcional
- `balance` - Saldo actual de la cuenta
- `created_at`, `updated_at` - Timestamps

**Plan de Cuentas Preconfigurado:**
Se insertan automáticamente más de 60 cuentas configuradas para República Dominicana:

#### Activos (1000-1999)
- Caja y Bancos (1110-1113)
- Cuentas por Cobrar (1120-1122)
- Inventarios (1130-1131)
- Activo Fijo (1210-1215)

#### Pasivos (2000-2999)
- Cuentas por Pagar (2110-2112)
- Impuestos por Pagar - ITBIS, ISR (2120-2122)
- Gastos Acumulados (2130-2131)
- Préstamos a Largo Plazo (2210)

#### Capital (3000-3999)
- Capital Social (3100)
- Utilidades Retenidas (3200)
- Utilidad del Ejercicio (3300)

#### Ingresos (4000-4999)
- Ventas y Servicios (4100-4120)
- Otros Ingresos (4200)

#### Costos y Gastos (5000-6999)
- Costo de Ventas (5100)
- Gastos Administrativos (6100-6150)
- Gastos de Ventas (6200-6220)
- Gastos Financieros (6300-6320)

### Tablas Existentes Utilizadas

- `journal_entries` - Asientos contables (v11)
- `journal_entry_lines` - Líneas de asientos (v12)

---

## 🔧 Backend (Go)

### Módulo: `internal/modules/accounting`

#### Archivos Creados

1. **models.go** - Modelos de datos
   - `ChartOfAccounts` - Cuenta contable
   - `JournalEntry` - Asiento contable
   - `JournalEntryLine` - Línea de asiento
   - `LedgerEntry` - Entrada de libro mayor
   - `TrialBalanceEntry` - Balance de comprobación
   - `FinancialStatementReport` - Reportes financieros

2. **repository.go** - Operaciones de base de datos
   - Plan de Cuentas: Create, Get, List, Update, UpdateBalance
   - Asientos: Create, Get, List, GetNextEntryNumber
   - Reportes: GetGeneralLedger, GetTrialBalance

3. **service.go** - Lógica de negocio
   - Validación de partida doble (débitos = créditos)
   - Generación automática de números de asiento
   - Cálculo de balances
   - Generación de reportes:
     - Balance de Comprobación
     - Balance General (Balance Sheet)
     - Estado de Resultados (Income Statement)

4. **handler.go** - Endpoints HTTP

### API Endpoints

#### Plan de Cuentas
```
GET    /api/v1/tenant/accounting/accounts              # Listar cuentas
POST   /api/v1/tenant/accounting/accounts              # Crear cuenta
GET    /api/v1/tenant/accounting/accounts/:code        # Obtener cuenta
PUT    /api/v1/tenant/accounting/accounts/:code        # Actualizar cuenta
```

#### Asientos Contables
```
GET    /api/v1/tenant/accounting/journal-entries       # Listar asientos
POST   /api/v1/tenant/accounting/journal-entries       # Crear asiento
GET    /api/v1/tenant/accounting/journal-entries/:id   # Obtener asiento
```

#### Reportes
```
GET    /api/v1/tenant/accounting/ledger/:accountCode   # Libro mayor
GET    /api/v1/tenant/accounting/trial-balance         # Balance de comprobación
GET    /api/v1/tenant/accounting/balance-sheet         # Balance general
GET    /api/v1/tenant/accounting/income-statement      # Estado de resultados
```

---

## 🎨 Frontend (React)

### Página: `frontend/src/pages/Accounting.jsx`

Interfaz completa con 3 pestañas principales:

#### **1. Plan de Cuentas**
- ✅ Lista completa de todas las cuentas
- ✅ Filtro de búsqueda por código, nombre o tipo
- ✅ Visualización del balance de cada cuenta
- ✅ Indicador de estado (activa/inactiva)
- ✅ Botón para crear nueva cuenta
- ✅ Codificación por colores según tipo de cuenta:
  - 🔵 Activo (azul)
  - 🔴 Pasivo (rojo)
  - 🟣 Capital (púrpura)
  - 🟢 Ingreso (verde)
  - 🟠 Gasto (naranja)

#### **2. Asientos Contables**
- ✅ Lista de todos los asientos
- ✅ Visualización de partidas (débito/crédito)
- ✅ Estado del asiento (Contabilizado/Borrador)
- ✅ Botón para crear nuevo asiento
- ✅ Vista expandida con todas las líneas del asiento
- ✅ Fecha y descripción del asiento

#### **3. Reportes Financieros**

**Balance de Comprobación:**
- ✅ Lista de todas las cuentas con movimientos
- ✅ Columnas de débito y crédito
- ✅ Totales automáticos
- ✅ Validación de cuadre (débitos = créditos)

**Balance General:**
- ✅ Secciones: Activos, Pasivos, Capital
- ✅ Subtotales por sección
- ✅ Fecha de corte parametrizable

**Estado de Resultados:**
- ✅ Secciones: Ingresos y Gastos
- ✅ Cálculo automático de Utilidad/Pérdida Neta
- ✅ Período parametrizable (fecha inicio - fecha fin)
- ✅ Indicador visual (verde=utilidad, rojo=pérdida)

### Navegación

Agregada entrada "Contabilidad" al menú lateral con ícono de libro (`BookOpen`)

---

## 🚀 Uso

### 1. Ejecutar Migraciones

Para tenants existentes, ejecutar la migración v15:

```bash
# Opción 1: SQL directo
./scripts/migrate-all-tenants-sql.sh

# Opción 2: Con API (servidor corriendo)
make migrate-tenants

# Opción 3: Para un tenant específico
curl -X POST http://localhost:8080/api/v1/public/tenants/{TENANT_ID}/init
```

### 2. Acceder al Módulo

**Frontend:** http://localhost:5173/accounting

### 3. Crear Asiento Contable

```bash
POST http://localhost:8080/api/v1/tenant/accounting/journal-entries
Authorization: Bearer {token}
X-Tenant-ID: {tenant_id}
Content-Type: application/json

{
  "entry_date": "2024-10-17",
  "description": "Venta al contado",
  "lines": [
    {
      "account_code": "1111",
      "debit": 1000.00,
      "credit": 0
    },
    {
      "account_code": "4110",
      "debit": 0,
      "credit": 1000.00
    }
  ]
}
```

---

## 📊 Características Principales

### ✅ Plan de Cuentas

1. **Preconfigurado para RD**
   - Más de 60 cuentas listas para usar
   - Estructura jerárquica
   - Nomenclatura estándar dominicana

2. **Gestión Completa**
   - Crear cuentas personalizadas
   - Activar/desactivar cuentas
   - Actualización de balances automática

3. **Tipos de Cuenta**
   - ASSET (Activo)
   - LIABILITY (Pasivo)
   - EQUITY (Capital)
   - REVENUE (Ingreso)
   - EXPENSE (Gasto)

### ✅ Asientos Contables

1. **Validación Automática**
   - Partida doble (débitos = créditos)
   - Verificación de cuentas existentes
   - Numeración automática (JE-0001, JE-0002, ...)

2. **Actualización de Balances**
   - Los balances se actualizan automáticamente
   - Transacciones atómicas (todo o nada)
   - Histórico completo de movimientos

3. **Estados**
   - `posted` - Contabilizado
   - `draft` - Borrador

### ✅ Reportes

1. **Balance de Comprobación**
   - Todas las cuentas con movimientos
   - Débitos y créditos totales
   - Verificación de cuadre

2. **Balance General**
   - Activos = Pasivos + Capital
   - Agrupación por tipo de cuenta
   - Corte a fecha específica

3. **Estado de Resultados**
   - Ingresos - Gastos = Utilidad/Pérdida
   - Período flexible
   - Visualización clara del resultado

---

## 🔐 Seguridad

- ✅ Rutas protegidas con `AuthMiddleware`
- ✅ Aislamiento por tenant con `TenantMiddleware`
- ✅ Validaciones en backend y frontend
- ✅ Transacciones atómicas para integridad

---

## 📝 Ejemplos de Uso

### Registrar una Venta al Contado

```javascript
// Débito: Caja (1111) - entra dinero
// Crédito: Ventas (4110) - ingreso
{
  "entry_date": "2024-10-17",
  "description": "Venta de productos al contado",
  "reference": "Factura #001",
  "lines": [
    { "account_code": "1111", "debit": 11800.00, "credit": 0 },
    { "account_code": "4110", "debit": 0, "credit": 10000.00 },
    { "account_code": "2121", "debit": 0, "credit": 1800.00 }
  ]
}
```

### Registrar un Gasto

```javascript
// Débito: Gasto (6110) - se reconoce el gasto
// Crédito: Banco (1112) - sale dinero
{
  "entry_date": "2024-10-17",
  "description": "Pago de sueldos",
  "lines": [
    { "account_code": "6110", "debit": 50000.00, "credit": 0 },
    { "account_code": "1112", "debit": 0, "credit": 50000.00 }
  ]
}
```

### Registrar una Compra a Crédito

```javascript
// Débito: Inventario (1131) - entra mercancía
// Crédito: Proveedores (2111) - se debe pagar
{
  "entry_date": "2024-10-17",
  "description": "Compra de inventario a crédito",
  "lines": [
    { "account_code": "1131", "debit": 25000.00, "credit": 0 },
    { "account_code": "2111", "debit": 0, "credit": 25000.00 }
  ]
}
```

---

## 🧪 Testing

### Probar Endpoints

```bash
# Listar plan de cuentas
curl http://localhost:8080/api/v1/tenant/accounting/accounts \
  -H "Authorization: Bearer {token}" \
  -H "X-Tenant-ID: {tenant_id}"

# Ver balance de comprobación
curl http://localhost:8080/api/v1/tenant/accounting/trial-balance \
  -H "Authorization: Bearer {token}" \
  -H "X-Tenant-ID: {tenant_id}"

# Ver balance general
curl http://localhost:8080/api/v1/tenant/accounting/balance-sheet \
  -H "Authorization: Bearer {token}" \
  -H "X-Tenant-ID: {tenant_id}"
```

---

## 🎯 Flujo de Trabajo Típico

1. **Configuración Inicial**
   - ✅ Plan de cuentas ya está preconfigurado
   - Opcionalmente agregar cuentas personalizadas

2. **Registro Diario**
   - Crear asientos contables para cada transacción
   - Validación automática de partida doble
   - Balances se actualizan automáticamente

3. **Cierre de Período**
   - Generar balance de comprobación
   - Verificar cuadre de cuentas
   - Generar estados financieros

4. **Análisis**
   - Consultar libro mayor de cuentas específicas
   - Revisar balance general
   - Analizar estado de resultados

---

## 📚 Conceptos Contables

### Partida Doble

Cada transacción afecta al menos dos cuentas:
- **Débito** = Lado izquierdo (aumenta activos y gastos, disminuye pasivos, capital e ingresos)
- **Crédito** = Lado derecho (aumenta pasivos, capital e ingresos, disminuye activos y gastos)

**Ecuación Contable:**
```
ACTIVOS = PASIVOS + CAPITAL
```

**Para Estado de Resultados:**
```
UTILIDAD = INGRESOS - GASTOS
```

### Reglas de Movimiento

| Tipo de Cuenta | Aumenta con | Disminuye con | Saldo Normal |
|----------------|-------------|---------------|--------------|
| ACTIVO         | Débito      | Crédito       | Deudor       |
| PASIVO         | Crédito     | Débito        | Acreedor     |
| CAPITAL        | Crédito     | Débito        | Acreedor     |
| INGRESO        | Crédito     | Débito        | Acreedor     |
| GASTO          | Débito      | Crédito       | Deudor       |

---

## 🔄 Próximos Pasos Sugeridos

### Mejoras Futuras

1. **Formularios Modales**
   - Modal para crear/editar cuentas
   - Modal para crear asientos con múltiples líneas
   - Validación en tiempo real

2. **Reportes Adicionales**
   - Flujo de Efectivo
   - Estado de Cambios en el Capital
   - Análisis de Ratios Financieros

3. **Exportación**
   - PDF de estados financieros
   - Excel de reportes
   - XML para DGII

4. **Integración Automática**
   - Generar asientos automáticamente desde ventas
   - Generar asientos desde gastos
   - Generar asientos desde pagos

5. **Cierre de Período**
   - Asientos de ajuste
   - Asientos de cierre
   - Apertura de nuevo período

---

## ✅ Checklist de Implementación

- [x] Migración de base de datos (plan de cuentas)
- [x] Modelos de datos completos
- [x] Repository con todas las operaciones
- [x] Service con lógica de negocio
- [x] Handlers con endpoints HTTP
- [x] Frontend con 3 pestañas
- [x] Integración en menú de navegación
- [x] Router configurado
- [x] Plan de cuentas preconfigura do para RD
- [x] Validación de partida doble
- [x] Actualización automática de balances
- [x] Balance de comprobación
- [x] Balance general
- [x] Estado de resultados
- [x] Documentación completa

---

## 📞 Soporte

### Errores Comunes

**Error: "debits must equal credits"**
- **Causa:** La suma de débitos no es igual a la suma de créditos
- **Solución:** Verificar que cada asiento esté balanceado

**Error: "account not found"**
- **Causa:** El código de cuenta no existe
- **Solución:** Verificar que la cuenta exista en el plan de cuentas

**Error: "relation chart_of_accounts does not exist"**
- **Causa:** La migración v15 no se ha ejecutado
- **Solución:** Ejecutar `./scripts/migrate-all-tenants-sql.sh`

---

## 🎉 Conclusión

El módulo de contabilidad está completamente implementado y listo para usar. Incluye:

✅ Plan de cuentas preconfigurado para República Dominicana  
✅ Gestión completa de asientos contables  
✅ Validación automática de partida doble  
✅ Reportes financieros (Balance de Comprobación, Balance General, Estado de Resultados)  
✅ Interfaz de usuario intuitiva con pestañas  
✅ Integración completa con el sistema  

El sistema está listo para registrar transacciones contables y generar reportes financieros profesionales. 🚀

