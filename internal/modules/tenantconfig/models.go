package tenantconfig

import (
	"time"

	"github.com/google/uuid"
)

// TenantConfig represents the tenant configuration settings
type TenantConfig struct {
	ID uuid.UUID `db:"id" json:"id"`

	// Información básica de la empresa
	CompanyLegalName string  `db:"company_legal_name" json:"company_legal_name"`
	CompanyTradeName *string `db:"company_trade_name" json:"company_trade_name,omitempty"`
	CompanyTaxID     *string `db:"company_tax_id" json:"company_tax_id,omitempty"`
	CompanyEmail     *string `db:"company_email" json:"company_email,omitempty"`
	CompanyPhone     *string `db:"company_phone" json:"company_phone,omitempty"`
	CompanyWebsite   *string `db:"company_website" json:"company_website,omitempty"`

	// Dirección de la empresa
	AddressLine1  *string `db:"address_line1" json:"address_line1,omitempty"`
	AddressLine2  *string `db:"address_line2" json:"address_line2,omitempty"`
	City          *string `db:"city" json:"city,omitempty"`
	StateProvince *string `db:"state_province" json:"state_province,omitempty"`
	PostalCode    *string `db:"postal_code" json:"postal_code,omitempty"`
	Country       string  `db:"country" json:"country"`

	// Configuración fiscal
	TaxRegime       *string `db:"tax_regime" json:"tax_regime,omitempty"`
	DGIIIRNC        *string `db:"dgii_rnc" json:"dgii_rnc,omitempty"`
	DGIIUser        *string `db:"dgii_user" json:"dgii_user,omitempty"`
	DGIIAPIKey      *string `db:"dgii_api_key" json:"dgii_api_key,omitempty"`
	DGIIEnvironment string  `db:"dgii_environment" json:"dgii_environment"`
	NCFEnabled      bool    `db:"ncf_enabled" json:"ncf_enabled"`

	// Configuración de facturación
	InvoicePrefix     string  `db:"invoice_prefix" json:"invoice_prefix"`
	InvoiceNextNumber int     `db:"invoice_next_number" json:"invoice_next_number"`
	InvoiceTerms      *string `db:"invoice_terms" json:"invoice_terms,omitempty"`
	InvoiceFooter     *string `db:"invoice_footer" json:"invoice_footer,omitempty"`
	PaymentTermsDays  int     `db:"payment_terms_days" json:"payment_terms_days"`

	// Configuración de numeración NCF
	NCFFiscalCreditPrefix   *string `db:"ncf_fiscal_credit_prefix" json:"ncf_fiscal_credit_prefix,omitempty"`
	NCFFiscalCreditSequence int     `db:"ncf_fiscal_credit_sequence" json:"ncf_fiscal_credit_sequence"`
	NCFConsumerPrefix       *string `db:"ncf_consumer_prefix" json:"ncf_consumer_prefix,omitempty"`
	NCFConsumerSequence     int     `db:"ncf_consumer_sequence" json:"ncf_consumer_sequence"`
	NCFDebitNotePrefix      *string `db:"ncf_debit_note_prefix" json:"ncf_debit_note_prefix,omitempty"`
	NCFDebitNoteSequence    int     `db:"ncf_debit_note_sequence" json:"ncf_debit_note_sequence"`
	NCFCreditNotePrefix     *string `db:"ncf_credit_note_prefix" json:"ncf_credit_note_prefix,omitempty"`
	NCFCreditNoteSequence   int     `db:"ncf_credit_note_sequence" json:"ncf_credit_note_sequence"`
	NCFGovPrefix            *string `db:"ncf_gov_prefix" json:"ncf_gov_prefix,omitempty"`         // Prefijo para NCF tipo 15 (Gubernamental)
	NCFGovSequence          int     `db:"ncf_gov_sequence" json:"ncf_gov_sequence"`               // Secuencia actual para NCF tipo 15
	NCFGovEndRange          int     `db:"ncf_gov_end_range" json:"ncf_gov_end_range"`             // Fin del rango para NCF tipo 15

	// Configuración monetaria
	DefaultCurrency string  `db:"default_currency" json:"default_currency"`
	DefaultTaxRate  float64 `db:"default_tax_rate" json:"default_tax_rate"`

	// Configuración de notificaciones
	NotificationEmail    *string `db:"notification_email" json:"notification_email,omitempty"`
	SendInvoiceEmails    bool    `db:"send_invoice_emails" json:"send_invoice_emails"`
	SendPaymentReminders bool    `db:"send_payment_reminders" json:"send_payment_reminders"`

	// Logo y branding
	LogoURL        *string `db:"logo_url" json:"logo_url,omitempty"`
	PrimaryColor   string  `db:"primary_color" json:"primary_color"`
	SecondaryColor string  `db:"secondary_color" json:"secondary_color"`

	// Metadata
	Timezone   string `db:"timezone" json:"timezone"`
	DateFormat string `db:"date_format" json:"date_format"`
	TimeFormat string `db:"time_format" json:"time_format"`
	Locale     string `db:"locale" json:"locale"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// UpdateTenantConfigRequest represents the request to update tenant configuration
type UpdateTenantConfigRequest struct {
	// Información básica de la empresa
	CompanyLegalName *string `json:"company_legal_name,omitempty"`
	CompanyTradeName *string `json:"company_trade_name,omitempty"`
	CompanyTaxID     *string `json:"company_tax_id,omitempty"`
	CompanyEmail     *string `json:"company_email,omitempty"`
	CompanyPhone     *string `json:"company_phone,omitempty"`
	CompanyWebsite   *string `json:"company_website,omitempty"`

	// Dirección de la empresa
	AddressLine1  *string `json:"address_line1,omitempty"`
	AddressLine2  *string `json:"address_line2,omitempty"`
	City          *string `json:"city,omitempty"`
	StateProvince *string `json:"state_province,omitempty"`
	PostalCode    *string `json:"postal_code,omitempty"`
	Country       *string `json:"country,omitempty"`

	// Configuración fiscal
	TaxRegime       *string `json:"tax_regime,omitempty"`
	DGIIIRNC        *string `json:"dgii_rnc,omitempty"`
	DGIIUser        *string `json:"dgii_user,omitempty"`
	DGIIAPIKey      *string `json:"dgii_api_key,omitempty"`
	DGIIEnvironment *string `json:"dgii_environment,omitempty"`
	NCFEnabled      *bool   `json:"ncf_enabled,omitempty"`

	// Configuración de facturación
	InvoicePrefix     *string `json:"invoice_prefix,omitempty"`
	InvoiceNextNumber *int    `json:"invoice_next_number,omitempty"`
	InvoiceTerms      *string `json:"invoice_terms,omitempty"`
	InvoiceFooter     *string `json:"invoice_footer,omitempty"`
	PaymentTermsDays  *int    `json:"payment_terms_days,omitempty"`

	// Configuración de numeración NCF
	NCFFiscalCreditPrefix   *string `json:"ncf_fiscal_credit_prefix,omitempty"`
	NCFFiscalCreditSequence *int    `json:"ncf_fiscal_credit_sequence,omitempty"`
	NCFConsumerPrefix       *string `json:"ncf_consumer_prefix,omitempty"`
	NCFConsumerSequence     *int    `json:"ncf_consumer_sequence,omitempty"`
	NCFDebitNotePrefix      *string `json:"ncf_debit_note_prefix,omitempty"`
	NCFDebitNoteSequence    *int    `json:"ncf_debit_note_sequence,omitempty"`
	NCFCreditNotePrefix     *string `json:"ncf_credit_note_prefix,omitempty"`
	NCFCreditNoteSequence   *int    `json:"ncf_credit_note_sequence,omitempty"`
	NCFGovPrefix            *string `json:"ncf_gov_prefix,omitempty"`         // Prefijo para NCF tipo 15
	NCFGovSequence          *int    `json:"ncf_gov_sequence,omitempty"`       // Secuencia para NCF tipo 15

	// Configuración monetaria
	DefaultCurrency *string  `json:"default_currency,omitempty"`
	DefaultTaxRate  *float64 `json:"default_tax_rate,omitempty"`

	// Configuración de notificaciones
	NotificationEmail    *string `json:"notification_email,omitempty"`
	SendInvoiceEmails    *bool   `json:"send_invoice_emails,omitempty"`
	SendPaymentReminders *bool   `json:"send_payment_reminders,omitempty"`

	// Logo y branding
	LogoURL        *string `json:"logo_url,omitempty"`
	PrimaryColor   *string `json:"primary_color,omitempty"`
	SecondaryColor *string `json:"secondary_color,omitempty"`

	// Metadata
	Timezone   *string `json:"timezone,omitempty"`
	DateFormat *string `json:"date_format,omitempty"`
	TimeFormat *string `json:"time_format,omitempty"`
	Locale     *string `json:"locale,omitempty"`
}
