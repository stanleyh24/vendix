#!/bin/bash

# Script para recrear el tenant demo desde cero

set -e

echo "🔄 Recreando tenant demo..."
echo "=================================================="
echo ""

# Configuración
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-vendix}"
API_URL="${API_URL:-http://localhost:8080}"

export PGPASSWORD="$DB_PASSWORD"

# 1. Eliminar tenant existente
echo "1️⃣ Eliminando tenant demo existente..."

# Obtener ID y schema del tenant demo
TENANT_INFO=$(psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -t -c "
    SELECT id, schema_name 
    FROM public.tenants 
    WHERE slug = 'demo' OR name LIKE '%Demo%'
    LIMIT 1;
" 2>/dev/null | sed 's/^[ \t]*//;s/[ \t]*$//')

if [ -n "$TENANT_INFO" ]; then
    TENANT_ID=$(echo "$TENANT_INFO" | awk '{print $1}')
    SCHEMA_NAME=$(echo "$TENANT_INFO" | awk '{print $3}')
    
    echo "   Encontrado: ID=$TENANT_ID, Schema=$SCHEMA_NAME"
    
    # Eliminar schema
    if [ -n "$SCHEMA_NAME" ]; then
        echo "   Eliminando schema: $SCHEMA_NAME"
        psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c "
            DROP SCHEMA IF EXISTS $SCHEMA_NAME CASCADE;
        " > /dev/null 2>&1
        echo "   ✅ Schema eliminado"
    fi
    
    # Eliminar tenant de la tabla public.tenants
    if [ -n "$TENANT_ID" ]; then
        echo "   Eliminando tenant de public.tenants"
        psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c "
            DELETE FROM public.tenants WHERE id = '$TENANT_ID';
        " > /dev/null 2>&1
        echo "   ✅ Tenant eliminado"
    fi
else
    echo "   ℹ️  No se encontró tenant demo existente"
fi

echo ""

# 2. Crear nuevo tenant demo
echo "2️⃣ Creando nuevo tenant demo..."

CREATE_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/public/tenants" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Demo Company",
        "slug": "demo",
        "domain": null
    }')

NEW_TENANT_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$NEW_TENANT_ID" ]; then
    echo "   ❌ Error al crear tenant"
    echo "   Respuesta: $CREATE_RESPONSE"
    exit 1
fi

echo "   ✅ Tenant creado"
echo "   ID: $NEW_TENANT_ID"
echo "   Slug: demo"
echo "   Schema: tenant_demo"

echo ""

# 3. Inicializar schema (ejecutar migraciones)
echo "3️⃣ Inicializando schema del tenant..."

INIT_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/public/tenants/$NEW_TENANT_ID/init")

if echo "$INIT_RESPONSE" | grep -q "successfully"; then
    echo "   ✅ Schema inicializado (todas las tablas creadas)"
else
    echo "   ⚠️  Respuesta: $INIT_RESPONSE"
fi

echo ""

# 4. Verificar cuentas contables
echo "4️⃣ Verificando cuentas contables..."

sleep 2  # Esperar a que se apliquen las migraciones

ACCOUNTS_COUNT=$(psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -t -c "
    SELECT count(*) FROM tenant_demo.chart_of_accounts;
" 2>/dev/null | sed 's/^[ \t]*//;s/[ \t]*$//')

if [ "$ACCOUNTS_COUNT" -gt "0" ]; then
    echo "   ✅ Cuentas insertadas: $ACCOUNTS_COUNT"
else
    echo "   ⚠️  No se insertaron cuentas automáticamente"
    echo "   Ejecutando inserción manual..."
    
    psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" \
        -v schema=tenant_demo \
        -f "$(dirname "$0")/insert-chart-of-accounts.sql" \
        > /dev/null 2>&1
    
    ACCOUNTS_COUNT=$(psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -t -c "
        SELECT count(*) FROM tenant_demo.chart_of_accounts;
    " 2>/dev/null | sed 's/^[ \t]*//;s/[ \t]*$//')
    
    echo "   ✅ Cuentas insertadas: $ACCOUNTS_COUNT"
fi

echo ""

# 5. Crear usuario admin
echo "5️⃣ Creando usuario administrador..."

# Generar hash para password demo123
PASSWORD_HASH="\$2a\$10\$Rmc4ikuiVZ6fhWBMCOWG3urzWIboPHYvgalGku8prJqpLOnT78ztK"

# Establecer search_path y crear usuario
psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" << EOF > /dev/null 2>&1
SET search_path TO tenant_demo, public;

INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
VALUES (
    'admin@demo.com',
    '$PASSWORD_HASH',
    'Admin',
    'Demo',
    'admin',
    true
) ON CONFLICT (email) DO UPDATE SET password_hash = EXCLUDED.password_hash;
EOF

if [ $? -eq 0 ]; then
    echo "   ✅ Usuario administrador creado"
    echo "      Email: admin@demo.com"
    echo "      Password: demo123"
else
    echo "   ⚠️  El usuario ya existe o hubo un error"
fi

echo ""

# 6. Resumen
echo "=================================================="
echo "✅ Tenant demo recreado exitosamente"
echo ""
echo "📋 Información del tenant:"
echo "   ID: $NEW_TENANT_ID"
echo "   Nombre: Demo Company"
echo "   Slug: demo"
echo "   Schema: tenant_demo"
echo "   Cuentas contables: $ACCOUNTS_COUNT"
echo ""
echo "👤 Credenciales de acceso:"
echo "   Email: admin@demo.com"
echo "   Password: demo123"
echo "   Tenant ID: demo"
echo ""
echo "🌐 Puedes iniciar sesión en:"
echo "   http://localhost:5173/login"
echo ""
echo "=================================================="

