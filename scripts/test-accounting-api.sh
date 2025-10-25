#!/bin/bash

# Script para probar el endpoint de cuentas contables

echo "🧪 Probando API de Contabilidad"
echo "================================"
echo ""

# Configuración
API_URL="${API_URL:-http://localhost:8080}"
TENANT_ID="${TENANT_ID:-demo}"

echo "📍 API URL: $API_URL"
echo "🏢 Tenant ID: $TENANT_ID"
echo ""

# 1. Verificar que el backend esté corriendo
echo "1️⃣ Verificando backend..."
if curl -s -f "$API_URL/health" > /dev/null 2>&1; then
    echo "   ✅ Backend está corriendo"
    curl -s "$API_URL/health" | head -3
else
    echo "   ❌ Backend NO está corriendo"
    echo "   💡 Inicia el backend con: make run-api"
    exit 1
fi

echo ""

# 2. Probar sin autenticación (debe fallar)
echo "2️⃣ Probando endpoint sin autenticación..."
RESPONSE=$(curl -s -w "\n%{http_code}" "$API_URL/api/v1/tenant/accounting/accounts" \
    -H "X-Tenant-ID: $TENANT_ID" 2>&1)

HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)
BODY=$(echo "$RESPONSE" | head -n -1)

if [ "$HTTP_CODE" = "401" ] || [ "$HTTP_CODE" = "403" ]; then
    echo "   ✅ Correctamente requiere autenticación (HTTP $HTTP_CODE)"
else
    echo "   ⚠️  Respuesta inesperada: HTTP $HTTP_CODE"
    echo "   Body: $BODY"
fi

echo ""

# 3. Instrucciones para probar con token
echo "3️⃣ Para probar con autenticación:"
echo ""
echo "   # Primero, obtén un token:"
echo "   curl -X POST $API_URL/api/v1/tenant/auth/login \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -H 'X-Tenant-ID: $TENANT_ID' \\"
echo "     -d '{\"email\":\"admin@demo.com\",\"password\":\"demo123\"}'"
echo ""
echo "   # Luego, usa el token para consultar cuentas:"
echo "   curl $API_URL/api/v1/tenant/accounting/accounts \\"
echo "     -H 'Authorization: Bearer {TU_TOKEN}' \\"
echo "     -H 'X-Tenant-ID: $TENANT_ID'"
echo ""

# 4. Verificar en base de datos
echo "4️⃣ Verificando base de datos..."
echo ""

DB_HOST="${DB_HOST:-localhost}"
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-vendix}"

export PGPASSWORD="${DB_PASSWORD:-postgres}"

# Obtener el schema del tenant
SCHEMA=$(psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -t -c "
    SELECT schema_name 
    FROM public.tenants 
    WHERE slug = '$TENANT_ID' OR slug LIKE '%$TENANT_ID%'
    LIMIT 1;
" 2>/dev/null | sed 's/^[ \t]*//;s/[ \t]*$//')

if [ -z "$SCHEMA" ]; then
    echo "   ❌ No se encontró el tenant '$TENANT_ID' en la base de datos"
    echo ""
    echo "   Tenants disponibles:"
    psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c "
        SELECT slug, schema_name FROM public.tenants WHERE deleted_at IS NULL;
    " 2>/dev/null || echo "   Error al conectar con la base de datos"
else
    echo "   ✅ Tenant encontrado: $SCHEMA"
    
    # Verificar tabla
    TABLE_EXISTS=$(psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -t -c "
        SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_schema = '$SCHEMA' 
            AND table_name = 'chart_of_accounts'
        );
    " 2>/dev/null | sed 's/^[ \t]*//;s/[ \t]*$//')
    
    if [ "$TABLE_EXISTS" = "t" ]; then
        echo "   ✅ Tabla chart_of_accounts existe"
        
        # Contar cuentas
        COUNT=$(psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -t -c "
            SELECT count(*) FROM $SCHEMA.chart_of_accounts;
        " 2>/dev/null | sed 's/^[ \t]*//;s/[ \t]*$//')
        
        if [ "$COUNT" -gt "0" ]; then
            echo "   ✅ Tiene $COUNT cuentas"
            echo ""
            echo "   Primeras 5 cuentas:"
            psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c "
                SELECT account_code, account_name, account_type 
                FROM $SCHEMA.chart_of_accounts 
                ORDER BY account_code 
                LIMIT 5;
            " 2>/dev/null
        else
            echo "   ❌ La tabla está VACÍA (0 cuentas)"
            echo ""
            echo "   💡 Solución: make insert-accounts"
        fi
    else
        echo "   ❌ Tabla chart_of_accounts NO existe"
        echo ""
        echo "   💡 Solución: make migrate-tenants"
    fi
fi

echo ""
echo "================================"
echo "🏁 Diagnóstico completo"
echo ""
echo "📋 Checklist:"
echo "   1. Backend corriendo: Verificar arriba"
echo "   2. Tenant existe: Verificar arriba"
echo "   3. Tabla existe: Verificar arriba"
echo "   4. Tabla tiene datos: Verificar arriba"
echo "   5. Autenticación: Necesitas token válido"
echo ""
echo "💡 Siguiente paso:"
echo "   - Si falta la tabla: make migrate-tenants"
echo "   - Si la tabla está vacía: make insert-accounts"
echo "   - Si todo está bien: Verifica la consola del frontend (F12)"

