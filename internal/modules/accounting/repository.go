package accounting

import (
	"context"
	"fmt"
	"time"

	"vendix/internal/database"

	"github.com/google/uuid"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// Chart of Accounts methods

func (r *Repository) CreateAccount(ctx context.Context, schema string, account *ChartOfAccounts) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.chart_of_accounts (id, account_code, account_name, account_type, parent_account_code, description, is_active, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		account.ID,
		account.AccountCode,
		account.AccountName,
		account.AccountType,
		account.ParentAccountCode,
		account.Description,
		account.IsActive,
		account.Balance,
		account.CreatedAt,
		account.UpdatedAt,
	)

	return err
}

func (r *Repository) GetAccountByCode(ctx context.Context, schema string, code string) (*ChartOfAccounts, error) {
	var account ChartOfAccounts
	query := fmt.Sprintf(`
		SELECT id, account_code, account_name, account_type, parent_account_code, is_active, description, balance, created_at, updated_at
		FROM %s.chart_of_accounts
		WHERE account_code = $1
	`, schema)

	err := r.db.GetContext(ctx, &account, query, code)
	return &account, err
}

func (r *Repository) ListAccounts(ctx context.Context, schema string, accountType *string) ([]*ChartOfAccounts, error) {
	var accounts []*ChartOfAccounts
	baseQuery := fmt.Sprintf(`
		SELECT id, account_code, account_name, account_type, parent_account_code, is_active, description, balance, created_at, updated_at
		FROM %s.chart_of_accounts
		WHERE 1=1
	`, schema)

	args := []interface{}{}
	argCount := 1

	if accountType != nil {
		baseQuery += fmt.Sprintf(" AND account_type = $%d", argCount)
		args = append(args, *accountType)
		argCount++
	}

	baseQuery += " ORDER BY account_code"

	err := r.db.SelectContext(ctx, &accounts, baseQuery, args...)
	if err != nil {
		return nil, err
	}

	// Return empty slice if no accounts found instead of nil
	if accounts == nil {
		accounts = []*ChartOfAccounts{}
	}

	return accounts, nil
}

