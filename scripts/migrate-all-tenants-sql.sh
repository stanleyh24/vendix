#!/bin/bash

# Script para ejecutar migraciones directamente con SQL en todos los tenants
# Este script NO requiere que el servidor esté corriendo

set -e

# Configuración
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-vendix}"

echo "🔄 Ejecutando migraciones SQL en todos los tenants..."
echo "=================================================="
echo "Host: $DB_HOST"
echo "Puerto: $DB_PORT"
echo "Base de datos: $DB_NAME"
echo ""

# Exportar password para psql
export PGPASSWORD="$DB_PASSWORD"

# Obtener lista de schemas de tenants
SCHEMAS=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
    SELECT schema_name 
    FROM public.tenants 
    WHERE deleted_at IS NULL 
    ORDER BY created_at;
" | sed 's/^[ \t]*//;s/[ \t]*$//' | grep -v '^$')

if [ -z "$SCHEMAS" ]; then
    echo "❌ No se encontraron tenants en la base de datos"
    exit 1
fi

echo "📋 Schemas de tenants encontrados:"
echo "$SCHEMAS"
echo ""

# Contador
TOTAL=0
SUCCESS=0
FAILED=0

# Ejecutar migraciones para cada schema
while IFS= read -r SCHEMA; do
    if [ -n "$SCHEMA" ]; then
        TOTAL=$((TOTAL + 1))
        echo "🔄 Migrando schema: $SCHEMA"
        
        # Ejecutar el script SQL
        if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" \
            -v schema="$SCHEMA" \
            -f "$(dirname "$0")/migrate-tenant-direct.sql" \
            > /dev/null 2>&1; then
            echo "   ✅ Migración exitosa"
            SUCCESS=$((SUCCESS + 1))
        else
            echo "   ❌ Error en la migración"
            FAILED=$((FAILED + 1))
        fi
        
        echo ""
    fi
done <<< "$SCHEMAS"

# Resumen
echo "=================================================="
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

