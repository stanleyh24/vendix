#!/bin/bash
# Script para agregar campo ncf_type a invoices
# Ejecuta la migración v19

set -e

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Migración: Tipo de Comprobante NCF${NC}"
echo -e "${BLUE}  Versión: 19${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Configuración de base de datos
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-vendix}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"

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
-- Agregar campo para tipo de comprobante fiscal (NCF)
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS ncf_type VARCHAR(2) DEFAULT '02' NOT NULL;

-- Crear índice para búsquedas por tipo de NCF
CREATE INDEX IF NOT EXISTS idx_invoices_ncf_type ON invoices(ncf_type);

-- Agregar comentario explicativo
COMMENT ON COLUMN invoices.ncf_type IS 'Tipo de comprobante fiscal: 01=Crédito Fiscal, 02=Consumidor Final';
EOF
)

# Contador
SUCCESS_COUNT=0
FAILED_COUNT=0
ALREADY_EXISTS=0
TOTAL_COUNT=$(echo "$TENANTS" | wc -l)

echo -e "${YELLOW}Aplicando migración a $TOTAL_COUNT tenant(s)...${NC}"
echo ""

# Iterar sobre cada tenant
echo "$TENANTS" | while read tenant; do
    if [ -z "$tenant" ]; then
        continue
    fi
    
    echo -e "${BLUE}Procesando: $tenant${NC}"
    
    # Verificar si ya existe el campo
    HAS_COLUMN=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
        SET search_path TO $tenant;
        SELECT EXISTS (
            SELECT 1 
            FROM information_schema.columns 
            WHERE table_schema = '$tenant' 
            AND table_name = 'invoices' 
            AND column_name = 'ncf_type'
        );
    " | sed 's/^[ \t]*//;s/[ \t]*$//')
    
    if [ "$HAS_COLUMN" = "t" ]; then
        echo -e "  ${YELLOW}⊙ Ya tiene la columna ncf_type - Saltando${NC}"
        ALREADY_EXISTS=$((ALREADY_EXISTS + 1))
    else
        # Aplicar migración
        if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
            SET search_path TO $tenant;
            $MIGRATION_SQL
        " > /dev/null 2>&1; then
            echo -e "  ${GREEN}✓ Migración aplicada exitosamente${NC}"
            SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
        else
            echo -e "  ${RED}✗ Error al aplicar migración${NC}"
            FAILED_COUNT=$((FAILED_COUNT + 1))
        fi
    fi
    
    echo ""
done

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Resumen de Migración${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "Total de tenants: $TOTAL_COUNT"
echo -e "${GREEN}Exitosos: $SUCCESS_COUNT${NC}"
echo -e "${YELLOW}Ya existían: $ALREADY_EXISTS${NC}"
if [ $FAILED_COUNT -gt 0 ]; then
    echo -e "${RED}Fallidos: $FAILED_COUNT${NC}"
fi
echo ""

if [ $FAILED_COUNT -eq 0 ]; then
    echo -e "${GREEN}✓ Migración completada exitosamente${NC}"
    echo ""
    echo -e "${YELLOW}Tipo de Comprobante Fiscal (NCF):${NC}"
    echo "  • NCF 01 = Crédito Fiscal (empresas con RNC)"
    echo "  • NCF 02 = Consumidor Final (default)"
    echo ""
    echo -e "${YELLOW}Próximos pasos:${NC}"
    echo "1. Reiniciar el servidor backend"
    echo "2. Probar selector de NCF en punto de venta"
    echo "3. Verificar que facturas se crean con el tipo correcto"
    echo ""
else
    echo -e "${YELLOW}⚠ Algunos tenants fallaron. Revisa los errores arriba.${NC}"
    exit 1
fi

