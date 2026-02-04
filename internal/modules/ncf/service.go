package ncf

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/stanleyh24/vendix/internal/database"
	"github.com/stanleyh24/vendix/internal/logger"
	"github.com/stanleyh24/vendix/internal/modules/tenantconfig"
)

// NCFType represents the type of fiscal document
type NCFType string

const (
	NCFTypeFiscalCredit NCFType = "01" // Crédito Fiscal
	NCFTypeConsumer     NCFType = "02" // Consumidor Final
	NCFTypeDebitNote    NCFType = "03" // Nota de Débito
	NCFTypeCreditNote   NCFType = "04" // Nota de Crédito
	NCFTypeSupplier     NCFType = "11" // Regímenes Especiales
	NCFTypeGov          NCFType = "15" // Gubernamental (Facturas Gubernamentales)
	NCFTypeExport       NCFType = "13" // Exportaciones
	NCFTypeImport       NCFType = "14" // Importaciones
)

// NCFInfo contains information about an NCF
type NCFInfo struct {
	Type      NCFType
	Prefix    string
	Sequence  int
	FullNCF   string
	Formatted string
}

// Service provides NCF generation and validation
type Service struct {
	db         *database.DB
	configRepo *tenantconfig.Repository
}

// NewService creates a new NCF service
func NewService(db *database.DB) *Service {
	return &Service{
		db:         db,
		configRepo: tenantconfig.NewRepository(db),
	}
}

// GenerateNCF generates the next NCF for a given type
func (s *Service) GenerateNCF(ctx context.Context, schema string, ncfType NCFType) (*NCFInfo, error) {
	// Get tenant config
	configService := tenantconfig.NewService(s.db)
	config, err := configService.Get(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant config: %w", err)
	}

	if !config.NCFEnabled {
		return nil, fmt.Errorf("NCF is not enabled for this tenant")
	}

	var prefix *string
	var sequence int
	var endRange int

	// Get prefix, sequence, and range based on type
	switch ncfType {
	case NCFTypeFiscalCredit:
		prefix = config.NCFFiscalCreditPrefix
		sequence = config.NCFFiscalCreditSequence
		endRange = getNCFEndRange(config, "fiscal_credit")
	case NCFTypeConsumer:
		prefix = config.NCFConsumerPrefix
		sequence = config.NCFConsumerSequence
		endRange = getNCFEndRange(config, "consumer")
	case NCFTypeDebitNote:
		prefix = config.NCFDebitNotePrefix
		sequence = config.NCFDebitNoteSequence
		endRange = getNCFEndRange(config, "debit_note")
	case NCFTypeCreditNote:
		prefix = config.NCFCreditNotePrefix
		sequence = config.NCFCreditNoteSequence
		endRange = getNCFEndRange(config, "credit_note")
	case NCFTypeGov:
		prefix = config.NCFGovPrefix
		sequence = config.NCFGovSequence
		endRange = config.NCFGovEndRange
		if endRange == 0 {
			endRange = getNCFEndRange(config, "gov")
		}
	default:
		return nil, fmt.Errorf("unsupported NCF type: %s", ncfType)
	}

	if prefix == nil || *prefix == "" {
		return nil, fmt.Errorf("NCF prefix not configured for type %s", ncfType)
	}

	// Validate range
	if sequence >= endRange {
		return nil, fmt.Errorf("NCF range exhausted for type %s: sequence %d >= end range %d. Please contact DGII for a new range", ncfType, sequence, endRange)
	}

	if sequence >= int(float64(endRange)*0.9) {
		logger.Warn("NCF range approaching limit",
			"type", ncfType,
			"sequence", sequence,
			"end_range", endRange,
			"percentage", float64(sequence)/float64(endRange)*100)
	}

	// Generate full NCF: Prefix + Type + 11 digits (sequence padded)
	fullNCF := s.buildNCF(*prefix, string(ncfType), sequence)
	
	// Increment sequence in database (within a transaction to ensure atomicity)
	newSequence := sequence + 1
	if err := s.incrementNCFSequence(ctx, schema, ncfType, newSequence); err != nil {
		return nil, fmt.Errorf("failed to increment NCF sequence: %w", err)
	}

	return &NCFInfo{
		Type:      ncfType,
		Prefix:    *prefix,
		Sequence:  sequence,
		FullNCF:   fullNCF,
		Formatted: s.formatNCF(fullNCF),
	}, nil
}

