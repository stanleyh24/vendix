#!/bin/bash

# Script para configurar el entorno de desarrollo la primera vez

set -e

echo "🚀 Configurando entorno de desarrollo..."
echo ""

# Verificar que Docker esté corriendo
if ! docker info > /dev/null 2>&1; then
    echo "❌ Error: Docker no está corriendo"
    echo "   Por favor, inicia Docker y vuelve a intentar"
    exit 1
fi

echo "✅ Docker está corriendo"

# Crear archivo .env si no existe
if [ ! -f .env ]; then
    echo "📝 Creando archivo .env desde env.example..."
    cp env.example .env
    echo "✅ Archivo .env creado"
    echo "   Puedes editarlo si necesitas cambiar alguna configuración"
else
    echo "✅ Archivo .env ya existe"
fi

echo ""
echo "🔧 Iniciando servicios de desarrollo..."
echo "   Esto puede tomar unos minutos la primera vez..."
echo ""

# Iniciar docker-compose en modo desarrollo
docker-compose -f docker-compose.dev.yml up -d

echo ""
echo "⏳ Esperando a que los servicios estén listos..."
sleep 5

# Verificar que Postgres esté listo
echo "   Verificando PostgreSQL..."
until docker exec vendix-postgres-dev pg_isready -U postgres > /dev/null 2>&1; do
    echo "   PostgreSQL aún no está listo, esperando..."
    sleep 2
done
echo "   ✅ PostgreSQL listo"

# Verificar que Redis esté listo
echo "   Verificando Redis..."
until docker exec vendix-redis-dev redis-cli ping > /dev/null 2>&1; do
    echo "   Redis aún no está listo, esperando..."
    sleep 2
done
echo "   ✅ Redis listo"

# Verificar que el API esté listo
echo "   Verificando API..."
retries=0
max_retries=30
until curl -s http://localhost:8080/health > /dev/null 2>&1; do
    if [ $retries -eq $max_retries ]; then
        echo "   ⚠️  El API está tardando más de lo esperado"
        echo "   Puedes verificar los logs con: make dev-logs"
        break
    fi
    echo "   API aún no está listo, esperando..."
    sleep 2
    retries=$((retries + 1))
done

if [ $retries -lt $max_retries ]; then
    echo "   ✅ API listo"
fi

echo ""
echo "✨ ¡Entorno de desarrollo configurado!"
echo ""
echo "📍 URLs disponibles:"
echo "   Frontend:  http://localhost:3000"
echo "   API:       http://localhost:8080"
echo "   MinIO:     http://localhost:9001 (minioadmin/minioadmin)"
echo ""
echo "📋 Próximos pasos:"
echo "   1. Ver logs:           make dev-logs"
echo "   2. Crear tenant:       bash scripts/create-test-tenant.sh"
echo "   3. Detener servicios:  make dev-down"
echo ""
echo "📚 Documentación:"
echo "   - Guía rápida:    QUICK_START.md"
echo "   - Desarrollo:     DEV_SETUP.md"
echo "   - README:         README.md"
echo ""

