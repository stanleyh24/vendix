# ✅ Tenant Demo Recreado

## 🎉 El tenant demo ha sido recreado exitosamente

---

## 📋 Información del Tenant

**Detalles:**
- 🏢 **Nombre:** Demo Company
- 🔖 **Slug:** demo
- 🗄️ **Schema:** tenant_demo
- 📊 **ID:** 1ddd937e-75bf-450e-b4f2-02173edae5ba
- ✅ **Estado:** Activo

---

## 👤 Credenciales de Acceso

### Usuario Administrador

```
Email: admin@demo.com
Password: demo123
Tenant ID: demo
```

---

## ✅ Lo que se Creó

### 1. Todas las Tablas (16 migraciones)

- ✅ users
- ✅ refresh_tokens
- ✅ customers
- ✅ products
- ✅ tax_rates
- ✅ invoices
- ✅ invoice_lines
- ✅ payments
- ✅ payment_allocations
- ✅ credit_notes
- ✅ journal_entries
- ✅ journal_entry_lines
- ✅ audit_logs
- ✅ settings
- ✅ **chart_of_accounts** (57 cuentas)
- ✅ **tenant_config**

### 2. Datos Iniciales

**Plan de Cuentas:** 57 cuentas preconfiguradas
- 18 cuentas de Activos (1000-1999)
- 10 cuentas de Pasivos (2000-2999)
- 4 cuentas de Capital (3000-3999)
- 5 cuentas de Ingresos (4000-4999)
- 20 cuentas de Gastos (5000-6999)

**Tasas de Impuesto:**
- ITBIS 18%
- ITBIS 16%
- Exento

**Usuario:** 1 administrador

**Configuración:** 1 registro con valores por defecto

---

## 🚀 Cómo Usar

### 1. Iniciar Sesión

```
1. Ve a: http://localhost:5173/login
2. Ingresa:
   - Email: admin@demo.com
   - Password: demo123
3. Click en "Iniciar Sesión"
```

### 2. Explorar Módulos

Una vez dentro, puedes acceder a:

- **Dashboard** - `/` - Vista general
- **Ventas** - `/sales` - Registro de ventas
- **Clientes** - `/customers` - Gestión de clientes
- **Inventario** - `/products` - Productos y servicios
- **Gastos** - `/expenses` - Control de gastos
- **Facturación** - `/invoices` - Facturación electrónica
- **Contabilidad** - `/accounting` - Plan de cuentas y reportes ← **NUEVO**
- **Configuración** - `/settings` - Configuración del tenant ← **NUEVO**

---

## 📊 Plan de Cuentas Incluido

### Activos (1000-1999)
```
1000 - ACTIVOS
1100 - Activo Circulante
1110 - Caja y Bancos
1111 - Caja General
1112 - Banco - Cuenta Corriente
1113 - Banco - Cuenta de Ahorros
1120 - Cuentas por Cobrar
1121 - Clientes
1122 - Documentos por Cobrar
1130 - Inventarios
1131 - Inventario de Mercancías
1200 - Activo Fijo
1210 - Propiedad, Planta y Equipo
1211 - Edificios
1212 - Maquinaria y Equipo
1213 - Equipo de Oficina
1214 - Vehículos
1215 - Equipo de Computación
```

### Pasivos (2000-2999)
```
2000 - PASIVOS
2100 - Pasivo Circulante
2110 - Cuentas por Pagar
2111 - Proveedores
2112 - Documentos por Pagar
2120 - Impuestos por Pagar
2121 - ITBIS por Pagar
2122 - ISR por Pagar
2130 - Gastos Acumulados
2131 - Sueldos por Pagar
2200 - Pasivo a Largo Plazo
2210 - Préstamos Bancarios
```

### Capital (3000-3999)
```
3000 - CAPITAL
3100 - Capital Social
3200 - Utilidades Retenidas
3300 - Utilidad del Ejercicio
```

### Ingresos (4000-4999)
```
4000 - INGRESOS
4100 - Ingresos por Ventas
4110 - Ventas
4120 - Servicios
4200 - Otros Ingresos
```

### Gastos (5000-6999)
```
5000 - COSTOS
5100 - Costo de Ventas
6000 - GASTOS
6100 - Gastos de Administración
6110 - Sueldos y Salarios
6120 - Alquileres
6130 - Servicios Públicos
6131 - Energía Eléctrica
6132 - Agua
6133 - Teléfono e Internet
6140 - Papelería y Útiles
6150 - Combustible y Lubricantes
6200 - Gastos de Ventas
6210 - Publicidad y Mercadeo
6220 - Comisiones de Ventas
6300 - Gastos Financieros
6310 - Intereses Bancarios
6320 - Comisiones Bancarias
```

---

## 🔄 Recrear en el Futuro

Si necesitas recrear el tenant demo nuevamente:

```bash
# Opción 1: Con make
make recreate-demo

# Opción 2: Directamente
./scripts/recreate-demo-tenant.sh
```

---

## 🧪 Verificación

### Verificar en el Frontend:

1. **Login:**
   - http://localhost:5173/login
   - Email: admin@demo.com
   - Password: demo123

2. **Contabilidad:**
   - http://localhost:5173/accounting
   - Deberías ver 57 cuentas en "Plan de Cuentas"

3. **Configuración:**
   - http://localhost:5173/settings
   - Puedes configurar los datos de la empresa

### Verificar con API:

```bash
# Login
curl -X POST http://localhost:8080/api/v1/tenant/auth/login \
  -H 'Content-Type: application/json' \
  -H 'X-Tenant-ID: demo' \
  -d '{"email":"admin@demo.com","password":"demo123"}'

# Copiar el access_token y usarlo:
curl http://localhost:8080/api/v1/tenant/accounting/accounts \
  -H "Authorization: Bearer {TOKEN}" \
  -H "X-Tenant-ID: demo" \
  | jq '.[0:5]'  # Muestra las primeras 5 cuentas
```

---

## 📚 Módulos Disponibles

Con el tenant demo recién creado, puedes probar:

### ✅ Contabilidad
- Ver plan de cuentas (57 cuentas)
- Crear asientos contables
- Generar reportes:
  - Balance de Comprobación
  - Balance General
  - Estado de Resultados

### ✅ Configuración
- Datos de la empresa
- Configuración fiscal (DGII, NCF)
- Preferencias de facturación
- Branding (colores, logo)

### ✅ Otros Módulos
- Clientes
- Productos
- Ventas
- Gastos
- Facturación
- Reportes

---

## 🎯 Próximos Pasos

1. **Inicia sesión** con las credenciales arriba
2. **Configura tu empresa** en `/settings`
3. **Agrega clientes** en `/customers`
4. **Agrega productos** en `/products`
5. **Registra ventas** en `/sales`
6. **Consulta contabilidad** en `/accounting`

---

## 🛠️ Comando Creado

Se agregó al Makefile:

```bash
make recreate-demo
```

Este comando:
1. Elimina el tenant demo existente (si existe)
2. Crea un tenant nuevo
3. Ejecuta todas las migraciones
4. Inserta el plan de cuentas
5. Crea el usuario administrador
6. Muestra las credenciales

---

## ✅ Todo Listo

El tenant demo está completamente configurado y listo para usar. 

**¡Ve a http://localhost:5173/login e inicia sesión!** 🚀

**Credenciales:**
- Email: `admin@demo.com`
- Password: `demo123`
- Tenant: `demo`

El módulo de contabilidad con las 57 cuentas ya está funcionando. 🎉
