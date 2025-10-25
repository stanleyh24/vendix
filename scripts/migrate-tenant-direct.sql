-- Script para crear la tabla tenant_config en un tenant específico
-- Uso: psql -h localhost -U postgres -d vendix -v schema=tenant_demo -f scripts/migrate-tenant-direct.sql

\echo 'Creando tabla tenant_config en el schema :schema'

SET search_path TO :schema, public;

-- Migración v15: create_tenant_config_table
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

-- Insertar configuración por defecto si no existe
INSERT INTO tenant_config (company_legal_name) 
SELECT 'Mi Empresa'
WHERE NOT EXISTS (SELECT 1 FROM tenant_config);

-- Registrar la migración si existe la tabla schema_migrations
INSERT INTO schema_migrations (version, name)
SELECT 15, 'create_tenant_config_table'
WHERE EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = 'schema_migrations')
  AND NOT EXISTS (SELECT 1 FROM schema_migrations WHERE version = 15);

\echo 'Tabla tenant_config creada exitosamente en :schema'

