# Cambio de Nombre del Proyecto

## 📋 Resumen del Cambio

**Fecha:** 2025-10-16
**De:** saas-billing / saas_billing
**A:** vendix

---

## ✅ Cambios Realizados

### 1. Archivos de Configuración

- ✅ `go.mod` - Nombre del módulo actualizado a `vendix`
- ✅ `frontend/package.json` - Nombre del paquete actualizado a `vendix-frontend`
- ✅ `env.example` - Variables de base de datos actualizadas
- ✅ `docker-compose.yml` - Nombres de contenedores actualizados
- ✅ `docker-compose.dev.yml` - Nombres de contenedores de desarrollo actualizados

### 2. Nombres de Contenedores Docker

**Producción (`docker-compose.yml`):**
- `saas-billing-postgres` → `vendix-postgres`
- `saas-billing-redis` → `vendix-redis`
- `saas-billing-minio` → `vendix-minio`
- `saas-billing-api` → `vendix-api`
- `saas-billing-worker` → `vendix-worker`
- `saas-billing-frontend` → `vendix-frontend`

**Desarrollo (`docker-compose.dev.yml`):**
- `saas-billing-postgres-dev` → `vendix-postgres-dev`
- `saas-billing-redis-dev` → `vendix-redis-dev`
- `saas-billing-minio-dev` → `vendix-minio-dev`
- `saas-billing-api-dev` → `vendix-api-dev`
- `saas-billing-worker-dev` → `vendix-worker-dev`
- `saas-billing-frontend-dev` → `vendix-frontend-dev`

### 3. Base de Datos

- `saas_billing` → `vendix`

### 4. Documentación

- ✅ README.md - Título y referencias actualizadas
- ✅ START_HERE.md - Referencias actualizadas
- ✅ QUICK_START.md - Referencias actualizadas
- ✅ DEV_SETUP.md - Referencias actualizadas
- ✅ SETUP_COMPLETO.md - Referencias actualizadas
- ✅ CHANGELOG_DEV_SETUP.md - Referencias actualizadas
- ✅ .github/DEVELOPMENT.md - Referencias actualizadas
- ✅ ARCHITECTURE.md - Referencias actualizadas
- ✅ PROJECT_SUMMARY.md - Referencias actualizadas
- ✅ docs/ - Todos los archivos de documentación actualizados

### 5. Scripts

- ✅ `scripts/dev-setup.sh` - Nombres de contenedores actualizados
- ✅ `scripts/verify-setup.sh` - Nombres de contenedores actualizados
- ✅ `scripts/README.md` - Referencias actualizadas

### 6. Código Fuente

- ✅ Todos los archivos `.go` - Importaciones y referencias actualizadas
- ✅ Comentarios y documentación en el código actualizada

---

## 🔄 Acciones Necesarias

### 1. Si tienes contenedores corriendo

Detén los contenedores antiguos:

```bash
# Detener contenedores de desarrollo
docker-compose -f docker-compose.dev.yml down

# Detener contenedores de producción
docker-compose down

# (Opcional) Eliminar contenedores antiguos
docker rm -f saas-billing-postgres saas-billing-redis saas-billing-minio \
  saas-billing-api saas-billing-worker saas-billing-frontend \
  saas-billing-postgres-dev saas-billing-redis-dev saas-billing-minio-dev \
  saas-billing-api-dev saas-billing-worker-dev saas-billing-frontend-dev
```

### 2. Actualizar archivo .env

Si tienes un archivo `.env` personalizado, actualízalo:

```bash
# Cambiar nombre de base de datos
DB_NAME=vendix  # antes: saas_billing

# Cambiar bucket de S3 (opcional)
S3_BUCKET=vendix  # antes: saas-billing
```

### 3. Reconstruir imágenes

```bash
# Para desarrollo
make dev-rebuild

# Para producción
make docker-rebuild
```

### 4. Actualizar dependencias Go

```bash
go mod tidy
```

---

## 🚀 Iniciar el Proyecto con el Nuevo Nombre

### Desarrollo

```bash
make dev-up
```

### Producción

```bash
make docker-up
```

---

## 🔍 Verificar los Cambios

```bash
# Ver contenedores con el nuevo nombre
docker ps

# Verificar módulo Go
cat go.mod | head -1

# Verificar frontend
cat frontend/package.json | grep name

# Verificar base de datos
cat env.example | grep DB_NAME
```

---

## 📝 Notas Importantes

1. **Base de Datos Existente:**
   Si tienes datos en una base de datos llamada `saas_billing`, necesitarás:
   - Renombrar la base de datos: `ALTER DATABASE saas_billing RENAME TO vendix;`
   - O actualizar tu `.env` para usar `vendix` como nombre de base

2. **Volúmenes Docker:**
   Los volúmenes existentes seguirán funcionando, pero tienen el nombre antiguo.
   Si quieres limpiarlos:
   ```bash
   docker volume ls | grep saas-billing
   docker volume rm <volume-name>
   ```

3. **IDE y Editor:**
   - VSCode: Puede necesitar recargar la ventana
   - GoLand: Invalidar cachés y reiniciar
   - Otros: Recargar proyecto

4. **Git:**
   No es necesario renombrar el repositorio Git, pero si quieres hacerlo:
   - En GitHub/GitLab: Renombrar en Settings
   - Local: `git remote set-url origin <new-url>`

---

## ✅ Checklist Post-Cambio

- [ ] Contenedores antiguos detenidos
- [ ] Archivo `.env` actualizado
- [ ] Dependencias Go actualizadas (`go mod tidy`)
- [ ] Imágenes Docker reconstruidas
- [ ] Proyecto inicia correctamente con `make dev-up`
- [ ] Frontend accesible en http://localhost:3000
- [ ] API accesible en http://localhost:8080
- [ ] Base de datos renombrada o nueva creada

---

## 🆘 Problemas Comunes

### Error: "container name already exists"

```bash
# Eliminar contenedores antiguos
docker rm -f $(docker ps -aq --filter "name=saas-billing")
```

### Error: "database does not exist"

```bash
# Crear nueva base de datos
docker exec -it vendix-postgres-dev psql -U postgres -c "CREATE DATABASE vendix;"
```

### Error: "module not found"

```bash
# Actualizar dependencias
cd /home/scorpion/desa/Projects/saas-billing
go mod tidy
```

### El IDE no reconoce el módulo

```bash
# Limpiar caché y reiniciar IDE
# VSCode: Recargar ventana (Ctrl+Shift+P -> Reload Window)
# GoLand: File -> Invalidate Caches -> Invalidate and Restart
```

---

## 📊 Estadísticas del Cambio

- **Archivos modificados:** ~50 archivos
- **Tipos de archivo:** .go, .md, .yml, .json, .sh, .txt
- **Referencias cambiadas:** 150+ ocurrencias
- **Contenedores renombrados:** 12 (6 prod + 6 dev)

---

## 🎉 ¡Cambio Completado!

El proyecto ahora se llama **Vendix** en lugar de **SaaS Billing**.

Todos los archivos, configuraciones, documentación y código han sido actualizados.

**Siguiente paso:** Ejecuta `make dev-up` para iniciar el proyecto con el nuevo nombre.

