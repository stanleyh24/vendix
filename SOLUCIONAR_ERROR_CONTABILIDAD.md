# 🔧 Solucionar Error al Cargar Cuentas

## ❌ Problema

Al acceder a la página de Contabilidad, aparece un error al intentar cargar las cuentas.

---

## 🔍 Diagnóstico

El error más común es que **la tabla `chart_of_accounts` no existe** en el schema del tenant.

### Paso 1: Verificar si la Tabla Existe

Ejecuta el script de verificación:

```bash
./scripts/check-accounting-tables.sh
```

**Salida esperada si HAY problema:**
```
❌ tenant_demo - Tabla NO existe (necesita migración)

⚠️  ACCIÓN REQUERIDA:
   Ejecuta: ./scripts/migrate-all-tenants-sql.sh
```

**Salida esperada si TODO está bien:**
```
✅ tenant_demo - Tabla existe con 60+ cuentas
✅ Todos los tenants tienen la tabla de contabilidad
```

---

## ✅ Solución

### Opción 1: Migrar TODOS los Tenants (Recomendado)

```bash
./scripts/migrate-all-tenants-sql.sh
```

**Esto va a:**
1. Encontrar todos los tenants en la base de datos
2. Crear la tabla `chart_of_accounts` en cada uno
3. Insertar 60+ cuentas preconfiguradas
4. Mostrar resumen de éxito/fallo

**Salida esperada:**
```
🔄 Ejecutando migraciones SQL en todos los tenants...
==================================================
📋 Schemas de tenants encontrados:
tenant_demo

🔄 Migrando schema: tenant_demo
   ✅ Migración exitosa

==================================================
📊 Resumen:
   Total: 1
   Exitosos: 1
   Fallidos: 0

✅ Todas las migraciones se ejecutaron correctamente
```

### Opción 2: Migrar UN Tenant Específico

```bash
# Identificar el schema del tenant
psql -h localhost -U postgres -d vendix -c "
  SELECT id, name, schema_name FROM public.tenants;
"

# Migrar ese tenant específico
psql -h localhost -U postgres -d vendix \
  -v schema=tenant_demo \
  -f scripts/migrate-tenant-direct.sql
```

### Opción 3: Con API (si el servidor está corriendo)

```bash
# 1. Listar tenants para obtener el ID
curl http://localhost:8080/api/v1/public/tenants

# 2. Ejecutar migraciones para ese tenant
curl -X POST http://localhost:8080/api/v1/public/tenants/{TENANT_ID}/init
```

---

## 🧪 Verificación

Después de ejecutar las migraciones:

### 1. Verificar en Base de Datos

```bash
# Listar todas las cuentas
psql -h localhost -U postgres -d vendix -c "
  SELECT count(*) FROM tenant_demo.chart_of_accounts;
"
```

**Debe mostrar:** Al menos 60 cuentas

### 2. Verificar con API

```bash
curl http://localhost:8080/api/v1/tenant/accounting/accounts \
  -H "Authorization: Bearer {tu_token}" \
  -H "X-Tenant-ID: {tu_tenant_slug}"
```

**Debe retornar:** JSON con array de cuentas

### 3. Verificar en Frontend

1. Ir a http://localhost:5173/accounting
2. Seleccionar la pestaña "Plan de Cuentas"
3. Debe mostrar una lista con más de 60 cuentas

---

## 🐛 Otros Errores Posibles

### Error: "Permission denied"

**Causa:** El script no tiene permisos de ejecución

**Solución:**
```bash
chmod +x scripts/*.sh
./scripts/migrate-all-tenants-sql.sh
```

### Error: "Connection refused"

**Causa:** PostgreSQL no está corriendo

**Solución:**
```bash
# Iniciar servicios con Docker
make dev-up

# O solo la base de datos
docker-compose up -d postgres
```

### Error: "No rows in result set"

**Causa:** El tenant existe pero la tabla está vacía

**Solución:**
```bash
# La migración inserta datos automáticamente
./scripts/migrate-all-tenants-sql.sh
```

### Error: "X-Tenant-ID header required"

**Causa:** No se está enviando el header del tenant

**Solución:** Verificar que estás autenticado y el authStore tiene el tenantId

```javascript
// En el frontend, verificar:
console.log(useAuthStore.getState())
// Debe tener: { tenantId: 'demo', accessToken: '...' }
```