func (r *Repository) UpdateAccount(ctx context.Context, schema string, account *ChartOfAccounts) error {
	query := fmt.Sprintf(`
		UPDATE %s.chart_of_accounts
		SET account_name = $1, account_type = $2, parent_account_code = $3, 
		    description = $4, is_active = $5, updated_at = $6
		WHERE account_code = $7
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		account.AccountName,
		account.AccountType,
		account.ParentAccountCode,
		account.Description,
		account.IsActive,
		account.UpdatedAt,
		account.AccountCode,
	)

	return err
}

func (r *Repository) UpdateAccountBalance(ctx context.Context, schema string, code string, balance float64) error {
	query := fmt.Sprintf(`UPDATE %s.chart_of_accounts SET balance = $1 WHERE account_code = $2`, schema)
	_, err := r.db.ExecContext(ctx, query, balance, code)
	return err
}

// AccountHasJournalEntries verifica si una cuenta tiene asientos contables
func (r *Repository) AccountHasJournalEntries(ctx context.Context, schema string, accountCode string) (bool, error) {
	var count int
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s.journal_entry_lines
		WHERE account_code = $1
	`, schema)
	
	err := r.db.GetContext(ctx, &count, query, accountCode)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// DeleteAccount elimina una cuenta del plan de cuentas
func (r *Repository) DeleteAccount(ctx context.Context, schema string, accountCode string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.chart_of_accounts
		WHERE account_code = $1
	`, schema)
	
	_, err := r.db.ExecContext(ctx, query, accountCode)
	return err
}

// Journal Entry methods

func (r *Repository) CreateJournalEntry(ctx context.Context, schema string, entry *JournalEntry) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert journal entry
	query := fmt.Sprintf(`
		INSERT INTO %s.journal_entries (id, entry_number, entry_date, description, reference, status, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, schema)

	_, err = tx.ExecContext(ctx, query,
		entry.ID,
		entry.EntryNumber,
		entry.EntryDate,
		entry.Description,
		entry.Reference,
		entry.Status,
		entry.CreatedBy,
		entry.CreatedAt,
		entry.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Insert journal entry lines
	lineQuery := fmt.Sprintf(`
		INSERT INTO %s.journal_entry_lines (id, journal_entry_id, account_code, account_name, debit, credit, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, schema)

	for _, line := range entry.Lines {
		_, err = tx.ExecContext(ctx, lineQuery,
			line.ID,
			line.JournalEntryID,
			line.AccountCode,
			line.AccountName,
			line.Debit,
			line.Credit,
			line.Description,
			line.CreatedAt,
		)
		if err != nil {
			return err
		}

		// Update account balance
		var account ChartOfAccounts
		balanceQuery := fmt.Sprintf("SELECT balance FROM %s.chart_of_accounts WHERE account_code = $1", schema)
		err = tx.GetContext(ctx, &account, balanceQuery, line.AccountCode)
		if err != nil {
			return err
		}

		newBalance := account.Balance + line.Debit - line.Credit
		updateBalanceQuery := fmt.Sprintf("UPDATE %s.chart_of_accounts SET balance = $1 WHERE account_code = $2", schema)
		_, err = tx.ExecContext(ctx, updateBalanceQuery, newBalance, line.AccountCode)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) GetJournalEntryByID(ctx context.Context, schema string, id uuid.UUID) (*JournalEntry, error) {
	var entry JournalEntry
	query := fmt.Sprintf(`
		SELECT id, entry_number, entry_date, description, reference, status, created_by, created_at, updated_at
		FROM %s.journal_entries
		WHERE id = $1
	`, schema)

	err := r.db.GetContext(ctx, &entry, query, id)
	if err != nil {
		return nil, err
	}

	// Get lines
	lines, err := r.getJournalEntryLines(ctx, schema, id)
	if err != nil {
		return nil, err
	}

	entry.Lines = lines
	return &entry, nil
}

func (r *Repository) getJournalEntryLines(ctx context.Context, schema string, entryID uuid.UUID) ([]JournalEntryLine, error) {
	var lines []JournalEntryLine
	query := fmt.Sprintf(`
		SELECT id, journal_entry_id, account_code, account_name, debit, credit, description, created_at
		FROM %s.journal_entry_lines
		WHERE journal_entry_id = $1
		ORDER BY created_at
	`, schema)

	err := r.db.SelectContext(ctx, &lines, query, entryID)
	return lines, err
}

func (r *Repository) ListJournalEntries(ctx context.Context, schema string, startDate, endDate *string, status *string) ([]*JournalEntry, error) {
	var entries []*JournalEntry
	baseQuery := fmt.Sprintf(`
		SELECT id, entry_number, entry_date, description, reference, status, created_by, created_at, updated_at
		FROM %s.journal_entries
		WHERE 1=1
	`, schema)

	args := []interface{}{}
	argCount := 1

	if startDate != nil {
		baseQuery += fmt.Sprintf(" AND entry_date >= $%d", argCount)
		args = append(args, *startDate)
		argCount++
	}

	if endDate != nil {
		baseQuery += fmt.Sprintf(" AND entry_date <= $%d", argCount)
		args = append(args, *endDate)
		argCount++
	}

	if status != nil {
		baseQuery += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *status)
		argCount++
	}

	baseQuery += " ORDER BY entry_date DESC, entry_number DESC"

	err := r.db.SelectContext(ctx, &entries, baseQuery, args...)
	if err != nil {
		return nil, err
	}

	// Get lines for each entry
	for i, entry := range entries {
		lines, err := r.getJournalEntryLines(ctx, schema, entry.ID)
		if err != nil {
			return nil, err
		}
		entries[i].Lines = lines
	}

	return entries, nil
}

func (r *Repository) GetNextEntryNumber(ctx context.Context, schema string) (string, error) {
	var lastNumber string
	query := fmt.Sprintf(`
		SELECT entry_number 
		FROM %s.journal_entries 
		ORDER BY created_at DESC 
		LIMIT 1
	`, schema)

	err := r.db.GetContext(ctx, &lastNumber, query)
	if err != nil {
		// If no entries exist, start with 1
		return "JE-0001", nil
	}

	// Parse and increment (format: JE-0001)
	var num int
	fmt.Sscanf(lastNumber, "JE-%d", &num)
	num++

	return fmt.Sprintf("JE-%04d", num), nil
}

// Reporting methods

func (r *Repository) GetGeneralLedger(ctx context.Context, schema string, accountCode string, startDate, endDate *string) ([]LedgerEntry, error) {
	var entries []LedgerEntry
	baseQuery := fmt.Sprintf(`
		SELECT 
			je.entry_date as date,
			jel.account_code,
			ca.account_name,
			je.description,
			COALESCE(je.reference, '') as reference,
			jel.debit,
			jel.credit,
			je.id as journal_entry_id
		FROM %s.journal_entry_lines jel
		INNER JOIN %s.journal_entries je ON jel.journal_entry_id = je.id
		INNER JOIN %s.chart_of_accounts ca ON jel.account_code = ca.account_code
		WHERE jel.account_code = $1 AND je.status = 'posted'
	`, schema, schema, schema)

	args := []interface{}{accountCode}
	argCount := 2

	if startDate != nil {
		baseQuery += fmt.Sprintf(" AND je.entry_date >= $%d", argCount)
		args = append(args, *startDate)
		argCount++
	}

	if endDate != nil {
		baseQuery += fmt.Sprintf(" AND je.entry_date <= $%d", argCount)
		args = append(args, *endDate)
		argCount++
	}

	baseQuery += " ORDER BY je.entry_date, je.entry_number"

	err := r.db.SelectContext(ctx, &entries, baseQuery, args...)
	if err != nil {
		return nil, err
	}

	// Calculate running balance
	balance := 0.0
	for i := range entries {
		balance += entries[i].Debit - entries[i].Credit
		entries[i].Balance = balance
	}

	return entries, nil
}

func (r *Repository) GetTrialBalance(ctx context.Context, schema string, asOfDate *string) ([]TrialBalanceEntry, error) {
	var entries []TrialBalanceEntry
	baseQuery := fmt.Sprintf(`
		SELECT 
			ca.account_code,
			ca.account_name,
			ca.account_type,
			COALESCE(SUM(jel.debit), 0) as debit,
			COALESCE(SUM(jel.credit), 0) as credit
		FROM %s.chart_of_accounts ca
		LEFT JOIN %s.journal_entry_lines jel ON ca.account_code = jel.account_code
		LEFT JOIN %s.journal_entries je ON jel.journal_entry_id = je.id
		WHERE ca.is_active = true AND (je.status = 'posted' OR je.status IS NULL)
	`, schema, schema, schema)

	args := []interface{}{}
	if asOfDate != nil {
		baseQuery += " AND je.entry_date <= $1"
		args = append(args, *asOfDate)
	}

	baseQuery += " GROUP BY ca.account_code, ca.account_name, ca.account_type ORDER BY ca.account_code"

	err := r.db.SelectContext(ctx, &entries, baseQuery, args...)
	return entries, err
}

// Account Mapping methods

// GetAccountMapping obtiene el mapeo de cuenta para un tipo de transacción
func (r *Repository) GetAccountMapping(ctx context.Context, schema string, transactionType string) (*AccountMapping, error) {
	var mapping AccountMapping
	query := fmt.Sprintf(`
		SELECT id, transaction_type, account_code, account_name, description, is_active, created_at, updated_at
		FROM %s.account_mappings
		WHERE transaction_type = $1 AND is_active = true
	`, schema)
	
	err := r.db.GetContext(ctx, &mapping, query, transactionType)
	if err != nil {
		return nil, err
	}
	return &mapping, nil
}

// GetAccountMappingWithFallback obtiene el mapeo o usa un valor por defecto
func (r *Repository) GetAccountMappingWithFallback(ctx context.Context, schema string, transactionType string, defaultCode string, defaultName string) (*AccountMapping, error) {
	mapping, err := r.GetAccountMapping(ctx, schema, transactionType)
	if err != nil {
		// Si no existe, crear uno temporal con valores por defecto
		return &AccountMapping{
			TransactionType: transactionType,
			AccountCode:     defaultCode,
			AccountName:     defaultName,
			IsActive:        true,
		}, nil
	}
	return mapping, nil
}

// ListAccountMappings lista todos los mapeos
func (r *Repository) ListAccountMappings(ctx context.Context, schema string) ([]*AccountMapping, error) {
	var mappings []*AccountMapping
	query := fmt.Sprintf(`
		SELECT id, transaction_type, account_code, account_name, description, is_active, created_at, updated_at
		FROM %s.account_mappings
		ORDER BY transaction_type ASC
	`, schema)
	
	err := r.db.SelectContext(ctx, &mappings, query)
	return mappings, err
}

// CreateAccountMapping crea un nuevo mapeo
func (r *Repository) CreateAccountMapping(ctx context.Context, schema string, mapping *AccountMapping) error {
	// Verificar que la cuenta existe
	account, err := r.GetAccountByCode(ctx, schema, mapping.AccountCode)
	if err != nil {
		return fmt.Errorf("account %s not found", mapping.AccountCode)
	}
	
	mapping.AccountName = account.AccountName
	
	query := fmt.Sprintf(`
		INSERT INTO %s.account_mappings (id, transaction_type, account_code, account_name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		ON CONFLICT (transaction_type) 
		DO UPDATE SET 
			account_code = EXCLUDED.account_code,
			account_name = EXCLUDED.account_name,
			description = EXCLUDED.description,
			is_active = EXCLUDED.is_active,
			updated_at = NOW()
	`, schema)
	
	_, err = r.db.ExecContext(ctx, query, 
		mapping.ID, mapping.TransactionType, mapping.AccountCode, 
		mapping.AccountName, mapping.Description, mapping.IsActive)
	return err
}

// UpdateAccountMapping actualiza un mapeo existente
func (r *Repository) UpdateAccountMapping(ctx context.Context, schema string, transactionType string, mapping *AccountMapping) error {
	// Si se actualiza el código de cuenta, verificar que existe
	if mapping.AccountCode != "" {
		account, err := r.GetAccountByCode(ctx, schema, mapping.AccountCode)
		if err != nil {
			return fmt.Errorf("account %s not found", mapping.AccountCode)
		}
		mapping.AccountName = account.AccountName
	}
	
	query := fmt.Sprintf(`
		UPDATE %s.account_mappings
		SET account_code = COALESCE(NULLIF($1, ''), account_code),
			account_name = COALESCE(NULLIF($2, ''), account_name),
			description = COALESCE($3, description),
			is_active = COALESCE($4, is_active),
			updated_at = NOW()
		WHERE transaction_type = $5
	`, schema)
	
	_, err := r.db.ExecContext(ctx, query, 
		mapping.AccountCode, mapping.AccountName, mapping.Description, 
		mapping.IsActive, transactionType)
	return err
}

// GetPaymentByID obtiene información de un pago usando sqlx
func (r *Repository) GetPaymentByID(ctx context.Context, schema string, paymentID string) (map[string]interface{}, error) {
	query := fmt.Sprintf(`
		SELECT id, payment_number, customer_id, payment_method, payment_date, amount, reference, notes, category, supplier, description
		FROM %s.payments
		WHERE id = $1
	`, schema)
	
	var result struct {
		ID            uuid.UUID  `db:"id"`
		PaymentNumber string    `db:"payment_number"`
		CustomerID    *uuid.UUID `db:"customer_id"`
		PaymentMethod string    `db:"payment_method"`
		PaymentDate   time.Time `db:"payment_date"`
		Amount        float64   `db:"amount"`
		Reference     *string   `db:"reference"`
		Notes         *string   `db:"notes"`
		Category      string    `db:"category"`
		Supplier      string    `db:"supplier"`
		Description   string    `db:"description"`
	}
	
	if err := r.db.GetContext(ctx, &result, query, paymentID); err != nil {
		return nil, err
	}
	
	// Convertir a mapa
	payment := make(map[string]interface{})
	payment["id"] = result.ID
	payment["payment_number"] = result.PaymentNumber
	payment["customer_id"] = result.CustomerID
	payment["payment_method"] = result.PaymentMethod
	payment["payment_date"] = result.PaymentDate
	payment["amount"] = result.Amount
	payment["reference"] = result.Reference
	payment["notes"] = result.Notes
	payment["category"] = result.Category
	payment["supplier"] = result.Supplier
	payment["description"] = result.Description
	
	return payment, nil
}
