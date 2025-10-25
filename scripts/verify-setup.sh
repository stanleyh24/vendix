#!/bin/bash

# Script para verificar que la configuración de desarrollo está correcta

set -e

echo "🔍 Verificando configuración de desarrollo..."
echo ""

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Función para verificar archivos
check_file() {
    if [ -f "$1" ]; then
        echo -e "${GREEN}✅${NC} $1 existe"
        return 0
    else
        echo -e "${RED}❌${NC} $1 NO existe"
        return 1
    fi
}

# Función para verificar ejecutables
check_executable() {
    if [ -x "$1" ]; then
        echo -e "${GREEN}✅${NC} $1 es ejecutable"
        return 0
    else
        echo -e "${YELLOW}⚠️${NC}  $1 no es ejecutable, corrigiendo..."
        chmod +x "$1"
        if [ -x "$1" ]; then
            echo -e "${GREEN}✅${NC} $1 ahora es ejecutable"
            return 0
        else
            echo -e "${RED}❌${NC} No se pudo hacer ejecutable $1"
            return 1
        fi
    fi
}

# Verificar que estamos en el directorio correcto
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌${NC} Error: No estás en el directorio raíz del proyecto"
    exit 1
fi

echo "📁 Verificando archivos de configuración..."
echo ""

# Archivos Docker
check_file "docker-compose.dev.yml"
check_file "Dockerfile.dev"
check_file "frontend/Dockerfile.dev"
check_file ".dockerignore"

echo ""
echo "⚙️  Verificando archivos de configuración de Air..."
echo ""

check_file ".air.toml"
check_file ".air.worker.toml"

echo ""
echo "📚 Verificando documentación..."
echo ""

check_file "START_HERE.md"
check_file "QUICK_START.md"
check_file "DEV_SETUP.md"
check_file "CHANGELOG_DEV_SETUP.md"
check_file "SETUP_COMPLETO.md"
check_file ".github/DEVELOPMENT.md"

echo ""
echo "🛠️  Verificando scripts..."
echo ""

check_file "scripts/dev-setup.sh"
check_file "scripts/create-test-tenant.sh"
check_file "scripts/README.md"

echo ""
echo "🔧 Verificando que scripts sean ejecutables..."
echo ""

check_executable "scripts/dev-setup.sh"
check_executable "scripts/create-test-tenant.sh"

echo ""
echo "🐳 Verificando sintaxis de Docker Compose..."
echo ""

if docker-compose -f docker-compose.dev.yml config > /dev/null 2>&1; then
    echo -e "${GREEN}✅${NC} docker-compose.dev.yml es válido"
else
    echo -e "${RED}❌${NC} docker-compose.dev.yml tiene errores de sintaxis"
    exit 1
fi

if docker-compose -f docker-compose.yml config > /dev/null 2>&1; then
    echo -e "${GREEN}✅${NC} docker-compose.yml es válido"
else
    echo -e "${RED}❌${NC} docker-compose.yml tiene errores de sintaxis"
    exit 1
fi

echo ""
echo "🔍 Verificando que Docker esté disponible..."
echo ""

if command -v docker &> /dev/null; then
    echo -e "${GREEN}✅${NC} Docker está instalado"
    
    if docker info > /dev/null 2>&1; then
        echo -e "${GREEN}✅${NC} Docker está corriendo"
    else
        echo -e "${YELLOW}⚠️${NC}  Docker está instalado pero no está corriendo"
        echo "   Por favor, inicia Docker Desktop"
    fi
else
    echo -e "${RED}❌${NC} Docker no está instalado"
fi

if command -v docker-compose &> /dev/null; then
    echo -e "${GREEN}✅${NC} Docker Compose está instalado"
    COMPOSE_VERSION=$(docker-compose --version)
    echo "   Versión: $COMPOSE_VERSION"
else
    echo -e "${RED}❌${NC} Docker Compose no está instalado"
fi

echo ""
echo "📋 Verificando Makefile..."
echo ""

if grep -q "dev-up" Makefile; then
    echo -e "${GREEN}✅${NC} Comando 'make dev-up' está disponible"
else
    echo -e "${RED}❌${NC} Comando 'make dev-up' NO está disponible"
fi

if grep -q "dev-down" Makefile; then
    echo -e "${GREEN}✅${NC} Comando 'make dev-down' está disponible"
else
    echo -e "${RED}❌${NC} Comando 'make dev-down' NO está disponible"
fi

if grep -q "dev-logs" Makefile; then
    echo -e "${GREEN}✅${NC} Comando 'make dev-logs' está disponible"
else
    echo -e "${RED}❌${NC} Comando 'make dev-logs' NO está disponible"
fi

echo ""
echo "📦 Verificando archivos de proyecto..."
echo ""

check_file "go.mod"
check_file "frontend/package.json"
check_file "frontend/vite.config.js"

echo ""
echo "🎯 Resumen de Verificación"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo -e "${GREEN}✨ Configuración de desarrollo verificada correctamente${NC}"
echo ""
echo "📖 Próximos pasos:"
echo ""
echo "   1. Si aún no has iniciado el entorno:"
echo "      ${GREEN}bash scripts/dev-setup.sh${NC}"
echo ""
echo "   2. O inicia manualmente:"
echo "      ${GREEN}make dev-up${NC}"
echo ""
echo "   3. Ver logs:"
echo "      ${GREEN}make dev-logs${NC}"
echo ""
echo "   4. Crear tenant de prueba:"
echo "      ${GREEN}bash scripts/create-test-tenant.sh${NC}"
echo ""
echo "📚 Documentación:"
echo "   - Inicio rápido: ${GREEN}START_HERE.md${NC}"
echo "   - Guía completa:  ${GREEN}DEV_SETUP.md${NC}"
echo ""

