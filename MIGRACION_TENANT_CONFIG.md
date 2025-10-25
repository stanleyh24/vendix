# 🔄 Migración de Tabla tenant_config

## 📋 Resumen

La migración v15 agrega la tabla `tenant_config` a cada schema de tenant para gestionar la configuración de la empresa.

---

## ⚠️ Importante

Las migraciones de tenant **NO se ejecutan automáticamente** al iniciar el servidor. Solo las migraciones del schema `public` se ejecutan automáticamente.

Para cada tenant (nuevo o existente), debes ejecutar las migraciones manualmente.

---

## 🚀 Métodos de Ejecución

### **Método 1: API (Recomendado) - Con Servidor Corriendo**

#### Para UN tenant específico:

```bash
# Obtener el ID del tenant
curl -X GET http://localhost:8080/api/v1/public/tenants

# Ejecutar migraciones para ese tenant
curl -X POST http://localhost:8080/api/v1/public/tenants/{TENANT_ID}/init
```

#### Con Make (para un tenant):

```bash
make init-tenant ID={TENANT_ID}
```

#### Para TODOS los tenants:

```bash
# Opción 1: Usar el script
./scripts/migrate-tenants.sh

# Opción 2: Usar Make
make migrate-tenants
```

---

### **Método 2: SQL Directo - Sin Servidor**

Si el servidor no está corriendo, puedes ejecutar el SQL directamente:

#### Para UN tenant específico:

```bash
# Reemplaza 'tenant_demo' con el schema_name de tu tenant
psql -h localhost -U postgres -d vendix \
  -v schema=tenant_demo \
  -f scripts/migrate-tenant-direct.sql
```

#### Para TODOS los tenants:

```bash
#!/bin/bash
# Listar todos los tenants
psql -h localhost -U postgres -d vendix -t -c "
  SELECT schema_name FROM public.tenants WHERE deleted_at IS NULL;
" | while read schema; do
  if [ -n "$schema" ]; then
    echo "Migrando $schema..."
    psql -h localhost -U postgres -d vendix \
      -v schema=$schema \
      -f scripts/migrate-tenant-direct.sql
  fi
done
```

---

### **Método 3: Desde Docker**

Si estás usando Docker:

```bash
# Con el script
docker-compose exec api ./scripts/migrate-tenants.sh

# O ejecutar migraciones para un tenant específico
docker-compose exec api curl -X POST \
  http://localhost:8080/api/v1/public/tenants/{TENANT_ID}/init
```

---

## 🔍 Verificar Migración

### Verificar que la tabla existe:

```bash
psql -h localhost -U postgres -d vendix -c "
  SELECT EXISTS (
    SELECT FROM information_schema.tables 
    WHERE table_schema = 'tenant_demo' 
    AND table_name = 'tenant_config'
  );
"
```

Debería retornar `t` (true).

### Verificar que hay una fila de configuración:

```bash
psql -h localhost -U postgres -d vendix -c "
  SELECT count(*) FROM tenant_demo.tenant_config;
"
```

Debería retornar `1`.

### Verificar la migración registrada:

```bash
psql -h localhost -U postgres -d vendix -c "
  SELECT version, name, applied_at 
  FROM tenant_demo.schema_migrations 
  WHERE version = 15;
"
```

Debería mostrar la migración v15.

---

## 📝 Pasos Completos para Migración

### Si el servidor está corriendo:

```bash
# 1. Listar tenants existentes
curl http://localhost:8080/api/v1/public/tenants

# 2. Para cada tenant, ejecutar migraciones
# Opción A: Uno por uno
curl -X POST http://localhost:8080/api/v1/public/tenants/{TENANT_ID_1}/init
curl -X POST http://localhost:8080/api/v1/public/tenants/{TENANT_ID_2}/init

# Opción B: Todos a la vez
make migrate-tenants
```

### Si el servidor NO está corriendo:

```bash
# 1. Identificar schemas de tenants
psql -h localhost -U postgres -d vendix -c "
  SELECT id, name, schema_name FROM public.tenants WHERE deleted_at IS NULL;
"

# 2. Migrar cada uno
psql -h localhost -U postgres -d vendix \
  -v schema=tenant_demo \
  -f scripts/migrate-tenant-direct.sql

psql -h localhost -U postgres -d vendix \
  -v schema=tenant_testcompany \
  -f scripts/migrate-tenant-direct.sql

# ... etc
```

