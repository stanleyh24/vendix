package accounting

import (
	"context"
	"fmt"
	"time"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/logger"
	"vendix/internal/modules/sales"

	"github.com/google/uuid"
)

type Service struct {
	repo         *Repository
	salesRepo    *sales.Repository
	cfg          *config.Config
}

func NewService(db *database.DB, cfg *config.Config) *Service {
	return &Service{
		repo:      NewRepository(db),
		salesRepo: sales.NewRepository(db),
		cfg:       cfg,
	}
}

// Chart of Accounts methods

func (s *Service) CreateAccount(ctx context.Context, schema string, req *CreateAccountRequest) (*ChartOfAccounts, error) {
	// Validate account type
	validTypes := map[string]bool{
		AccountTypeAsset:     true,
		AccountTypeLiability: true,
		AccountTypeEquity:    true,
		AccountTypeRevenue:   true,
		AccountTypeExpense:   true,
	}

	if !validTypes[req.AccountType] {
		return nil, fmt.Errorf("invalid account type: %s", req.AccountType)
	}

	// Check if account code already exists
	existing, err := s.repo.GetAccountByCode(ctx, schema, req.AccountCode)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("account code already exists: %s", req.AccountCode)
	}

	account := &ChartOfAccounts{
		ID:                uuid.New(),
		AccountCode:       req.AccountCode,
		AccountName:       req.AccountName,
		AccountType:       req.AccountType,
		ParentAccountCode: req.ParentAccountCode,
		Description:       req.Description,
		IsActive:          true,
		Balance:           0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.repo.CreateAccount(ctx, schema, account); err != nil {
		logger.Error("Failed to create account", "error", err)
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	logger.Info("Account created", "code", account.AccountCode)
	return account, nil
}

func (s *Service) GetAccountByCode(ctx context.Context, schema string, code string) (*ChartOfAccounts, error) {
	account, err := s.repo.GetAccountByCode(ctx, schema, code)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	return account, nil
}

func (s *Service) ListAccounts(ctx context.Context, schema string, accountType *string) ([]*ChartOfAccounts, error) {
	accounts, err := s.repo.ListAccounts(ctx, schema, accountType)

	if err != nil {
		logger.Error("Failed to list accounts", "error", err)
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	// Return empty slice instead of nil
	if accounts == nil {
		accounts = []*ChartOfAccounts{}
	}

	return accounts, nil
}

func (s *Service) UpdateAccount(ctx context.Context, schema string, code string, req *UpdateAccountRequest) (*ChartOfAccounts, error) {
	// Get existing account
	account, err := s.repo.GetAccountByCode(ctx, schema, code)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	// Update fields
	if req.AccountName != nil {
		account.AccountName = *req.AccountName
	}
	if req.AccountType != nil {
		account.AccountType = *req.AccountType
	}
	if req.ParentAccountCode != nil {
		account.ParentAccountCode = req.ParentAccountCode
	}
	if req.Description != nil {
		account.Description = req.Description
	}
	if req.IsActive != nil {
		account.IsActive = *req.IsActive
	}

	account.UpdatedAt = time.Now()

	if err := s.repo.UpdateAccount(ctx, schema, account); err != nil {
		logger.Error("Failed to update account", "error", err)
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	logger.Info("Account updated", "code", code)
	return account, nil
}

// Journal Entry methods

func (s *Service) CreateJournalEntry(ctx context.Context, schema string, req *CreateJournalEntryRequest) (*JournalEntry, error) {
	// Validate that debits equal credits
	totalDebit := 0.0
	totalCredit := 0.0

	for _, line := range req.Lines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}

	if totalDebit != totalCredit {
		return nil, fmt.Errorf("debits (%.2f) must equal credits (%.2f)", totalDebit, totalCredit)
	}

	// Parse entry date
	entryDate, err := time.Parse("2006-01-02", req.EntryDate)
	if err != nil {
		return nil, fmt.Errorf("invalid entry date format: %w", err)
	}

	// Generate entry number
	entryNumber, err := s.repo.GetNextEntryNumber(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to generate entry number: %w", err)
	}

	// Create journal entry
	entry := &JournalEntry{
		ID:          uuid.New(),
		EntryNumber: entryNumber,
		EntryDate:   entryDate,
		Description: req.Description,
		Reference:   req.Reference,
		Status:      "posted",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create lines
	entry.Lines = make([]JournalEntryLine, len(req.Lines))
	for i, lineReq := range req.Lines {
		// Get account name
		account, err := s.repo.GetAccountByCode(ctx, schema, lineReq.AccountCode)
		if err != nil {
			return nil, fmt.Errorf("account not found: %s", lineReq.AccountCode)
		}

		entry.Lines[i] = JournalEntryLine{
			ID:             uuid.New(),
			JournalEntryID: entry.ID,
			AccountCode:    lineReq.AccountCode,
			AccountName:    account.AccountName,
			Debit:          lineReq.Debit,
			Credit:         lineReq.Credit,
			Description:    lineReq.Description,
			CreatedAt:      time.Now(),
		}
	}

	// Save to database
	if err := s.repo.CreateJournalEntry(ctx, schema, entry); err != nil {
		logger.Error("Failed to create journal entry", "error", err)
		return nil, fmt.Errorf("failed to create journal entry: %w", err)
	}

	logger.Info("Journal entry created", "number", entry.EntryNumber)
	return entry, nil
}

func (s *Service) GetJournalEntry(ctx context.Context, schema string, id string) (*JournalEntry, error) {
	entryID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid entry ID: %w", err)
	}

	entry, err := s.repo.GetJournalEntryByID(ctx, schema, entryID)
	if err != nil {
		return nil, fmt.Errorf("journal entry not found: %w", err)
	}

	return entry, nil
}

func (s *Service) ListJournalEntries(ctx context.Context, schema string, startDate, endDate, status *string) ([]*JournalEntry, error) {
	entries, err := s.repo.ListJournalEntries(ctx, schema, startDate, endDate, status)
	if err != nil {
		logger.Error("Failed to list journal entries", "error", err)
		return nil, fmt.Errorf("failed to list journal entries: %w", err)
	}

	return entries, nil
}

// Reporting methods

func (s *Service) GetGeneralLedger(ctx context.Context, schema string, accountCode string, startDate, endDate *string) ([]LedgerEntry, error) {
	// Verify account exists
	_, err := s.repo.GetAccountByCode(ctx, schema, accountCode)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	ledger, err := s.repo.GetGeneralLedger(ctx, schema, accountCode, startDate, endDate)
	if err != nil {
		logger.Error("Failed to get general ledger", "error", err)
		return nil, fmt.Errorf("failed to get general ledger: %w", err)
	}

	return ledger, nil
}

func (s *Service) GetTrialBalance(ctx context.Context, schema string, asOfDate *string) ([]TrialBalanceEntry, error) {
	balance, err := s.repo.GetTrialBalance(ctx, schema, asOfDate)
	if err != nil {
		logger.Error("Failed to get trial balance", "error", err)
		return nil, fmt.Errorf("failed to get trial balance: %w", err)
	}

	return balance, nil
}

func (s *Service) GetBalanceSheet(ctx context.Context, schema string, asOfDate string) (*FinancialStatementReport, error) {
	balance, err := s.repo.GetTrialBalance(ctx, schema, &asOfDate)
	if err != nil {
		return nil, err
	}

	report := &FinancialStatementReport{
		StatementType: "Balance Sheet",
		PeriodEnd:     asOfDate,
		Sections:      []FinancialStatementSection{},
	}

	// Group by account type
	assets := []FinancialStatementItem{}
	liabilities := []FinancialStatementItem{}
	equity := []FinancialStatementItem{}

	for _, entry := range balance {
		amount := entry.Debit - entry.Credit

		item := FinancialStatementItem{
			AccountCode: entry.AccountCode,
			AccountName: entry.AccountName,
			Amount:      amount,
		}

		switch entry.AccountType {
		case AccountTypeAsset:
			assets = append(assets, item)
		case AccountTypeLiability:
			liabilities = append(liabilities, item)
		case AccountTypeEquity:
			equity = append(equity, item)
		}
	}

	// Add sections
	if len(assets) > 0 {
		assetsTotal := 0.0
		for _, item := range assets {
			assetsTotal += item.Amount
		}
		report.Sections = append(report.Sections, FinancialStatementSection{
			Title:    "Activos",
			Accounts: assets,
			Subtotal: assetsTotal,
		})
	}

	if len(liabilities) > 0 {
		liabilitiesTotal := 0.0
		for _, item := range liabilities {
			liabilitiesTotal += item.Amount
		}
		report.Sections = append(report.Sections, FinancialStatementSection{
			Title:    "Pasivos",
			Accounts: liabilities,
			Subtotal: liabilitiesTotal,
		})
	}

	if len(equity) > 0 {
		equityTotal := 0.0
		for _, item := range equity {
			equityTotal += item.Amount
		}
		report.Sections = append(report.Sections, FinancialStatementSection{
			Title:    "Capital",
			Accounts: equity,
			Subtotal: equityTotal,
		})
	}

	return report, nil
}

func (s *Service) GetIncomeStatement(ctx context.Context, schema string, startDate, endDate string) (*FinancialStatementReport, error) {
	balance, err := s.repo.GetTrialBalance(ctx, schema, &endDate)
	if err != nil {
		return nil, err
	}

	report := &FinancialStatementReport{
		StatementType: "Income Statement",
		PeriodStart:   startDate,
		PeriodEnd:     endDate,
		Sections:      []FinancialStatementSection{},
	}

	// Group by account type
	revenues := []FinancialStatementItem{}
	expenses := []FinancialStatementItem{}

	for _, entry := range balance {
		amount := entry.Credit - entry.Debit // For revenue and expenses

		item := FinancialStatementItem{
			AccountCode: entry.AccountCode,
			AccountName: entry.AccountName,
			Amount:      amount,
		}

		switch entry.AccountType {
		case AccountTypeRevenue:
			revenues = append(revenues, item)
		case AccountTypeExpense:
			expenses = append(expenses, item)
		}
	}

	// Add sections
	revenueTotal := 0.0
	for _, item := range revenues {
		revenueTotal += item.Amount
	}
	report.Sections = append(report.Sections, FinancialStatementSection{
		Title:    "Ingresos",
		Accounts: revenues,
		Subtotal: revenueTotal,
	})

	expenseTotal := 0.0
	for _, item := range expenses {
		expenseTotal += item.Amount
	}
	report.Sections = append(report.Sections, FinancialStatementSection{
		Title:    "Gastos",
		Accounts: expenses,
		Subtotal: expenseTotal,
	})

	report.Total = revenueTotal - expenseTotal

	return report, nil
}

// GenerateDailySalesJournalEntry generates a journal entry for all sales of a specific day
// It groups sales by payment method (cash, card, credit) and creates a single journal entry
func (s *Service) GenerateDailySalesJournalEntry(ctx context.Context, schema string, date string, userID *uuid.UUID) (*JournalEntry, error) {
	// Get daily sales summary
	summary, err := s.salesRepo.GetDailySalesSummary(ctx, schema, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily sales summary: %w", err)
	}

	// Check if there are any sales
	if summary.TotalTransactions == 0 {
		return nil, fmt.Errorf("no sales found for date %s", date)
	}

	// Check if journal entry already exists for this date
	existingEntries, err := s.repo.ListJournalEntries(ctx, schema, &date, &date, stringPtr("posted"))
	if err == nil {
		for _, entry := range existingEntries {
			if entry.Description == fmt.Sprintf("Asiento diario de ventas - %s", date) {
				return nil, fmt.Errorf("journal entry already exists for date %s", date)
			}
		}
	}

	// Generate entry number
	entryNumber, err := s.repo.GetNextEntryNumber(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to generate entry number: %w", err)
	}

	// Parse entry date
	entryDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	// Create journal entry
	entry := &JournalEntry{
		ID:          uuid.New(),
		EntryNumber: entryNumber,
		EntryDate:   entryDate,
		Description: fmt.Sprintf("Asiento diario de ventas - %s", date),
		Reference:   stringPtr(fmt.Sprintf("VENTAS-%s", date)),
		Status:      "posted",
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Lines:       []JournalEntryLine{},
	}

	// Build journal entry lines
	// DEBE (Debit):
	// - Caja General (1111) - ventas en efectivo
	// - Banco - Cuenta Corriente (1112) - ventas con tarjeta y transferencias
	// - Clientes (1121) - ventas a crédito (CXC)

	// HABER (Credit):
	// - Ventas (4110) - subtotal de todas las ventas
	// - ITBIS por Pagar (2121) - total de impuestos
	// - Descuentos en Ventas (4111) - si hay descuentos (por ahora 0)

	// Verify accounts exist
	accountsToCheck := []string{"1111", "1112", "1121", "4110", "2121"}
	for _, code := range accountsToCheck {
		_, err := s.repo.GetAccountByCode(ctx, schema, code)
		if err != nil {
			return nil, fmt.Errorf("account %s not found in chart of accounts", code)
		}
	}

	// Get account names
	cashAccount, _ := s.repo.GetAccountByCode(ctx, schema, "1111")
	bankAccount, _ := s.repo.GetAccountByCode(ctx, schema, "1112")
	arAccount, _ := s.repo.GetAccountByCode(ctx, schema, "1121")
	salesAccount, _ := s.repo.GetAccountByCode(ctx, schema, "4110")
	taxAccount, _ := s.repo.GetAccountByCode(ctx, schema, "2121")

	// DEBIT LINES (DEBE)
	// 1. Caja General - ventas en efectivo (subtotal + tax)
	if summary.CashTotal > 0 {
		entry.Lines = append(entry.Lines, JournalEntryLine{
			ID:             uuid.New(),
			JournalEntryID: entry.ID,
			AccountCode:    "1111",
			AccountName:    cashAccount.AccountName,
			Debit:          summary.CashTotal,
			Credit:         0,
			Description:    stringPtr(fmt.Sprintf("Ventas en efectivo - %d transacciones", summary.TotalTransactions)),
			CreatedAt:      time.Now(),
		})
	}

	// 2. Banco - Cuenta Corriente - ventas con tarjeta y transferencias
	cardAndTransferTotal := summary.CardTotal + summary.TransferTotal
	if cardAndTransferTotal > 0 {
		entry.Lines = append(entry.Lines, JournalEntryLine{
			ID:             uuid.New(),
			JournalEntryID: entry.ID,
			AccountCode:    "1112",
			AccountName:    bankAccount.AccountName,
			Debit:          cardAndTransferTotal,
			Credit:         0,
			Description:    stringPtr(fmt.Sprintf("Ventas con tarjeta y transferencias")),
			CreatedAt:      time.Now(),
		})
	}

	// 3. Clientes (CXC) - ventas a crédito
	if summary.CreditTotal > 0 {
		entry.Lines = append(entry.Lines, JournalEntryLine{
			ID:             uuid.New(),
			JournalEntryID: entry.ID,
			AccountCode:    "1121",
			AccountName:    arAccount.AccountName,
			Debit:          summary.CreditTotal,
			Credit:         0,
			Description:    stringPtr(fmt.Sprintf("Ventas a crédito (CXC)")),
			CreatedAt:      time.Now(),
		})
	}

	// CREDIT LINES (HABER)
	// 4. Ventas - subtotal de todas las ventas
	if summary.TotalSubtotal > 0 {
		entry.Lines = append(entry.Lines, JournalEntryLine{
			ID:             uuid.New(),
			JournalEntryID: entry.ID,
			AccountCode:    "4110",
			AccountName:    salesAccount.AccountName,
			Debit:          0,
			Credit:         summary.TotalSubtotal,
			Description:    stringPtr(fmt.Sprintf("Ventas del día - %d transacciones", summary.TotalTransactions)),
			CreatedAt:      time.Now(),
		})
	}

	// 5. ITBIS por Pagar - total de impuestos
	if summary.TotalTax > 0 {
		entry.Lines = append(entry.Lines, JournalEntryLine{
			ID:             uuid.New(),
			JournalEntryID: entry.ID,
			AccountCode:    "2121",
			AccountName:    taxAccount.AccountName,
			Debit:          0,
			Credit:         summary.TotalTax,
			Description:    stringPtr("ITBIS por Pagar"),
			CreatedAt:      time.Now(),
		})
	}

	// 6. Impuesto Selectivo (si aplica) - por ahora 0, pero estructura lista
	if summary.TotalSelective > 0 {
		// Verificar si existe la cuenta 2123 para Impuesto Selectivo
		selectiveAccount, err := s.repo.GetAccountByCode(ctx, schema, "2123")
		if err == nil {
			entry.Lines = append(entry.Lines, JournalEntryLine{
				ID:             uuid.New(),
				JournalEntryID: entry.ID,
				AccountCode:    "2123",
				AccountName:    selectiveAccount.AccountName,
				Debit:          0,
				Credit:         summary.TotalSelective,
				Description:    stringPtr("Impuesto Selectivo por Pagar"),
				CreatedAt:      time.Now(),
			})
		}
	}

	// Validate that debits equal credits
	totalDebit := 0.0
	totalCredit := 0.0
	for _, line := range entry.Lines {
		totalDebit += line.Debit
		totalCredit += line.Credit
	}

	if totalDebit != totalCredit {
		return nil, fmt.Errorf("debits (%.2f) must equal credits (%.2f)", totalDebit, totalCredit)
	}

	// Save journal entry
	if err := s.repo.CreateJournalEntry(ctx, schema, entry); err != nil {
		logger.Error("Failed to create daily sales journal entry", "error", err)
		return nil, fmt.Errorf("failed to create journal entry: %w", err)
	}

	logger.Info("Daily sales journal entry created", "entry_number", entry.EntryNumber, "date", date)
	return entry, nil
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
