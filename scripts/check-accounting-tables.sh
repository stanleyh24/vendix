#!/bin/bash

# Script para verificar si las tablas de contabilidad existen en los tenants

set -e

# Configuración
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-vendix}"

echo "🔍 Verificando tablas de contabilidad en tenants..."
echo "=================================================="
echo "Host: $DB_HOST"
echo "Puerto: $DB_PORT"
echo "Base de datos: $DB_NAME"
echo ""

# Exportar password para psql
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
WITH_TABLE=0
WITHOUT_TABLE=0

# Verificar cada tenant
while IFS= read -r SCHEMA; do
    if [ -n "$SCHEMA" ]; then
        TOTAL=$((TOTAL + 1))
        
        # Verificar si existe la tabla chart_of_accounts
        EXISTS=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
            SELECT EXISTS (
                SELECT FROM information_schema.tables 
                WHERE table_schema = '$SCHEMA' 
                AND table_name = 'chart_of_accounts'
            );
        " | sed 's/^[ \t]*//;s/[ \t]*$//')
        
        if [ "$EXISTS" = "t" ]; then
            # Contar cuentas
            COUNT=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
                SELECT COUNT(*) FROM $SCHEMA.chart_of_accounts;
            " | sed 's/^[ \t]*//;s/[ \t]*$//')
            
            echo "✅ $SCHEMA - Tabla existe con $COUNT cuentas"
            WITH_TABLE=$((WITH_TABLE + 1))
        else
            echo "❌ $SCHEMA - Tabla NO existe (necesita migración)"
            WITHOUT_TABLE=$((WITHOUT_TABLE + 1))
        fi
    fi
done <<< "$TENANTS"

echo ""
echo "=================================================="
echo "📊 Resumen:"
echo "   Total tenants: $TOTAL"
echo "   Con tabla: $WITH_TABLE"
echo "   Sin tabla: $WITHOUT_TABLE"
echo ""

if [ $WITHOUT_TABLE -gt 0 ]; then
    echo "⚠️  ACCIÓN REQUERIDA:"
    echo "   Ejecuta: ./scripts/migrate-all-tenants-sql.sh"
    echo ""
    exit 1
else
    echo "✅ Todos los tenants tienen la tabla de contabilidad"
    exit 0
fi

