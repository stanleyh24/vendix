package database

// getTenantMigrations returns migrations for tenant schemas
// Each tenant has their own isolated database schema
func getTenantMigrations() []Migration {
	return []Migration{
		{
			Version: 1,
			Name:    "create_users_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS users (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					email VARCHAR(255) UNIQUE NOT NULL,
					password_hash VARCHAR(255) NOT NULL,
					first_name VARCHAR(100),
					last_name VARCHAR(100),
					role VARCHAR(50) NOT NULL DEFAULT 'viewer',
					is_active BOOLEAN DEFAULT true,
					last_login_at TIMESTAMP,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_users_email ON users(email);
				CREATE INDEX idx_users_role ON users(role);
			`,
		},
		{
			Version: 2,
			Name:    "create_refresh_tokens_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS refresh_tokens (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
					token VARCHAR(500) UNIQUE NOT NULL,
					expires_at TIMESTAMP NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
				CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
			`,
		},
		{
			Version: 3,
			Name:    "create_customers_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS customers (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					customer_type VARCHAR(50) NOT NULL DEFAULT 'individual',
					tax_id VARCHAR(50) UNIQUE,
					name VARCHAR(255) NOT NULL,
					email VARCHAR(255),
					phone VARCHAR(50),
					address TEXT,
					city VARCHAR(100),
					state VARCHAR(100),
					postal_code VARCHAR(20),
					country VARCHAR(2) DEFAULT 'DO',
					is_active BOOLEAN DEFAULT true,
					metadata JSONB DEFAULT '{}',
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_customers_tax_id ON customers(tax_id);
				CREATE INDEX idx_customers_email ON customers(email);
				CREATE INDEX idx_customers_name ON customers(name);
			`,
		},
		{
			Version: 4,
			Name:    "create_products_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS products (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					code VARCHAR(100) UNIQUE NOT NULL,
					name VARCHAR(255) NOT NULL,
					description TEXT,
					product_type VARCHAR(50) NOT NULL DEFAULT 'product',
					unit VARCHAR(50) NOT NULL DEFAULT 'unit',
					price DECIMAL(10, 2) NOT NULL DEFAULT 0,
					cost DECIMAL(10, 2),
					tax_rate DECIMAL(5, 2) DEFAULT 0,
					is_active BOOLEAN DEFAULT true,
					metadata JSONB DEFAULT '{}',
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_products_code ON products(code);
				CREATE INDEX idx_products_product_type ON products(product_type);
			`,
		},
		{
			Version: 5,
			Name:    "create_tax_rates_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS tax_rates (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					name VARCHAR(255) NOT NULL,
					code VARCHAR(50) UNIQUE NOT NULL,
					rate DECIMAL(5, 2) NOT NULL,
					is_compound BOOLEAN DEFAULT false,
					is_active BOOLEAN DEFAULT true,
					description TEXT,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_tax_rates_code ON tax_rates(code);

				-- Insert default Dominican Republic tax rates
				INSERT INTO tax_rates (name, code, rate, is_active) VALUES
					('ITBIS 18%', 'ITBIS18', 18.00, true),
					('ITBIS 16%', 'ITBIS16', 16.00, true),
					('Exento', 'EXEMPT', 0.00, true);
			`,
		},
		{
			Version: 6,
			Name:    "create_invoices_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS invoices (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					invoice_number VARCHAR(100) UNIQUE NOT NULL,
					ncf VARCHAR(50),
					customer_id UUID NOT NULL REFERENCES customers(id),
					issue_date DATE NOT NULL,
					due_date DATE NOT NULL,
					status VARCHAR(50) NOT NULL DEFAULT 'draft',
					subtotal DECIMAL(10, 2) NOT NULL DEFAULT 0,
					tax_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
					total DECIMAL(10, 2) NOT NULL DEFAULT 0,
					paid_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
					currency VARCHAR(3) NOT NULL DEFAULT 'DOP',
					notes TEXT,
					terms TEXT,
					dgii_status VARCHAR(50),
					dgii_response JSONB,
					signed_at TIMESTAMP,
					sent_at TIMESTAMP,
					metadata JSONB DEFAULT '{}',
					created_by UUID REFERENCES users(id),
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_invoices_customer_id ON invoices(customer_id);
				CREATE INDEX idx_invoices_status ON invoices(status);
				CREATE INDEX idx_invoices_issue_date ON invoices(issue_date);
				CREATE INDEX idx_invoices_invoice_number ON invoices(invoice_number);
			`,
		},
		{
			Version: 7,
			Name:    "create_invoice_lines_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS invoice_lines (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
					product_id UUID REFERENCES products(id),
					line_number INTEGER NOT NULL,
					description TEXT NOT NULL,
					quantity DECIMAL(10, 2) NOT NULL DEFAULT 1,
					unit_price DECIMAL(10, 2) NOT NULL,
					tax_rate DECIMAL(5, 2) DEFAULT 0,
					tax_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
					line_total DECIMAL(10, 2) NOT NULL,
					metadata JSONB DEFAULT '{}',
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(invoice_id, line_number)
				);

				CREATE INDEX idx_invoice_lines_invoice_id ON invoice_lines(invoice_id);
			`,
		},
		{
			Version: 8,
			Name:    "create_payments_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS payments (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					payment_number VARCHAR(100) UNIQUE NOT NULL,
					customer_id UUID NOT NULL REFERENCES customers(id),
					payment_method VARCHAR(50) NOT NULL,
					payment_date DATE NOT NULL,
					amount DECIMAL(10, 2) NOT NULL,
					currency VARCHAR(3) NOT NULL DEFAULT 'DOP',
					reference VARCHAR(255),
					notes TEXT,
					stripe_payment_id VARCHAR(255),
					status VARCHAR(50) NOT NULL DEFAULT 'completed',
					metadata JSONB DEFAULT '{}',
					created_by UUID REFERENCES users(id),
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_payments_customer_id ON payments(customer_id);
				CREATE INDEX idx_payments_payment_date ON payments(payment_date);
				CREATE INDEX idx_payments_stripe_payment_id ON payments(stripe_payment_id);
			`,
		},
		{
			Version: 9,
			Name:    "create_payment_allocations_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS payment_allocations (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
					invoice_id UUID NOT NULL REFERENCES invoices(id),
					amount DECIMAL(10, 2) NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_payment_allocations_payment_id ON payment_allocations(payment_id);
				CREATE INDEX idx_payment_allocations_invoice_id ON payment_allocations(invoice_id);
			`,
		},
		{
			Version: 10,
			Name:    "create_credit_notes_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS credit_notes (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					credit_note_number VARCHAR(100) UNIQUE NOT NULL,
					invoice_id UUID NOT NULL REFERENCES invoices(id),
					ncf VARCHAR(50),
					customer_id UUID NOT NULL REFERENCES customers(id),
					issue_date DATE NOT NULL,
					reason TEXT NOT NULL,
					subtotal DECIMAL(10, 2) NOT NULL DEFAULT 0,
					tax_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
					total DECIMAL(10, 2) NOT NULL DEFAULT 0,
					status VARCHAR(50) NOT NULL DEFAULT 'draft',
					dgii_status VARCHAR(50),
					dgii_response JSONB,
					metadata JSONB DEFAULT '{}',
					created_by UUID REFERENCES users(id),
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_credit_notes_invoice_id ON credit_notes(invoice_id);
				CREATE INDEX idx_credit_notes_customer_id ON credit_notes(customer_id);
			`,
		},
		{
			Version: 11,
			Name:    "create_journal_entries_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS journal_entries (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					entry_number VARCHAR(100) UNIQUE NOT NULL,
					entry_date DATE NOT NULL,
					description TEXT NOT NULL,
					reference VARCHAR(255),
					status VARCHAR(50) NOT NULL DEFAULT 'draft',
					created_by UUID REFERENCES users(id),
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_journal_entries_entry_date ON journal_entries(entry_date);
			`,
		},
		{
			Version: 12,
			Name:    "create_journal_entry_lines_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS journal_entry_lines (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					journal_entry_id UUID NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
					account_code VARCHAR(50) NOT NULL,
					account_name VARCHAR(255) NOT NULL,
					debit DECIMAL(10, 2) NOT NULL DEFAULT 0,
					credit DECIMAL(10, 2) NOT NULL DEFAULT 0,
					description TEXT,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_journal_entry_lines_journal_entry_id ON journal_entry_lines(journal_entry_id);
			`,
		},
		{
			Version: 13,
			Name:    "create_audit_logs_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS audit_logs (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					user_id UUID REFERENCES users(id),
					entity_type VARCHAR(100) NOT NULL,
					entity_id UUID NOT NULL,
					action VARCHAR(50) NOT NULL,
					old_values JSONB,
					new_values JSONB,
					ip_address VARCHAR(45),
					user_agent TEXT,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
				CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
				CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
			`,
		},
		{
			Version: 14,
			Name:    "create_settings_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS settings (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					key VARCHAR(255) UNIQUE NOT NULL,
					value JSONB NOT NULL,
					description TEXT,
					is_public BOOLEAN DEFAULT false,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_settings_key ON settings(key);

				-- Insert default settings
				INSERT INTO settings (key, value, description, is_public) VALUES
					('company_name', '"Mi Empresa"', 'Nombre de la empresa', true),
					('company_tax_id', '""', 'RNC de la empresa', true),
					('company_address', '""', 'Dirección de la empresa', true),
					('company_phone', '""', 'Teléfono de la empresa', true),
					('company_email', '""', 'Email de la empresa', true),
					('invoice_prefix', '"INV-"', 'Prefijo de facturas', false),
					('next_invoice_number', '1', 'Siguiente número de factura', false),
					('default_currency', '"DOP"', 'Moneda por defecto', true),
					('default_tax_rate', '18.00', 'Tasa de impuesto por defecto (ITBIS)', true),
					('ncf_enabled', 'true', 'NCF habilitado para DGII', false),
					('fiscal_year_start', '"01-01"', 'Inicio del año fiscal (MM-DD)', false);
			`,
		},
		{
			Version: 15,
			Name:    "create_chart_of_accounts_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS chart_of_accounts (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					account_code VARCHAR(50) UNIQUE NOT NULL,
					account_name VARCHAR(255) NOT NULL,
					account_type VARCHAR(50) NOT NULL,
					parent_account_code VARCHAR(50),
					is_active BOOLEAN DEFAULT true,
					description TEXT,
					balance DECIMAL(15, 2) DEFAULT 0,
					metadata JSONB DEFAULT '{}',
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_chart_of_accounts_code ON chart_of_accounts(account_code);
				CREATE INDEX idx_chart_of_accounts_type ON chart_of_accounts(account_type);
				CREATE INDEX idx_chart_of_accounts_parent ON chart_of_accounts(parent_account_code);

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
			`,
		},
		{
			Version: 16,
			Name:    "create_tenant_config_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS tenant_config (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					-- Información básica de la empresa
					company_legal_name VARCHAR(255) NOT NULL DEFAULT '',
					company_trade_name VARCHAR(255),
					company_tax_id VARCHAR(50),
					company_email VARCHAR(255),
					company_phone VARCHAR(50),
					company_website VARCHAR(255),
					
					-- Dirección de la empresa
					address_line1 TEXT,
					address_line2 TEXT,
					city VARCHAR(100),
					state_province VARCHAR(100),
					postal_code VARCHAR(20),
					country VARCHAR(2) DEFAULT 'DO',
					
					-- Configuración fiscal
					tax_regime VARCHAR(100),
					dgii_rnc VARCHAR(50),
					dgii_user VARCHAR(100),
					dgii_api_key VARCHAR(255),
					dgii_environment VARCHAR(20) DEFAULT 'sandbox',
					ncf_enabled BOOLEAN DEFAULT true,
					
					-- Configuración de facturación
					invoice_prefix VARCHAR(20) DEFAULT 'INV',
					invoice_next_number INTEGER DEFAULT 1,
					invoice_terms TEXT,
					invoice_footer TEXT,
					payment_terms_days INTEGER DEFAULT 30,
					
					-- Configuración de numeración NCF
					ncf_fiscal_credit_prefix VARCHAR(20),
					ncf_fiscal_credit_sequence INTEGER DEFAULT 1,
					ncf_consumer_prefix VARCHAR(20),
					ncf_consumer_sequence INTEGER DEFAULT 1,
					ncf_debit_note_prefix VARCHAR(20),
					ncf_debit_note_sequence INTEGER DEFAULT 1,
					ncf_credit_note_prefix VARCHAR(20),
					ncf_credit_note_sequence INTEGER DEFAULT 1,
					
					-- Configuración monetaria
					default_currency VARCHAR(3) DEFAULT 'DOP',
					default_tax_rate DECIMAL(5, 2) DEFAULT 18.00,
					
					-- Configuración de notificaciones
					notification_email VARCHAR(255),
					send_invoice_emails BOOLEAN DEFAULT true,
					send_payment_reminders BOOLEAN DEFAULT true,
					
					-- Logo y branding
					logo_url TEXT,
					primary_color VARCHAR(7) DEFAULT '#3B82F6',
					secondary_color VARCHAR(7) DEFAULT '#10B981',
					
					-- Metadata
					timezone VARCHAR(50) DEFAULT 'America/Santo_Domingo',
					date_format VARCHAR(20) DEFAULT 'DD/MM/YYYY',
					time_format VARCHAR(20) DEFAULT 'HH:mm',
					locale VARCHAR(10) DEFAULT 'es_DO',
					
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				-- Insert default configuration (one row per tenant)
				INSERT INTO tenant_config (company_legal_name) VALUES ('Mi Empresa');
			`,
		},
		{
			Version: 17,
			Name:    "add_stock_quantity_to_products",
			SQL: `
				-- Agregar campo de cantidad en stock a la tabla de productos
				ALTER TABLE products ADD COLUMN IF NOT EXISTS stock_quantity DECIMAL(10, 2) DEFAULT 0 NOT NULL;
				
				-- Agregar índice para búsquedas rápidas de productos con bajo stock
				CREATE INDEX IF NOT EXISTS idx_products_stock_quantity ON products(stock_quantity);
				
				-- Crear tabla para registro de movimientos de inventario (auditoría)
				CREATE TABLE IF NOT EXISTS inventory_movements (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
					movement_type VARCHAR(50) NOT NULL, -- 'sale', 'purchase', 'adjustment', 'return'
					quantity DECIMAL(10, 2) NOT NULL,
					previous_stock DECIMAL(10, 2) NOT NULL,
					new_stock DECIMAL(10, 2) NOT NULL,
					reference_type VARCHAR(50), -- 'invoice', 'credit_note', 'manual'
					reference_id UUID,
					notes TEXT,
					created_by UUID REFERENCES users(id),
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_inventory_movements_product_id ON inventory_movements(product_id);
				CREATE INDEX idx_inventory_movements_created_at ON inventory_movements(created_at);
				CREATE INDEX idx_inventory_movements_reference ON inventory_movements(reference_type, reference_id);
			`,
		},
		{
			Version: 18,
			Name:    "create_generic_customer",
			SQL: `
				-- Crear cliente genérico para ventas sin cliente específico
				-- Este cliente se usa como "Consumidor Final" o "Cliente Genérico"
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
			`,
		},
		{
			Version: 19,
			Name:    "add_ncf_type_to_invoices",
			SQL: `
				-- Agregar campo para tipo de comprobante fiscal (NCF)
				-- En República Dominicana existen diferentes tipos de NCF:
				-- 01 = Crédito Fiscal (para empresas con RNC)
				-- 02 = Consumidor Final (para personas/ventas generales)
				-- 03 = Nota de Débito
				-- 04 = Nota de Crédito
				-- 11 = Proveedores Informales
				-- 12 = Único Ingreso
				-- 13 = Menor Cuantía
				-- 14 = Régimen Especial
				-- 15 = Gubernamental
				
				ALTER TABLE invoices ADD COLUMN IF NOT EXISTS ncf_type VARCHAR(2) DEFAULT '02' NOT NULL;
				
				-- Crear índice para búsquedas por tipo de NCF
				CREATE INDEX IF NOT EXISTS idx_invoices_ncf_type ON invoices(ncf_type);
				
				-- Agregar comentario explicativo
				COMMENT ON COLUMN invoices.ncf_type IS 'Tipo de comprobante fiscal: 01=Crédito Fiscal, 02=Consumidor Final';
			`,
		},
		{
			Version: 20,
			Name:    "create_sales_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS sales (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					sale_number VARCHAR(100) UNIQUE NOT NULL,
					customer_id UUID REFERENCES customers(id),
					invoice_id UUID REFERENCES invoices(id),
					subtotal DECIMAL(10, 2) NOT NULL DEFAULT 0,
					tax_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
					total DECIMAL(10, 2) NOT NULL DEFAULT 0,
					payment_type VARCHAR(50) NOT NULL,
					status VARCHAR(50) NOT NULL DEFAULT 'completed',
					notes TEXT,
					metadata JSONB DEFAULT '{}',
					created_by UUID REFERENCES users(id),
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_sales_customer_id ON sales(customer_id);
				CREATE INDEX idx_sales_invoice_id ON sales(invoice_id);
				CREATE INDEX idx_sales_status ON sales(status);
				CREATE INDEX idx_sales_created_at ON sales(created_at);
				CREATE INDEX idx_sales_sale_number ON sales(sale_number);
			`,
		},
		{
			Version: 21,
			Name:    "create_sale_lines_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS sale_lines (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					sale_id UUID NOT NULL REFERENCES sales(id) ON DELETE CASCADE,
					product_id UUID REFERENCES products(id),
					line_number INTEGER NOT NULL,
					description TEXT NOT NULL,
					quantity DECIMAL(10, 2) NOT NULL DEFAULT 1,
					unit_price DECIMAL(10, 2) NOT NULL,
					tax_rate DECIMAL(5, 2) DEFAULT 0,
					tax_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
					line_total DECIMAL(10, 2) NOT NULL,
					metadata JSONB DEFAULT '{}',
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(sale_id, line_number)
				);

				CREATE INDEX idx_sale_lines_sale_id ON sale_lines(sale_id);
				CREATE INDEX idx_sale_lines_product_id ON sale_lines(product_id);
			`,
		},
		{
			Version: 22,
			Name:    "create_suppliers_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS suppliers (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					name VARCHAR(255) NOT NULL,
					tax_id VARCHAR(50),
					email VARCHAR(255),
					phone VARCHAR(50),
					address TEXT,
					city VARCHAR(100),
					state VARCHAR(100),
					postal_code VARCHAR(20),
					country VARCHAR(2) DEFAULT 'DO',
					is_active BOOLEAN DEFAULT true,
					metadata JSONB DEFAULT '{}',
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX IF NOT EXISTS idx_suppliers_name ON suppliers(name);
				CREATE INDEX IF NOT EXISTS idx_suppliers_tax_id ON suppliers(tax_id);
			`,
		},
		{
			Version: 23,
			Name:    "add_supplier_id_to_products",
			SQL: `
				-- Agregar relación opcional de proveedor a productos
				ALTER TABLE products ADD COLUMN IF NOT EXISTS supplier_id UUID REFERENCES suppliers(id);
				CREATE INDEX IF NOT EXISTS idx_products_supplier_id ON products(supplier_id);
			`,
		},
	}
}
