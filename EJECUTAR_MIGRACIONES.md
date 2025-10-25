# 🚀 Cómo Ejecutar las Migraciones de tenant_config

## ⚡ Solución Rápida

Elige UNA de estas opciones según tu situación:

---

## 📍 Opción 1: Con Servidor Corriendo (API)

Si tu servidor backend está corriendo en http://localhost:8080:

```bash
# Ejecutar migraciones en TODOS los tenants
make migrate-tenants
```

**O si prefieres hacerlo manualmente para cada tenant:**

```bash
# 1. Ver lista de tenants
curl http://localhost:8080/api/v1/public/tenants | jq

# 2. Para cada tenant, ejecutar migraciones (reemplaza {ID} con el UUID del tenant)
curl -X POST http://localhost:8080/api/v1/public/tenants/{ID}/init

# Ejemplo:
# curl -X POST http://localhost:8080/api/v1/public/tenants/123e4567-e89b-12d3-a456-426614174000/init
```

---

## 📍 Opción 2: Sin Servidor (SQL Directo)

Si el servidor NO está corriendo o prefieres ejecutar SQL directamente:

```bash
# Ejecutar migraciones SQL en TODOS los tenants
./scripts/migrate-all-tenants-sql.sh
```

**O si prefieres hacerlo manualmente para cada tenant:**

```bash
# 1. Ver schemas de tenants
psql -h localhost -U postgres -d vendix -c "
  SELECT schema_name FROM public.tenants WHERE deleted_at IS NULL;
"

# 2. Para cada schema, ejecutar:
psql -h localhost -U postgres -d vendix \
  -v schema=tenant_demo \
  -f scripts/migrate-tenant-direct.sql

# Ejemplo para otro tenant:
# psql -h localhost -U postgres -d vendix \
#   -v schema=tenant_testcompany \
#   -f scripts/migrate-tenant-direct.sql
```

---

## 📍 Opción 3: Con Docker

Si estás usando Docker:

```bash
# Entrar al contenedor
docker-compose exec api bash

# Dentro del contenedor, ejecutar:
./scripts/migrate-all-tenants-sql.sh
```

**O desde fuera del contenedor:**

```bash
docker-compose exec api ./scripts/migrate-all-tenants-sql.sh
```

---

## ✅ Verificar que Funcionó

Después de ejecutar las migraciones, verifica que todo está bien:

### 1. Verificar en la base de datos:

```bash
# Ver si la tabla existe en un tenant (reemplaza tenant_demo con tu schema)
psql -h localhost -U postgres -d vendix -c "
  SELECT table_name 
  FROM information_schema.tables 
  WHERE table_schema = 'tenant_demo' 
    AND table_name = 'tenant_config';
"
```

Debería mostrar: `tenant_config`

### 2. Verificar con la API:

```bash
# Obtener configuración (necesitas un token válido)
curl http://localhost:8080/api/v1/tenant/config \
  -H "Authorization: Bearer TU_TOKEN_AQUI" \
  -H "X-Tenant-ID: SLUG_DEL_TENANT"
```

Debería retornar un JSON con la configuración.

### 3. Verificar en el Frontend:

1. Inicia sesión en http://localhost:5173/login
2. Ve a "Configuración" en el menú lateral
3. Deberías ver el formulario con las pestañas de configuración

---

## 🎯 Ejemplo Completo Paso a Paso

### Escenario: Tienes un tenant llamado "demo"

#### Con API (servidor corriendo):

```bash
# 1. Iniciar el servidor si no está corriendo
make dev-up

# 2. Esperar unos segundos a que inicie
sleep 5

# 3. Listar tenants para obtener el ID
curl http://localhost:8080/api/v1/public/tenants

# Respuesta ejemplo:
# [
#   {
#     "id": "123e4567-e89b-12d3-a456-426614174000",
#     "name": "Demo Company",
#     "slug": "demo",
#     "schema_name": "tenant_demo",
#     ...
#   }
# ]

# 4. Ejecutar migración (usar el ID de arriba)
curl -X POST http://localhost:8080/api/v1/public/tenants/123e4567-e89b-12d3-a456-426614174000/init

# Respuesta exitosa:
# {"message":"Tenant schema initialized successfully"}

# 5. Verificar accediendo a la configuración en el frontend
# http://localhost:5173/settings
```

#### Con SQL (servidor NO necesario):

```bash
# 1. Ejecutar migración directamente
psql -h localhost -U postgres -d vendix \
  -v schema=tenant_demo \
  -f scripts/migrate-tenant-direct.sql

# Respuesta:
# Creando tabla tenant_config en el schema tenant_demo
# ...
# Tabla tenant_config creada exitosamente en tenant_demo

# 2. Verificar
psql -h localhost -U postgres -d vendix -c "
  SELECT count(*) FROM tenant_demo.tenant_config;
"

# Debería mostrar: 1
```

---

## 🆘 Si Algo Sale Mal

### Error: "could not connect to database"

```bash
# Asegurarte que PostgreSQL está corriendo
make dev-up
# O
docker-compose up -d postgres
```

### Error: "relation tenant_config already exists"

✅ **Esto es bueno!** Significa que la migración ya se ejecutó antes. No necesitas hacer nada.

### Error: "no such file or directory: ./scripts/..."

```bash
# Asegurarte de estar en el directorio raíz del proyecto
cd /home/scorpion/desa/Projects/vendix

# Ejecutar desde ahí
./scripts/migrate-all-tenants-sql.sh
```

### Error: "Permission denied"

```bash
# Dar permisos de ejecución
chmod +x scripts/*.sh

# Intentar de nuevo
./scripts/migrate-all-tenants-sql.sh
```

---

## 📝 Resumen de Comandos

```bash
# Ver todos los comandos disponibles
make help

# Migrar todos los tenants (servidor corriendo)
make migrate-tenants

# Migrar todos los tenants (sin servidor)
./scripts/migrate-all-tenants-sql.sh

# Migrar un tenant específico
make init-tenant ID={tenant_id}

# O con SQL directo
psql -h localhost -U postgres -d vendix \
  -v schema=tenant_demo \
  -f scripts/migrate-tenant-direct.sql

# Ver logs si algo falla
make dev-logs
```

---

## 🎉 ¡Listo!

Después de ejecutar las migraciones:

1. ✅ La tabla `tenant_config` existe en cada tenant
2. ✅ Hay una fila con configuración por defecto
3. ✅ El endpoint `/api/v1/tenant/config` funciona
4. ✅ La página `/settings` en el frontend carga correctamente

**Ahora puedes usar la sección de configuración!** 🚀

---

## 📚 Documentación Adicional

- `MIGRACION_TENANT_CONFIG.md` - Detalles técnicos completos
- `FIX_ERROR_CONFIGURACION.md` - Solución al error de carga
- `CONFIGURACION_TENANT_IMPLEMENTADA.md` - Funcionalidad implementada

