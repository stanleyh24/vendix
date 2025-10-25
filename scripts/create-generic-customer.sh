#!/bin/bash
# Script para crear cliente genérico en tenants que no lo tengan
# Ejecuta la migración v18

set -e

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Migración: Cliente Genérico${NC}"
echo -e "${BLUE}  Versión: 18${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Configuración de base de datos
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-vendix_db}"
DB_USER="${DB_USER:-vendix_user}"

echo -e "${YELLOW}Configuración:${NC}"
echo "  Host: $DB_HOST"
echo "  Puerto: $DB_PORT"
echo "  Base de datos: $DB_NAME"
echo "  Usuario: $DB_USER"
echo ""

# Verificar conexión
echo -e "${YELLOW}Verificando conexión a PostgreSQL...${NC}"
if ! PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\q' 2>/dev/null; then
    echo -e "${RED}Error: No se puede conectar a PostgreSQL${NC}"
    echo "Verifica las variables de entorno DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME"
    exit 1
fi
echo -e "${GREEN}✓ Conexión exitosa${NC}"
echo ""

# Obtener lista de tenants
echo -e "${YELLOW}Obteniendo lista de tenants...${NC}"
TENANTS=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT schema_name FROM tenants WHERE status = 'active'" | sed 's/^[ \t]*//;s/[ \t]*$//' | grep -v '^$')

if [ -z "$TENANTS" ]; then
    echo -e "${YELLOW}No se encontraron tenants activos${NC}"
    exit 0
fi

echo -e "${GREEN}Tenants encontrados:${NC}"
echo "$TENANTS" | while read tenant; do
    echo "  - $tenant"
done
echo ""

# SQL de migración
MIGRATION_SQL=$(cat <<'EOF'
-- Crear cliente genérico para ventas sin cliente específico
INSERT INTO customers (
    id,
    customer_type,
    tax_id,
    name,
    email,
    phone,
    address,
    city,
    state,
    country,
    is_active,
    metadata,
    created_at,
    updated_at
) 
SELECT 
    '00000000-0000-0000-0000-000000000001'::UUID,
    'individual',
    '000000000',
    'Cliente Genérico',
    'generico@sistema.local',
    NULL,
    NULL,
    NULL,
    NULL,
    'DO',
    true,
    '{"is_generic": true, "description": "Cliente genérico para ventas sin identificación específica"}'::JSONB,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
WHERE NOT EXISTS (
    SELECT 1 FROM customers WHERE id = '00000000-0000-0000-0000-000000000001'::UUID
);
EOF
)

# Contador
SUCCESS_COUNT=0
FAILED_COUNT=0
ALREADY_EXISTS=0
TOTAL_COUNT=$(echo "$TENANTS" | wc -l)

echo -e "${YELLOW}Creando cliente genérico en $TOTAL_COUNT tenant(s)...${NC}"
echo ""

# Iterar sobre cada tenant
echo "$TENANTS" | while read tenant; do
    if [ -z "$tenant" ]; then
        continue
    fi
    
    echo -e "${BLUE}Procesando: $tenant${NC}"
    
    # Verificar si ya existe el cliente genérico
    HAS_CUSTOMER=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
        SET search_path TO $tenant;
        SELECT EXISTS (
            SELECT 1 
            FROM customers 
            WHERE id = '00000000-0000-0000-0000-000000000001'::UUID
        );
    " | sed 's/^[ \t]*//;s/[ \t]*$//')
    
    if [ "$HAS_CUSTOMER" = "t" ]; then
        echo -e "  ${YELLOW}⊙ Cliente genérico ya existe - Saltando${NC}"
        ALREADY_EXISTS=$((ALREADY_EXISTS + 1))
    else
        # Crear cliente genérico
        if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
            SET search_path TO $tenant;
            $MIGRATION_SQL
        " > /dev/null 2>&1; then
            echo -e "  ${GREEN}✓ Cliente genérico creado exitosamente${NC}"
            SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
        else
            echo -e "  ${RED}✗ Error al crear cliente genérico${NC}"
            FAILED_COUNT=$((FAILED_COUNT + 1))
        fi
    fi
    
    echo ""
done

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Resumen de Migración${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "Total de tenants: $TOTAL_COUNT"
echo -e "${GREEN}Creados: $SUCCESS_COUNT${NC}"
echo -e "${YELLOW}Ya existían: $ALREADY_EXISTS${NC}"
if [ $FAILED_COUNT -gt 0 ]; then
    echo -e "${RED}Fallidos: $FAILED_COUNT${NC}"
fi
echo ""

if [ $FAILED_COUNT -eq 0 ]; then
    echo -e "${GREEN}✓ Migración completada exitosamente${NC}"
    echo ""
    echo -e "${YELLOW}Cliente Genérico:${NC}"
    echo "  ID: 00000000-0000-0000-0000-000000000001"
    echo "  Nombre: Cliente Genérico"
    echo "  Uso: Ventas sin identificación específica del cliente"
    echo ""
    echo -e "${YELLOW}Próximos pasos:${NC}"
    echo "1. Reiniciar el servidor backend"
    echo "2. Probar ventas con cliente genérico desde el frontend"
    echo "3. Verificar que las facturas se crean correctamente"
    echo ""
else
    echo -e "${YELLOW}⚠ Algunos tenants fallaron. Revisa los errores arriba.${NC}"
    exit 1
fi

