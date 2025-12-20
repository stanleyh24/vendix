-- Script para poblar datos de ejemplo en el tenant demo
-- Ejecutar: docker exec -i vendix-postgres-dev psql -U postgres -d vendix -f /tmp/seed-demo-data.sql

-- Cambiar al schema del tenant demo
SET search_path TO tenant_demo;

-- Limpiar datos existentes (en orden debido a foreign keys)
DELETE FROM payroll_entries;
DELETE FROM payroll_periods;
DELETE FROM employees;
DELETE FROM invoice_lines;
DELETE FROM invoices;
DELETE FROM payment_allocations;
DELETE FROM payments;
DELETE FROM products;
DELETE FROM customers;
DELETE FROM suppliers;

-- ====================
-- EMPLEADOS
-- ====================

INSERT INTO employees (
  id, employee_code, first_name, last_name, email, phone, tax_id,
  address, city, state, postal_code, country, department, position,
  hire_date, salary, salary_type, is_active, metadata, created_at, updated_at
)
VALUES
  (gen_random_uuid(), 'EMP-001', 'Carolina', 'Martínez', 'carolina.martinez@vendixdemo.do', '809-555-3001', '001-8899001-2', 'Av. Gustavo M. Ricart #45', 'Santo Domingo', 'DN', '10129', 'DO', 'Operaciones', 'Gerente de Operaciones', CURRENT_DATE - INTERVAL '720 days', 68000.00, 'monthly', true, jsonb_build_object('contract_type', 'indefinido', 'bank_account', '001-456789-1'), NOW() - INTERVAL '720 days', NOW()),
  (gen_random_uuid(), 'EMP-002', 'Luis', 'Fernández', 'luis.fernandez@vendixdemo.do', '809-555-3002', '001-2233445-6', 'Calle José Amado Soler #12', 'Santo Domingo', 'DN', '10127', 'DO', 'Ventas', 'Ejecutivo Comercial', CURRENT_DATE - INTERVAL '540 days', 45000.00, 'monthly', true, jsonb_build_object('commission_rate', 0.03, 'bank_account', '001-112233-4'), NOW() - INTERVAL '540 days', NOW()),
  (gen_random_uuid(), 'EMP-003', 'Mariela', 'Suárez', 'mariela.suarez@vendixdemo.do', '809-555-3003', '001-4455667-8', 'Residencial Don Honorio #8', 'Santo Domingo Oeste', 'SDO', '10903', 'DO', 'Soporte', 'Especialista de Soporte', CURRENT_DATE - INTERVAL '365 days', 320.00, 'hourly', true, jsonb_build_object('shift', 'nocturno', 'bank_account', '001-665544-7'), NOW() - INTERVAL '365 days', NOW()),
  (gen_random_uuid(), 'EMP-004', 'Raúl', 'Peña', 'raul.pena@vendixdemo.do', '829-555-3004', '002-5566778-9', 'Residencial Praderas del Este', 'Santo Domingo Este', 'SDE', '11517', 'DO', 'Logística', 'Coordinador de Despacho', CURRENT_DATE - INTERVAL '410 days', 26800.00, 'daily', true, jsonb_build_object('transport_allowance', true, 'bank_account', '001-778899-0'), NOW() - INTERVAL '410 days', NOW());

-- ====================
-- NOMINA (PERIODOS Y ENTRADAS)
-- ====================

DO $$
DECLARE
  period_jan UUID := gen_random_uuid();
  period_feb UUID := gen_random_uuid();
  emp_ops employees%ROWTYPE;
  emp_sales employees%ROWTYPE;
  emp_support employees%ROWTYPE;
  emp_log employees%ROWTYPE;
  current_emp employees%ROWTYPE;
  hours_worked NUMERIC;
  days_worked NUMERIC;
  overtime_hours NUMERIC;
  bonuses NUMERIC;
  commissions NUMERIC;
  other_deductions NUMERIC;
  base_salary NUMERIC;
  entry_gross_salary NUMERIC;
  entry_overtime_pay NUMERIC;
  entry_total_gross NUMERIC;
  entry_tax NUMERIC;
  entry_social NUMERIC;
  entry_total_deductions NUMERIC;
  entry_net NUMERIC;
  entry_id UUID;
