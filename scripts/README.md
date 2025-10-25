# Scripts de Utilidades

Este directorio contiene scripts útiles para desarrollo y despliegue.

## 🚀 Scripts de Desarrollo

### `dev-setup.sh`

Configura automáticamente el entorno de desarrollo la primera vez.

**Uso:**
```bash
bash scripts/dev-setup.sh
```

**Qué hace:**
- ✅ Verifica que Docker esté corriendo
- ✅ Crea `.env` desde `env.example` si no existe
- ✅ Inicia todos los servicios en modo desarrollo
- ✅ Espera a que todos los servicios estén listos
- ✅ Verifica conectividad de servicios

**Cuándo usar:**
- Primera vez que clonas el repositorio
- Después de hacer `git clean` o borrar configuración
- Para resetear el entorno de desarrollo

---

### `create-test-tenant.sh`

Crea un tenant de prueba completo con usuario admin.

**Uso:**
```bash
bash scripts/create-test-tenant.sh
```

**Qué hace:**
- Crea un nuevo tenant en la base de datos
- Inicializa el schema del tenant
- Crea un usuario administrador
- Muestra las credenciales

**Cuándo usar:**
- Para testing local
- Para demos
- Para desarrollo de features

---

### `verify-setup.sh`

Verifica que toda la configuración de desarrollo esté correcta.

**Uso:**
```bash
bash scripts/verify-setup.sh
```

**Qué hace:**
- ✅ Verifica que todos los archivos necesarios existen
- ✅ Valida sintaxis de docker-compose files
- ✅ Verifica que Docker esté instalado y corriendo
- ✅ Comprueba que los scripts sean ejecutables
- ✅ Verifica comandos Make disponibles
- ✅ Muestra resumen con próximos pasos

**Cuándo usar:**
- Después de clonar el repositorio
- Para diagnosticar problemas de configuración
- Antes de reportar un issue
- Para verificar que todo esté listo para desarrollar

---

## 📝 Convenciones

### Estructura de Scripts

Todos los scripts deben:
- Tener `#!/bin/bash` al inicio
- Usar `set -e` para fallar en errores
- Incluir mensajes descriptivos con emojis
- Ser ejecutables: `chmod +x script.sh`
- Incluir comentarios explicativos

### Ejemplo de estructura:

```bash
#!/bin/bash

# Descripción del script

set -e

echo "🚀 Iniciando proceso..."

# Verificaciones
if [ condición ]; then
    echo "❌ Error: mensaje"
    exit 1
fi

echo "✅ Proceso completado"
```

### Mensajes

- ✅ Éxito
- ❌ Error
- ⚠️  Advertencia
- 🚀 Inicio de proceso
- 📝 Creación/Escritura
- 🔧 Configuración
- ⏳ Espera/Carga
- 📍 Información de ubicación
- 💡 Tip/Sugerencia

---

## 🆕 Agregar Nuevos Scripts

1. Crear el script en este directorio
2. Hacerlo ejecutable: `chmod +x nuevo-script.sh`
3. Documentarlo en este README
4. Seguir las convenciones de mensajes
5. Incluir manejo de errores

---

## 🔍 Scripts Planeados

Ideas para futuros scripts:

- `backup-database.sh` - Backup de la base de datos
- `restore-database.sh` - Restaurar desde backup
- `seed-demo-data.sh` - Llenar con datos de demo
- `check-health.sh` - Verificar salud de servicios
- `generate-ssl-certs.sh` - Generar certificados SSL para desarrollo
- `cleanup-dev.sh` - Limpiar datos de desarrollo
- `migrate-tenant.sh` - Migrar schema de tenant específico

---

## 💡 Tips

### Ejecutar scripts desde cualquier directorio

```bash
# Desde la raíz del proyecto
bash scripts/dev-setup.sh

# Desde el directorio scripts
./dev-setup.sh

# Usando Make (si está configurado)
make setup-dev
```

### Debug de scripts

```bash
# Modo verbose
bash -x scripts/dev-setup.sh

# Ver solo errores
bash scripts/dev-setup.sh 2>&1 | grep "Error"
```

### Variables de entorno

Los scripts pueden usar variables de entorno:

```bash
# Usar puerto personalizado
API_PORT=9000 bash scripts/dev-setup.sh

# Usar base de datos diferente
DB_NAME=testing bash scripts/create-test-tenant.sh
```

---

## 📚 Recursos

- [Bash Scripting Guide](https://www.gnu.org/software/bash/manual/bash.html)
- [Shell Check](https://www.shellcheck.net/) - Linter para bash
- [Google Shell Style Guide](https://google.github.io/styleguide/shellguide.html)

