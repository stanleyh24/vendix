package tenantconfig

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"vendix/internal/database"
	"vendix/internal/logger"
)

type Service struct {
	repo *Repository
}

// NewService creates a new tenant config service
func NewService(db *database.DB) *Service {
	return &Service{
		repo: NewRepository(db),
	}
}

// Get retrieves the tenant configuration, creating it if it doesn't exist
func (s *Service) Get(ctx context.Context, schema string) (*TenantConfig, error) {
	config, err := s.repo.Get(ctx, schema)
	if err != nil {
		// If no config exists, create a default one
		if errors.Is(err, sql.ErrNoRows) {
			logger.Info("No tenant config found, creating default configuration")
			config, err = s.repo.Create(ctx, schema)
			if err != nil {
				logger.Error("Failed to create default tenant config", "error", err)
				return nil, fmt.Errorf("failed to create tenant configuration: %w", err)
			}
			return config, nil
		}
		logger.Error("Failed to get tenant config", "error", err)
		return nil, fmt.Errorf("failed to get tenant configuration: %w", err)
	}

	return config, nil
}

// Update updates the tenant configuration
func (s *Service) Update(ctx context.Context, schema string, req *UpdateTenantConfigRequest) (*TenantConfig, error) {
	// Get current configuration, creating if it doesn't exist
	config, err := s.Get(ctx, schema)
	if err != nil {
		logger.Error("Failed to get current config for update", "error", err)
		return nil, fmt.Errorf("failed to get current configuration: %w", err)
	}

	// Update fields if provided
	if req.CompanyLegalName != nil {
		config.CompanyLegalName = *req.CompanyLegalName
	}
	if req.CompanyTradeName != nil {
		config.CompanyTradeName = req.CompanyTradeName
	}
	if req.CompanyTaxID != nil {
		config.CompanyTaxID = req.CompanyTaxID
	}
	if req.CompanyEmail != nil {
		config.CompanyEmail = req.CompanyEmail
	}
	if req.CompanyPhone != nil {
		config.CompanyPhone = req.CompanyPhone
	}
	if req.CompanyWebsite != nil {
		config.CompanyWebsite = req.CompanyWebsite
	}

	// Address fields
	if req.AddressLine1 != nil {
		config.AddressLine1 = req.AddressLine1
	}
	if req.AddressLine2 != nil {
		config.AddressLine2 = req.AddressLine2
	}
	if req.City != nil {
		config.City = req.City
	}
	if req.StateProvince != nil {
		config.StateProvince = req.StateProvince
	}
	if req.PostalCode != nil {
		config.PostalCode = req.PostalCode
	}
	if req.Country != nil {
		config.Country = *req.Country
	}

	// Fiscal fields
	if req.TaxRegime != nil {
		config.TaxRegime = req.TaxRegime
	}
	if req.DGIIIRNC != nil {
		config.DGIIIRNC = req.DGIIIRNC
	}
	if req.DGIIUser != nil {
		config.DGIIUser = req.DGIIUser
	}
	if req.DGIIAPIKey != nil {
		config.DGIIAPIKey = req.DGIIAPIKey
	}
	if req.DGIIEnvironment != nil {
		config.DGIIEnvironment = *req.DGIIEnvironment
	}
	if req.NCFEnabled != nil {
		config.NCFEnabled = *req.NCFEnabled
	}

	// Invoice fields
	if req.InvoicePrefix != nil {
		config.InvoicePrefix = *req.InvoicePrefix
	}
	if req.InvoiceNextNumber != nil {
		config.InvoiceNextNumber = *req.InvoiceNextNumber
	}
	if req.InvoiceTerms != nil {
		config.InvoiceTerms = req.InvoiceTerms
	}
	if req.InvoiceFooter != nil {
		config.InvoiceFooter = req.InvoiceFooter
	}
	if req.PaymentTermsDays != nil {
		config.PaymentTermsDays = *req.PaymentTermsDays
	}

	// NCF fields
	if req.NCFFiscalCreditPrefix != nil {
		config.NCFFiscalCreditPrefix = req.NCFFiscalCreditPrefix
	}
	if req.NCFFiscalCreditSequence != nil {
		config.NCFFiscalCreditSequence = *req.NCFFiscalCreditSequence
	}
	if req.NCFConsumerPrefix != nil {
		config.NCFConsumerPrefix = req.NCFConsumerPrefix
	}
	if req.NCFConsumerSequence != nil {
		config.NCFConsumerSequence = *req.NCFConsumerSequence
	}
	if req.NCFDebitNotePrefix != nil {
		config.NCFDebitNotePrefix = req.NCFDebitNotePrefix
	}
	if req.NCFDebitNoteSequence != nil {
		config.NCFDebitNoteSequence = *req.NCFDebitNoteSequence
	}
	if req.NCFCreditNotePrefix != nil {
		config.NCFCreditNotePrefix = req.NCFCreditNotePrefix
	}
	if req.NCFCreditNoteSequence != nil {
		config.NCFCreditNoteSequence = *req.NCFCreditNoteSequence
	}
	if req.NCFGovPrefix != nil {
		config.NCFGovPrefix = req.NCFGovPrefix
	}
	if req.NCFGovSequence != nil {
		config.NCFGovSequence = *req.NCFGovSequence
	}

	// Currency fields
	if req.DefaultCurrency != nil {
		config.DefaultCurrency = *req.DefaultCurrency
	}
	if req.DefaultTaxRate != nil {
		config.DefaultTaxRate = *req.DefaultTaxRate
	}

	// Notification fields
	if req.NotificationEmail != nil {
		config.NotificationEmail = req.NotificationEmail
	}
	if req.SendInvoiceEmails != nil {
		config.SendInvoiceEmails = *req.SendInvoiceEmails
	}
	if req.SendPaymentReminders != nil {
		config.SendPaymentReminders = *req.SendPaymentReminders
	}

	// Branding fields
	if req.LogoURL != nil {
		config.LogoURL = req.LogoURL
	}
	if req.PrimaryColor != nil {
		config.PrimaryColor = *req.PrimaryColor
	}
	if req.SecondaryColor != nil {
		config.SecondaryColor = *req.SecondaryColor
	}

	// Metadata fields
	if req.Timezone != nil {
		config.Timezone = *req.Timezone
	}
	if req.DateFormat != nil {
		config.DateFormat = *req.DateFormat
	}
	if req.TimeFormat != nil {
		config.TimeFormat = *req.TimeFormat
	}
	if req.Locale != nil {
		config.Locale = *req.Locale
	}

	config.UpdatedAt = time.Now()

	// Update in database
	if err := s.repo.Update(ctx, schema, config); err != nil {
		logger.Error("Failed to update tenant config", "error", err)
		return nil, fmt.Errorf("failed to update tenant configuration: %w", err)
	}

	logger.Info("Tenant config updated successfully")

	return config, nil
}
