package tenantconfig

import (
	"context"
	"fmt"
	"time"

	"vendix/internal/database"
)

type Repository struct {
	db *database.DB
}

// NewRepository creates a new tenant config repository
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// Get retrieves the tenant configuration (there should be only one row per tenant)
func (r *Repository) Get(ctx context.Context, schema string) (*TenantConfig, error) {
	var config TenantConfig
	query := fmt.Sprintf(`
		SELECT 
			id, company_legal_name, company_trade_name, company_tax_id, company_email, 
			company_phone, company_website, address_line1, address_line2, city, 
			state_province, postal_code, country, tax_regime, dgii_rnc, dgii_user, 
			dgii_api_key, dgii_environment, ncf_enabled, invoice_prefix, 
			invoice_next_number, invoice_terms, invoice_footer, payment_terms_days,
			ncf_fiscal_credit_prefix, ncf_fiscal_credit_sequence, ncf_consumer_prefix,
			ncf_consumer_sequence, ncf_debit_note_prefix, ncf_debit_note_sequence,
			ncf_credit_note_prefix, ncf_credit_note_sequence, 
			ncf_gov_prefix, ncf_gov_sequence, ncf_gov_end_range,
			default_currency,
			default_tax_rate, notification_email, send_invoice_emails, 
			send_payment_reminders, logo_url, primary_color, secondary_color,
			timezone, date_format, time_format, locale, created_at, updated_at
		FROM %s.tenant_config
		LIMIT 1
	`, schema)

	err := r.db.GetContext(ctx, &config, query)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Create creates the initial tenant configuration
func (r *Repository) Create(ctx context.Context, schema string) (*TenantConfig, error) {
	config := &TenantConfig{
		CompanyLegalName:        "Mi Empresa",
		Country:                 "DO",
		DGIIEnvironment:         "sandbox",
		NCFEnabled:              true,
		InvoicePrefix:           "INV",
		InvoiceNextNumber:       1,
		PaymentTermsDays:        30,
		DefaultCurrency:         "DOP",
		DefaultTaxRate:          18.00,
		SendInvoiceEmails:       true,
		SendPaymentReminders:    true,
		PrimaryColor:            "#3B82F6",
		SecondaryColor:          "#10B981",
		Timezone:                "America/Santo_Domingo",
		DateFormat:              "DD/MM/YYYY",
		TimeFormat:              "HH:mm",
		Locale:                  "es_DO",
		NCFFiscalCreditSequence: 1,
		NCFConsumerSequence:     1,
		NCFDebitNoteSequence:    1,
		NCFCreditNoteSequence:   1,
		NCFGovSequence:          1,
		NCFGovEndRange:          999999999,
	}

	query := fmt.Sprintf(`
		INSERT INTO %s.tenant_config (
			company_legal_name, country, dgii_environment, ncf_enabled,
			invoice_prefix, invoice_next_number, payment_terms_days,
			ncf_fiscal_credit_sequence, ncf_consumer_sequence,
			ncf_debit_note_sequence, ncf_credit_note_sequence,
			ncf_gov_sequence, ncf_gov_end_range,
			default_currency, default_tax_rate, send_invoice_emails,
			send_payment_reminders, primary_color, secondary_color,
			timezone, date_format, time_format, locale
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
		RETURNING id, created_at, updated_at
	`, schema)

	err := r.db.QueryRowContext(ctx, query,
		config.CompanyLegalName,
		config.Country,
		config.DGIIEnvironment,
		config.NCFEnabled,
		config.InvoicePrefix,
		config.InvoiceNextNumber,
		config.PaymentTermsDays,
		config.NCFFiscalCreditSequence,
		config.NCFConsumerSequence,
		config.NCFDebitNoteSequence,
		config.NCFCreditNoteSequence,
		config.NCFGovSequence,
		config.NCFGovEndRange,
		config.DefaultCurrency,
		config.DefaultTaxRate,
		config.SendInvoiceEmails,
		config.SendPaymentReminders,
		config.PrimaryColor,
		config.SecondaryColor,
		config.Timezone,
		config.DateFormat,
		config.TimeFormat,
		config.Locale,
	).Scan(&config.ID, &config.CreatedAt, &config.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return config, nil
}

// Update updates the tenant configuration
func (r *Repository) Update(ctx context.Context, schema string, config *TenantConfig) error {
	query := fmt.Sprintf(`
		UPDATE %s.tenant_config
		SET 
			company_legal_name = $1,
			company_trade_name = $2,
			company_tax_id = $3,
			company_email = $4,
			company_phone = $5,
			company_website = $6,
			address_line1 = $7,
			address_line2 = $8,
			city = $9,
			state_province = $10,
			postal_code = $11,
			country = $12,
			tax_regime = $13,
			dgii_rnc = $14,
			dgii_user = $15,
			dgii_api_key = $16,
			dgii_environment = $17,
			ncf_enabled = $18,
			invoice_prefix = $19,
			invoice_next_number = $20,
			invoice_terms = $21,
			invoice_footer = $22,
			payment_terms_days = $23,
			ncf_fiscal_credit_prefix = $24,
			ncf_fiscal_credit_sequence = $25,
			ncf_consumer_prefix = $26,
			ncf_consumer_sequence = $27,
			ncf_debit_note_prefix = $28,
			ncf_debit_note_sequence = $29,
			ncf_credit_note_prefix = $30,
			ncf_credit_note_sequence = $31,
			ncf_gov_prefix = $32,
			ncf_gov_sequence = $33,
			ncf_gov_end_range = $34,
			default_currency = $35,
			default_tax_rate = $36,
			notification_email = $37,
			send_invoice_emails = $38,
			send_payment_reminders = $39,
			logo_url = $40,
			primary_color = $41,
			secondary_color = $42,
			timezone = $43,
			date_format = $44,
			time_format = $45,
			locale = $46,
			updated_at = $47
		WHERE id = $48
	`, schema)

	result, err := r.db.ExecContext(ctx, query,
		config.CompanyLegalName,
		config.CompanyTradeName,
		config.CompanyTaxID,
		config.CompanyEmail,
		config.CompanyPhone,
		config.CompanyWebsite,
		config.AddressLine1,
		config.AddressLine2,
		config.City,
		config.StateProvince,
		config.PostalCode,
		config.Country,
		config.TaxRegime,
		config.DGIIIRNC,
		config.DGIIUser,
		config.DGIIAPIKey,
		config.DGIIEnvironment,
		config.NCFEnabled,
		config.InvoicePrefix,
		config.InvoiceNextNumber,
		config.InvoiceTerms,
		config.InvoiceFooter,
		config.PaymentTermsDays,
		config.NCFFiscalCreditPrefix,
		config.NCFFiscalCreditSequence,
		config.NCFConsumerPrefix,
		config.NCFConsumerSequence,
		config.NCFDebitNotePrefix,
		config.NCFDebitNoteSequence,
		config.NCFCreditNotePrefix,
		config.NCFCreditNoteSequence,
		config.NCFGovPrefix,
		config.NCFGovSequence,
		config.NCFGovEndRange,
		config.DefaultCurrency,
		config.DefaultTaxRate,
		config.NotificationEmail,
		config.SendInvoiceEmails,
		config.SendPaymentReminders,
		config.LogoURL,
		config.PrimaryColor,
		config.SecondaryColor,
		config.Timezone,
		config.DateFormat,
		config.TimeFormat,
		config.Locale,
		time.Now(),
		config.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no configuration found to update")
	}

	return nil
}
