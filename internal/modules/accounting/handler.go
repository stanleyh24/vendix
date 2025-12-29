package accounting

import (
	"strings"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(db *database.DB, cfg *config.Config) *Handler {
	return &Handler{service: NewService(db, cfg)}
}

// RegisterRoutes registers accounting routes
func RegisterRoutes(router fiber.Router, db *database.DB, cfg *config.Config) {
	h := NewHandler(db, cfg)
	accounting := router.Group("/accounting")

	// Chart of Accounts
	accounting.Get("/accounts", h.ListAccounts)
	accounting.Post("/accounts", h.CreateAccount)
	accounting.Get("/accounts/:code", h.GetAccount)
	accounting.Put("/accounts/:code", h.UpdateAccount)
	accounting.Delete("/accounts/:code", h.DeleteAccount)

	// Journal Entries
	accounting.Get("/journal-entries", h.ListJournalEntries)
	accounting.Post("/journal-entries", h.CreateJournalEntry)
	// Daily Sales Journal Entry - must be registered before parameterized routes
	accounting.Post("/journal-entries/daily-sales", h.GenerateDailySalesJournalEntry)
	accounting.Get("/journal-entries/:id", h.GetJournalEntry)

	// Reports
	accounting.Get("/ledger/:accountCode", h.GetGeneralLedger)
	accounting.Get("/trial-balance", h.GetTrialBalance)
	accounting.Get("/balance-sheet", h.GetBalanceSheet)
	accounting.Get("/income-statement", h.GetIncomeStatement)

	// Account Mappings
	accounting.Get("/account-mappings", h.ListAccountMappings)
	accounting.Post("/account-mappings", h.CreateOrUpdateAccountMapping)
	accounting.Get("/account-mappings/:type", h.GetAccountMapping)
	accounting.Put("/account-mappings/:type", h.UpdateAccountMapping)

	// Payment and Expense Journal Entries
	accounting.Post("/journal-entries/payment/:id", h.GeneratePaymentJournalEntry)
	accounting.Post("/journal-entries/expense/:id", h.GenerateExpenseJournalEntry)
}

// Chart of Accounts handlers

// ListAccounts lists all accounts
// @Summary List accounts
// @Tags accounting
// @Produce json
// @Param account_type query string false "Filter by account type"
// @Success 200 {array} ChartOfAccounts
// @Router /api/v1/tenant/accounting/accounts [get]
func (h *Handler) ListAccounts(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	if schema == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant schema not found in context",
		})
	}

	accountType := c.Query("account_type")
	var accountTypePtr *string
	if accountType != "" {
		accountTypePtr = &accountType
	}

	accounts, err := h.service.ListAccounts(c.Context(), schema, accountTypePtr)
	if err != nil {
		// Check if it's a "table does not exist" error
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "does not exist") || strings.Contains(errMsg, "no existe") {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "La tabla de contabilidad no existe. Por favor ejecuta las migraciones: ./scripts/migrate-all-tenants-sql.sh",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(accounts)
}

// CreateAccount creates a new account
// @Summary Create account
// @Tags accounting
// @Accept json
// @Produce json
// @Param account body CreateAccountRequest true "Account data"
// @Success 201 {object} ChartOfAccounts
// @Router /api/v1/tenant/accounting/accounts [post]
func (h *Handler) CreateAccount(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	var req CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	account, err := h.service.CreateAccount(c.Context(), schema, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(account)
}

// GetAccount retrieves an account by code
// @Summary Get account
// @Tags accounting
// @Produce json
// @Param code path string true "Account code"
// @Success 200 {object} ChartOfAccounts
// @Router /api/v1/tenant/accounting/accounts/{code} [get]
func (h *Handler) GetAccount(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	code := c.Params("code")

	account, err := h.service.GetAccountByCode(c.Context(), schema, code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account not found",
		})
	}

	return c.JSON(account)
}

