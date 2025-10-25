-- Script para limpiar facturas duplicadas en el tenant demo
-- USO: psql vendix -U vendix -f scripts/clean-invoices.sql

-- Cambiar al schema demo
SET search_path TO demo;

-- Mostrar facturas actuales
\echo '============================================'
\echo 'FACTURAS ACTUALES EN EL SISTEMA'
\echo '============================================'
SELECT 
    id,
    invoice_number,
    ncf_type,
    total,
    status,
    created_at
FROM invoices 
ORDER BY created_at DESC;

\echo ''
\echo '============================================'
\echo 'VERIFICANDO DUPLICADOS'
\echo '============================================'
SELECT 
    invoice_number, 
    COUNT(*) as count,
    STRING_AGG(id::text, ', ') as invoice_ids
FROM invoices 
GROUP BY invoice_number 
HAVING COUNT(*) > 1;

-- Descomentar las siguientes líneas para ELIMINAR TODAS LAS FACTURAS
-- ADVERTENCIA: Esto eliminará TODAS las facturas en el tenant demo
-- Solo usar en ambiente de desarrollo

-- \echo ''
-- \echo '============================================'
-- \echo 'ELIMINANDO TODAS LAS FACTURAS...'
-- \echo '============================================'
-- DELETE FROM invoice_lines;
-- DELETE FROM invoices;
-- \echo 'Facturas eliminadas correctamente.'

\echo ''
\echo '============================================'
\echo 'Para eliminar TODAS las facturas ejecuta:'
\echo '  psql vendix -U vendix -c "SET search_path TO demo; DELETE FROM invoice_lines; DELETE FROM invoices;"'
\echo '============================================'

