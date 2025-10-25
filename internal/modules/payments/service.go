package payments

import (
	"context"
	"fmt"
	"time"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/logger"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
	cfg  *config.Config
}

func NewService(db *database.DB, cfg *config.Config) *Service {
	return &Service{
		repo: NewRepository(db),
		cfg:  cfg,
	}
}

func (s *Service) Create(ctx context.Context, schema string, req *CreatePaymentRequest) (*Payment, error) {
	// Parse date
	paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
	if err != nil {
		return nil, fmt.Errorf("invalid payment date: %w", err)
	}

	// Generate payment number
	prefix := "PAY-"
	if req.Supplier != "" {
		prefix = "EXP-" // Expense prefix
	}

	paymentNumber, err := s.repo.GetNextPaymentNumber(ctx, schema, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to generate payment number: %w", err)
	}

	currency := req.Currency
	if currency == "" {
		currency = "DOP"
	}

	payment := &Payment{
		ID:            uuid.New(),
		PaymentNumber: paymentNumber,
		CustomerID:    req.CustomerID,
		PaymentMethod: req.PaymentMethod,
		PaymentDate:   paymentDate,
		Amount:        req.Amount,
		Currency:      currency,
		Reference:     req.Reference,
		Notes:         req.Notes,
		Status:        "completed",
		Category:      req.Category,
		Supplier:      req.Supplier,
		Description:   req.Description,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.Create(ctx, schema, payment); err != nil {
		logger.Error("Failed to create payment", "error", err)
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	logger.Info("Payment created", "number", payment.PaymentNumber)
	return payment, nil
}

func (s *Service) GetByID(ctx context.Context, schema string, id string) (*Payment, error) {
	paymentID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid payment ID: %w", err)
	}

	payment, err := s.repo.GetByID(ctx, schema, paymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	return payment, nil
}

func (s *Service) List(ctx context.Context, schema string) ([]*Payment, error) {
	payments, err := s.repo.List(ctx, schema)
	if err != nil {
		logger.Error("Failed to list payments", "error", err)
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}

	return payments, nil
}

func (s *Service) Update(ctx context.Context, schema string, id string, req *UpdatePaymentRequest) (*Payment, error) {
	paymentID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid payment ID: %w", err)
	}

	payment, err := s.repo.GetByID(ctx, schema, paymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	if req.Status != nil {
		payment.Status = *req.Status
	}
	if req.Notes != nil {
		payment.Notes = req.Notes
	}
	if req.Reference != nil {
		payment.Reference = req.Reference
	}

	payment.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, schema, payment); err != nil {
		logger.Error("Failed to update payment", "error", err)
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	logger.Info("Payment updated", "id", id)
	return payment, nil
}

func (s *Service) Delete(ctx context.Context, schema string, id string) error {
	paymentID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid payment ID: %w", err)
	}

	if err := s.repo.Delete(ctx, schema, paymentID); err != nil {
		logger.Error("Failed to delete payment", "error", err)
		return fmt.Errorf("failed to delete payment: %w", err)
	}

	logger.Info("Payment deleted", "id", id)
	return nil
}
