# Guía de Desarrollo

Esta guía está dirigida a desarrolladores que trabajan en este proyecto.

## 🎯 Configuración Inicial

### Opción 1: Script Automático (Recomendado)

```bash
bash scripts/dev-setup.sh
```

Este script:
- ✅ Verifica que Docker esté corriendo
- ✅ Crea el archivo `.env` si no existe
- ✅ Inicia todos los servicios en modo desarrollo
- ✅ Verifica que todo esté funcionando

### Opción 2: Manual

```bash
# 1. Copiar configuración
cp env.example .env

# 2. Iniciar servicios
make dev-up

# 3. Ver logs
make dev-logs
```

## 🔄 Flujo de Trabajo Diario

### Iniciar el Día

```bash
# Iniciar todos los servicios
make dev-up

# En otra terminal, ver logs
make dev-logs
```

### Durante el Desarrollo

- **Backend (Go)**: Edita archivos `.go`, Air recompilará automáticamente
- **Frontend (React)**: Edita archivos en `frontend/src/`, Vite aplicará cambios instantáneamente
- **SQL/Migraciones**: Si cambias migraciones, reinicia el API: `make dev-restart-api`

### Finalizar el Día

```bash
make dev-down
```

## 📝 Convenciones de Código

### Backend (Go)

- **Formato**: Ejecuta `make fmt` antes de commit
- **Linter**: Ejecuta `make lint` para verificar problemas
- **Tests**: Ejecuta `make test` antes de crear PR
- **Comentarios**: Documenta funciones exportadas con GoDoc

Estructura de archivos:
```
internal/modules/<modulo>/
├── handler.go      # HTTP handlers
├── service.go      # Lógica de negocio
├── repository.go   # Acceso a datos (SQL)
├── models.go       # Estructuras de datos
└── *_test.go       # Tests
```

### Frontend (React + TypeScript)

- **Formato**: ESLint + Prettier (configurado en `frontend/`)
- **Componentes**: Usar functional components con hooks
- **Estado**: React Query para datos remotos, Zustand para estado local
- **Tipos**: Definir interfaces en `frontend/src/types/`

Estructura de archivos:
```
frontend/src/
├── modules/<modulo>/    # Módulos por funcionalidad
│   ├── components/      # Componentes del módulo
│   ├── hooks/           # Hooks personalizados
│   └── types/           # Tipos TypeScript
├── components/          # Componentes compartidos
├── lib/                 # Utilidades y configuración
└── pages/               # Páginas principales
```

### Base de Datos (SQL)

- **Formato**: snake_case para tablas y columnas
- **Migraciones**: Nombrar con timestamp: `001_nombre_descriptivo.sql`
- **Comentarios**: Documentar tablas y columnas importantes
- **Transacciones**: Usar para operaciones críticas

## 🧪 Testing

### Backend

```bash
# Todos los tests
make test

# Con coverage
make test-coverage

# Módulo específico
go test -v ./internal/modules/invoices/...

# Test específico
go test -v -run TestCreateInvoice ./internal/modules/invoices/
```

### Frontend

```bash
cd frontend

# Todos los tests
npm test

# Con coverage
npm test -- --coverage

# Modo watch
npm test -- --watch
```

## 🐛 Debugging

### Backend

1. **Logs estructurados**: Ya están configurados con zerolog
   ```go
   log.Info().
       Str("tenant_id", tenantID).
       Str("invoice_id", invoiceID).
       Msg("Invoice created")
   ```

2. **Delve (debugger)**: Puedes adjuntar Delve al contenedor
   ```bash
   # Modificar Dockerfile.dev para incluir delve
   # Luego adjuntar con tu IDE
   ```

3. **Imprimir SQL**: Activa logs de SQL en modo desarrollo

### Frontend

1. **React DevTools**: Instalar extensión de navegador
2. **Console logs**: Usar durante desarrollo, remover antes de commit
3. **Redux DevTools**: Para inspeccionar Zustand

