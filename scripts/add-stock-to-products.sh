#!/bin/bash
# Script para agregar campo stock_quantity a productos en tenants existentes
# Ejecuta la migración v17 que agrega control de inventario

set -e

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Migración: Sistema de Inventario${NC}"
echo -e "${BLUE}  Versión: 17${NC}"
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
-- Agregar campo de cantidad en stock a la tabla de productos
ALTER TABLE products ADD COLUMN IF NOT EXISTS stock_quantity DECIMAL(10, 2) DEFAULT 0 NOT NULL;

-- Agregar índice para búsquedas rápidas de productos con bajo stock
CREATE INDEX IF NOT EXISTS idx_products_stock_quantity ON products(stock_quantity);

-- Crear tabla para registro de movimientos de inventario (auditoría)
CREATE TABLE IF NOT EXISTS inventory_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    movement_type VARCHAR(50) NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL,
    previous_stock DECIMAL(10, 2) NOT NULL,
    new_stock DECIMAL(10, 2) NOT NULL,
    reference_type VARCHAR(50),
    reference_id UUID,
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_inventory_movements_product_id ON inventory_movements(product_id);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_created_at ON inventory_movements(created_at);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_reference ON inventory_movements(reference_type, reference_id);
EOF
)

# Contador
SUCCESS_COUNT=0
FAILED_COUNT=0
TOTAL_COUNT=$(echo "$TENANTS" | wc -l)

echo -e "${YELLOW}Aplicando migración a $TOTAL_COUNT tenant(s)...${NC}"
echo ""

# Iterar sobre cada tenant
echo "$TENANTS" | while read tenant; do
    if [ -z "$tenant" ]; then
        continue
    fi
    
    echo -e "${BLUE}Procesando: $tenant${NC}"
    
    # Verificar si la migración ya está aplicada
    HAS_COLUMN=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
        SET search_path TO $tenant;
        SELECT EXISTS (
            SELECT 1 
            FROM information_schema.columns 
            WHERE table_schema = '$tenant' 
            AND table_name = 'products' 
            AND column_name = 'stock_quantity'
        );
    " | sed 's/^[ \t]*//;s/[ \t]*$//')
    
    if [ "$HAS_COLUMN" = "t" ]; then
        echo -e "  ${YELLOW}⊙ Ya tiene la columna stock_quantity - Saltando${NC}"
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
if [ $FAILED_COUNT -gt 0 ]; then
    echo -e "${RED}Fallidos: $FAILED_COUNT${NC}"
fi
echo ""

if [ $FAILED_COUNT -eq 0 ]; then
    echo -e "${GREEN}✓ Migración completada exitosamente en todos los tenants${NC}"
    echo ""
    echo -e "${YELLOW}Próximos pasos:${NC}"
    echo "1. Actualizar el stock inicial de los productos existentes"
    echo "2. Reiniciar el servidor backend para cargar los nuevos cambios"
    echo "3. Probar creando productos con stock y realizando ventas"
    echo ""
else
    echo -e "${YELLOW}⚠ Algunos tenants fallaron. Revisa los errores arriba.${NC}"
    exit 1
fi