BEGIN
  SELECT * INTO emp_ops FROM employees WHERE employee_code = 'EMP-001' LIMIT 1;
  SELECT * INTO emp_sales FROM employees WHERE employee_code = 'EMP-002' LIMIT 1;
  SELECT * INTO emp_support FROM employees WHERE employee_code = 'EMP-003' LIMIT 1;
  SELECT * INTO emp_log FROM employees WHERE employee_code = 'EMP-004' LIMIT 1;

  INSERT INTO payroll_periods (id, period_code, period_start, period_end, status, total_gross, total_deductions, total_net, notes, created_by, created_at, updated_at)
  VALUES
    (period_jan, '2025-01', DATE '2025-01-01', DATE '2025-01-31', 'completed', 0, 0, 0, 'Nómina Enero 2025', NULL, NOW() - INTERVAL '60 days', NOW() - INTERVAL '60 days'),
    (period_feb, '2025-02', DATE '2025-02-01', DATE '2025-02-28', 'processing', 0, 0, 0, 'Nómina Febrero 2025 (en progreso)', NULL, NOW() - INTERVAL '30 days', NOW() - INTERVAL '5 days');

  -- Período enero 2025 - Carolina Martínez
  current_emp := emp_ops;
  hours_worked := NULL;
  days_worked := NULL;
  overtime_hours := 10;
  bonuses := 7000;
  commissions := 1500;
  other_deductions := 2800;
  base_salary := COALESCE(current_emp.salary, 0);
  IF current_emp.salary_type = 'monthly' THEN
    entry_gross_salary := base_salary;
    entry_overtime_pay := (base_salary / 160.0) * 1.5 * COALESCE(overtime_hours, 0);
  ELSIF current_emp.salary_type = 'hourly' THEN
    entry_gross_salary := base_salary * COALESCE(hours_worked, 0);
    entry_overtime_pay := base_salary * 1.5 * COALESCE(overtime_hours, 0);
  ELSIF current_emp.salary_type = 'daily' THEN
    entry_gross_salary := (base_salary / 22.0) * COALESCE(days_worked, 0);
    entry_overtime_pay := ((base_salary / 22.0) / 8.0) * 1.5 * COALESCE(overtime_hours, 0);
  ELSE
    RAISE EXCEPTION 'Tipo de salario no soportado: %', current_emp.salary_type;
  END IF;
  entry_total_gross := entry_gross_salary + entry_overtime_pay + bonuses + commissions;
  IF entry_total_gross <= 416220 THEN
    entry_tax := entry_total_gross * 0.15;
  ELSE
    entry_tax := 416220 * 0.15 + (entry_total_gross - 416220) * 0.20;
  END IF;
  entry_social := entry_total_gross * 0.0304;
  entry_total_deductions := entry_tax + entry_social + other_deductions;
  entry_net := entry_total_gross - entry_total_deductions;

  entry_gross_salary := ROUND(entry_gross_salary, 2);
  entry_overtime_pay := ROUND(entry_overtime_pay, 2);
  entry_total_gross := ROUND(entry_total_gross, 2);
  entry_tax := ROUND(entry_tax, 2);
  entry_social := ROUND(entry_social, 2);
  entry_total_deductions := ROUND(entry_total_deductions, 2);
  entry_net := ROUND(entry_net, 2);

  entry_id := gen_random_uuid();
  INSERT INTO payroll_entries (
    id, payroll_period_id, employee_id, employee_code, employee_name,
    department, position, salary, salary_type, hours_worked, days_worked,
    gross_salary, overtime_hours, overtime_pay, bonuses, commissions,
    total_gross, tax_deduction, social_security, other_deductions,
    total_deductions, net_salary, notes, created_at, updated_at
  )
  VALUES (
    entry_id, period_jan, current_emp.id, current_emp.employee_code,
    current_emp.first_name || ' ' || current_emp.last_name, current_emp.department, current_emp.position,
    entry_gross_salary, current_emp.salary_type, hours_worked, days_worked,
    entry_gross_salary, overtime_hours, entry_overtime_pay, bonuses, commissions,
    entry_total_gross, entry_tax, entry_social, other_deductions,
    entry_total_deductions, entry_net, 'Incluye bono por metas y horas extras de soporte a clientes',
    NOW() - INTERVAL '45 days', NOW() - INTERVAL '45 days'
  );

  UPDATE payroll_periods
  SET total_gross = ROUND(total_gross + entry_total_gross, 2),
      total_deductions = ROUND(total_deductions + entry_total_deductions, 2),
      total_net = ROUND(total_net + entry_net, 2),
      updated_at = NOW()
  WHERE id = period_jan;

  -- Período enero 2025 - Luis Fernández
  current_emp := emp_sales;
  hours_worked := NULL;
  days_worked := NULL;
  overtime_hours := 0;
  bonuses := 3500;
  commissions := 12000;
  other_deductions := 1800;
  base_salary := COALESCE(current_emp.salary, 0);
  IF current_emp.salary_type = 'monthly' THEN
    entry_gross_salary := base_salary;
    entry_overtime_pay := (base_salary / 160.0) * 1.5 * COALESCE(overtime_hours, 0);
  ELSIF current_emp.salary_type = 'hourly' THEN
    entry_gross_salary := base_salary * COALESCE(hours_worked, 0);
    entry_overtime_pay := base_salary * 1.5 * COALESCE(overtime_hours, 0);
  ELSIF current_emp.salary_type = 'daily' THEN
    entry_gross_salary := (base_salary / 22.0) * COALESCE(days_worked, 0);
    entry_overtime_pay := ((base_salary / 22.0) / 8.0) * 1.5 * COALESCE(overtime_hours, 0);
  ELSE
    RAISE EXCEPTION 'Tipo de salario no soportado: %', current_emp.salary_type;
  END IF;
  entry_total_gross := entry_gross_salary + entry_overtime_pay + bonuses + commissions;
  IF entry_total_gross <= 416220 THEN
    entry_tax := entry_total_gross * 0.15;
  ELSE
    entry_tax := 416220 * 0.15 + (entry_total_gross - 416220) * 0.20;
  END IF;
  entry_social := entry_total_gross * 0.0304;
  entry_total_deductions := entry_tax + entry_social + other_deductions;
  entry_net := entry_total_gross - entry_total_deductions;

  entry_gross_salary := ROUND(entry_gross_salary, 2);
  entry_overtime_pay := ROUND(entry_overtime_pay, 2);
  entry_total_gross := ROUND(entry_total_gross, 2);
  entry_tax := ROUND(entry_tax, 2);
  entry_social := ROUND(entry_social, 2);
  entry_total_deductions := ROUND(entry_total_deductions, 2);
  entry_net := ROUND(entry_net, 2);

  entry_id := gen_random_uuid();
  INSERT INTO payroll_entries (
    id, payroll_period_id, employee_id, employee_code, employee_name,
    department, position, salary, salary_type, hours_worked, days_worked,
    gross_salary, overtime_hours, overtime_pay, bonuses, commissions,
    total_gross, tax_deduction, social_security, other_deductions,
    total_deductions, net_salary, notes, created_at, updated_at
  )
  VALUES (
    entry_id, period_jan, current_emp.id, current_emp.employee_code,
    current_emp.first_name || ' ' || current_emp.last_name, current_emp.department, current_emp.position,
    entry_gross_salary, current_emp.salary_type, hours_worked, days_worked,
    entry_gross_salary, overtime_hours, entry_overtime_pay, bonuses, commissions,
    entry_total_gross, entry_tax, entry_social, other_deductions,
    entry_total_deductions, entry_net, 'Comisiones por nuevas cuentas corporativas',
    NOW() - INTERVAL '45 days', NOW() - INTERVAL '45 days'
  );

  UPDATE payroll_periods
  SET total_gross = ROUND(total_gross + entry_total_gross, 2),
      total_deductions = ROUND(total_deductions + entry_total_deductions, 2),
      total_net = ROUND(total_net + entry_net, 2),
      updated_at = NOW()
  WHERE id = period_jan;

  -- Período enero 2025 - Mariela Suárez
  current_emp := emp_support;
  hours_worked := 168;
  days_worked := NULL;
  overtime_hours := 6;
  bonuses := 2500;
  commissions := 0;
  other_deductions := 900;
  base_salary := COALESCE(current_emp.salary, 0);
  IF current_emp.salary_type = 'monthly' THEN
    entry_gross_salary := base_salary;
    entry_overtime_pay := (base_salary / 160.0) * 1.5 * COALESCE(overtime_hours, 0);
  ELSIF current_emp.salary_type = 'hourly' THEN
    entry_gross_salary := base_salary * COALESCE(hours_worked, 0);
    entry_overtime_pay := base_salary * 1.5 * COALESCE(overtime_hours, 0);
  ELSIF current_emp.salary_type = 'daily' THEN
    entry_gross_salary := (base_salary / 22.0) * COALESCE(days_worked, 0);
    entry_overtime_pay := ((base_salary / 22.0) / 8.0) * 1.5 * COALESCE(overtime_hours, 0);
  ELSE
    RAISE EXCEPTION 'Tipo de salario no soportado: %', current_emp.salary_type;
  END IF;
  entry_total_gross := entry_gross_salary + entry_overtime_pay + bonuses + commissions;
  IF entry_total_gross <= 416220 THEN
    entry_tax := entry_total_gross * 0.15;
  ELSE
    entry_tax := 416220 * 0.15 + (entry_total_gross - 416220) * 0.20;
  END IF;
  entry_social := entry_total_gross * 0.0304;
  entry_total_deductions := entry_tax + entry_social + other_deductions;
  entry_net := entry_total_gross - entry_total_deductions;

  entry_gross_salary := ROUND(entry_gross_salary, 2);
  entry_overtime_pay := ROUND(entry_overtime_pay, 2);
  entry_total_gross := ROUND(entry_total_gross, 2);
  entry_tax := ROUND(entry_tax, 2);
  entry_social := ROUND(entry_social, 2);
  entry_total_deductions := ROUND(entry_total_deductions, 2);
  entry_net := ROUND(entry_net, 2);

  entry_id := gen_random_uuid();
  INSERT INTO payroll_entries (
    id, payroll_period_id, employee_id, employee_code, employee_name,
    department, position, salary, salary_type, hours_worked, days_worked,
    gross_salary, overtime_hours, overtime_pay, bonuses, commissions,
    total_gross, tax_deduction, social_security, other_deductions,
    total_deductions, net_salary, notes, created_at, updated_at
  )
  VALUES (
    entry_id, period_jan, current_emp.id, current_emp.employee_code,
    current_emp.first_name || ' ' || current_emp.last_name, current_emp.department, current_emp.position,
    entry_gross_salary, current_emp.salary_type, hours_worked, days_worked,
    entry_gross_salary, overtime_hours, entry_overtime_pay, bonuses, commissions,
    entry_total_gross, entry_tax, entry_social, other_deductions,
    entry_total_deductions, entry_net, 'Guardias nocturnas y soporte crítico a clientes',
    NOW() - INTERVAL '45 days', NOW() - INTERVAL '45 days'
  );

  UPDATE payroll_periods
  SET total_gross = ROUND(total_gross + entry_total_gross, 2),
      total_deductions = ROUND(total_deductions + entry_total_deductions, 2),
      total_net = ROUND(total_net + entry_net, 2),
      updated_at = NOW()
  WHERE id = period_jan;

  -- Período enero 2025 - Raúl Peña
  current_emp := emp_log;
  hours_worked := NULL;
  days_worked := 24;
  overtime_hours := 4;
  bonuses := 1800;
  commissions := 0;
  other_deductions := 600;
  base_salary := COALESCE(current_emp.salary, 0);
  IF current_emp.salary_type = 'monthly' THEN
    entry_gross_salary := base_salary;
    entry_overtime_pay := (base_salary / 160.0) * 1.5 * COALESCE(overtime_hours, 0);
  ELSIF current_emp.salary_type = 'hourly' THEN
    entry_gross_salary := base_salary * COALESCE(hours_worked, 0);
    entry_overtime_pay := base_salary * 1.5 * COALESCE(overtime_hours, 0);
  ELSIF current_emp.salary_type = 'daily' THEN
    entry_gross_salary := (base_salary / 22.0) * COALESCE(days_worked, 0);
    entry_overtime_pay := ((base_salary / 22.0) / 8.0) * 1.5 * COALESCE(overtime_hours, 0);
  ELSE
    RAISE EXCEPTION 'Tipo de salario no soportado: %', current_emp.salary_type;
  END IF;
  entry_total_gross := entry_gross_salary + entry_overtime_pay + bonuses + commissions;
  IF entry_total_gross <= 416220 THEN
    entry_tax := entry_total_gross * 0.15;
  ELSE
    entry_tax := 416220 * 0.15 + (entry_total_gross - 416220) * 0.20;
  END IF;
  entry_social := entry_total_gross * 0.0304;
  entry_total_deductions := entry_tax + entry_social + other_deductions;
  entry_net := entry_total_gross - entry_total_deductions;

  entry_gross_salary := ROUND(entry_gross_salary, 2);
  entry_overtime_pay := ROUND(entry_overtime_pay, 2);
  entry_total_gross := ROUND(entry_total_gross, 2);
  entry_tax := ROUND(entry_tax, 2);
  entry_social := ROUND(entry_social, 2);
  entry_total_deductions := ROUND(entry_total_deductions, 2);
  entry_net := ROUND(entry_net, 2);

  entry_id := gen_random_uuid();
  INSERT INTO payroll_entries (
    id, payroll_period_id, employee_id, employee_code, employee_name,
    department, position, salary, salary_type, hours_worked, days_worked,
    gross_salary, overtime_hours, overtime_pay, bonuses, commissions,
    total_gross, tax_deduction, social_security, other_deductions,
    total_deductions, net_salary, notes, created_at, updated_at
  )
  VALUES (
    entry_id, period_jan, current_emp.id, current_emp.employee_code,
    current_emp.first_name || ' ' || current_emp.last_name, current_emp.department, current_emp.position,
    entry_gross_salary, current_emp.salary_type, hours_worked, days_worked,
    entry_gross_salary, overtime_hours, entry_overtime_pay, bonuses, commissions,
    entry_total_gross, entry_tax, entry_social, other_deductions,
    entry_total_deductions, entry_net, 'Turnos extendidos por incremento en despachos',
    NOW() - INTERVAL '45 days', NOW() - INTERVAL '45 days'
  );

  UPDATE payroll_periods
  SET total_gross = ROUND(total_gross + entry_total_gross, 2),
      total_deductions = ROUND(total_deductions + entry_total_deductions, 2),
      total_net = ROUND(total_net + entry_net, 2),
      updated_at = NOW()
  WHERE id = period_jan;

  -- Período febrero 2025 - Carolina Martínez
  current_emp := emp_ops;
  hours_worked := NULL;
  days_worked := NULL;
  overtime_hours := 6;
  bonuses := 4500;
  commissions := 0;
  other_deductions := 2500;
  base_salary := COALESCE(current_emp.salary, 0);
  IF current_emp.salary_type = 'monthly' THEN
    entry_gross_salary := base_salary;
    entry_overtime_pay := (base_salary / 160.0) * 1.5 * COALESCE(overtime_hours, 0);
  ELSIF current_emp.salary_type = 'hourly' THEN
    entry_gross_salary := base_salary * COALESCE(hours_worked, 0);
    entry_overtime_pay := base_salary * 1.5 * COALESCE(overtime_hours, 0);
  ELSIF current_emp.salary_type = 'daily' THEN
    entry_gross_salary := (base_salary / 22.0) * COALESCE(days_worked, 0);
    entry_overtime_pay := ((base_salary / 22.0) / 8.0) * 1.5 * COALESCE(overtime_hours, 0);
  ELSE
    RAISE EXCEPTION 'Tipo de salario no soportado: %', current_emp.salary_type;
  END IF;
  entry_total_gross := entry_gross_salary + entry_overtime_pay + bonuses + commissions;
  IF entry_total_gross <= 416220 THEN
    entry_tax := entry_total_gross * 0.15;
  ELSE
    entry_tax := 416220 * 0.15 + (entry_total_gross - 416220) * 0.20;
  END IF;
  entry_social := entry_total_gross * 0.0304;
  entry_total_deductions := entry_tax + entry_social + other_deductions;
  entry_net := entry_total_gross - entry_total_deductions;

  entry_gross_salary := ROUND(entry_gross_salary, 2);
  entry_overtime_pay := ROUND(entry_overtime_pay, 2);
  entry_total_gross := ROUND(entry_total_gross, 2);
  entry_tax := ROUND(entry_tax, 2);
  entry_social := ROUND(entry_social, 2);
  entry_total_deductions := ROUND(entry_total_deductions, 2);
  entry_net := ROUND(entry_net, 2);

  entry_id := gen_random_uuid();
  INSERT INTO payroll_entries (
    id, payroll_period_id, employee_id, employee_code, employee_name,
    department, position, salary, salary_type, hours_worked, days_worked,
    gross_salary, overtime_hours, overtime_pay, bonuses, commissions,
    total_gross, tax_deduction, social_security, other_deductions,
    total_deductions, net_salary, notes, created_at, updated_at
  )
  VALUES (
    entry_id, period_feb, current_emp.id, current_emp.employee_code,
    current_emp.first_name || ' ' || current_emp.last_name, current_emp.department, current_emp.position,
    entry_gross_salary, current_emp.salary_type, hours_worked, days_worked,
    entry_gross_salary, overtime_hours, entry_overtime_pay, bonuses, commissions,
    entry_total_gross, entry_tax, entry_social, other_deductions,
    entry_total_deductions, entry_net, 'Planificación de expansión logística',
    NOW() - INTERVAL '15 days', NOW() - INTERVAL '15 days'
  );

  UPDATE payroll_periods
  SET total_gross = ROUND(total_gross + entry_total_gross, 2),
      total_deductions = ROUND(total_deductions + entry_total_deductions, 2),
      total_net = ROUND(total_net + entry_net, 2),
      updated_at = NOW()
  WHERE id = period_feb;

  -- Período febrero 2025 - Luis Fernández
  current_emp := emp_sales;
  hours_worked := NULL;
  days_worked := NULL;
  overtime_hours := 0;
  bonuses := 2500;
  commissions := 8500;
  other_deductions := 1800;
  base_salary := COALESCE(current_emp.salary, 0);
  IF current_emp.salary_type = 'monthly' THEN
    entry_gross_salary := base_salary;
    entry_overtime_pay := (base_salary / 160.0) * 1.5 * COALESCE(overtime_hours, 0);
  ELSIF current_emp.salary_type = 'hourly' THEN
    entry_gross_salary := base_salary * COALESCE(hours_worked, 0);
    entry_overtime_pay := base_salary * 1.5 * COALESCE(overtime_hours, 0);
  ELSIF current_emp.salary_type = 'daily' THEN
    entry_gross_salary := (base_salary / 22.0) * COALESCE(days_worked, 0);
    entry_overtime_pay := ((base_salary / 22.0) / 8.0) * 1.5 * COALESCE(overtime_hours, 0);
  ELSE
    RAISE EXCEPTION 'Tipo de salario no soportado: %', current_emp.salary_type;
  END IF;
  entry_total_gross := entry_gross_salary + entry_overtime_pay + bonuses + commissions;
  IF entry_total_gross <= 416220 THEN
    entry_tax := entry_total_gross * 0.15;
  ELSE
    entry_tax := 416220 * 0.15 + (entry_total_gross - 416220) * 0.20;
  END IF;
  entry_social := entry_total_gross * 0.0304;
  entry_total_deductions := entry_tax + entry_social + other_deductions;
  entry_net := entry_total_gross - entry_total_deductions;

  entry_gross_salary := ROUND(entry_gross_salary, 2);
  entry_overtime_pay := ROUND(entry_overtime_pay, 2);
  entry_total_gross := ROUND(entry_total_gross, 2);
  entry_tax := ROUND(entry_tax, 2);
  entry_social := ROUND(entry_social, 2);
  entry_total_deductions := ROUND(entry_total_deductions, 2);
  entry_net := ROUND(entry_net, 2);

  entry_id := gen_random_uuid();
  INSERT INTO payroll_entries (
    id, payroll_period_id, employee_id, employee_code, employee_name,
    department, position, salary, salary_type, hours_worked, days_worked,
    gross_salary, overtime_hours, overtime_pay, bonuses, commissions,
    total_gross, tax_deduction, social_security, other_deductions,
    total_deductions, net_salary, notes, created_at, updated_at
  )
  VALUES (
    entry_id, period_feb, current_emp.id, current_emp.employee_code,
    current_emp.first_name || ' ' || current_emp.last_name, current_emp.department, current_emp.position,
    entry_gross_salary, current_emp.salary_type, hours_worked, days_worked,
    entry_gross_salary, overtime_hours, entry_overtime_pay, bonuses, commissions,
    entry_total_gross, entry_tax, entry_social, other_deductions,
    entry_total_deductions, entry_net, 'Campaña de relanzamiento de servicios premium',
    NOW() - INTERVAL '15 days', NOW() - INTERVAL '15 days'
  );

  UPDATE payroll_periods
  SET total_gross = ROUND(total_gross + entry_total_gross, 2),
      total_deductions = ROUND(total_deductions + entry_total_deductions, 2),
      total_net = ROUND(total_net + entry_net, 2),
      updated_at = NOW()
  WHERE id = period_feb;
