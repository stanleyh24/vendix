package accounting

import (
	"context"
	"fmt"
	"time"

	"github.com/stanleyh24/vendix/internal/config"
	"github.com/stanleyh24/vendix/internal/database"
	"github.com/stanleyh24/vendix/internal/logger"
	"github.com/stanleyh24/vendix/internal/modules/sales"

	"github.com/google/uuid"
)

type Service struct {
	repo      *Repository
	salesRepo *sales.Repository
	cfg       *config.Config
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

func (s *Service) DeleteAccount(ctx context.Context, schema string, code string) error {
	// Verificar que la cuenta existe
	account, err := s.repo.GetAccountByCode(ctx, schema, code)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	
	// Verificar si la cuenta tiene asientos contables
	hasEntries, err := s.repo.AccountHasJournalEntries(ctx, schema, code)
	if err != nil {
		logger.Error("Failed to check journal entries", "error", err)
		return fmt.Errorf("failed to check journal entries: %w", err)
	}
	
	if hasEntries {
		return fmt.Errorf("cannot delete account %s: account has journal entries. Only accounts without journal entries can be deleted", code)
	}
	
	// Verificar si hay cuentas hijas (subcuentas)
	childAccounts, err := s.repo.ListAccounts(ctx, schema, nil)
	if err == nil {
		for _, child := range childAccounts {
			if child.ParentAccountCode != nil && *child.ParentAccountCode == code {
				return fmt.Errorf("cannot delete account %s: account has child accounts. Please delete or reassign child accounts first", code)
			}
		}
	}
	
	// Eliminar la cuenta
	if err := s.repo.DeleteAccount(ctx, schema, code); err != nil {
		logger.Error("Failed to delete account", "error", err, "code", code)
		return fmt.Errorf("failed to delete account: %w", err)
	}
	
	logger.Info("Account deleted", "code", code, "name", account.AccountName)
	return nil
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

	// Build journal entry lines usando cuentas configuradas
	// DEBE (Debit):
	// - Caja General - ventas en efectivo
	// - Banco - Cuenta Corriente - ventas con tarjeta y transferencias
	// - Clientes - ventas a crédito (CXC)

	// HABER (Credit):
	// - Ventas - subtotal de todas las ventas
	// - ITBIS por Pagar - total de impuestos
	// - Descuentos en Ventas - si hay descuentos (por ahora 0)

	// Obtener cuentas configuradas (con fallback a valores por defecto)
	cashAccount, err := s.GetAccountForTransaction(ctx, schema, TransactionTypeSalesCash, "1111", "Caja General")
	if err != nil {
		return nil, err
	}
	
	bankAccount, err := s.GetAccountForTransaction(ctx, schema, TransactionTypeSalesCard, "1112", "Banco - Cuenta Corriente")
	if err != nil {
		return nil, err
	}
	
	arAccount, err := s.GetAccountForTransaction(ctx, schema, TransactionTypeSalesCredit, "1121", "Clientes")
	if err != nil {
		return nil, err
	}
	
	salesAccount, err := s.GetAccountForTransaction(ctx, schema, TransactionTypeSalesRevenue, "4110", "Ventas")
	if err != nil {
		return nil, err
	}
	
	taxAccount, err := s.GetAccountForTransaction(ctx, schema, TransactionTypeTaxPayable, "2121", "ITBIS por Pagar")
	if err != nil {
		return nil, err
	}

	// DEBIT LINES (DEBE) - usando cuentas configuradas
	// 1. Caja General - ventas en efectivo (subtotal + tax)
	if summary.CashTotal > 0 {
		entry.Lines = append(entry.Lines, JournalEntryLine{
			ID:             uuid.New(),
			JournalEntryID: entry.ID,
			AccountCode:    cashAccount.AccountCode,
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
			AccountCode:    bankAccount.AccountCode,
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
			AccountCode:    arAccount.AccountCode,
			AccountName:    arAccount.AccountName,
			Debit:          summary.CreditTotal,
			Credit:         0,
			Description:    stringPtr(fmt.Sprintf("Ventas a crédito (CXC)")),
			CreatedAt:      time.Now(),
		})
	}

	// CREDIT LINES (HABER) - usando cuentas configuradas
	// 4. Ventas - subtotal de todas las ventas
	if summary.TotalSubtotal > 0 {
		entry.Lines = append(entry.Lines, JournalEntryLine{
			ID:             uuid.New(),
			JournalEntryID: entry.ID,
			AccountCode:    salesAccount.AccountCode,
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
			AccountCode:    taxAccount.AccountCode,
			AccountName:    taxAccount.AccountName,
			Debit:          0,
			Credit:         summary.TotalTax,
			Description:    stringPtr("ITBIS por Pagar"),
			CreatedAt:      time.Now(),
		})
	}

	// 6. Impuesto Selectivo (si aplica)
	if summary.TotalSelective > 0 {
		selectiveAccount, err := s.GetAccountForTransaction(ctx, schema, TransactionTypeTaxSelective, "2123", "Impuesto Selectivo por Pagar")
		if err == nil {
			entry.Lines = append(entry.Lines, JournalEntryLine{
				ID:             uuid.New(),
				JournalEntryID: entry.ID,
				AccountCode:    selectiveAccount.AccountCode,
				AccountName:    selectiveAccount.AccountName,
				Debit:          0,
				Credit:         summary.TotalSelective,
				Description:    stringPtr("Impuesto Selectivo por Pagar"),
				CreatedAt:      time.Now(),
			})
		}
	}

	// 7. COSTO DE VENTAS (COGS) - Registrar el costo de los productos vendidos
	// Solo se registra si hay costos definidos en los productos
	if summary.TotalCost > 0 {
		// DEBE: Costo de Ventas (5100) - aumenta el gasto
		cogsAccount, err := s.GetAccountForTransaction(ctx, schema, TransactionTypePurchaseExpense, "5100", "Costo de Ventas")
		if err != nil {
			logger.Warn("Failed to get COGS account, skipping cost entry", "error", err)
		} else {
			entry.Lines = append(entry.Lines, JournalEntryLine{
				ID:             uuid.New(),
				JournalEntryID: entry.ID,
				AccountCode:    cogsAccount.AccountCode,
				AccountName:    cogsAccount.AccountName,
				Debit:          summary.TotalCost,
				Credit:         0,
				Description:    stringPtr(fmt.Sprintf("Costo de productos vendidos - %d transacciones", summary.TotalTransactions)),
				CreatedAt:      time.Now(),
			})

			// HABER: Inventario (1131) - disminuye el inventario
			inventoryAccount, err := s.GetAccountForTransaction(ctx, schema, TransactionTypePurchaseInventory, "1131", "Inventario de Mercancías")
			if err != nil {
				logger.Warn("Failed to get Inventory account, skipping cost entry", "error", err)
			} else {
				entry.Lines = append(entry.Lines, JournalEntryLine{
					ID:             uuid.New(),
					JournalEntryID: entry.ID,
					AccountCode:    inventoryAccount.AccountCode,
					AccountName:    inventoryAccount.AccountName,
					Debit:          0,
					Credit:         summary.TotalCost,
					Description:    stringPtr(fmt.Sprintf("Salida de inventario por ventas - %d transacciones", summary.TotalTransactions)),
					CreatedAt:      time.Now(),
				})
			}
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

// Account Mapping methods

// GetAccountForTransaction obtiene la cuenta configurada para un tipo de transacción
func (s *Service) GetAccountForTransaction(ctx context.Context, schema string, transactionType string, defaultCode string, defaultName string) (*ChartOfAccounts, error) {
	mapping, err := s.repo.GetAccountMappingWithFallback(ctx, schema, transactionType, defaultCode, defaultName)
	if err != nil {
		return nil, fmt.Errorf("failed to get account mapping: %w", err)
	}
	
	account, err := s.repo.GetAccountByCode(ctx, schema, mapping.AccountCode)
	if err != nil {
		return nil, fmt.Errorf("account %s not found for transaction type %s", mapping.AccountCode, transactionType)
	}
	
	return account, nil
}

// ListAccountMappings lista todos los mapeos de cuentas
func (s *Service) ListAccountMappings(ctx context.Context, schema string) ([]*AccountMapping, error) {
	mappings, err := s.repo.ListAccountMappings(ctx, schema)
	if err != nil {
		logger.Error("Failed to list account mappings", "error", err)
		return nil, fmt.Errorf("failed to list account mappings: %w", err)
	}
	
	if mappings == nil {
		mappings = []*AccountMapping{}
	}
	
	return mappings, nil
}

// GetAccountMapping obtiene el mapeo para un tipo de transacción
func (s *Service) GetAccountMapping(ctx context.Context, schema string, transactionType string) (*AccountMapping, error) {
	mapping, err := s.repo.GetAccountMapping(ctx, schema, transactionType)
	if err != nil {
		return nil, fmt.Errorf("account mapping not found: %w", err)
	}
	return mapping, nil
}

// CreateOrUpdateAccountMapping crea o actualiza un mapeo
func (s *Service) CreateOrUpdateAccountMapping(ctx context.Context, schema string, req *CreateAccountMappingRequest) (*AccountMapping, error) {
	// Verificar que la cuenta existe
	account, err := s.repo.GetAccountByCode(ctx, schema, req.AccountCode)
	if err != nil {
		return nil, fmt.Errorf("account %s not found", req.AccountCode)
	}
	
	mapping := &AccountMapping{
		ID:              uuid.New(),
		TransactionType: req.TransactionType,
		AccountCode:     req.AccountCode,
		AccountName:     account.AccountName,
		Description:     req.Description,
		IsActive:        true,
	}
	
	if err := s.repo.CreateAccountMapping(ctx, schema, mapping); err != nil {
		logger.Error("Failed to create account mapping", "error", err)
		return nil, fmt.Errorf("failed to create account mapping: %w", err)
	}
	
	logger.Info("Account mapping created/updated", "transaction_type", req.TransactionType, "account_code", req.AccountCode)
	return mapping, nil
}

// UpdateAccountMapping actualiza un mapeo existente
func (s *Service) UpdateAccountMapping(ctx context.Context, schema string, transactionType string, req *UpdateAccountMappingRequest) (*AccountMapping, error) {
	// Obtener mapeo existente
	existing, err := s.repo.GetAccountMapping(ctx, schema, transactionType)
	if err != nil {
		return nil, fmt.Errorf("account mapping not found: %w", err)
	}
	
	// Actualizar campos
	if req.AccountCode != nil {
		// Verificar que la nueva cuenta existe
		account, err := s.repo.GetAccountByCode(ctx, schema, *req.AccountCode)
		if err != nil {
			return nil, fmt.Errorf("account %s not found", *req.AccountCode)
		}
		existing.AccountCode = *req.AccountCode
		existing.AccountName = account.AccountName
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	
	if err := s.repo.UpdateAccountMapping(ctx, schema, transactionType, existing); err != nil {
		logger.Error("Failed to update account mapping", "error", err)
		return nil, fmt.Errorf("failed to update account mapping: %w", err)
	}
	
	logger.Info("Account mapping updated", "transaction_type", transactionType)
	return existing, nil
}

// GeneratePaymentJournalEntry genera un asiento contable para un pago
// DEBE: Cuenta de pago (Caja/Banco) según método de pago
// HABER: Cuenta de cliente (CXC) o cuenta de gasto según el tipo
func (s *Service) GeneratePaymentJournalEntry(ctx context.Context, schema string, paymentID string, userID *uuid.UUID) (*JournalEntry, error) {
	// Obtener información del pago
	paymentData, err := s.repo.GetPaymentByID(ctx, schema, paymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}
	
	// Extraer valores del mapa
	paymentNumber, _ := paymentData["payment_number"].(string)
	paymentMethod, _ := paymentData["payment_method"].(string)
	amount, _ := paymentData["amount"].(float64)
	
	var customerID *uuid.UUID
	if cid, ok := paymentData["customer_id"].(uuid.UUID); ok {
		customerID = &cid
	} else if cidStr, ok := paymentData["customer_id"].(string); ok && cidStr != "" {
		if parsed, err := uuid.Parse(cidStr); err == nil {
			customerID = &parsed
		}
	}
	
	var paymentDate time.Time
	if pd, ok := paymentData["payment_date"].(time.Time); ok {
		paymentDate = pd
	} else if pdStr, ok := paymentData["payment_date"].(string); ok {
		if parsed, err := time.Parse("2006-01-02", pdStr); err == nil {
			paymentDate = parsed
		} else {
			paymentDate = time.Now()
		}
	} else {
		paymentDate = time.Now()
	}
	
	var reference *string
	if ref, ok := paymentData["reference"].(*string); ok && ref != nil {
		reference = ref
	} else if refStr, ok := paymentData["reference"].(string); ok && refStr != "" {
		reference = &refStr
	}
	
	var notes *string
	if n, ok := paymentData["notes"].(*string); ok && n != nil {
		notes = n
	} else if nStr, ok := paymentData["notes"].(string); ok && nStr != "" {
		notes = &nStr
	}
	
	// Generar número de asiento
	entryNumber, err := s.repo.GetNextEntryNumber(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to generate entry number: %w", err)
	}
	
	entry := &JournalEntry{
		ID:          uuid.New(),
		EntryNumber: entryNumber,
		EntryDate:   paymentDate,
		Description: fmt.Sprintf("Pago %s", paymentNumber),
		Reference:   reference,
		Status:      "posted",
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Lines:       []JournalEntryLine{},
	}
	
	// Determinar cuenta de pago según método
	var paymentAccount *ChartOfAccounts
	var paymentTransactionType string
	switch paymentMethod {
	case "cash":
		paymentTransactionType = TransactionTypePaymentCash
	case "card", "transfer", "check":
		paymentTransactionType = TransactionTypePaymentBank
	default:
		paymentTransactionType = TransactionTypePaymentCash // Por defecto
	}
	
	paymentAccount, err = s.GetAccountForTransaction(ctx, schema, paymentTransactionType, "1111", "Caja General")
	if err != nil {
		return nil, fmt.Errorf("failed to get payment account: %w", err)
	}
	
	// Determinar cuenta de crédito (HABER)
	var creditAccount *ChartOfAccounts
	if customerID != nil {
		// Es un pago de cliente - reducir CXC (Clientes)
		creditAccount, err = s.GetAccountForTransaction(ctx, schema, TransactionTypeSalesCredit, "1121", "Clientes")
		if err != nil {
			return nil, fmt.Errorf("failed to get accounts receivable account: %w", err)
		}
	} else {
		// Es un gasto - usar cuenta de gastos generales
		creditAccount, err = s.GetAccountForTransaction(ctx, schema, TransactionTypeExpenseGeneral, "6100", "Gastos de Administración")
		if err != nil {
			return nil, fmt.Errorf("failed to get expense account: %w", err)
		}
	}
	
	// Crear líneas del asiento
	// DEBE: Cuenta de pago
	entry.Lines = append(entry.Lines, JournalEntryLine{
		ID:             uuid.New(),
		JournalEntryID: entry.ID,
		AccountCode:    paymentAccount.AccountCode,
		AccountName:    paymentAccount.AccountName,
		Debit:          amount,
		Credit:         0,
		Description:    notes,
		CreatedAt:      time.Now(),
	})
	
	// HABER: Cuenta de cliente o gasto
	entry.Lines = append(entry.Lines, JournalEntryLine{
		ID:             uuid.New(),
		JournalEntryID: entry.ID,
		AccountCode:    creditAccount.AccountCode,
		AccountName:    creditAccount.AccountName,
		Debit:          0,
		Credit:         amount,
		Description:    notes,
		CreatedAt:      time.Now(),
	})
	
	// Guardar asiento
	if err := s.repo.CreateJournalEntry(ctx, schema, entry); err != nil {
		logger.Error("Failed to create payment journal entry", "error", err)
		return nil, fmt.Errorf("failed to create journal entry: %w", err)
	}
	
	logger.Info("Payment journal entry created", "entry_number", entry.EntryNumber, "payment_id", paymentID)
	return entry, nil
}

// GenerateExpenseJournalEntry genera un asiento contable para un gasto
// DEBE: Cuenta de gasto según categoría
// HABER: Cuenta de pago (Caja/Banco) según método de pago
func (s *Service) GenerateExpenseJournalEntry(ctx context.Context, schema string, paymentID string, userID *uuid.UUID) (*JournalEntry, error) {
	// Obtener información del gasto
	expenseData, err := s.repo.GetPaymentByID(ctx, schema, paymentID)
	if err != nil {
		return nil, fmt.Errorf("expense not found: %w", err)
	}
	
	// Verificar que es un gasto (tiene prefijo EXP-)
	paymentNumber, _ := expenseData["payment_number"].(string)
	if len(paymentNumber) < 4 || paymentNumber[:4] != "EXP-" {
		return nil, fmt.Errorf("payment is not an expense")
	}
	
	// Extraer valores del mapa
	paymentMethod, _ := expenseData["payment_method"].(string)
	amount, _ := expenseData["amount"].(float64)
	category, _ := expenseData["category"].(string)
	supplier, _ := expenseData["supplier"].(string)
	description, _ := expenseData["description"].(string)
	
	var paymentDate time.Time
	if pd, ok := expenseData["payment_date"].(time.Time); ok {
		paymentDate = pd
	} else if pdStr, ok := expenseData["payment_date"].(string); ok {
		if parsed, err := time.Parse("2006-01-02", pdStr); err == nil {
			paymentDate = parsed
		}
	}
	
	var reference *string
	if ref, ok := expenseData["reference"].(string); ok && ref != "" {
		reference = &ref
	}
	
	var notes *string
	if n, ok := expenseData["notes"].(string); ok && n != "" {
		notes = &n
	}
	
	// Generar número de asiento
	entryNumber, err := s.repo.GetNextEntryNumber(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to generate entry number: %w", err)
	}
	
	desc := fmt.Sprintf("Gasto %s", paymentNumber)
	if description != "" {
		desc = fmt.Sprintf("Gasto %s - %s", paymentNumber, description)
	}
	if supplier != "" {
		desc = fmt.Sprintf("%s - %s", desc, supplier)
	}
	
	entry := &JournalEntry{
		ID:          uuid.New(),
		EntryNumber: entryNumber,
		EntryDate:   paymentDate,
		Description: desc,
		Reference:   reference,
		Status:      "posted",
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Lines:       []JournalEntryLine{},
	}
	
	// Determinar cuenta de gasto según categoría
	var expenseAccount *ChartOfAccounts
	var expenseTransactionType string
	switch category {
	case "utilities", "servicios":
		expenseTransactionType = TransactionTypeExpenseUtilities
	case "rent", "alquiler":
		expenseTransactionType = TransactionTypeExpenseRent
	case "supplies", "suministros":
		expenseTransactionType = TransactionTypeExpenseSupplies
	default:
		expenseTransactionType = TransactionTypeExpenseGeneral
	}
	
	expenseAccount, err = s.GetAccountForTransaction(ctx, schema, expenseTransactionType, "6100", "Gastos de Administración")
	if err != nil {
		return nil, fmt.Errorf("failed to get expense account: %w", err)
	}
	
	// Determinar cuenta de pago según método
	var paymentAccount *ChartOfAccounts
	var paymentTransactionType string
	switch paymentMethod {
	case "cash":
		paymentTransactionType = TransactionTypePaymentCash
	case "card", "transfer", "check":
		paymentTransactionType = TransactionTypePaymentBank
	default:
		paymentTransactionType = TransactionTypePaymentCash // Por defecto
	}
	
	paymentAccount, err = s.GetAccountForTransaction(ctx, schema, paymentTransactionType, "1111", "Caja General")
	if err != nil {
		return nil, fmt.Errorf("failed to get payment account: %w", err)
	}
	
	// Crear líneas del asiento
	// DEBE: Cuenta de gasto
	descPtr := &description
	if description == "" {
		descPtr = notes
	}
	entry.Lines = append(entry.Lines, JournalEntryLine{
		ID:             uuid.New(),
		JournalEntryID: entry.ID,
		AccountCode:    expenseAccount.AccountCode,
		AccountName:    expenseAccount.AccountName,
		Debit:          amount,
		Credit:         0,
		Description:    descPtr,
		CreatedAt:      time.Now(),
	})
	
	// HABER: Cuenta de pago
	entry.Lines = append(entry.Lines, JournalEntryLine{
		ID:             uuid.New(),
		JournalEntryID: entry.ID,
		AccountCode:    paymentAccount.AccountCode,
		AccountName:    paymentAccount.AccountName,
		Debit:          0,
		Credit:         amount,
		Description:    notes,
		CreatedAt:      time.Now(),
	})
	
	// Guardar asiento
	if err := s.repo.CreateJournalEntry(ctx, schema, entry); err != nil {
		logger.Error("Failed to create expense journal entry", "error", err)
		return nil, fmt.Errorf("failed to create journal entry: %w", err)
	}
	
	logger.Info("Expense journal entry created", "entry_number", entry.EntryNumber, "expense_id", paymentID)
	return entry, nil
}

// GetAccountsPayableBalance calculates the current accounts payable balance
// based on journal entries
func (s *Service) GetAccountsPayableBalance(ctx context.Context, schema string) (float64, error) {
	// Get the accounts payable account code
	apAccount, err := s.GetAccountForTransaction(ctx, schema, TransactionTypePurchaseAP, "2111", "Proveedores")
	if err != nil {
		return 0, fmt.Errorf("failed to get accounts payable account: %w", err)
	}

	// Calculate balance from journal entries
	balance, err := s.repo.GetAccountBalanceFromEntries(ctx, schema, apAccount.AccountCode)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate balance: %w", err)
	}

	return balance, nil
}

// RecalculateAccountsPayableBalance recalculates the accounts payable balance
// based on journal entries and updates the account balance
func (s *Service) RecalculateAccountsPayableBalance(ctx context.Context, schema string) error {
	// Get the accounts payable account code
	apAccount, err := s.GetAccountForTransaction(ctx, schema, TransactionTypePurchaseAP, "2111", "Proveedores")
	if err != nil {
		return fmt.Errorf("failed to get accounts payable account: %w", err)
	}

	// Calculate balance from journal entries
	balance, err := s.repo.GetAccountBalanceFromEntries(ctx, schema, apAccount.AccountCode)
	if err != nil {
		return fmt.Errorf("failed to calculate balance: %w", err)
	}

	// Update the account balance
	if err := s.repo.UpdateAccountBalance(ctx, schema, apAccount.AccountCode, balance); err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	logger.Info("Accounts payable balance recalculated", "account_code", apAccount.AccountCode, "balance", balance)
	return nil
}

// GetAccountsPayableSummary returns a summary of accounts payable
// including total outstanding, number of unpaid purchases, etc.
// Note: This method now only returns the account balance from the chart of accounts.
// For detailed aging report, use the reports module.
func (s *Service) GetAccountsPayableSummary(ctx context.Context, schema string) (*AccountsPayableSummary, error) {
	// Get the accounts payable account balance
	balance, err := s.GetAccountsPayableBalance(ctx, schema)
	if err != nil {
		logger.Warn("Failed to get accounts payable balance", "error", err)
		balance = 0
	}

	summary := &AccountsPayableSummary{
		AccountBalance:   balance,
		TotalOutstanding: 0, // Use reports module for detailed aging report
		UnpaidPurchases:  0, // Use reports module for detailed aging report
		Total0_30:        0, // Use reports module for detailed aging report
		Total31_60:       0, // Use reports module for detailed aging report
		Total61_90:       0, // Use reports module for detailed aging report
		TotalOver90:      0, // Use reports module for detailed aging report
		LastCalculated:   time.Now(),
	}

	return summary, nil
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
