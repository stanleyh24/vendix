package dgii

import (
	"context"
	"fmt"

	"vendix/internal/config"
	"vendix/internal/logger"
)

// Service provides DGII compliance operations
type Service struct {
	adapter DGIIAdapter
	cfg     *config.Config
}

// NewService creates a new DGII service
func NewService(cfg *config.Config) *Service {
	var adapter DGIIAdapter

	if cfg.DGIIEnabled {
		// In production, this would initialize the real DGII adapter
		// For now, we use the mock adapter
		logger.Info("Initializing DGII service with mock adapter", "enabled", cfg.DGIIEnabled)
		adapter = NewMockAdapter(true)
	} else {
		logger.Info("DGII service disabled")
		adapter = NewMockAdapter(false)
	}

	return &Service{
		adapter: adapter,
		cfg:     cfg,
	}
}

// ProcessInvoiceResponse includes both the DGII response and the signed document
type ProcessInvoiceResponse struct {
	Response     *DGIIResponse
	SignedDoc    *SignedDocument
}

// ProcessInvoice signs and sends an invoice to DGII
func (s *Service) ProcessInvoice(ctx context.Context, document *ElectronicDocument) (*DGIIResponse, error) {
	// Sign document
	signed, err := s.adapter.SignDocument(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("failed to sign document: %w", err)
	}

	// Send to DGII
	response, err := s.adapter.SendDocument(ctx, signed)
	if err != nil {
		return nil, fmt.Errorf("failed to send document to DGII: %w", err)
	}

	return response, nil
}

// ProcessInvoiceWithSignedDoc signs and sends an invoice to DGII, returning both response and signed document
func (s *Service) ProcessInvoiceWithSignedDoc(ctx context.Context, document *ElectronicDocument) (*ProcessInvoiceResponse, error) {
	// Sign document
	signed, err := s.adapter.SignDocument(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("failed to sign document: %w", err)
	}

	// Send to DGII
	response, err := s.adapter.SendDocument(ctx, signed)
	if err != nil {
		return nil, fmt.Errorf("failed to send document to DGII: %w", err)
	}

	return &ProcessInvoiceResponse{
		Response:  response,
		SignedDoc: signed,
	}, nil
}

// CancelInvoice cancels an invoice in DGII
func (s *Service) CancelInvoice(ctx context.Context, documentID string, reason string) error {
	return s.adapter.CancelDocument(ctx, documentID, reason)
}

// GetStatus retrieves the status of a document from DGII
func (s *Service) GetStatus(ctx context.Context, documentID string) (*DocumentStatus, error) {
	return s.adapter.GetDocumentStatus(ctx, documentID)
}
