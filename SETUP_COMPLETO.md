# ✅ Configuración de Desarrollo con Hot-Reload - COMPLETADO

## 🎉 ¡Tu entorno de desarrollo está listo!

Se ha configurado un entorno de desarrollo completo con **hot-reload** para que no tengas que reconstruir las imágenes Docker en cada cambio.

---

## 📦 ¿Qué se ha creado?

### 🐳 Docker para Desarrollo

| Archivo | Propósito |
|---------|-----------|
| `docker-compose.dev.yml` | Configuración Docker para desarrollo |
| `Dockerfile.dev` | Dockerfile del backend con Air (hot-reload) |
| `frontend/Dockerfile.dev` | Dockerfile del frontend con Vite dev server |
| `.air.toml` | Configuración de Air para el API |
| `.air.worker.toml` | Configuración de Air para el Worker |

### 📚 Documentación Nueva

| Archivo | Descripción |
|---------|-------------|
| `START_HERE.md` | 🚀 **EMPIEZA AQUÍ** - Punto de entrada principal |
| `QUICK_START.md` | Guía de inicio rápido (5 minutos) |
| `DEV_SETUP.md` | Guía completa de configuración de desarrollo |
| `.github/DEVELOPMENT.md` | Guía para desarrolladores (convenciones, testing, etc) |
| `CHANGELOG_DEV_SETUP.md` | Registro de cambios del entorno de desarrollo |
| `scripts/README.md` | Documentación de scripts disponibles |

### 🛠️ Scripts Nuevos

| Script | Función |
|--------|---------|
| `scripts/dev-setup.sh` | Configuración automática del entorno |

### 🔧 Archivos Actualizados

| Archivo | Cambios |
|---------|---------|
| `Makefile` | Agregados comandos `dev-*` para desarrollo |
| `README.md` | Nueva sección de desarrollo con hot-reload |
| `.gitignore` | Actualizado con `tmp/` y otros temporales |
| `.dockerignore` | Creado para optimizar builds |

---

## 🚀 Cómo Usar - Inicio Rápido

### 1️⃣ Primera Vez (Setup Automático)

```bash
# Ejecuta el script de configuración
bash scripts/dev-setup.sh
```

Este script:
- ✅ Verifica Docker
- ✅ Crea archivo `.env`
- ✅ Inicia todos los servicios
- ✅ Verifica que todo funcione

### 2️⃣ Crear Tenant de Prueba

```bash
bash scripts/create-test-tenant.sh
```

### 3️⃣ Abrir la Aplicación

- **Frontend:** http://localhost:3000
- **API Docs:** http://localhost:8080/swagger/index.html
- **MinIO Console:** http://localhost:9001

### 4️⃣ Desarrollar

¡Eso es todo! Ahora puedes:
- ✏️ Editar archivos en `internal/` (Go)
- ✏️ Editar archivos en `frontend/src/` (React)
- 🔄 Los cambios se aplicarán **automáticamente**
- ⚡ Sin necesidad de reconstruir Docker

---

## 📋 Comandos Disponibles

### Desarrollo Diario

```bash
# Iniciar entorno de desarrollo
make dev-up

# Ver logs en tiempo real
make dev-logs

# Detener entorno
make dev-down
```

### Comandos Útiles

```bash
# Reiniciar solo el API
make dev-restart-api

# Reiniciar solo el Frontend
make dev-restart-frontend

# Reconstruir (solo si cambias dependencias)
make dev-rebuild

# Limpiar todo
make clean
```

### Testing

```bash
# Ejecutar tests
make test

# Tests con coverage
make test-coverage

# Linter
make lint

# Formatear código
make fmt
```

### Ver Todos los Comandos

```bash
make help
```

---

## 🔥 Hot-Reload en Acción

### Backend (Go)

