# 🚀 Empieza Aquí

> **¿Primera vez en el proyecto? Este es el lugar correcto.**

## ⚡ Inicio Ultra-Rápido (5 minutos)

```bash
# 1. Configura el entorno automáticamente
bash scripts/dev-setup.sh

# 2. Crea un tenant de prueba
bash scripts/create-test-tenant.sh

# 3. Abre en tu navegador
# Frontend: http://localhost:3000
```

**¡Eso es todo!** Ahora puedes desarrollar con hot-reload. ✨

---

## 📚 Documentación por Rol

### 👨‍💻 Soy Desarrollador

**Desarrollo Local (Recomendado)**
```bash
make dev-up      # Iniciar con hot-reload
make dev-logs    # Ver logs
# ... edita código, los cambios se aplican automáticamente
make dev-down    # Detener
```

📖 **Lee:** [Guía de Desarrollo](.github/DEVELOPMENT.md)

---

### 🆕 Soy Nuevo en el Proyecto

**Paso a paso:**
1. 📖 Lee [QUICK_START.md](QUICK_START.md) (5 min)
2. 🔧 Ejecuta el setup: `bash scripts/dev-setup.sh`
3. 📚 Explora [ARCHITECTURE.md](ARCHITECTURE.md)
4. 💻 Lee [Guía de Desarrollo](.github/DEVELOPMENT.md)

---

### 🚀 Quiero Deployear

**Producción:**
```bash
make docker-up    # Iniciar en modo producción
make docker-logs  # Ver logs
```

📖 **Lee:** [DEPLOYMENT.md](docs/DEPLOYMENT.md)

---

### 📖 Quiero Entender la Arquitectura

**Documentos principales:**
- [ARCHITECTURE.md](ARCHITECTURE.md) - Arquitectura del sistema
- [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) - Resumen del proyecto
- [README.md](README.md) - Documentación completa

---

### 🐛 Tengo un Problema

**Soluciones rápidas:**

```bash
# Ver qué está corriendo
docker ps

# Ver logs de todos los servicios
make dev-logs

# Reiniciar un servicio específico
make dev-restart-api
make dev-restart-frontend

# Limpiar y empezar de cero
make dev-down
make clean
make dev-up
```

📖 **Lee:** [DEV_SETUP.md - Troubleshooting](DEV_SETUP.md#-solución-de-problemas)

---

## 🎯 Modos de Ejecución

### 🔥 Modo Desarrollo (Recomendado para desarrollo)

```bash
make dev-up
```

**Características:**
- ✅ Hot-reload en Backend (Go)
- ✅ Hot-reload en Frontend (React)
- ✅ No necesitas rebuild
- ✅ Cambios instantáneos

### 🚀 Modo Producción

```bash
make docker-up
```

**Características:**
- ✅ Optimizado para rendimiento
- ✅ Imágenes más pequeñas
- ✅ Build multi-stage
- ❌ Sin hot-reload

---

## 📋 Comandos Más Usados

```bash
# DESARROLLO
make dev-up              # Iniciar desarrollo
make dev-down            # Detener desarrollo
make dev-logs            # Ver logs
make dev-rebuild         # Reconstruir si es necesario

# TESTING
make test                # Ejecutar tests
make test-coverage       # Tests con cobertura
make lint                # Ejecutar linter

# UTILIDADES
make help                # Ver todos los comandos
make clean               # Limpiar archivos temporales
```

---

## 🌐 URLs Importantes

Después de ejecutar `make dev-up`:

| Servicio | URL | Credenciales |
|----------|-----|--------------|
| **Frontend** | http://localhost:3000 | - |
| **API** | http://localhost:8080 | - |
| **Swagger** | http://localhost:8080/swagger/index.html | - |
| **MinIO** | http://localhost:9001 | minioadmin / minioadmin |
| **PostgreSQL** | localhost:5432 | postgres / postgres |
| **Redis** | localhost:6379 | - |

---

## 📖 Documentación Completa

### Guías de Desarrollo
- 🚀 [QUICK_START.md](QUICK_START.md) - Inicio en 5 minutos
- 🔧 [DEV_SETUP.md](DEV_SETUP.md) - Configuración detallada
- 👨‍💻 [DEVELOPMENT.md](.github/DEVELOPMENT.md) - Guía completa de desarrollo
- 📜 [CHANGELOG_DEV_SETUP.md](CHANGELOG_DEV_SETUP.md) - Cambios del entorno

### Documentación Técnica
- 🏗️ [ARCHITECTURE.md](ARCHITECTURE.md) - Arquitectura del sistema
- 📖 [README.md](README.md) - Documentación principal
- 📊 [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) - Resumen del proyecto
- 🌐 [API_EXAMPLES.md](docs/API_EXAMPLES.md) - Ejemplos de API

### Deployment
- 🚀 [DEPLOYMENT.md](docs/DEPLOYMENT.md) - Guía de deployment

---

## 🎓 Flujo de Trabajo Típico

```bash
# 1️⃣ Lunes por la mañana
make dev-up              # Iniciar servicios
make dev-logs            # Ver logs (en otra terminal)

# 2️⃣ Durante la semana
# ... editas código ...
# ... los cambios se aplican automáticamente ...
# ... desarrollas features ...
# ... haces commits ...

# 3️⃣ Antes de hacer Push/PR
make test                # Ejecutar tests
make lint                # Verificar código
make fmt                 # Formatear código

# 4️⃣ Viernes al terminar
make dev-down            # Detener servicios
```

---

## 🆘 ¿Necesitas Ayuda?

1. **Problemas técnicos:** Abre un Issue en GitHub
2. **Preguntas de código:** Lee la [Guía de Desarrollo](.github/DEVELOPMENT.md)
3. **Setup no funciona:** Lee [DEV_SETUP.md - Troubleshooting](DEV_SETUP.md#-solución-de-problemas)
4. **API no funciona:** Verifica logs con `make dev-logs`

---

## ✨ Tips Pro

### Para Máxima Productividad

1. **Usa dos terminales:**
   - Terminal 1: `make dev-logs` (déjala corriendo)
   - Terminal 2: Comandos y desarrollo

2. **Atajos útiles:**
   ```bash
   alias dev-up='make dev-up'
   alias dev-down='make dev-down'
   alias dev-logs='make dev-logs'
   ```

3. **Verifica salud:**
   ```bash
   curl http://localhost:8080/health
   ```

4. **Acceso directo a servicios:**
   ```bash
   # PostgreSQL
   docker exec -it vendix-postgres-dev psql -U postgres -d vendix
   
   # Redis
   docker exec -it vendix-redis-dev redis-cli
   ```

---

## 🎯 Siguiente Paso

Elige tu siguiente paso según tu objetivo:

- 🏃 **Quiero empezar YA:** → `bash scripts/dev-setup.sh`
- 📖 **Quiero entender primero:** → [QUICK_START.md](QUICK_START.md)
- 👨‍💻 **Voy a desarrollar:** → [DEVELOPMENT.md](.github/DEVELOPMENT.md)
- 🏗️ **Quiero entender la arquitectura:** → [ARCHITECTURE.md](ARCHITECTURE.md)
- 🚀 **Voy a deployear:** → [DEPLOYMENT.md](docs/DEPLOYMENT.md)

---

**¡Bienvenido al proyecto!** 🎉

Si tienes dudas, revisa la documentación o abre un Issue.