// UpdateAccount updates an account
// @Summary Update account
// @Tags accounting
// @Accept json
// @Produce json
// @Param code path string true "Account code"
// @Param account body UpdateAccountRequest true "Account data"
// @Success 200 {object} ChartOfAccounts
// @Router /api/v1/tenant/accounting/accounts/{code} [put]
func (h *Handler) UpdateAccount(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	code := c.Params("code")

	var req UpdateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	account, err := h.service.UpdateAccount(c.Context(), schema, code, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(account)
}

// DeleteAccount elimina una cuenta del plan de cuentas
// @Summary Delete account
// @Tags accounting
// @Produce json
// @Param code path string true "Account code"
// @Success 204 "Account deleted successfully"
// @Failure 400 {object} map[string]string "Account has journal entries or child accounts"
// @Failure 404 {object} map[string]string "Account not found"
// @Router /api/v1/tenant/accounting/accounts/{code} [delete]
func (h *Handler) DeleteAccount(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	code := c.Params("code")

	err := h.service.DeleteAccount(c.Context(), schema, code)
	if err != nil {
		errMsg := err.Error()
		// Verificar si es un error de validación (tiene asientos o cuentas hijas)
		if strings.Contains(errMsg, "has journal entries") || strings.Contains(errMsg, "has child accounts") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": errMsg,
			})
		}
		// Si contiene "not found", es un error 404
		if strings.Contains(errMsg, "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": errMsg,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errMsg,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Journal Entry handlers

// ListJournalEntries lists all journal entries
// @Summary List journal entries
// @Tags accounting
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param status query string false "Filter by status"
// @Success 200 {array} JournalEntry
// @Router /api/v1/tenant/accounting/journal-entries [get]
func (h *Handler) ListJournalEntries(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")

	var startDatePtr, endDatePtr, statusPtr *string
	if startDate != "" {
		startDatePtr = &startDate
	}
	if endDate != "" {
		endDatePtr = &endDate
	}
	if status != "" {
		statusPtr = &status
	}

	entries, err := h.service.ListJournalEntries(c.Context(), schema, startDatePtr, endDatePtr, statusPtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(entries)
}

// CreateJournalEntry creates a new journal entry
// @Summary Create journal entry
// @Tags accounting
// @Accept json
// @Produce json
// @Param entry body CreateJournalEntryRequest true "Journal entry data"
// @Success 201 {object} JournalEntry
// @Router /api/v1/tenant/accounting/journal-entries [post]
func (h *Handler) CreateJournalEntry(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	var req CreateJournalEntryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	entry, err := h.service.CreateJournalEntry(c.Context(), schema, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(entry)
}

// GetJournalEntry retrieves a journal entry by ID
// @Summary Get journal entry
// @Tags accounting
// @Produce json
// @Param id path string true "Journal entry ID"
// @Success 200 {object} JournalEntry
// @Router /api/v1/tenant/accounting/journal-entries/{id} [get]
func (h *Handler) GetJournalEntry(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	entry, err := h.service.GetJournalEntry(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Journal entry not found",
		})
	}

	return c.JSON(entry)
}

// Report handlers

// GetGeneralLedger retrieves the general ledger for an account
// @Summary Get general ledger
// @Tags accounting
// @Produce json
// @Param accountCode path string true "Account code"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} LedgerEntry
// @Router /api/v1/tenant/accounting/ledger/{accountCode} [get]
func (h *Handler) GetGeneralLedger(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	accountCode := c.Params("accountCode")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var startDatePtr, endDatePtr *string
	if startDate != "" {
		startDatePtr = &startDate
	}
	if endDate != "" {
		endDatePtr = &endDate
	}

	ledger, err := h.service.GetGeneralLedger(c.Context(), schema, accountCode, startDatePtr, endDatePtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(ledger)
}

// GetTrialBalance retrieves the trial balance
// @Summary Get trial balance
// @Tags accounting
// @Produce json
// @Param as_of_date query string false "As of date (YYYY-MM-DD)"
// @Success 200 {array} TrialBalanceEntry
// @Router /api/v1/tenant/accounting/trial-balance [get]
func (h *Handler) GetTrialBalance(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	asOfDate := c.Query("as_of_date")

	var asOfDatePtr *string
	if asOfDate != "" {
		asOfDatePtr = &asOfDate
	}

	balance, err := h.service.GetTrialBalance(c.Context(), schema, asOfDatePtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(balance)
}

// GetBalanceSheet retrieves the balance sheet
// @Summary Get balance sheet
// @Tags accounting
// @Produce json
// @Param as_of_date query string false "As of date (YYYY-MM-DD)"
// @Success 200 {object} FinancialStatementReport
// @Router /api/v1/tenant/accounting/balance-sheet [get]
func (h *Handler) GetBalanceSheet(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	asOfDate := c.Query("as_of_date")
	if asOfDate == "" {
		asOfDate = "2024-12-31" // Default to end of year
	}

	report, err := h.service.GetBalanceSheet(c.Context(), schema, asOfDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(report)
}

// GetIncomeStatement retrieves the income statement
// @Summary Get income statement
// @Tags accounting
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} FinancialStatementReport
// @Router /api/v1/tenant/accounting/income-statement [get]
func (h *Handler) GetIncomeStatement(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" {
		startDate = "2024-01-01"
	}
	if endDate == "" {
		endDate = "2024-12-31"
	}

	report, err := h.service.GetIncomeStatement(c.Context(), schema, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(report)
}

// GenerateDailySalesJournalEntry generates a journal entry for all sales of a specific day
// @Summary Generate daily sales journal entry
// @Tags accounting
// @Accept json
// @Produce json
// @Param request body GenerateDailySalesJournalEntryRequest true "Date for journal entry"
// @Success 201 {object} JournalEntry
// @Router /api/v1/tenant/accounting/journal-entries/daily-sales [post]
func (h *Handler) GenerateDailySalesJournalEntry(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	if schema == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant schema not found in context",
		})
	}

	var req GenerateDailySalesJournalEntryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user ID from context (set by auth middleware)
	var userIDPtr *uuid.UUID
	if userIDStr, ok := c.Locals("user_id").(string); ok && userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err == nil {
			userIDPtr = &userID
		}
	}

	entry, err := h.service.GenerateDailySalesJournalEntry(c.Context(), schema, req.Date, userIDPtr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(entry)
}