// buildNCF constructs the full NCF string
// Format: {PREFIX}{TYPE}{SEQUENCE} (e.g., B010100000000001 for prefix B01, type 01, sequence 1)
func (s *Service) buildNCF(prefix, ncfType string, sequence int) string {
	// Format sequence as 11 digits with leading zeros
	sequenceStr := fmt.Sprintf("%011d", sequence)
	return prefix + ncfType + sequenceStr
}

// formatNCF formats NCF with dashes for readability
// Example: B01-01-00000000001
func (s *Service) formatNCF(ncf string) string {
	if len(ncf) < 4 {
		return ncf
	}
	// Insert dashes: prefix-type-sequence
	prefix := ncf[0:3]
	typ := ncf[3:5]
	seq := ncf[5:]
	return fmt.Sprintf("%s-%s-%s", prefix, typ, seq)
}

// ValidateNCF validates an NCF format
func (s *Service) ValidateNCF(ncf string) error {
	// Remove dashes and spaces
	cleaned := strings.ReplaceAll(strings.ReplaceAll(ncf, "-", ""), " ", "")
	
	// Basic format validation: should be at least 16 characters (prefix + type + sequence)
	if len(cleaned) < 16 {
		return fmt.Errorf("NCF format invalid: too short (expected at least 16 characters)")
	}

	// Validate prefix (3 characters, alphanumeric)
	prefix := cleaned[0:3]
	if matched, _ := regexp.MatchString("^[A-Z0-9]{3}$", prefix); !matched {
		return fmt.Errorf("NCF prefix invalid: must be 3 alphanumeric characters")
	}

	// Validate type (2 digits, 01-15)
	ncfType := cleaned[3:5]
	if matched, _ := regexp.MatchString("^(01|02|03|04|11|12|13|14|15)$", ncfType); !matched {
		return fmt.Errorf("NCF type invalid: must be 01-15")
	}

	// Validate sequence (11 digits)
	sequence := cleaned[5:]
	if matched, _ := regexp.MatchString("^[0-9]{11}$", sequence); !matched {
		return fmt.Errorf("NCF sequence invalid: must be 11 digits")
	}

	return nil
}

// ValidateRNC validates a Dominican RNC (Registro Nacional del Contribuyente)
func ValidateRNC(rnc string) error {
	// Remove dashes and spaces
	cleaned := strings.ReplaceAll(strings.ReplaceAll(rnc, "-", ""), " ", "")

	// RNC must be 9 or 11 digits
	if len(cleaned) != 9 && len(cleaned) != 11 {
		return fmt.Errorf("RNC must be 9 or 11 digits, got %d", len(cleaned))
	}

	// Must be all digits
	if matched, _ := regexp.MatchString("^[0-9]+$", cleaned); !matched {
		return fmt.Errorf("RNC must contain only digits")
	}

	// Validate checksum (algoritmo de validación de RNC)
	if len(cleaned) == 9 || len(cleaned) == 11 {
		if !validateRNCChecksum(cleaned) {
			return fmt.Errorf("RNC checksum validation failed")
		}
	}

	return nil
}

// ValidateCedula validates a Dominican Cédula (personal ID)
func ValidateCedula(cedula string) error {
	// Remove dashes and spaces
	cleaned := strings.ReplaceAll(strings.ReplaceAll(cedula, "-", ""), " ", "")

	// Cédula must be 11 digits
	if len(cleaned) != 11 {
		return fmt.Errorf("Cédula must be 11 digits, got %d", len(cleaned))
	}

	// Must be all digits
	if matched, _ := regexp.MatchString("^[0-9]{11}$", cleaned); !matched {
		return fmt.Errorf("Cédula must contain only 11 digits")
	}

	// Validate checksum (algoritmo de validación de cédula dominicana)
	if !validateCedulaChecksum(cleaned) {
		return fmt.Errorf("Cédula checksum validation failed")
	}

	return nil
}

