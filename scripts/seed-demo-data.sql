-- Script para poblar datos de ejemplo en el tenant demo
-- Ejecutar: docker exec -i vendix-postgres-dev psql -U postgres -d vendix -f /tmp/seed-demo-data.sql

-- Cambiar al schema del tenant demo
SET search_path TO tenant_demo;

-- Limpiar datos existentes (en orden debido a foreign keys)
DELETE FROM invoice_lines;
DELETE FROM invoices;
DELETE FROM payment_allocations;
DELETE FROM payments;
DELETE FROM products;
DELETE FROM customers;

-- ====================
-- CLIENTES
-- ====================

INSERT INTO customers (id, name, email, phone, tax_id, address, customer_type, is_active, created_at, updated_at)
VALUES
  (gen_random_uuid(), 'Juan Pérez', 'juan.perez@email.com', '809-555-0101', '001-1234567-8', 'Av. Winston Churchill #45, Santo Domingo', 'individual', true, NOW(), NOW()),
  (gen_random_uuid(), 'María García', 'maria.garcia@email.com', '809-555-0102', '001-9876543-2', 'Calle El Conde #123, Zona Colonial', 'individual', true, NOW(), NOW()),
  (gen_random_uuid(), 'Supermercado La Económica SRL', 'contacto@laeconomica.do', '809-555-0201', '130-56789-0', 'Av. 27 de Febrero #890, Santo Domingo', 'business', true, NOW(), NOW()),
  (gen_random_uuid(), 'Ferretería El Constructor SA', 'ventas@elconstructor.do', '809-555-0202', '130-11223-3', 'Autopista Duarte Km 8, Santo Domingo Norte', 'business', true, NOW(), NOW()),
  (gen_random_uuid(), 'Restaurante Don Pepe', 'admin@donpepe.do', '809-555-0301', '131-99887-7', 'Calle Hostos #56, Zona Colonial', 'business', true, NOW(), NOW()),
  (gen_random_uuid(), 'Carlos Rodríguez', 'carlos.rod@email.com', '809-555-0103', '001-5555666-7', 'Los Prados, Santo Domingo', 'individual', true, NOW(), NOW()),
  (gen_random_uuid(), 'Farmacia Moderna', 'info@farmaciamoderna.do', '809-555-0401', '130-77665-5', 'Av. Abraham Lincoln #234', 'business', true, NOW(), NOW()),
  (gen_random_uuid(), 'Ana Martínez', 'ana.martinez@email.com', '809-555-0104', '001-3332221-1', 'Naco, Santo Domingo', 'individual', true, NOW(), NOW()),
  (gen_random_uuid(), 'Gastos Generales', 'gastos@interno.local', NULL, '999-9999999-9', 'Interno', 'business', true, NOW(), NOW());

-- ====================
-- PRODUCTOS
-- ====================

