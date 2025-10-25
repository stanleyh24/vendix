#!/bin/bash

# Script para cargar datos de ejemplo en el tenant demo
# Ejecuta el archivo SQL dentro del contenedor de PostgreSQL

echo "🌱 Poblando datos de ejemplo en el tenant demo..."

# Copiar el archivo SQL al contenedor
docker cp /home/scorpion/desa/Projects/vendix/scripts/seed-demo-data.sql vendix-postgres-dev:/tmp/seed-demo-data.sql

# Ejecutar el script SQL
docker exec -i vendix-postgres-dev psql -U postgres -d vendix -f /tmp/seed-demo-data.sql

# Limpiar archivo temporal
docker exec vendix-postgres-dev rm /tmp/seed-demo-data.sql

echo ""
echo "✅ Datos de ejemplo cargados exitosamente!"
echo ""
echo "📊 Resumen de datos generados:"
echo "   - 8 Clientes (individuales y empresas)"
echo "   - 13 Productos (electrónica, papelería, servicios, alimentos)"
echo "   - 7 Facturas con diferentes estados:"
echo "     • 2 Pagadas"
echo "     • 3 Pendientes"
echo "     • 1 Borrador"
echo "     • 1 Anulada"
echo "   - 10 Gastos de diferentes categorías"
echo ""
echo "🌐 Puedes verificar en:"
echo "   - http://localhost:5173/customers (Clientes)"
echo "   - http://localhost:5173/products (Productos)"
echo "   - http://localhost:5173/invoices (Facturas)"
echo "   - http://localhost:5173/expenses (Gastos)"
echo "   - http://localhost:5173/sales (Punto de Venta)"
echo ""