// validateRNCChecksum validates RNC using DGII algorithm
func validateRNCChecksum(rnc string) bool {
	// Simplified validation - in production, implement full DGII algorithm
	// For now, just check it's numeric and has correct length
	digits := strings.Split(rnc, "")
	
	// Calculate weighted sum
	weights := []int{7, 9, 8, 6, 5, 4, 3, 2}
	sum := 0
	
	for i := 0; i < len(digits)-1 && i < len(weights); i++ {
		digit, _ := strconv.Atoi(digits[i])
		sum += digit * weights[i]
	}
	
	checkDigit, _ := strconv.Atoi(digits[len(digits)-1])
	remainder := sum % 11
	
	// Special case: if remainder is 0 or 1, check digit should be 0 or 1
	if remainder <= 1 {
		return checkDigit == remainder
	}
	
	expectedCheck := 11 - remainder
	return checkDigit == expectedCheck
}

// validateCedulaChecksum validates Dominican cédula using standard algorithm
func validateCedulaChecksum(cedula string) bool {
	if len(cedula) != 11 {
		return false
	}

	digits := make([]int, 11)
	for i, char := range cedula {
		digits[i], _ = strconv.Atoi(string(char))
	}

	// Multiply odd positions (1, 3, 5, 7, 9) by 1, even (2, 4, 6, 8) by 2
	sum := 0
	for i := 0; i < 10; i++ {
		multiplier := 1
		if (i+1)%2 == 0 {
			multiplier = 2
		}
		product := digits[i] * multiplier
		if product > 9 {
			product = product/10 + product%10
		}
		sum += product
	}

	checkDigit := (10 - (sum % 10)) % 10
	return checkDigit == digits[10]
}

// incrementNCFSequence increments the NCF sequence in the database
func (s *Service) incrementNCFSequence(ctx context.Context, schema string, ncfType NCFType, newSequence int) error {
	updateReq := &tenantconfig.UpdateTenantConfigRequest{}

	switch ncfType {
	case NCFTypeFiscalCredit:
		updateReq.NCFFiscalCreditSequence = &newSequence
	case NCFTypeConsumer:
		updateReq.NCFConsumerSequence = &newSequence
	case NCFTypeDebitNote:
		updateReq.NCFDebitNoteSequence = &newSequence
	case NCFTypeCreditNote:
		updateReq.NCFCreditNoteSequence = &newSequence
	case NCFTypeGov:
		updateReq.NCFGovSequence = &newSequence
	default:
		return fmt.Errorf("unsupported NCF type: %s", ncfType)
	}

	configService := tenantconfig.NewService(s.db)
	_, err := configService.Update(ctx, schema, updateReq)
	return err
}


// getNCFEndRange gets the end range for an NCF type  
func getNCFEndRange(config *tenantconfig.TenantConfig, ncfField string) int {
	// Read from database directly since these fields may not be in the model yet
	// For now, default to 999999999 - the actual implementation should read from tenant_config table
	// TODO: Add start/end range fields to TenantConfig model
	return 999999999
}

// GetNCFTypeFromString converts string to NCFType
func GetNCFTypeFromString(s string) (NCFType, error) {
	switch s {
	case "01":
		return NCFTypeFiscalCredit, nil
	case "02":
		return NCFTypeConsumer, nil
	case "03":
		return NCFTypeDebitNote, nil
	case "04":
		return NCFTypeCreditNote, nil
	case "11":
		return NCFTypeSupplier, nil
	case "12":
		return NCFTypeGov, nil // Mantener compatibilidad con tipo 12
	case "13":
		return NCFTypeExport, nil
	case "14":
		return NCFTypeImport, nil
	case "15":
		return NCFTypeGov, nil // Tipo 15 es Gubernamental (más común)
	default:
		return "", fmt.Errorf("invalid NCF type: %s", s)
	}
}

