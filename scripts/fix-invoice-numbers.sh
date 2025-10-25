#!/bin/bash

# Script para verificar y corregir números de facturas duplicados

# Colores
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}======================================${NC}"
echo -e "${GREEN}  Verificación de Números de Facturas${NC}"
echo -e "${GREEN}======================================${NC}"
echo ""

# Verificar que existan las variables de entorno
if [ -z "$DATABASE_URL" ]; then
    echo -e "${YELLOW}DATABASE_URL no está configurada. Usando valores por defecto...${NC}"
    DB_HOST="localhost"
    DB_PORT="5432"
    DB_USER="postgres"
    DB_PASS="postgres"
    DB_NAME="vendix"
else
    echo -e "${GREEN}Usando DATABASE_URL de las variables de entorno${NC}"
fi

# Función para ejecutar queries
run_query() {
    local schema=$1
    local query=$2
    
    if [ -n "$DATABASE_URL" ]; then
        psql "$DATABASE_URL" -c "SET search_path TO $schema; $query"
    else
        PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SET search_path TO $schema; $query"
    fi
}

echo -e "${YELLOW}1. Verificando facturas en schema 'demo'...${NC}"
echo ""

# Listar facturas existentes
echo "Facturas existentes:"
run_query "demo" "SELECT id, invoice_number, ncf_type, customer_id, total, created_at FROM invoices ORDER BY created_at DESC LIMIT 10;"

echo ""
echo -e "${YELLOW}2. Verificando duplicados...${NC}"
echo ""

# Verificar si hay duplicados
run_query "demo" "SELECT invoice_number, COUNT(*) as count FROM invoices GROUP BY invoice_number HAVING COUNT(*) > 1;"

echo ""
echo -e "${YELLOW}3. Solución: Si hay duplicados, puedes eliminar las facturas de prueba${NC}"
echo -e "${YELLOW}   Ejecuta el siguiente comando para eliminar TODAS las facturas (CUIDADO):${NC}"
echo ""
echo -e "${RED}   SET search_path TO demo;${NC}"
echo -e "${RED}   DELETE FROM invoice_lines;${NC}"
echo -e "${RED}   DELETE FROM invoices;${NC}"
echo ""

echo -e "${GREEN}======================================${NC}"
echo -e "${GREEN}  Verificación Completada${NC}"
echo -e "${GREEN}======================================${NC}"

