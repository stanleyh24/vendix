#!/bin/bash
# Script para agregar campos de inicio y fin de rango NCF
# Ejecuta la migración v20

set -e

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Migración: Rangos de NCF${NC}"
echo -e "${BLUE}  Versión: 20${NC}"
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
-- Agregar campos de inicio y fin de secuencia NCF para control de rangos DGII

-- NCF Crédito Fiscal (01)
ALTER TABLE tenant_config ADD COLUMN IF NOT EXISTS ncf_fiscal_credit_start INTEGER DEFAULT 1;
ALTER TABLE tenant_config ADD COLUMN IF NOT EXISTS ncf_fiscal_credit_end INTEGER DEFAULT 999999999;

-- NCF Consumidor Final (02)
ALTER TABLE tenant_config ADD COLUMN IF NOT EXISTS ncf_consumer_start INTEGER DEFAULT 1;
ALTER TABLE tenant_config ADD COLUMN IF NOT EXISTS ncf_consumer_end INTEGER DEFAULT 999999999;

-- NCF Nota de Débito (03)
ALTER TABLE tenant_config ADD COLUMN IF NOT EXISTS ncf_debit_note_start INTEGER DEFAULT 1;
ALTER TABLE tenant_config ADD COLUMN IF NOT EXISTS ncf_debit_note_end INTEGER DEFAULT 999999999;

-- NCF Nota de Crédito (04)
ALTER TABLE tenant_config ADD COLUMN IF NOT EXISTS ncf_credit_note_start INTEGER DEFAULT 1;
ALTER TABLE tenant_config ADD COLUMN IF NOT EXISTS ncf_credit_note_end INTEGER DEFAULT 999999999;

-- Agregar comentarios explicativos
COMMENT ON COLUMN tenant_config.ncf_fiscal_credit_start IS 'Inicio del rango de NCF tipo 01 asignado por DGII';
COMMENT ON COLUMN tenant_config.ncf_fiscal_credit_end IS 'Fin del rango de NCF tipo 01 asignado por DGII';
COMMENT ON COLUMN tenant_config.ncf_consumer_start IS 'Inicio del rango de NCF tipo 02 asignado por DGII';
COMMENT ON COLUMN tenant_config.ncf_consumer_end IS 'Fin del rango de NCF tipo 02 asignado por DGII';
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
            AND table_name = 'tenant_config' 
            AND column_name = 'ncf_fiscal_credit_start'
        );
    " | sed 's/^[ \t]*//;s/[ \t]*$//')
    
    if [ "$HAS_COLUMN" = "t" ]; then
        echo -e "  ${YELLOW}⊙ Ya tiene los campos de rango NCF - Saltando${NC}"
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
    echo -e "${YELLOW}Rangos de NCF:${NC}"
    echo "  Ahora cada tipo de NCF tiene:"
    echo "  • Inicio de rango (asignado por DGII)"
    echo "  • Fin de rango (asignado por DGII)"
    echo "  • Secuencia actual (se incrementa automáticamente)"
    echo ""
    echo -e "${YELLOW}Próximos pasos:${NC}"
    echo "1. Configurar los rangos asignados por DGII en Configuración → NCF"
    echo "2. El sistema alertará al llegar al 90% del rango"
    echo "3. Bloqueará al alcanzar el 100% del rango"
    echo ""
else
    echo -e "${YELLOW}⚠ Algunos tenants fallaron. Revisa los errores arriba.${NC}"
    exit 1
fi













