package cashregisters

import (
	"context"
	"fmt"
	"time"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/logger"
	"vendix/internal/modules/accounting"

	"github.com/google/uuid"
)

type Service struct {
	repo         *Repository
	accountingSvc *accounting.Service
	cfg          *config.Config
}

func NewService(db *database.DB, cfg *config.Config) *Service {
	return &Service{
		repo:         NewRepository(db),
		accountingSvc: accounting.NewService(db, cfg),
		cfg:          cfg,
	}
}

// CreateCashRegister creates a new cash register
func (s *Service) CreateCashRegister(ctx context.Context, schema string, userID uuid.UUID, req *CreateCashRegisterRequest) (*CashRegister, error) {
	register := &CashRegister{
		ID:            uuid.New(),
		Name:          req.Name,
		Location:      req.Location,
		IsActive:      true,
		InitialBalance: req.InitialBalance,
		CreatedBy:     &userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.CreateCashRegister(ctx, schema, register); err != nil {
		logger.Error("Failed to create cash register", "error", err)
		return nil, fmt.Errorf("failed to create cash register: %w", err)
	}

	logger.Info("Cash register created", "id", register.ID)
	return register, nil
}

// GetCashRegisterByID gets a cash register by ID
func (s *Service) GetCashRegisterByID(ctx context.Context, schema string, id uuid.UUID) (*CashRegister, error) {
	register, err := s.repo.GetCashRegisterByID(ctx, schema, id)
	if err != nil {
		return nil, fmt.Errorf("cash register not found: %w", err)
	}
	return register, nil
}

// ListCashRegisters lists all cash registers
func (s *Service) ListCashRegisters(ctx context.Context, schema string, activeOnly bool) ([]*CashRegister, error) {
	registers, err := s.repo.ListCashRegisters(ctx, schema, activeOnly)
	if err != nil {
		logger.Error("Failed to list cash registers", "error", err)
		return nil, fmt.Errorf("failed to list cash registers: %w", err)
	}
	return registers, nil
}

// OpenSession opens a new cash register session
func (s *Service) OpenSession(ctx context.Context, schema string, userID uuid.UUID, req *OpenSessionRequest) (*CashRegisterSession, error) {
	// Verify cash register exists
	_, err := s.repo.GetCashRegisterByID(ctx, schema, req.CashRegisterID)
	if err != nil {
		return nil, fmt.Errorf("cash register not found: %w", err)
	}

	// Check if there's already an open session
	openSession, err := s.repo.GetOpenSession(ctx, schema, req.CashRegisterID)
	if err == nil && openSession != nil {
		return nil, fmt.Errorf("ya existe una sesión abierta para esta caja")
	}

	session := &CashRegisterSession{
		ID:              uuid.New(),
		CashRegisterID:  req.CashRegisterID,
		OpenedBy:        userID,
		OpeningBalance:  req.OpeningBalance,
		ExpectedCash:    0,
		ExpectedCard:    0,
		ExpectedTransfer: 0,
		Status:          "open",
		Notes:           req.Notes,
		OpenedAt:        time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.CreateSession(ctx, schema, session); err != nil {
		logger.Error("Failed to create session", "error", err)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	logger.Info("Cash register session opened", "session_id", session.ID)
	return session, nil
}

// CloseSession closes a cash register session and generates journal entry automatically
func (s *Service) CloseSession(ctx context.Context, schema string, sessionID uuid.UUID, userID uuid.UUID, req *CloseSessionRequest) (*SessionSummary, error) {
	// Get session
	session, err := s.repo.GetSessionByID(ctx, schema, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if session.Status != "open" {
		return nil, fmt.Errorf("la sesión ya está cerrada")
	}

	// Calculate sales from opening time to now
	salesSummary, err := s.repo.GetSalesSummaryByDate(ctx, schema, session.OpenedAt, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get sales summary: %w", err)
	}

	// Calculate expected cash (sales in cash + opening balance)
	expectedCashTotal := salesSummary["cash"] + session.OpeningBalance
	difference := req.CountedCash - expectedCashTotal

	// Update session
	now := time.Now()
	session.CountedCash = &req.CountedCash
	session.ExpectedCash = salesSummary["cash"] + session.OpeningBalance // Total expected including opening balance
	session.ExpectedCard = salesSummary["card"]
	session.ExpectedTransfer = salesSummary["transfer"]
	session.Difference = &difference
	session.Status = "closed"
	session.ClosedBy = &userID
	session.ClosedAt = &now
	session.Notes = req.Notes
	session.UpdatedAt = now

	// Save count details
	if len(req.CountDetails) > 0 {
		for _, detail := range req.CountDetails {
			detail.ID = uuid.New()
			detail.SessionID = sessionID
			detail.Subtotal = detail.Denomination * float64(detail.Quantity)
			if err := s.repo.CreateCountDetail(ctx, schema, &detail); err != nil {
				return nil, fmt.Errorf("failed to save count detail: %w", err)
			}
		}
	}

	// Generate journal entry automatically for the day
	dateStr := session.OpenedAt.Format("2006-01-02")
	journalEntry, err := s.accountingSvc.GenerateDailySalesJournalEntry(ctx, schema, dateStr, &userID)
	if err != nil {
		// Log error but don't fail the close operation
		logger.Error("Failed to generate journal entry on cash close", "error", err, "date", dateStr)
		// Continue with closing the session even if journal entry fails
	} else {
		// Link journal entry to session
		session.JournalEntryID = &journalEntry.ID
		logger.Info("Journal entry generated on cash close", "entry_id", journalEntry.ID, "entry_number", journalEntry.EntryNumber)
	}

	// Update session with journal entry ID
	if err := s.repo.UpdateSession(ctx, schema, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	// Prepare summary
	summary := &SessionSummary{
		Session: session,
		TotalSales: salesSummary["cash"] + salesSummary["card"] + salesSummary["transfer"],
		SalesByPaymentType: salesSummary,
		TotalTransactions: int(salesSummary["count"]),
	}

	logger.Info("Cash register session closed", "session_id", sessionID, "difference", difference)
	return summary, nil
}

// GetSessionSummary gets the summary of an open session
func (s *Service) GetSessionSummary(ctx context.Context, schema string, sessionID uuid.UUID) (*SessionSummary, error) {
	session, err := s.repo.GetSessionByID(ctx, schema, sessionID)
	if err != nil {
		logger.Error("Failed to get session by ID", "error", err, "session_id", sessionID, "schema", schema)
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Get sales from opening time to now
	salesSummary, err := s.repo.GetSalesSummaryByDate(ctx, schema, session.OpenedAt, time.Now())
	if err != nil {
		logger.Error("Failed to get sales summary by date", "error", err, "schema", schema, "start", session.OpenedAt, "end", time.Now())
		return nil, fmt.Errorf("failed to get sales summary: %w", err)
	}

	// Ensure all payment types are initialized (in case query returns no rows)
	if salesSummary == nil {
		salesSummary = make(map[string]float64)
	}
	if _, ok := salesSummary["cash"]; !ok {
		salesSummary["cash"] = 0
	}
	if _, ok := salesSummary["card"]; !ok {
		salesSummary["card"] = 0
	}
	if _, ok := salesSummary["transfer"]; !ok {
		salesSummary["transfer"] = 0
	}
	if _, ok := salesSummary["count"]; !ok {
		salesSummary["count"] = 0
	}

	// Calculate expected amounts (cash includes opening balance + sales)
	session.ExpectedCash = salesSummary["cash"] + session.OpeningBalance
	session.ExpectedCard = salesSummary["card"]
	session.ExpectedTransfer = salesSummary["transfer"]

	summary := &SessionSummary{
		Session: session,
		TotalSales: salesSummary["cash"] + salesSummary["card"] + salesSummary["transfer"],
		SalesByPaymentType: salesSummary,
		TotalTransactions: int(salesSummary["count"]),
	}

	logger.Info("Session summary generated", "session_id", sessionID, "total_sales", summary.TotalSales, "cash_sales", salesSummary["cash"], "card_sales", salesSummary["card"], "transfer_sales", salesSummary["transfer"])
	return summary, nil
}

// GetOpenSession gets the open session for a cash register
func (s *Service) GetOpenSession(ctx context.Context, schema string, cashRegisterID uuid.UUID) (*CashRegisterSession, error) {
	session, err := s.repo.GetOpenSession(ctx, schema, cashRegisterID)
	if err != nil {
		return nil, fmt.Errorf("no open session found: %w", err)
	}
	return session, nil
}

// ListSessions lists sessions with filters
func (s *Service) ListSessions(ctx context.Context, schema string, cashRegisterID *uuid.UUID, status *string, startDate, endDate *time.Time) ([]*CashRegisterSession, error) {
	sessions, err := s.repo.ListSessions(ctx, schema, cashRegisterID, status, startDate, endDate)
	if err != nil {
		logger.Error("Failed to list sessions", "error", err)
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	return sessions, nil
}

