package dgii

import (
	"context"
	"fmt"
	"time"

	"github.com/stanleyh24/vendix/internal/logger"

	"github.com/google/uuid"
)

// MockAdapter is a mock implementation of DGIIAdapter for testing
type MockAdapter struct {
	enabled bool
}

// NewMockAdapter creates a new mock DGII adapter
func NewMockAdapter(enabled bool) DGIIAdapter {
	return &MockAdapter{enabled: enabled}
}

// SignDocument mocks the digital signature of a document
func (m *MockAdapter) SignDocument(ctx context.Context, document *ElectronicDocument) (*SignedDocument, error) {
	if !m.enabled {
		return nil, fmt.Errorf("DGII integration is disabled")
	}

	logger.Info("Mock: Signing document",
		"document_type", document.DocumentType,
		"ncf", document.NCF,
	)

	// Simulate signing delay
	time.Sleep(100 * time.Millisecond)

	// Generate a more realistic XML structure (simplified eCF format)
	signedXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<eCF>
	<Encabezado>
		<IDVersion>1.0</IDVersion>
		<eNCF>%s</eNCF>
		<TipoCF>%s</TipoCF>
		<FechaEmision>%s</FechaEmision>
	</Encabezado>
	<DGII>
		<TrackingCode>PENDING</TrackingCode>
	</DGII>
</eCF>`, document.NCF, document.DocumentType, document.IssueDate)

	signed := &SignedDocument{
		OriginalDocument: document,
		Signature:        fmt.Sprintf("MOCK_SIGNATURE_%s", uuid.New().String()),
		Certificate:      "MOCK_CERTIFICATE",
		SignedAt:         time.Now().Format(time.RFC3339),
		SignedXML:        signedXML,
	}

	logger.Info("Mock: Document signed successfully", "ncf", document.NCF)

	return signed, nil
}

// SendDocument mocks sending a signed document to DGII
func (m *MockAdapter) SendDocument(ctx context.Context, signed *SignedDocument) (*DGIIResponse, error) {
	if !m.enabled {
		return nil, fmt.Errorf("DGII integration is disabled")
	}

	logger.Info("Mock: Sending document to DGII",
		"ncf", signed.OriginalDocument.NCF,
	)

	// Simulate network delay
	time.Sleep(200 * time.Millisecond)

	response := &DGIIResponse{
		Success:      true,
		TrackingCode: fmt.Sprintf("DGII-TRACK-%s", uuid.New().String()[:8]),
		Message:      "Documento recibido y procesado exitosamente (MOCK)",
		ReceivedAt:   time.Now().Format(time.RFC3339),
		Metadata: map[string]interface{}{
			"mock": true,
			"ncf":  signed.OriginalDocument.NCF,
		},
	}

	logger.Info("Mock: Document sent successfully",
		"tracking_code", response.TrackingCode,
	)

	return response, nil
}

// CancelDocument mocks cancelling a document in DGII
func (m *MockAdapter) CancelDocument(ctx context.Context, documentID string, reason string) error {
	if !m.enabled {
		return fmt.Errorf("DGII integration is disabled")
	}

	logger.Info("Mock: Cancelling document",
		"document_id", documentID,
		"reason", reason,
	)

	// Simulate network delay
	time.Sleep(150 * time.Millisecond)

	logger.Info("Mock: Document cancelled successfully",
		"document_id", documentID,
	)

	return nil
}

// GetDocumentStatus mocks retrieving document status from DGII
func (m *MockAdapter) GetDocumentStatus(ctx context.Context, documentID string) (*DocumentStatus, error) {
	if !m.enabled {
		return nil, fmt.Errorf("DGII integration is disabled")
	}

	logger.Info("Mock: Getting document status",
		"document_id", documentID,
	)

	// Simulate network delay
	time.Sleep(100 * time.Millisecond)

	status := &DocumentStatus{
		DocumentID:    documentID,
		Status:        "approved",
		TrackingCode:  fmt.Sprintf("DGII-TRACK-%s", uuid.New().String()[:8]),
		Message:       "Documento aprobado por DGII (MOCK)",
		LastUpdatedAt: time.Now().Format(time.RFC3339),
	}

	logger.Info("Mock: Document status retrieved",
		"status", status.Status,
	)

	return status, nil
}
