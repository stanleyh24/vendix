#!/bin/bash

# Script para agregar cash register por defecto a todos los tenants existentes
# que no tengan uno. Primero ejecuta las migraciones pendientes si es necesario.

set -e

# Configuración
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-vendix}"
API_URL="${API_URL:-http://localhost:8080}"

echo "🏦 Agregando cash register por defecto a tenants..."
echo "==================================================="
echo "Host: $DB_HOST"
echo "Puerto: $DB_PORT"
echo "Base de datos: $DB_NAME"
echo "API URL: $API_URL"
echo ""

# Exportar password para psql
export PGPASSWORD="$DB_PASSWORD"

# Verificar si la API está disponible
API_AVAILABLE=false
if curl -s -f "${API_URL}/health" > /dev/null 2>&1; then
    API_AVAILABLE=true
    echo "✅ API disponible"
else
    echo "⚠️  API no disponible - solo se ejecutarán migraciones SQL directas si la tabla no existe"
fi
echo ""

# Obtener lista de tenants con sus IDs y schemas
TENANTS_QUERY="SELECT id::text, schema_name FROM public.tenants WHERE deleted_at IS NULL ORDER BY created_at;"

TENANTS=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -A -F'|' -c "$TENANTS_QUERY" | sed 's/^[ \t]*//;s/[ \t]*$//' | grep -v '^$')

if [ -z "$TENANTS" ]; then
    echo "❌ No se encontraron tenants en la base de datos"
    exit 1
fi

TENANT_COUNT=$(echo "$TENANTS" | wc -l)
echo "📋 Tenants encontrados: $TENANT_COUNT"
echo ""

# Contador
TOTAL=0
SUCCESS=0
SKIPPED=0
FAILED=0
MIGRATED=0

# Procesar cada tenant
while IFS='|' read -r TENANT_ID SCHEMA; do
    TENANT_ID=$(echo "$TENANT_ID" | sed 's/^[ \t]*//;s/[ \t]*$//')
    SCHEMA=$(echo "$SCHEMA" | sed 's/^[ \t]*//;s/[ \t]*$//')
    
    if [ -n "$SCHEMA" ] && [ -n "$TENANT_ID" ]; then
        TOTAL=$((TOTAL + 1))
        echo "🔄 Procesando: $SCHEMA (ID: $TENANT_ID)"
        
        # Verificar si existe la tabla cash_registers
        TABLE_EXISTS=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
            SELECT EXISTS (
                SELECT FROM information_schema.tables 
                WHERE table_schema = '$SCHEMA' 
                AND table_name = 'cash_registers'
            );
        " | sed 's/^[ \t]*//;s/[ \t]*$//')
        
        if [ "$TABLE_EXISTS" != "t" ]; then
            echo "   ⚠️  Tabla cash_registers no existe - ejecutando migraciones..."
            
            # Intentar ejecutar migraciones a través de la API si está disponible
            if [ "$API_AVAILABLE" = true ]; then
                RESPONSE=$(curl -s -X POST "${API_URL}/api/v1/public/tenants/${TENANT_ID}/init" \
                    -H "Content-Type: application/json" \
                    -w "\n%{http_code}")
                
                HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)
                BODY=$(echo "$RESPONSE" | head -n -1)
                
                if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "201" ]; then
                    echo "   ✅ Migraciones ejecutadas exitosamente vía API"
                    MIGRATED=$((MIGRATED + 1))
                    
                    # Verificar nuevamente si la tabla existe después de las migraciones
                    sleep 1
                    TABLE_EXISTS=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
                        SELECT EXISTS (
                            SELECT FROM information_schema.tables 
                            WHERE table_schema = '$SCHEMA' 
                            AND table_name = 'cash_registers'
                        );
                    " | sed 's/^[ \t]*//;s/[ \t]*$//')
                    
                    if [ "$TABLE_EXISTS" != "t" ]; then
                        echo "   ⚠️  Tabla aún no existe después de migraciones - omitiendo"
                        SKIPPED=$((SKIPPED + 1))
                        echo ""
                        continue
                    fi
                else
                    echo "   ❌ Error al ejecutar migraciones: HTTP $HTTP_CODE"
                    if [ -n "$BODY" ]; then
                        echo "   Mensaje: $BODY"
                    fi
                    echo "   ⚠️  Omitiendo este tenant"
                    SKIPPED=$((SKIPPED + 1))
                    echo ""
                    continue
                fi
            else
                echo "   ❌ API no disponible - no se pueden ejecutar migraciones automáticamente"
                echo "   💡 Ejecuta primero: make migrate-tenants"
                echo "   O ejecuta: curl -X POST ${API_URL}/api/v1/public/tenants/${TENANT_ID}/init"
                SKIPPED=$((SKIPPED + 1))
                echo ""
                continue
            fi
        fi
        
        # Contar cash registers existentes
        COUNT=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
            SELECT count(*) FROM $SCHEMA.cash_registers;
        " | sed 's/^[ \t]*//;s/[ \t]*$//')
        
        if [ "$COUNT" -gt "0" ]; then
            echo "   ℹ️  Ya tiene $COUNT cash register(s) - omitiendo"
            SKIPPED=$((SKIPPED + 1))
            echo ""
            continue
        fi
        
        # Crear cash register por defecto
        echo "   📝 Creando cash register por defecto..."
        
        if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
            INSERT INTO $SCHEMA.cash_registers (
                id,
                name,
                location,
                is_active,
                initial_balance,
                created_by,
                created_at,
                updated_at
            ) VALUES (
                gen_random_uuid(),
                'Caja Principal',
                NULL,
                true,
                0.00,
                NULL,
                CURRENT_TIMESTAMP,
                CURRENT_TIMESTAMP
            );
        " > /dev/null 2>&1; then
            
            # Verificar que se creó correctamente
            NEW_COUNT=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
                SELECT count(*) FROM $SCHEMA.cash_registers;
            " | sed 's/^[ \t]*//;s/[ \t]*$//')
            
            if [ "$NEW_COUNT" -gt "0" ]; then
                echo "   ✅ Cash register creado exitosamente (total: $NEW_COUNT)"
                SUCCESS=$((SUCCESS + 1))
            else
                echo "   ⚠️  No se pudo verificar la creación"
                FAILED=$((FAILED + 1))
            fi
        else
            echo "   ❌ Error al crear cash register"
            FAILED=$((FAILED + 1))
        fi
        
        echo ""
    fi
done <<< "$TENANTS"

# Resumen
echo "==================================================="
echo "📊 Resumen:"
echo "   Total procesados: $TOTAL"
echo "   Migrados: $MIGRATED"
echo "   Exitosos: $SUCCESS"
echo "   Omitidos: $SKIPPED"
echo "   Fallidos: $FAILED"
echo ""

if [ $FAILED -eq 0 ]; then
    if [ $SUCCESS -gt 0 ]; then
        echo "✅ Cash registers creados correctamente en $SUCCESS tenant(s)"
    else
        if [ $MIGRATED -gt 0 ]; then
            echo "ℹ️  Se ejecutaron migraciones en $MIGRATED tenant(s), pero ya tenían cash registers o no se pudo crear"
        else
            echo "ℹ️  No había tenants sin cash registers"
        fi
    fi
    exit 0
else
    echo "⚠️  Algunos tenants fallaron al crear cash registers"
    exit 1
fi

