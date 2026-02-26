# Guía de Despliegue en Producción - Vendix

Esta guía detalla los pasos para desplegar Vendix en un entorno de producción utilizando Docker.

## 📋 Requisitos Previos

- Servidor con Linux (Ubuntu 22.04+ recomendado).
- Docker y Docker Compose instalados.
- Un dominio apuntando a la IP del servidor.
- Credenciales de la DGII (opcional para pruebas, requerido para producción real).

## 🚀 Pasos para el Despliegue

### 1. Preparación del Entorno

Clona el repositorio en el servidor de producción:

```bash
git clone <repository-url>
cd vendix
```

Crea el archivo de variables de entorno de producción:

```bash
cp env.example .env
```

### 2. Configuración Crítica del `.env`

Edita el archivo `.env` y asegúrate de configurar los siguientes valores:

```ini
# App
APP_ENV=production
JWT_SECRET=generar-un-secreto-muy-largo-y-seguro

# Database
DB_PASSWORD=contrasena-fuerte-para-postgres

# DGII
DGII_ENABLED=true
DGII_API_URL=https://api.dgii.gov.do
DGII_CERT_PATH=/app/certs/certificado.pem
DGII_KEY_PATH=/app/certs/certificado.key

# Storage (MinIO)
S3_ACCESS_KEY=usuario-minio
S3_SECRET_KEY=contrasena-minio
```

### 3. Certificados DGII

Crea una carpeta `certs` en la raíz del proyecto y coloca tus certificados de firma digital (si tienes `DGII_ENABLED=true`):

```bash
mkdir certs
# Copia tus archivos .pem y .key aquí
```

### 4. Lanzamiento con Docker Compose

Utiliza el archivo específico de producción:

```bash
docker-compose -f docker-compose.prod.yml up -d --build
```

### 5. Configuración de SSL (Opcional pero Recomendado)

Aunque el `docker-compose.prod.yml` expone el puerto 80, se recomienda usar un reverse proxy como **Nginx** o **Caddy** para habilitar HTTPS.

#### Ejemplo de Caddy (más sencillo):

```caddyfile
vendix.tudominio.com {
    reverse_proxy localhost:80
}
```

### 6. Inicialización del Sistema

Una vez los contenedores estén corriendo, debes crear el primer tenant:

```bash
# Cambia 'demo' por el nombre deseado
bash scripts/create-test-tenant.sh
```

## 🛠 Comandos de Mantenimiento

- **Ver logs:** `docker-compose -f docker-compose.prod.yml logs -f`
- **Actualizar sistema:** 
  ```bash
  git pull
  docker-compose -f docker-compose.prod.yml up -d --build
  ```
- **Backup de DB:**
  ```bash
  docker exec vendix-postgres-prod pg_dump -U postgres vendix > backup.sql
  ```

## 🔐 Seguridad

1. **Puertos:** Solo el puerto 80 (y 443 si usas SSL) debe estar abierto al público. Los puertos de PostgreSQL (5432) y Redis (6379) están protegidos internamente por Docker.
2. **JWT:** Asegúrate de que `JWT_SECRET` sea único por cada instalación.
3. **Passwords:** Cambia todas las contraseñas por defecto en el archivo `.env`.
