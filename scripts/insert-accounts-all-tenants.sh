#!/bin/bash

# Script para insertar cuentas en todos los tenants que tienen la tabla pero sin datos

set -e

# Configuración
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-vendix}"

echo "💰 Insertando cuentas en tenants sin datos..."
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
SKIPPED=0
FAILED=0

# Procesar cada schema
while IFS= read -r SCHEMA; do
    if [ -n "$SCHEMA" ]; then
        TOTAL=$((TOTAL + 1))
        echo "🔄 Procesando: $SCHEMA"
        
        # Verificar si existe la tabla
        TABLE_EXISTS=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
            SELECT EXISTS (
                SELECT FROM information_schema.tables 
                WHERE table_schema = '$SCHEMA' 
                AND table_name = 'chart_of_accounts'
            );
        " | sed 's/^[ \t]*//;s/[ \t]*$//')
        
        if [ "$TABLE_EXISTS" != "t" ]; then
            echo "   ⚠️  Tabla no existe - necesita ejecutar migración primero"
            SKIPPED=$((SKIPPED + 1))
            echo ""
            continue
        fi
        
        # Contar cuentas existentes
        COUNT=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
            SELECT count(*) FROM $SCHEMA.chart_of_accounts;
        " | sed 's/^[ \t]*//;s/[ \t]*$//')
        
        if [ "$COUNT" -gt "0" ]; then
            echo "   ℹ️  Ya tiene $COUNT cuentas - omitiendo"
            SKIPPED=$((SKIPPED + 1))
            echo ""
            continue
        fi
        
        # Insertar cuentas
        if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" \
            -v schema="$SCHEMA" \
            -f "$(dirname "$0")/insert-chart-of-accounts.sql" \
            > /dev/null 2>&1; then
            
            # Verificar cuántas se insertaron
            NEW_COUNT=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
                SELECT count(*) FROM $SCHEMA.chart_of_accounts;
            " | sed 's/^[ \t]*//;s/[ \t]*$//')
            
            echo "   ✅ Insertadas $NEW_COUNT cuentas exitosamente"
            SUCCESS=$((SUCCESS + 1))
        else
            echo "   ❌ Error al insertar cuentas"
            FAILED=$((FAILED + 1))
        fi
        
        echo ""
    fi
done <<< "$SCHEMAS"

# Resumen
echo "=================================================="
echo "📊 Resumen:"
echo "   Total procesados: $TOTAL"
echo "   Exitosos: $SUCCESS"
echo "   Omitidos: $SKIPPED"
echo "   Fallidos: $FAILED"
echo ""

if [ $FAILED -eq 0 ]; then
    if [ $SUCCESS -gt 0 ]; then
        echo "✅ Cuentas insertadas correctamente en $SUCCESS tenant(s)"
    else
        echo "ℹ️  No había tenants sin cuentas"
    fi
    exit 0
else
    echo "⚠️  Algunos tenants fallaron al insertar cuentas"
    exit 1
fi