### Error en Frontend: "Network Error"

**Causa:** El backend no está corriendo

**Solución:**
```bash
# Iniciar el backend
make run-api

# O con Docker
make dev-up
```

---

## 📋 Checklist de Solución

Marca cada paso a medida que lo completes:

- [ ] 1. Verificar que PostgreSQL está corriendo
- [ ] 2. Ejecutar `./scripts/check-accounting-tables.sh`
- [ ] 3. Si falta la tabla, ejecutar `./scripts/migrate-all-tenants-sql.sh`
- [ ] 4. Verificar que la tabla tiene datos (60+ cuentas)
- [ ] 5. Reiniciar el backend (si está corriendo)
- [ ] 6. Refrescar el frontend
- [ ] 7. Acceder a /accounting y verificar que carga

---

## 🔄 Flujo Completo de Solución

```bash
# 1. Verificar el estado actual
./scripts/check-accounting-tables.sh

# 2. Si hay problemas, ejecutar migraciones
./scripts/migrate-all-tenants-sql.sh

# 3. Verificar que funcionó
psql -h localhost -U postgres -d vendix -c "
  SELECT schema_name, 
         (SELECT count(*) FROM chart_of_accounts) as num_cuentas
  FROM information_schema.schemata 
  WHERE schema_name LIKE 'tenant_%'
  LIMIT 5;
"

# 4. Reiniciar backend (si está corriendo)
# Ctrl+C para detener
# make run-api para reiniciar

# 5. Probar en el frontend
# Ir a http://localhost:5173/accounting
```

---

## 💡 Prevención

Para evitar este problema en el futuro:

### 1. Siempre Inicializar Nuevos Tenants

Cuando crees un nuevo tenant:

```bash
# Opción A: Con API
curl -X POST http://localhost:8080/api/v1/public/tenants \
  -H "Content-Type: application/json" \
  -d '{"name":"Mi Empresa","slug":"miempresa"}'

# Obtener el ID del tenant creado, luego:
curl -X POST http://localhost:8080/api/v1/public/tenants/{TENANT_ID}/init

# Opción B: Comando make
make create-tenant SLUG=miempresa
make init-tenant ID={tenant_id}
```

### 2. Ejecutar Migraciones Después de Actualizar

Después de hacer `git pull` o actualizar el código:

```bash
# Verificar si hay nuevas migraciones
./scripts/check-accounting-tables.sh

# Ejecutar migraciones si es necesario
./scripts/migrate-all-tenants-sql.sh
```

---

## 📚 Documentación Relacionada

- `EJECUTAR_MIGRACIONES.md` - Guía completa de migraciones
- `CONTABILIDAD_IMPLEMENTADA.md` - Documentación del módulo
- `MIGRACION_TENANT_CONFIG.md` - Cómo funcionan las migraciones

---

## 🆘 Si Nada Funciona

### Opción Nuclear: Reinicializar Todo

⚠️ **ADVERTENCIA:** Esto borrará todos los datos del tenant

```bash
# 1. Identificar el tenant
psql -h localhost -U postgres -d vendix -c "
  SELECT id, schema_name FROM public.tenants WHERE slug = 'demo';
"

# 2. Eliminar el schema
psql -h localhost -U postgres -d vendix -c "
  DROP SCHEMA IF EXISTS tenant_demo CASCADE;
"

# 3. Reinicializar
curl -X POST http://localhost:8080/api/v1/public/tenants/{TENANT_ID}/init
```

---

## ✅ Verificación Final

Después de solucionar, deberías poder:

1. ✅ Ver la lista de cuentas en `/accounting`
2. ✅ Filtrar cuentas por tipo
3. ✅ Ver balances de cada cuenta
4. ✅ Crear nuevos asientos contables
5. ✅ Generar reportes (Balance de Comprobación, etc.)

---

## 📞 Resumen Rápido

**Problema:** Error al cargar cuentas  
**Causa más común:** Tabla `chart_of_accounts` no existe  
**Solución rápida:**
```bash
./scripts/migrate-all-tenants-sql.sh
```

**Si persiste el error:**
1. Verificar que PostgreSQL está corriendo
2. Verificar que el backend está corriendo
3. Verificar que estás autenticado en el frontend
4. Revisar los logs del backend para más detalles

---

¡El error debería estar solucionado! Si continúa, revisa los logs del servidor para más detalles sobre el error específico.

