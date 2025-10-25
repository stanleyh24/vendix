package database

// getPublicMigrations returns migrations for the public schema
// These are shared tables across all tenants
func getPublicMigrations() []Migration {
	return []Migration{
		{
			Version: 1,
			Name:    "create_tenants_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS tenants (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					name VARCHAR(255) NOT NULL,
					slug VARCHAR(100) UNIQUE NOT NULL,
					schema_name VARCHAR(63) UNIQUE NOT NULL,
					domain VARCHAR(255),
					status VARCHAR(50) NOT NULL DEFAULT 'active',
					plan_id UUID,
					settings JSONB DEFAULT '{}',
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					deleted_at TIMESTAMP
				);

				CREATE INDEX idx_tenants_slug ON tenants(slug);
				CREATE INDEX idx_tenants_domain ON tenants(domain);
				CREATE INDEX idx_tenants_status ON tenants(status);
			`,
		},
		{
			Version: 2,
			Name:    "create_plans_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS plans (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					name VARCHAR(255) NOT NULL,
					slug VARCHAR(100) UNIQUE NOT NULL,
					description TEXT,
					price DECIMAL(10, 2) NOT NULL,
					currency VARCHAR(3) NOT NULL DEFAULT 'DOP',
					billing_period VARCHAR(50) NOT NULL,
					max_invoices_per_month INTEGER,
					max_users INTEGER,
					features JSONB DEFAULT '[]',
					is_active BOOLEAN DEFAULT true,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_plans_slug ON plans(slug);
				CREATE INDEX idx_plans_is_active ON plans(is_active);
			`,
		},
		{
			Version: 3,
			Name:    "create_plan_modules_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS plan_modules (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
					module_name VARCHAR(100) NOT NULL,
					is_enabled BOOLEAN DEFAULT true,
					config JSONB DEFAULT '{}',
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(plan_id, module_name)
				);

				CREATE INDEX idx_plan_modules_plan_id ON plan_modules(plan_id);
			`,
		},
		{
			Version: 4,
			Name:    "create_subscriptions_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS subscriptions (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
					plan_id UUID NOT NULL REFERENCES plans(id),
					stripe_subscription_id VARCHAR(255),
					status VARCHAR(50) NOT NULL DEFAULT 'active',
					current_period_start TIMESTAMP NOT NULL,
					current_period_end TIMESTAMP NOT NULL,
					cancel_at TIMESTAMP,
					canceled_at TIMESTAMP,
					trial_start TIMESTAMP,
					trial_end TIMESTAMP,
					metadata JSONB DEFAULT '{}',
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX idx_subscriptions_tenant_id ON subscriptions(tenant_id);
				CREATE INDEX idx_subscriptions_status ON subscriptions(status);
				CREATE INDEX idx_subscriptions_stripe_subscription_id ON subscriptions(stripe_subscription_id);
			`,
		},
		{
			Version: 5,
			Name:    "create_tenant_usage_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS tenant_usage (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
					metric_name VARCHAR(100) NOT NULL,
					metric_value INTEGER NOT NULL DEFAULT 0,
					period_start TIMESTAMP NOT NULL,
					period_end TIMESTAMP NOT NULL,
					metadata JSONB DEFAULT '{}',
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(tenant_id, metric_name, period_start)
				);

				CREATE INDEX idx_tenant_usage_tenant_id ON tenant_usage(tenant_id);
				CREATE INDEX idx_tenant_usage_period ON tenant_usage(period_start, period_end);
			`,
		},
		{
			Version: 6,
			Name:    "insert_default_plans",
			SQL: `
				INSERT INTO plans (name, slug, description, price, currency, billing_period, max_invoices_per_month, max_users, features)
				VALUES 
					('Starter', 'starter', 'Plan básico para pequeños negocios', 29.99, 'USD', 'monthly', 50, 2, '["invoicing", "customers", "basic_reports"]'),
					('Professional', 'professional', 'Plan profesional para negocios en crecimiento', 79.99, 'USD', 'monthly', 200, 10, '["invoicing", "customers", "payments", "reports", "accounting", "pos"]'),
					('Enterprise', 'enterprise', 'Plan empresarial con todas las funcionalidades', 199.99, 'USD', 'monthly', -1, -1, '["invoicing", "customers", "payments", "reports", "accounting", "pos", "api_access", "priority_support"]')
				ON CONFLICT (slug) DO NOTHING;
			`,
		},
	}
}
