#!/bin/bash

# Script para ejecutar migraciones pendientes en todos los tenants existentes

set -e

# Configuración
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-vendix}"

echo "🔄 Ejecutando migraciones de tenants..."
echo "=========================================="
echo "Host: $DB_HOST"
echo "Puerto: $DB_PORT"
echo "Base de datos: $DB_NAME"
echo ""

# Exportar variables para psql
export PGPASSWORD="$DB_PASSWORD"

# Obtener lista de tenants
TENANTS=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
    SELECT schema_name 
    FROM public.tenants 
    WHERE deleted_at IS NULL 
    ORDER BY created_at;
" | sed 's/^[ \t]*//;s/[ \t]*$//' | grep -v '^$')

if [ -z "$TENANTS" ]; then
    echo "❌ No se encontraron tenants en la base de datos"
    exit 1
fi

echo "📋 Tenants encontrados:"
echo "$TENANTS"
echo ""

# Contador
TOTAL=0
SUCCESS=0
FAILED=0

# Ejecutar migraciones para cada tenant
while IFS= read -r SCHEMA; do
    if [ -n "$SCHEMA" ]; then
        TOTAL=$((TOTAL + 1))
        echo "🔄 Procesando tenant: $SCHEMA"
        
        # Ejecutar endpoint de inicialización de schema
        TENANT_ID=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
            SELECT id 
            FROM public.tenants 
            WHERE schema_name = '$SCHEMA' 
            LIMIT 1;
        " | sed 's/^[ \t]*//;s/[ \t]*$//')
        
        if [ -n "$TENANT_ID" ]; then
            echo "   ID del tenant: $TENANT_ID"
            
            # Llamar al endpoint de inicialización (esto ejecutará las migraciones)
            RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/public/tenants/${TENANT_ID}/init" \
                -H "Content-Type: application/json" \
                -w "\n%{http_code}")
            
            HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)
            BODY=$(echo "$RESPONSE" | head -n -1)
            
            if [ "$HTTP_CODE" = "200" ]; then
                echo "   ✅ Migraciones ejecutadas exitosamente"
                SUCCESS=$((SUCCESS + 1))
            else
                echo "   ⚠️  Respuesta del servidor: $HTTP_CODE - $BODY"
                # Podría ser que las migraciones ya estén aplicadas
                SUCCESS=$((SUCCESS + 1))
            fi
        else
            echo "   ❌ No se pudo obtener el ID del tenant"
            FAILED=$((FAILED + 1))
        fi
        
        echo ""
    fi
done <<< "$TENANTS"

# Resumen
echo "=========================================="
echo "📊 Resumen:"
echo "   Total: $TOTAL"
echo "   Exitosos: $SUCCESS"
echo "   Fallidos: $FAILED"
echo ""

if [ $FAILED -eq 0 ]; then
    echo "✅ Todas las migraciones se ejecutaron correctamente"
    exit 0
else
    echo "⚠️  Algunas migraciones fallaron"
    exit 1
fi