INSERT INTO products (id, code, name, description, product_type, unit, price, cost, tax_rate, is_active, created_at, updated_at)
VALUES
  -- Electrónica
  (gen_random_uuid(), 'LAPTOP-001', 'Laptop Dell XPS 15', 'Laptop Dell XPS 15, Intel i7, 16GB RAM, 512GB SSD', 'product', 'unit', 75000.00, 60000.00, 0.18, true, NOW(), NOW()),
  (gen_random_uuid(), 'LAPTOP-002', 'Laptop HP Pavilion', 'Laptop HP Pavilion, Intel i5, 8GB RAM, 256GB SSD', 'product', 'unit', 45000.00, 35000.00, 0.18, true, NOW(), NOW()),
  (gen_random_uuid(), 'MONITOR-001', 'Monitor LG 27"', 'Monitor LG 27 pulgadas, Full HD, IPS', 'product', 'unit', 18000.00, 14000.00, 0.18, true, NOW(), NOW()),
  (gen_random_uuid(), 'MOUSE-001', 'Mouse Logitech MX Master', 'Mouse inalámbrico ergonómico', 'product', 'unit', 4500.00, 3000.00, 0.18, true, NOW(), NOW()),
  (gen_random_uuid(), 'KEYBOARD-001', 'Teclado Mecánico Corsair', 'Teclado mecánico RGB, switches Cherry MX', 'product', 'unit', 8500.00, 6000.00, 0.18, true, NOW(), NOW()),
  
  -- Papelería
  (gen_random_uuid(), 'PAPER-001', 'Resma de Papel A4', 'Resma 500 hojas, 75g/m²', 'product', 'unit', 350.00, 250.00, 0.18, true, NOW(), NOW()),
  (gen_random_uuid(), 'PEN-001', 'Bolígrafos BIC (Caja 50)', 'Caja de 50 bolígrafos azules', 'product', 'unit', 450.00, 300.00, 0.18, true, NOW(), NOW()),
  (gen_random_uuid(), 'FOLDER-001', 'Folder Manila (Paquete 25)', 'Paquete de 25 folders tamaño carta', 'product', 'unit', 280.00, 180.00, 0.18, true, NOW(), NOW()),
  
  -- Servicios
  (gen_random_uuid(), 'SERV-001', 'Consultoría IT (Hora)', 'Servicio de consultoría en tecnología', 'service', 'hour', 3500.00, 2000.00, 0.18, true, NOW(), NOW()),
  (gen_random_uuid(), 'SERV-002', 'Desarrollo Web', 'Desarrollo de sitio web corporativo', 'service', 'project', 85000.00, 50000.00, 0.18, true, NOW(), NOW()),
  (gen_random_uuid(), 'SERV-003', 'Soporte Técnico Mensual', 'Plan de soporte técnico mensual', 'service', 'month', 15000.00, 8000.00, 0.18, true, NOW(), NOW()),
  
  -- Alimentos
  (gen_random_uuid(), 'FOOD-001', 'Café Premium (Libra)', 'Café dominicano de montaña', 'product', 'lb', 450.00, 280.00, 0.18, true, NOW(), NOW()),
  (gen_random_uuid(), 'FOOD-002', 'Galletas Surtidas', 'Paquete de galletas variadas', 'product', 'unit', 185.00, 120.00, 0.18, true, NOW(), NOW());

-- ====================
-- FACTURAS CON LÍNEAS
-- ====================

DO $$
DECLARE
  cliente1_id UUID;
  cliente2_id UUID;
  cliente3_id UUID;
  cliente4_id UUID;
  cliente5_id UUID;
  gastos_id UUID;
  
  laptop_dell_id UUID;
  monitor_lg_id UUID;
  mouse_id UUID;
  consulta_id UUID;
  desarrollo_id UUID;
  soporte_id UUID;
  teclado_id UUID;
  resma_id UUID;
  
  factura1_id UUID;
  factura2_id UUID;
  factura3_id UUID;
  factura4_id UUID;
  factura5_id UUID;
  factura6_id UUID;
  factura7_id UUID;
  