END;
$$;

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
-- PROVEEDORES
-- ====================

INSERT INTO suppliers (id, name, tax_id, email, phone, address, city, state, postal_code, country, is_active, created_at, updated_at)
VALUES
  (gen_random_uuid(), 'Suministros Oficina SRL', '1-31-000001', 'contacto@suministrosoficina.do', '809-555-1001', 'Av. 27 de Febrero #100', 'Santo Domingo', 'DN', '10109', 'DO', true, NOW(), NOW()),
  (gen_random_uuid(), 'Tecnología Caribe SAS', '1-30-000002', 'ventas@tec-caribe.do', '809-555-1002', 'Av. Churchill #250', 'Santo Domingo', 'DN', '10147', 'DO', true, NOW(), NOW()),
  (gen_random_uuid(), 'Servicios de Limpieza CleanPro', '1-32-000003', 'info@cleanpro.do', '809-555-1003', 'Calle Máximo Gómez #45', 'Santo Domingo', 'DN', '10112', 'DO', true, NOW(), NOW()),
  (gen_random_uuid(), 'Redes y Telecom SRL', '1-33-000004', 'soporte@redestelecom.do', '809-555-1004', 'Av. Sarasota #350', 'Santo Domingo', 'DN', '10148', 'DO', true, NOW(), NOW());

