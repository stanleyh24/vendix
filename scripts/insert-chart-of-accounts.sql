-- Script para insertar el plan de cuentas en la tabla chart_of_accounts
-- Uso: psql -h localhost -U postgres -d vendix -v schema=tenant_demo -f scripts/insert-chart-of-accounts.sql

\echo 'Insertando plan de cuentas en el schema :schema'

SET search_path TO :schema, public;

-- Verificar si ya hay cuentas
DO $$
DECLARE
    cuenta_count INTEGER;
BEGIN
    SELECT count(*) INTO cuenta_count FROM chart_of_accounts;
    
    IF cuenta_count > 0 THEN
        RAISE NOTICE 'Ya existen % cuentas en el plan. Omitiendo inserción.', cuenta_count;
    ELSE
        RAISE NOTICE 'Insertando plan de cuentas...';
        
        -- Insertar plan de cuentas básico para República Dominicana
        INSERT INTO chart_of_accounts (account_code, account_name, account_type, description) VALUES
            -- ACTIVOS
            ('1000', 'ACTIVOS', 'ASSET', 'Activos totales'),
            ('1100', 'Activo Circulante', 'ASSET', 'Activos corrientes'),
            ('1110', 'Caja y Bancos', 'ASSET', 'Efectivo y equivalentes'),
            ('1111', 'Caja General', 'ASSET', 'Efectivo en caja'),
            ('1112', 'Banco - Cuenta Corriente', 'ASSET', 'Cuenta corriente bancaria'),
            ('1113', 'Banco - Cuenta de Ahorros', 'ASSET', 'Cuenta de ahorros'),
            ('1120', 'Cuentas por Cobrar', 'ASSET', 'Clientes y deudores'),
            ('1121', 'Clientes', 'ASSET', 'Cuentas por cobrar a clientes'),
            ('1122', 'Documentos por Cobrar', 'ASSET', 'Pagarés y documentos'),
            ('1130', 'Inventarios', 'ASSET', 'Mercancías y productos'),
            ('1131', 'Inventario de Mercancías', 'ASSET', 'Productos para la venta'),
            ('1200', 'Activo Fijo', 'ASSET', 'Activos no corrientes'),
            ('1210', 'Propiedad, Planta y Equipo', 'ASSET', 'Activos tangibles'),
            ('1211', 'Edificios', 'ASSET', 'Construcciones e inmuebles'),
            ('1212', 'Maquinaria y Equipo', 'ASSET', 'Equipos de producción'),
            ('1213', 'Equipo de Oficina', 'ASSET', 'Mobiliario y equipos de oficina'),
            ('1214', 'Vehículos', 'ASSET', 'Transporte y vehículos'),
            ('1215', 'Equipo de Computación', 'ASSET', 'Computadoras y tecnología'),

            -- PASIVOS
            ('2000', 'PASIVOS', 'LIABILITY', 'Pasivos totales'),
            ('2100', 'Pasivo Circulante', 'LIABILITY', 'Pasivos corrientes'),
            ('2110', 'Cuentas por Pagar', 'LIABILITY', 'Proveedores y acreedores'),
            ('2111', 'Proveedores', 'LIABILITY', 'Cuentas por pagar a proveedores'),
            ('2112', 'Documentos por Pagar', 'LIABILITY', 'Pagarés y documentos'),
            ('2120', 'Impuestos por Pagar', 'LIABILITY', 'Obligaciones fiscales'),
            ('2121', 'ITBIS por Pagar', 'LIABILITY', 'Impuesto sobre transferencias'),
            ('2122', 'ISR por Pagar', 'LIABILITY', 'Impuesto sobre la renta'),
            ('2130', 'Gastos Acumulados', 'LIABILITY', 'Gastos devengados no pagados'),
            ('2131', 'Sueldos por Pagar', 'LIABILITY', 'Nómina pendiente de pago'),
            ('2200', 'Pasivo a Largo Plazo', 'LIABILITY', 'Pasivos no corrientes'),
            ('2210', 'Préstamos Bancarios', 'LIABILITY', 'Deudas a largo plazo'),

            -- CAPITAL
            ('3000', 'CAPITAL', 'EQUITY', 'Patrimonio neto'),
            ('3100', 'Capital Social', 'EQUITY', 'Aportaciones de socios'),
            ('3200', 'Utilidades Retenidas', 'EQUITY', 'Ganancias acumuladas'),
            ('3300', 'Utilidad del Ejercicio', 'EQUITY', 'Resultado del período'),

            -- INGRESOS
            ('4000', 'INGRESOS', 'REVENUE', 'Ingresos totales'),
            ('4100', 'Ingresos por Ventas', 'REVENUE', 'Ventas de productos y servicios'),
            ('4110', 'Ventas', 'REVENUE', 'Ventas al contado y crédito'),
            ('4120', 'Servicios', 'REVENUE', 'Ingresos por servicios'),
            ('4200', 'Otros Ingresos', 'REVENUE', 'Ingresos no operacionales'),

            -- COSTOS
            ('5000', 'COSTOS', 'EXPENSE', 'Costo de ventas'),
            ('5100', 'Costo de Ventas', 'EXPENSE', 'Costo de mercancías vendidas'),

            -- GASTOS
            ('6000', 'GASTOS', 'EXPENSE', 'Gastos operacionales'),
            ('6100', 'Gastos de Administración', 'EXPENSE', 'Gastos administrativos'),
            ('6110', 'Sueldos y Salarios', 'EXPENSE', 'Nómina administrativa'),
            ('6120', 'Alquileres', 'EXPENSE', 'Arrendamientos'),
            ('6130', 'Servicios Públicos', 'EXPENSE', 'Electricidad, agua, etc.'),
            ('6131', 'Energía Eléctrica', 'EXPENSE', 'Consumo eléctrico'),
            ('6132', 'Agua', 'EXPENSE', 'Servicio de agua'),
            ('6133', 'Teléfono e Internet', 'EXPENSE', 'Comunicaciones'),
            ('6140', 'Papelería y Útiles', 'EXPENSE', 'Material de oficina'),
            ('6150', 'Combustible y Lubricantes', 'EXPENSE', 'Combustible para vehículos'),
            ('6200', 'Gastos de Ventas', 'EXPENSE', 'Gastos del área comercial'),
            ('6210', 'Publicidad y Mercadeo', 'EXPENSE', 'Marketing y publicidad'),
            ('6220', 'Comisiones de Ventas', 'EXPENSE', 'Comisiones a vendedores'),
            ('6300', 'Gastos Financieros', 'EXPENSE', 'Intereses y cargos bancarios'),
            ('6310', 'Intereses Bancarios', 'EXPENSE', 'Intereses sobre préstamos'),
            ('6320', 'Comisiones Bancarias', 'EXPENSE', 'Cargos por servicios bancarios');
        
        RAISE NOTICE '✅ Se insertaron % cuentas exitosamente', (SELECT count(*) FROM chart_of_accounts);
    END IF;
END $$;

-- Mostrar resumen
SELECT 
    account_type,
    count(*) as num_cuentas
FROM chart_of_accounts
GROUP BY account_type
ORDER BY account_type;

\echo ''
\echo 'Total de cuentas:'
SELECT count(*) as total_cuentas FROM chart_of_accounts;