---

## ⚡ Migración Rápida (Desarrollo)

Para desarrollo local con Docker:

```bash
# 1. Asegurarse que los contenedores están corriendo
make dev-up

# 2. Esperar a que el API esté listo
sleep 5

# 3. Migrar todos los tenants
make migrate-tenants
```

---

## 🛠️ Troubleshooting

### Error: "relation tenant_config does not exist"

**Causa:** Las migraciones no se han ejecutado en ese tenant.

**Solución:**
```bash
# Ejecutar migraciones
curl -X POST http://localhost:8080/api/v1/public/tenants/{TENANT_ID}/init
```

### Error: "duplicate key value violates unique constraint"

**Causa:** Ya existe una fila en `tenant_config`.

**Solución:** No es un error crítico. La configuración ya existe.

### Error: "permission denied for schema"

**Causa:** El usuario de la base de datos no tiene permisos.

**Solución:**
```sql
GRANT ALL ON SCHEMA tenant_demo TO postgres;
GRANT ALL ON ALL TABLES IN SCHEMA tenant_demo TO postgres;
```

### Error: "could not connect to server"

**Causa:** La base de datos no está corriendo.

**Solución:**
```bash
# Iniciar base de datos
make dev-up

# O si usas producción
make docker-up
```

---

## 📊 Estado de Migraciones

### Ver todas las migraciones aplicadas en un tenant:

```sql
SELECT * FROM tenant_demo.schema_migrations ORDER BY version;
```

### Ver tenants sin migración v15:

```sql
SELECT t.id, t.name, t.schema_name
FROM public.tenants t
WHERE t.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1 
    FROM information_schema.tables 
    WHERE table_schema = t.schema_name 
      AND table_name = 'tenant_config'
  );
```

---

## 🔐 Consideraciones de Seguridad

1. **No ejecutar migraciones en producción sin backup**
   ```bash
   # Hacer backup antes
   pg_dump -h localhost -U postgres -d vendix > backup_$(date +%Y%m%d).sql
   ```

2. **Probar primero en desarrollo**
   ```bash
   # Probar en un tenant de prueba
   make create-tenant SLUG=test_migration
   make init-tenant ID={nuevo_tenant_id}
   ```

3. **Verificar después de ejecutar**
   ```bash
   # Verificar que todo funciona
   curl http://localhost:8080/api/v1/tenant/config \
     -H "Authorization: Bearer {token}" \
     -H "X-Tenant-ID: {tenant_id}"
   ```

---

## 📚 Comandos Útiles

```bash
# Ver todos los comandos disponibles
make help

# Crear tenant de prueba
make create-tenant SLUG=test1

# Inicializar schema de tenant
make init-tenant ID={tenant_id}

# Migrar todos los tenants
make migrate-tenants

# Ver logs del servidor
make dev-logs

# Reiniciar API
make dev-restart-api
```

---

## ✅ Checklist de Migración

Para cada tenant existente:

- [ ] Verificar que el tenant existe en `public.tenants`
- [ ] Obtener el `id` y `schema_name` del tenant
- [ ] Ejecutar migraciones (API o SQL directo)
- [ ] Verificar que la tabla `tenant_config` existe
- [ ] Verificar que hay una fila de configuración por defecto
- [ ] Probar el endpoint `/api/v1/tenant/config` (GET)
- [ ] Probar actualizar configuración (PUT)
- [ ] Verificar en el frontend que carga `/settings`

---

## 🎯 Resumen

### Para tenants NUEVOS:
✅ Las migraciones se ejecutan automáticamente al crear el tenant con `/init`

### Para tenants EXISTENTES:
⚠️ **Debes ejecutar manualmente:**
```bash
# Opción más fácil
make migrate-tenants
```

### Verificación:
```bash
# Debe retornar 200 con los datos de configuración
curl http://localhost:8080/api/v1/tenant/config \
  -H "Authorization: Bearer {token}" \
  -H "X-Tenant-ID: {tenant_id}"
```

---

## 📞 Soporte

Si tienes problemas:

1. Revisa los logs del servidor: `make dev-logs`
2. Verifica la conexión a la base de datos
3. Consulta este documento
4. Revisa `FIX_ERROR_CONFIGURACION.md` para soluciones comunes