-- Capturar IDs de proveedores para asociarlos a productos
DO $$
DECLARE
  prov_oficina UUID;
  prov_tecnologia UUID;
  prov_limpieza UUID;
  prov_telecom UUID;
BEGIN
  SELECT id INTO prov_oficina FROM suppliers WHERE name = 'Suministros Oficina SRL' LIMIT 1;
  SELECT id INTO prov_tecnologia FROM suppliers WHERE name = 'Tecnología Caribe SAS' LIMIT 1;
  SELECT id INTO prov_limpieza FROM suppliers WHERE name = 'Servicios de Limpieza CleanPro' LIMIT 1;
  SELECT id INTO prov_telecom FROM suppliers WHERE name = 'Redes y Telecom SRL' LIMIT 1;

  -- Asociaciones por categoría de producto (por código)
  UPDATE products SET supplier_id = prov_tecnologia WHERE code LIKE 'LAPTOP-%' OR code IN ('MONITOR-001','MOUSE-001','KEYBOARD-001');
  UPDATE products SET supplier_id = prov_oficina   WHERE code IN ('PAPER-001','PEN-001','FOLDER-001');
  UPDATE products SET supplier_id = prov_tecnologia WHERE code IN ('SERV-002'); -- desarrollo web
  UPDATE products SET supplier_id = prov_telecom   WHERE code IN ('SERV-003'); -- soporte técnico
  UPDATE products SET supplier_id = prov_oficina   WHERE code IN ('SERV-001'); -- consultoría (genérica)
  UPDATE products SET supplier_id = prov_oficina   WHERE code IN ('FOOD-001','FOOD-002');
END $$;

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
SELECT 'Proveedores insertados:', COUNT(*) FROM suppliers
UNION ALL
SELECT 'Productos insertados:', COUNT(*) FROM products
UNION ALL
SELECT 'Facturas insertadas:', COUNT(*) FROM invoices
UNION ALL
SELECT 'Líneas de factura:', COUNT(*) FROM invoice_lines
UNION ALL
SELECT 'Gastos insertados:', COUNT(*) FROM payments
UNION ALL
SELECT 'Empleados insertados:', COUNT(*) FROM employees
UNION ALL
SELECT 'Períodos de nómina:', COUNT(*) FROM payroll_periods
UNION ALL
SELECT 'Entradas de nómina:', COUNT(*) FROM payroll_entries;

-- Mostrar resumen de facturas por estado
SELECT 
  status,
  COUNT(*) as cantidad,
  SUM(total) as total
FROM invoices
GROUP BY status
ORDER BY status;

-- Resumen de nómina por período
SELECT
  period_code,
  status,
  total_gross,
  total_deductions,
  total_net
FROM payroll_periods
ORDER BY period_code;