// Account Mapping handlers

// ListAccountMappings lista todos los mapeos de cuentas
// @Summary List account mappings
// @Tags accounting
// @Produce json
// @Success 200 {array} AccountMapping
// @Security BearerAuth
// @Router /api/v1/tenant/accounting/account-mappings [get]
func (h *Handler) ListAccountMappings(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	mappings, err := h.service.ListAccountMappings(c.Context(), schema)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(mappings)
}

// GetAccountMapping obtiene el mapeo para un tipo de transacción
// @Summary Get account mapping
// @Tags accounting
// @Produce json
// @Param type path string true "Transaction type"
// @Success 200 {object} AccountMapping
// @Security BearerAuth
// @Router /api/v1/tenant/accounting/account-mappings/{type} [get]
func (h *Handler) GetAccountMapping(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	transactionType := c.Params("type")
	
	mapping, err := h.service.GetAccountMapping(c.Context(), schema, transactionType)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Account mapping not found",
		})
	}
	return c.JSON(mapping)
}

// CreateOrUpdateAccountMapping crea o actualiza un mapeo
// @Summary Create or update account mapping
// @Tags accounting
// @Accept json
// @Produce json
// @Param mapping body CreateAccountMappingRequest true "Account mapping data"
// @Success 201 {object} AccountMapping
// @Security BearerAuth
// @Router /api/v1/tenant/accounting/account-mappings [post]
func (h *Handler) CreateOrUpdateAccountMapping(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	var req CreateAccountMappingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	mapping, err := h.service.CreateOrUpdateAccountMapping(c.Context(), schema, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(mapping)
}

// UpdateAccountMapping actualiza un mapeo existente
// @Summary Update account mapping
// @Tags accounting
// @Accept json
// @Produce json
// @Param type path string true "Transaction type"
// @Param mapping body UpdateAccountMappingRequest true "Account mapping data"
// @Success 200 {object} AccountMapping
// @Security BearerAuth
// @Router /api/v1/tenant/accounting/account-mappings/{type} [put]
func (h *Handler) UpdateAccountMapping(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	transactionType := c.Params("type")
	var req UpdateAccountMappingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	mapping, err := h.service.UpdateAccountMapping(c.Context(), schema, transactionType, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(mapping)
}

// GeneratePaymentJournalEntry genera un asiento contable para un pago
// @Summary Generate payment journal entry
// @Tags accounting
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Success 201 {object} JournalEntry
// @Security BearerAuth
// @Router /api/v1/tenant/accounting/journal-entries/payment/{id} [post]
func (h *Handler) GeneratePaymentJournalEntry(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	paymentID := c.Params("id")
	
	var userIDPtr *uuid.UUID
	if userIDStr := c.Get("X-User-ID"); userIDStr != "" {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			userIDPtr = &userID
		}
	}
	
	entry, err := h.service.GeneratePaymentJournalEntry(c.Context(), schema, paymentID, userIDPtr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.Status(fiber.StatusCreated).JSON(entry)
}

// GenerateExpenseJournalEntry genera un asiento contable para un gasto
// @Summary Generate expense journal entry
// @Tags accounting
// @Accept json
// @Produce json
// @Param id path string true "Payment/Expense ID"
// @Success 201 {object} JournalEntry
// @Security BearerAuth
// @Router /api/v1/tenant/accounting/journal-entries/expense/{id} [post]
func (h *Handler) GenerateExpenseJournalEntry(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	paymentID := c.Params("id")
	
	var userIDPtr *uuid.UUID
	if userIDStr := c.Get("X-User-ID"); userIDStr != "" {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			userIDPtr = &userID
		}
	}
	
	entry, err := h.service.GenerateExpenseJournalEntry(c.Context(), schema, paymentID, userIDPtr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.Status(fiber.StatusCreated).JSON(entry)
}