```bash
# 1. Edita un archivo
vim internal/auth/service.go

# 2. Guarda el archivo
# ✨ Air detecta el cambio automáticamente
# ✨ Recompila en ~2-5 segundos
# ✨ Reinicia el servidor
# ✨ Listo para probar!
```

**Sin rebuild de Docker** ⚡

### Frontend (React)

```bash
# 1. Edita un componente
vim frontend/src/pages/Dashboard.jsx

# 2. Guarda el archivo
# ✨ Vite aplica el cambio instantáneamente
# ✨ HMR actualiza sin recargar la página
# ✨ El estado de React se preserva
# ✨ Listo!
```

**Sin recargar la página** ⚡

---

## 📖 Documentación por Situación

### 🆕 Primera Vez en el Proyecto

1. Lee: `START_HERE.md` (3 min)
2. Ejecuta: `bash scripts/dev-setup.sh`
3. Lee: `QUICK_START.md` (5 min)

### 👨‍💻 Voy a Desarrollar

1. Ejecuta: `make dev-up`
2. Lee: `.github/DEVELOPMENT.md`
3. Desarrolla normalmente

### 🐛 Tengo un Problema

1. Lee: `DEV_SETUP.md` → Sección "Troubleshooting"
2. Verifica logs: `make dev-logs`
3. Si persiste: `make dev-rebuild`

### 🏗️ Quiero Entender la Arquitectura

1. Lee: `ARCHITECTURE.md`
2. Lee: `PROJECT_SUMMARY.md`
3. Lee: `README.md`

---

## 🎯 Diferencias: Desarrollo vs Producción

| Aspecto | Desarrollo (`make dev-up`) | Producción (`make docker-up`) |
|---------|----------------------------|-------------------------------|
| **Hot-Reload** | ✅ Sí | ❌ No |
| **Rebuild necesario** | ❌ No | ✅ Sí |
| **Velocidad de cambios** | ⚡ Instantáneo | 🐌 2-5 minutos |
| **Tamaño de imagen** | ~800MB | ~50MB |
| **Optimización** | Debug | Producción |
| **Uso** | Desarrollo local | Deploy real |

---

## ✨ Ventajas del Nuevo Setup

### Antes ❌

```bash
# Hacer un cambio
vim internal/auth/service.go

# Reconstruir imagen
docker-compose up -d --build api  # ⏱️ 2-5 minutos

# Esperar...
# Probar cambio
# Repetir proceso...
```

**Tiempo por iteración: 2-5 minutos** 😫

### Ahora ✅

```bash
# Iniciar una vez
make dev-up

# Hacer cambio
vim internal/auth/service.go

# Guardar
# ✨ Cambio aplicado automáticamente

# Probar inmediatamente
```

**Tiempo por iteración: 2-5 segundos** 🚀

### 🎯 Resultado

**Mejora de velocidad: ~50-100x más rápido**

---

## 🔍 Estructura de Archivos de Desarrollo

```
vendix/
├── 🚀 START_HERE.md              ← Lee esto primero
├── 📖 QUICK_START.md             ← Inicio rápido
├── 🔧 DEV_SETUP.md               ← Setup completo
├── 📝 CHANGELOG_DEV_SETUP.md     ← Cambios realizados
│
├── 🐳 docker-compose.dev.yml     ← Docker para desarrollo
├── 🐳 Dockerfile.dev             ← Dockerfile backend dev
├── ⚙️  .air.toml                  ← Config Air (API)
├── ⚙️  .air.worker.toml           ← Config Air (Worker)
│
├── frontend/
│   └── 🐳 Dockerfile.dev         ← Dockerfile frontend dev
│
├── scripts/
│   ├── 🛠️ dev-setup.sh           ← Setup automático
│   ├── 🛠️ create-test-tenant.sh  ← Crear tenant
│   └── 📖 README.md              ← Docs de scripts
│
└── .github/
    └── 📖 DEVELOPMENT.md         ← Guía completa dev
```

---

## 💡 Consejos Pro