BEGIN
  -- Obtener IDs de clientes
  SELECT id INTO cliente1_id FROM customers WHERE tax_id = '001-1234567-8' LIMIT 1;
  SELECT id INTO cliente2_id FROM customers WHERE tax_id = '130-56789-0' LIMIT 1;
  SELECT id INTO cliente3_id FROM customers WHERE tax_id = '131-99887-7' LIMIT 1;
  SELECT id INTO cliente4_id FROM customers WHERE tax_id = '001-9876543-2' LIMIT 1;
  SELECT id INTO cliente5_id FROM customers WHERE tax_id = '130-11223-3' LIMIT 1;
  SELECT id INTO gastos_id FROM customers WHERE tax_id = '999-9999999-9' LIMIT 1;
  
  -- Obtener IDs de productos
  SELECT id INTO laptop_dell_id FROM products WHERE code = 'LAPTOP-001' LIMIT 1;
  SELECT id INTO monitor_lg_id FROM products WHERE code = 'MONITOR-001' LIMIT 1;
  SELECT id INTO mouse_id FROM products WHERE code = 'MOUSE-001' LIMIT 1;
  SELECT id INTO consulta_id FROM products WHERE code = 'SERV-001' LIMIT 1;
  SELECT id INTO desarrollo_id FROM products WHERE code = 'SERV-002' LIMIT 1;
  SELECT id INTO soporte_id FROM products WHERE code = 'SERV-003' LIMIT 1;
  SELECT id INTO teclado_id FROM products WHERE code = 'KEYBOARD-001' LIMIT 1;
  SELECT id INTO resma_id FROM products WHERE code = 'PAPER-001' LIMIT 1;
  
  -- FACTURA 1: Laptop + Mouse (PAGADA)
  factura1_id := gen_random_uuid();
  INSERT INTO invoices (id, invoice_number, customer_id, issue_date, due_date, status, subtotal, tax_amount, total, paid_amount, currency, notes, created_at, updated_at)
  VALUES (
    factura1_id, 'INV-00000001', cliente1_id, 
    CURRENT_DATE - INTERVAL '30 days', CURRENT_DATE - INTERVAL '0 days',
    'paid', 154500.00, 27810.00, 182310.00, 182310.00, 'DOP',
    'Equipo de trabajo para oficina', NOW() - INTERVAL '30 days', NOW()
  );
  
  INSERT INTO invoice_lines (id, invoice_id, product_id, line_number, description, quantity, unit_price, tax_rate, tax_amount, line_total, created_at)
  VALUES
    (gen_random_uuid(), factura1_id, laptop_dell_id, 1, 'Laptop Dell XPS 15', 2, 75000.00, 0.18, 27000.00, 177000.00, NOW() - INTERVAL '30 days'),
    (gen_random_uuid(), factura1_id, mouse_id, 2, 'Mouse Logitech MX Master', 1, 4500.00, 0.18, 810.00, 5310.00, NOW() - INTERVAL '30 days');
  
  -- FACTURA 2: Servicios de consultoría (ENVIADA)
  factura2_id := gen_random_uuid();
  INSERT INTO invoices (id, invoice_number, customer_id, issue_date, due_date, status, subtotal, tax_amount, total, paid_amount, currency, notes, sent_at, created_at, updated_at)
  VALUES (
    factura2_id, 'INV-00000002', cliente2_id,
    CURRENT_DATE - INTERVAL '15 days', CURRENT_DATE + INTERVAL '15 days',
    'sent', 35000.00, 6300.00, 41300.00, 0.00, 'DOP',
    'Servicios de consultoría IT - Octubre 2025',
    NOW() - INTERVAL '14 days',
    NOW() - INTERVAL '15 days', NOW() - INTERVAL '14 days'
  );
  
  INSERT INTO invoice_lines (id, invoice_id, product_id, line_number, description, quantity, unit_price, tax_rate, tax_amount, line_total, created_at)
  VALUES
    (gen_random_uuid(), factura2_id, consulta_id, 1, 'Consultoría IT', 10, 3500.00, 0.18, 6300.00, 41300.00, NOW() - INTERVAL '15 days');
  
  -- FACTURA 3: Desarrollo web (PENDIENTE)
  factura3_id := gen_random_uuid();
  INSERT INTO invoices (id, invoice_number, customer_id, issue_date, due_date, status, subtotal, tax_amount, total, paid_amount, currency, notes, created_at, updated_at)
  VALUES (
    factura3_id, 'INV-00000003', cliente3_id,
    CURRENT_DATE - INTERVAL '10 days', CURRENT_DATE + INTERVAL '20 days',
    'pending', 85000.00, 15300.00, 100300.00, 0.00, 'DOP',
    'Desarrollo de sitio web corporativo', NOW() - INTERVAL '10 days', NOW() - INTERVAL '10 days'
  );
  
  INSERT INTO invoice_lines (id, invoice_id, product_id, line_number, description, quantity, unit_price, tax_rate, tax_amount, line_total, created_at)
  VALUES
    (gen_random_uuid(), factura3_id, desarrollo_id, 1, 'Desarrollo Web Corporativo', 1, 85000.00, 0.18, 15300.00, 100300.00, NOW() - INTERVAL '10 days');
  
  -- FACTURA 4: Monitor (PAGADA)
  factura4_id := gen_random_uuid();
  INSERT INTO invoices (id, invoice_number, customer_id, issue_date, due_date, status, subtotal, tax_amount, total, paid_amount, currency, created_at, updated_at)
  VALUES (
    factura4_id, 'INV-00000004', cliente4_id,
    CURRENT_DATE - INTERVAL '5 days', CURRENT_DATE + INTERVAL '25 days',
    'paid', 18000.00, 3240.00, 21240.00, 21240.00, 'DOP',
    NOW() - INTERVAL '5 days', NOW() - INTERVAL '3 days'
  );
  
  INSERT INTO invoice_lines (id, invoice_id, product_id, line_number, description, quantity, unit_price, tax_rate, tax_amount, line_total, created_at)
  VALUES
    (gen_random_uuid(), factura4_id, monitor_lg_id, 1, 'Monitor LG 27" Full HD', 1, 18000.00, 0.18, 3240.00, 21240.00, NOW() - INTERVAL '5 days');
  
  -- FACTURA 5: Equipos de oficina (PENDIENTE)
  factura5_id := gen_random_uuid();
  INSERT INTO invoices (id, invoice_number, customer_id, issue_date, due_date, status, subtotal, tax_amount, total, paid_amount, currency, created_at, updated_at)
  VALUES (
    factura5_id, 'INV-00000005', cliente5_id,
    CURRENT_DATE - INTERVAL '3 days', CURRENT_DATE + INTERVAL '27 days',
    'pending', 56000.00, 10080.00, 66080.00, 0.00, 'DOP',
    NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days'
  );
  
  INSERT INTO invoice_lines (id, invoice_id, product_id, line_number, description, quantity, unit_price, tax_rate, tax_amount, line_total, created_at)
  VALUES
    (gen_random_uuid(), factura5_id, teclado_id, 1, 'Teclado Mecánico Corsair RGB', 5, 8500.00, 0.18, 7650.00, 50150.00, NOW() - INTERVAL '3 days'),
    (gen_random_uuid(), factura5_id, mouse_id, 2, 'Mouse Logitech MX Master', 1, 4500.00, 0.18, 810.00, 5310.00, NOW() - INTERVAL '3 days');
  
  -- FACTURA 6: Soporte técnico (BORRADOR)
  factura6_id := gen_random_uuid();
  INSERT INTO invoices (id, invoice_number, customer_id, issue_date, due_date, status, subtotal, tax_amount, total, paid_amount, currency, created_at, updated_at)
  VALUES (
    factura6_id, 'INV-00000006', cliente2_id,
    CURRENT_DATE, CURRENT_DATE + INTERVAL '30 days',
    'draft', 15000.00, 2700.00, 17700.00, 0.00, 'DOP',
    NOW(), NOW()
  );
  
  INSERT INTO invoice_lines (id, invoice_id, product_id, line_number, description, quantity, unit_price, tax_rate, tax_amount, line_total, created_at)
  VALUES
    (gen_random_uuid(), factura6_id, soporte_id, 1, 'Plan de Soporte Técnico Mensual', 1, 15000.00, 0.18, 2700.00, 17700.00, NOW());
  
  -- FACTURA 7: Laptop (ANULADA)
  factura7_id := gen_random_uuid();
  INSERT INTO invoices (id, invoice_number, customer_id, issue_date, due_date, status, subtotal, tax_amount, total, paid_amount, currency, notes, created_at, updated_at)
  VALUES (
    factura7_id, 'INV-00000007', cliente1_id,
    CURRENT_DATE - INTERVAL '20 days', CURRENT_DATE - INTERVAL '10 days',
    'cancelled', 75000.00, 13500.00, 88500.00, 0.00, 'DOP',
    'Factura anulada por error en especificaciones',
    NOW() - INTERVAL '20 days', NOW() - INTERVAL '18 days'
  );
  
  INSERT INTO invoice_lines (id, invoice_id, product_id, line_number, description, quantity, unit_price, tax_rate, tax_amount, line_total, created_at)
  VALUES
    (gen_random_uuid(), factura7_id, laptop_dell_id, 1, 'Laptop Dell XPS 15', 1, 75000.00, 0.18, 13500.00, 88500.00, NOW() - INTERVAL '20 days');

  -- ====================
  -- GASTOS (PAYMENTS)
  -- ====================

  INSERT INTO payments (id, payment_number, customer_id, payment_method, payment_date, amount, currency, reference, notes, status, created_at, updated_at)
  VALUES
    -- Gastos de octubre (usando cliente "Gastos Generales")
    (gen_random_uuid(), 'EXP-00000001', gastos_id, 'cash', CURRENT_DATE - INTERVAL '25 days', 8900.00, 'DOP', 'EDESUR-OCT2025', 'Factura de electricidad - Octubre', 'completed', NOW() - INTERVAL '25 days', NOW() - INTERVAL '25 days'),
    (gen_random_uuid(), 'EXP-00000002', gastos_id, 'transfer', CURRENT_DATE - INTERVAL '20 days', 5500.00, 'DOP', 'INV-PAP-001', 'Papelería y útiles de oficina', 'completed', NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days'),
    (gen_random_uuid(), 'EXP-00000003', gastos_id, 'card', CURRENT_DATE - INTERVAL '18 days', 15000.00, 'DOP', 'META-ADS-789', 'Campaña en redes sociales - Facebook e Instagram', 'completed', NOW() - INTERVAL '18 days', NOW() - INTERVAL '18 days'),
    (gen_random_uuid(), 'EXP-00000004', gastos_id, 'transfer', CURRENT_DATE - INTERVAL '15 days', 35000.00, 'DOP', 'ALQ-OCT-2025', 'Alquiler de oficina - Octubre 2025', 'completed', NOW() - INTERVAL '15 days', NOW() - INTERVAL '15 days'),
    (gen_random_uuid(), 'EXP-00000005', gastos_id, 'cash', CURRENT_DATE - INTERVAL '12 days', 2800.00, 'DOP', 'GASOLINA-001', 'Combustible para vehículo de la empresa', 'completed', NOW() - INTERVAL '12 days', NOW() - INTERVAL '12 days'),
    (gen_random_uuid(), 'EXP-00000006', gastos_id, 'transfer', CURRENT_DATE - INTERVAL '10 days', 4500.00, 'DOP', 'AGUA-OCT', 'Servicio de agua potable', 'completed', NOW() - INTERVAL '10 days', NOW() - INTERVAL '10 days'),
    (gen_random_uuid(), 'EXP-00000007', gastos_id, 'card', CURRENT_DATE - INTERVAL '8 days', 12000.00, 'DOP', 'INTERNET-FIB', 'Internet fibra óptica empresarial', 'completed', NOW() - INTERVAL '8 days', NOW() - INTERVAL '8 days'),
    (gen_random_uuid(), 'EXP-00000008', gastos_id, 'check', CURRENT_DATE - INTERVAL '5 days', 8500.00, 'DOP', 'LIMPIEZA-OCT', 'Servicio de limpieza mensual', 'completed', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days'),
    (gen_random_uuid(), 'EXP-00000009', gastos_id, 'transfer', CURRENT_DATE - INTERVAL '3 days', 3200.00, 'DOP', 'TELEFONIA', 'Servicios telefónicos empresariales', 'completed', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days'),
    (gen_random_uuid(), 'EXP-00000010', gastos_id, 'cash', CURRENT_DATE - INTERVAL '1 day', 1800.00, 'DOP', 'CAFE-SUP', 'Suministros de café y snacks para oficina', 'completed', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day');

END $$;

-- Verificar datos insertados
SELECT 'Clientes insertados:' as descripcion, COUNT(*) as cantidad FROM customers
UNION ALL
SELECT 'Productos insertados:', COUNT(*) FROM products
UNION ALL
SELECT 'Facturas insertadas:', COUNT(*) FROM invoices
UNION ALL
SELECT 'Líneas de factura:', COUNT(*) FROM invoice_lines
UNION ALL
SELECT 'Gastos insertados:', COUNT(*) FROM payments;

-- Mostrar resumen de facturas por estado
SELECT 
  status,
  COUNT(*) as cantidad,
  SUM(total) as total
FROM invoices
GROUP BY status
ORDER BY status;
