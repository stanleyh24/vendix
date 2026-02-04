# Inicio Rápido - 5 Minutos ⚡

Guía ultra-rápida para empezar a desarrollar.

## 1️⃣ Clonar e Instalar (1 min)

```bash
git clone <repository-url>
cd vendix
cp .env.dev.example .env
```

## 2️⃣ Iniciar Entorno de Desarrollo (2 min)

```bash
make dev-up
```

Esto iniciará:
- ✅ PostgreSQL
- ✅ Redis
- ✅ MinIO
- ✅ Backend API (Go con hot-reload)
- ✅ Worker
- ✅ Frontend (React con hot-reload)

## 3️⃣ Ver Logs (1 min)

En otra terminal:

```bash
make dev-logs
```

Espera a ver estos mensajes:
- `API server started on :8080`
- `VITE ready in XXX ms`

## 4️⃣ Crear Tenant de Prueba (1 min)

```bash
# Crear tenant
curl -X POST http://localhost:8080/api/v1/public/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mi Empresa Demo",
    "slug": "demo",
    "plan_id": null
  }'

# Guarda el ID que retorna, luego inicializa:
curl -X POST http://localhost:8080/api/v1/public/tenants/<TENANT_ID>/init
```

O usa el script:

```bash
bash scripts/create-test-tenant.sh
```

## 5️⃣ Registrar Usuario

```bash
curl -X POST http://localhost:8080/api/v1/tenant/auth/register \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: demo" \
  -d '{
    "email": "admin@demo.com",
    "password": "SecurePass123!",
    "first_name": "Admin",
    "last_name": "Demo",
    "role": "admin"
  }'
```

## 6️⃣ Abrir la Aplicación

🌐 **Frontend:** http://localhost:3000  
📦 **MinIO:** http://localhost:9001 (minioadmin / minioadmin)

## 🎯 Ahora Puedes Desarrollar

### Cambios en el Backend (Go)

1. Edita cualquier archivo `.go`
2. Air detectará el cambio y recompilará automáticamente
3. El servidor se reiniciará con tus cambios

### Cambios en el Frontend (React)

1. Edita cualquier archivo en `frontend/src/`
2. Vite aplicará el cambio instantáneamente (HMR)
3. No se recarga la página completa

## 📝 Comandos Útiles

```bash
# Ver logs en tiempo real
make dev-logs

# Reiniciar solo el API
make dev-restart-api

# Reiniciar solo el frontend
make dev-restart-frontend

# Detener todo
make dev-down

# Iniciar de nuevo
make dev-up
```

## 🐛 Si Algo Falla

### El contenedor no inicia

```bash
# Ver logs detallados
docker-compose -f docker-compose.dev.yml logs api

# Reconstruir
make dev-rebuild
```

### Cambios no se reflejan

```bash
# Backend: verifica que Air esté corriendo
docker-compose -f docker-compose.dev.yml logs api | grep "Building"

# Frontend: verifica que Vite esté en modo dev
docker-compose -f docker-compose.dev.yml logs frontend
```

### Limpiar todo y empezar de cero

```bash
make dev-down
docker-compose -f docker-compose.dev.yml down -v
make clean
make dev-up
```

## 🚀 Próximos Pasos

- Lee [DEV_SETUP.md](./DEV_SETUP.md) para guía completa de desarrollo
- Lee [README.md](./README.md) para arquitectura y características
- Explora [API_EXAMPLES.md](./docs/API_EXAMPLES.md) para ejemplos de API

## 💡 Tips Pro

1. **Usa 2 terminales:**
   - Terminal 1: `make dev-logs` (ver logs)
   - Terminal 2: comandos de desarrollo

2. **No reconstruyas innecesariamente:**
   - Solo usa `make dev-rebuild` si cambias `go.mod` o `package.json`
   - Los cambios en código se aplican automáticamente

3. **Verifica el health:**
   ```bash
   curl http://localhost:8080/health
   ```

4. **Acceso directo a la DB:**
   ```bash
   docker exec -it vendix-postgres-dev psql -U postgres -d vendix
   ```

¡Listo! Ahora estás desarrollando con hot-reload 🔥