### 1. Usa Dos Terminales

**Terminal 1:**
```bash
make dev-logs  # Déjala corriendo para ver logs
```

**Terminal 2:**
```bash
# Aquí trabajas normalmente
vim internal/auth/service.go
git commit -m "feat: nueva feature"
```

### 2. Aliases Útiles

Agrega a tu `~/.bashrc` o `~/.zshrc`:

```bash
alias dev-up='cd ~/path/to/vendix && make dev-up'
alias dev-down='cd ~/path/to/vendix && make dev-down'
alias dev-logs='cd ~/path/to/vendix && make dev-logs'
```

### 3. Verifica Salud de Servicios

```bash
# API
curl http://localhost:8080/health

# Frontend
curl http://localhost:3000

# Ver todos los contenedores
docker ps
```

### 4. Acceso Directo a Servicios

```bash
# PostgreSQL
docker exec -it vendix-postgres-dev psql -U postgres -d vendix

# Redis
docker exec -it vendix-redis-dev redis-cli

# Logs de un servicio específico
docker logs -f vendix-api-dev
```

---

## 🐛 Solución Rápida de Problemas

### Problema: Cambios no se reflejan

```bash
# Backend
make dev-restart-api

# Frontend
make dev-restart-frontend
```

### Problema: Error al iniciar

```bash
# Ver logs detallados
make dev-logs

# Reconstruir
make dev-rebuild
```

### Problema: Todo está roto

```bash
# Reset completo
make dev-down
docker-compose -f docker-compose.dev.yml down -v
make clean
make dev-up
```

---

## 📚 Recursos Adicionales

### Documentación del Proyecto
- [README.md](README.md) - Documentación completa
- [ARCHITECTURE.md](ARCHITECTURE.md) - Arquitectura
- [API_EXAMPLES.md](docs/API_EXAMPLES.md) - Ejemplos de API

### Guías de Desarrollo
- [START_HERE.md](START_HERE.md) - Punto de entrada
- [QUICK_START.md](QUICK_START.md) - Inicio rápido
- [DEV_SETUP.md](DEV_SETUP.md) - Setup detallado
- [DEVELOPMENT.md](.github/DEVELOPMENT.md) - Guía completa

### Herramientas Usadas
- [Air](https://github.com/cosmtrek/air) - Hot reload para Go
- [Vite](https://vitejs.dev/) - Build tool para frontend
- [Docker Compose](https://docs.docker.com/compose/) - Orquestación

---

## ✅ Checklist Final

Antes de empezar a desarrollar, verifica:

- [ ] Docker está instalado y corriendo
- [ ] Ejecutaste `bash scripts/dev-setup.sh`
- [ ] Los servicios están corriendo: `docker ps`
- [ ] Frontend accesible: http://localhost:3000
- [ ] API accesible: http://localhost:8080
- [ ] Logs visibles: `make dev-logs`
- [ ] Tenant de prueba creado
- [ ] Hot-reload funciona (prueba editando un archivo)

---

## 🎉 ¡Listo para Desarrollar!

Ahora tienes:
- ✅ Entorno de desarrollo configurado
- ✅ Hot-reload en backend (Go)
- ✅ Hot-reload en frontend (React)
- ✅ Documentación completa
- ✅ Scripts de utilidades
- ✅ Flujo de trabajo optimizado

**Siguiente paso:**

```bash
make dev-up
make dev-logs  # En otra terminal
# ... ¡A desarrollar! 🚀
```

---

## 📞 ¿Necesitas Ayuda?

1. **Lee la documentación:** Empieza con `START_HERE.md`
2. **Problemas técnicos:** `DEV_SETUP.md` → Troubleshooting
3. **Preguntas de desarrollo:** `.github/DEVELOPMENT.md`
4. **Issues:** Abre un issue en GitHub

---

**¡Feliz codificación!** 🎉🚀

*Última actualización: 2025-10-16*

