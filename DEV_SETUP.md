# Configuración de Entorno de Desarrollo

Este documento explica cómo configurar y usar el entorno de desarrollo local con hot-reload.

## 🚀 Inicio Rápido

### Desarrollo Local (con hot-reload)

```bash
# Iniciar todos los servicios en modo desarrollo
make dev-up

# Ver logs en tiempo real
make dev-logs

# Detener el entorno de desarrollo
make dev-down
```

### Producción (sin hot-reload)

```bash
# Iniciar en modo producción
make docker-up

# Ver logs
make docker-logs

# Detener
make docker-down
```

## 🔄 Diferencias entre Desarrollo y Producción

### Modo Desarrollo (`docker-compose.dev.yml`)

**Ventajas:**
- ✅ Hot-reload automático en backend (Go con Air)
- ✅ Hot-reload automático en frontend (Vite)
- ✅ No necesitas reconstruir imágenes en cada cambio
- ✅ Cambios instantáneos en el código
- ✅ Ideal para desarrollo diario

**Características:**
- Backend usa Air para detectar cambios y recompilar automáticamente
- Frontend usa Vite dev server con HMR (Hot Module Replacement)
- El código fuente se monta como volumen en los contenedores
- Los módulos de Go y node_modules se cachean para velocidad

### Modo Producción (`docker-compose.yml`)

**Ventajas:**
- ✅ Optimizado para producción
- ✅ Imágenes más pequeñas (multi-stage build)
- ✅ Mejor rendimiento
- ✅ Más seguro (sin código fuente montado)

**Características:**
- Backend compilado estáticamente
- Frontend con build optimizado servido por Nginx
- Sin dependencias de desarrollo
- Imágenes optimizadas con Alpine Linux

## 📋 Comandos Disponibles

### Desarrollo

```bash
make dev-up              # Iniciar entorno de desarrollo
make dev-down            # Detener entorno de desarrollo
make dev-logs            # Ver logs del entorno de desarrollo
make dev-rebuild         # Reconstruir contenedores de desarrollo
make dev-restart-api     # Reiniciar solo el API
make dev-restart-frontend # Reiniciar solo el frontend
```

### Producción

```bash
make docker-up           # Iniciar entorno de producción
make docker-down         # Detener entorno de producción
make docker-logs         # Ver logs de producción
make docker-rebuild      # Reconstruir contenedores de producción
```

### Desarrollo Local (sin Docker)

```bash
make run-api            # Ejecutar API directamente
make run-worker         # Ejecutar worker directamente
make test               # Ejecutar tests
make test-coverage      # Tests con coverage
```

## 🛠️ Tecnologías de Hot-Reload

### Backend: Air

Air es una herramienta de live-reload para Go que:
- Detecta cambios en archivos `.go`
- Recompila automáticamente
- Reinicia el servidor

**Configuración:** `.air.toml` y `.air.worker.toml`

### Frontend: Vite

Vite incluye HMR (Hot Module Replacement) nativo que:
- Actualiza módulos sin recargar la página completa
- Preserva el estado de React
- Es extremadamente rápido

## 📁 Estructura de Archivos de Desarrollo

```
.
├── docker-compose.yml       # Producción
├── docker-compose.dev.yml   # Desarrollo con hot-reload
├── Dockerfile               # Producción (multi-stage)
├── Dockerfile.dev           # Desarrollo (con Air)
├── .air.toml                # Configuración Air para API
├── .air.worker.toml         # Configuración Air para Worker
└── frontend/
    ├── Dockerfile           # Producción (Nginx)
    └── Dockerfile.dev       # Desarrollo (Vite dev server)
```

## 🔍 Solución de Problemas

### El código no se actualiza automáticamente

1. Verifica que los contenedores estén usando `docker-compose.dev.yml`:
   ```bash
   docker-compose -f docker-compose.dev.yml ps
   ```

2. Revisa los logs de Air:
   ```bash
   docker-compose -f docker-compose.dev.yml logs api
   ```

3. Si persiste, reconstruye:
   ```bash
   make dev-rebuild
   ```

### Errores de compilación en Go

Los errores se mostrarán en los logs de Air. Fija el código y Air recompilará automáticamente.

```bash
make dev-logs
```

### Problemas con node_modules

Si el frontend no encuentra módulos:

```bash
# Reconstruir el contenedor frontend
docker-compose -f docker-compose.dev.yml up -d --build frontend
```

### Limpiar todo y empezar de cero

```bash
# Detener y eliminar volúmenes
docker-compose -f docker-compose.dev.yml down -v

# Limpiar archivos temporales
make clean

# Reconstruir desde cero
make dev-rebuild
```

## 🌐 URLs en Desarrollo

- **Frontend:** http://localhost:3000
- **API:** http://localhost:8080
- **MinIO Console:** http://localhost:9001
- **PostgreSQL:** localhost:5432
- **Redis:** localhost:6379

## 💡 Consejos

1. **Usa `make dev-logs`** para ver todos los logs en tiempo real
2. **No reconstruyas** a menos que cambies dependencias (go.mod, package.json)
3. **Los cambios en código** se reflejan automáticamente
4. **Para cambios en dependencias:**
   ```bash
   # Backend (si cambias go.mod)
   docker-compose -f docker-compose.dev.yml restart api

   # Frontend (si cambias package.json)
   docker-compose -f docker-compose.dev.yml restart frontend
   ```

## 🎯 Flujo de Trabajo Recomendado

1. **Iniciar el entorno:**
   ```bash
   make dev-up
   ```

2. **Abrir logs en una terminal separada:**
   ```bash
   make dev-logs
   ```

3. **Desarrollar normalmente** - los cambios se aplicarán automáticamente

4. **Cuando termines:**
   ```bash
   make dev-down
   ```

## 📦 Gestión de Dependencias

### Agregar dependencia de Go

```bash
# 1. Agregar en tu código
import "github.com/nuevo/paquete"

# 2. Actualizar go.mod
go mod tidy

# 3. Reiniciar el contenedor para descargar
docker-compose -f docker-compose.dev.yml restart api
```

### Agregar dependencia de NPM

```bash
# 1. Entrar al contenedor
docker-compose -f docker-compose.dev.yml exec frontend sh

# 2. Instalar paquete
npm install nuevo-paquete

# 3. Salir y reiniciar
exit
docker-compose -f docker-compose.dev.yml restart frontend
```

## 🚀 Rendimiento

El modo desarrollo es más lento que ejecutar directamente en tu máquina, pero:
- ✅ Garantiza consistencia entre desarrolladores
- ✅ Evita problemas de "funciona en mi máquina"
- ✅ Incluye todos los servicios necesarios (Postgres, Redis, MinIO)
- ✅ Fácil de resetear y compartir

Si prefieres velocidad máxima, puedes ejecutar backend y frontend fuera de Docker y solo usar Docker para servicios (Postgres, Redis, MinIO).