## 📦 Gestión de Dependencias

### Backend (Go)

```bash
# Agregar dependencia
go get github.com/paquete/nuevo

# Actualizar go.mod
go mod tidy

# Reiniciar contenedor
make dev-restart-api
```

### Frontend (NPM)

```bash
# Entrar al contenedor
docker-compose -f docker-compose.dev.yml exec frontend sh

# Instalar paquete
npm install nuevo-paquete

# Salir y reiniciar
exit
make dev-restart-frontend
```

## 🔍 Troubleshooting

### "Los cambios no se reflejan"

**Backend:**
```bash
# Ver logs de Air
docker-compose -f docker-compose.dev.yml logs api | grep Building

# Reiniciar
make dev-restart-api
```

**Frontend:**
```bash
# Ver logs de Vite
docker-compose -f docker-compose.dev.yml logs frontend

# Limpiar y reconstruir
docker-compose -f docker-compose.dev.yml restart frontend
```

### "Error de compilación en Go"

1. Verifica que `go.mod` esté actualizado: `go mod tidy`
2. Revisa los logs: `make dev-logs`
3. Fija el error, Air recompilará automáticamente

### "node_modules no encontrado"

```bash
# Reconstruir contenedor frontend
docker-compose -f docker-compose.dev.yml up -d --build frontend
```

### "No puedo conectar a la base de datos"

```bash
# Verificar que Postgres esté corriendo
docker ps | grep postgres

# Ver logs de Postgres
docker logs vendix-postgres-dev

# Reiniciar Postgres
docker-compose -f docker-compose.dev.yml restart postgres
```

### "Limpiar todo y empezar de cero"

```bash
# Detener y eliminar volúmenes
make dev-down
docker-compose -f docker-compose.dev.yml down -v

# Limpiar archivos temporales
make clean
rm -rf frontend/node_modules

# Reconstruir
make dev-rebuild
```

## 🚀 Preparar para Producción

Antes de hacer merge a `main`:

```bash
# 1. Ejecutar todos los tests
make test
cd frontend && npm test

# 2. Verificar linter
make lint

# 3. Formatear código
make fmt
cd frontend && npm run lint

# 4. Actualizar documentación si es necesario

# 5. Probar build de producción
make docker-up
# Verificar que todo funcione
make docker-down
```

## 📚 Recursos Adicionales

- [Arquitectura del Sistema](../ARCHITECTURE.md)
- [Ejemplos de API](../docs/API_EXAMPLES.md)
- [Deployment](../docs/DEPLOYMENT.md)
- [Inicio Rápido](../QUICK_START.md)

## 💬 Comunicación

- Issues: Para bugs y features
- Pull Requests: Revisar código antes de merge
- Commits: Mensajes descriptivos en español
  - `feat: agregar módulo de reportes`
  - `fix: corregir cálculo de ITBIS`
  - `docs: actualizar README`
  - `refactor: mejorar estructura de handlers`

## ⚡ Tips de Productividad

1. **Alias útiles** (agregar a `.bashrc` o `.zshrc`):
   ```bash
   alias dev-up='make dev-up'
   alias dev-down='make dev-down'
   alias dev-logs='make dev-logs'
   alias dev-api='make dev-restart-api'
   ```

2. **VS Code**: Extensiones recomendadas
   - Go (oficial)
   - ESLint
   - Prettier
   - Docker
   - PostgreSQL

3. **Shortcuts de Docker**:
   ```bash
   # Ejecutar comando en contenedor
   docker exec -it vendix-api-dev sh
   
   # Ver recursos de Docker
   docker stats
   
   # Limpiar recursos no usados
   docker system prune
   ```

4. **Acceso directo a servicios**:
   ```bash
   # PostgreSQL
   docker exec -it vendix-postgres-dev psql -U postgres -d vendix
   
   # Redis
   docker exec -it vendix-redis-dev redis-cli
   ```

¡Feliz codificación! 🎉

