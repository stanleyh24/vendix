package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stanleyh24/vendix/internal/config"
	"github.com/stanleyh24/vendix/internal/database"
	"github.com/stanleyh24/vendix/internal/dgii"
	"github.com/stanleyh24/vendix/internal/logger"
	"github.com/stanleyh24/vendix/internal/modules/notifications"
	"github.com/stanleyh24/vendix/internal/modules/tenants"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

// Task types
const (
	TaskSendEmail                 = "task:send_email"
	TaskSignInvoice               = "task:sign_invoice"
	TaskSendInvoiceToDGII         = "task:send_invoice_dgii"
	TaskProcessRecurringBilling   = "task:process_recurring_billing"
	TaskGenerateMonthlyReports    = "task:generate_monthly_reports"
	TaskCheckPurchaseDueDates     = "task:check_purchase_due_dates"
	TaskSchedulePurchaseDueChecks = "task:schedule_purchase_due_checks" // Scheduler task that creates jobs for all tenants
)

type Handlers struct {
	db   *database.DB
	cfg  *config.Config
	dgii *dgii.Service
}

func NewHandlers(db *database.DB, cfg *config.Config) *Handlers {
	return &Handlers{
		db:   db,
		cfg:  cfg,
		dgii: dgii.NewService(cfg),
	}
}

// SendEmailPayload represents the payload for sending an email
type SendEmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// HandleSendEmail handles email sending tasks
func (h *Handlers) HandleSendEmail(ctx context.Context, task *asynq.Task) error {
	var payload SendEmailPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	logger.Info("Sending email",
		"to", payload.To,
		"subject", payload.Subject,
	)

	// TODO: Implement actual email sending logic
	// For now, just log it
	logger.Info("Email sent successfully", "to", payload.To)

	return nil
}

// SignInvoicePayload represents the payload for signing an invoice
type SignInvoicePayload struct {
	InvoiceID string `json:"invoice_id"`
	Schema    string `json:"schema"`
}

// HandleSignInvoice handles invoice signing tasks
func (h *Handlers) HandleSignInvoice(ctx context.Context, task *asynq.Task) error {
	var payload SignInvoicePayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	logger.Info("Signing invoice",
		"invoice_id", payload.InvoiceID,
		"schema", payload.Schema,
	)

	// TODO: Implement actual invoice signing logic
	// 1. Fetch invoice from database
	// 2. Create DGII electronic document
	// 3. Sign document using DGII service
	// 4. Update invoice with signature

	logger.Info("Invoice signed successfully", "invoice_id", payload.InvoiceID)

	return nil
}

// SendInvoiceToDGIIPayload represents the payload for sending invoice to DGII
type SendInvoiceToDGIIPayload struct {
	InvoiceID string `json:"invoice_id"`
	Schema    string `json:"schema"`
}

// HandleSendInvoiceToDGII handles sending invoices to DGII
func (h *Handlers) HandleSendInvoiceToDGII(ctx context.Context, task *asynq.Task) error {
	var payload SendInvoiceToDGIIPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	logger.Info("Sending invoice to DGII",
		"invoice_id", payload.InvoiceID,
		"schema", payload.Schema,
	)

	// TODO: Implement actual DGII submission logic
	// 1. Fetch signed invoice from database
	// 2. Create DGII electronic document
	// 3. Send to DGII using DGII service
	// 4. Update invoice with DGII response

	logger.Info("Invoice sent to DGII successfully", "invoice_id", payload.InvoiceID)

	return nil
}

// HandleProcessRecurringBilling handles recurring billing processing
func (h *Handlers) HandleProcessRecurringBilling(ctx context.Context, task *asynq.Task) error {
	logger.Info("Processing recurring billing")

	// TODO: Implement recurring billing logic
	// 1. Find all subscriptions that need to be billed
	// 2. Create invoices for each subscription
	// 3. Charge customers via Stripe
	// 4. Send notifications

	logger.Info("Recurring billing processed successfully")

	return nil
}

// HandleGenerateMonthlyReports handles monthly report generation
func (h *Handlers) HandleGenerateMonthlyReports(ctx context.Context, task *asynq.Task) error {
	logger.Info("Generating monthly reports")

	// TODO: Implement monthly report generation logic
	// 1. Generate sales reports for all tenants
	// 2. Generate tax reports
	// 3. Store reports in S3
	// 4. Send notifications to tenants

	logger.Info("Monthly reports generated successfully")

	return nil
}

// CheckPurchaseDueDatesPayload represents the payload for checking purchase due dates
type CheckPurchaseDueDatesPayload struct {
	TenantID   string `json:"tenant_id"`
	Schema     string `json:"schema"`
	DaysBefore int    `json:"days_before"` // Default: 7 days
}

// HandleCheckPurchaseDueDates checks for purchases due soon or overdue and creates notifications
func (h *Handlers) HandleCheckPurchaseDueDates(ctx context.Context, task *asynq.Task) error {
	var payload CheckPurchaseDueDatesPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if payload.DaysBefore == 0 {
		payload.DaysBefore = 7 // Default: check 7 days before
	}

	logger.Info("Checking purchase due dates",
		"tenant_id", payload.TenantID,
		"schema", payload.Schema,
		"days_before", payload.DaysBefore,
	)

	// Import notifications service
	// Note: We need to import it here to avoid circular dependencies
	// For now, we'll use the repository directly
	tenantID, err := uuid.Parse(payload.TenantID)
	if err != nil {
		return fmt.Errorf("invalid tenant ID: %w", err)
	}

	// Get notifications repository
	notificationsRepo := notifications.NewRepository(h.db)
	
	// Get purchases due soon
	dueSoon, err := notificationsRepo.GetPurchasesDueSoon(ctx, payload.Schema, payload.DaysBefore)
	if err != nil {
		logger.Error("Failed to get purchases due soon", "error", err)
		return fmt.Errorf("failed to get purchases due soon: %w", err)
	}

	// Get overdue purchases
	overdue, err := notificationsRepo.GetPurchasesOverdue(ctx, payload.Schema)
	if err != nil {
		logger.Error("Failed to get overdue purchases", "error", err)
		return fmt.Errorf("failed to get overdue purchases: %w", err)
	}

	// Create notifications service
	notificationsSvc := notifications.NewService(h.db, h.cfg)

	// Check and create notifications
	if err := notificationsSvc.CheckPurchaseDueNotifications(ctx, payload.Schema, tenantID, payload.DaysBefore); err != nil {
		logger.Error("Failed to create purchase due notifications", "error", err)
		return fmt.Errorf("failed to create notifications: %w", err)
	}

	logger.Info("Purchase due dates checked",
		"due_soon_count", len(dueSoon),
		"overdue_count", len(overdue),
	)

	return nil
}

// HandleSchedulePurchaseDueChecks is a scheduler task that creates individual jobs for each active tenant
func (h *Handlers) HandleSchedulePurchaseDueChecks(ctx context.Context, task *asynq.Task) error {
	logger.Info("Scheduling purchase due date checks for all active tenants")

	// Get all active tenants
	tenantRepo := tenants.NewRepository(h.db)
	activeTenants, err := tenantRepo.ListActive(ctx)
	if err != nil {
		logger.Error("Failed to get active tenants", "error", err)
		return fmt.Errorf("failed to get active tenants: %w", err)
	}

	if len(activeTenants) == 0 {
		logger.Info("No active tenants found, skipping purchase due date checks")
		return nil
	}

	// Create a job client to enqueue tasks
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     h.cfg.RedisAddr,
		Password: h.cfg.RedisPassword,
		DB:       h.cfg.RedisDB,
	})
	defer client.Close()

	// Create a job for each active tenant
	enqueuedCount := 0
	for _, tenant := range activeTenants {
		payload, err := json.Marshal(CheckPurchaseDueDatesPayload{
			TenantID:   tenant.ID.String(),
			Schema:     tenant.SchemaName,
			DaysBefore: 7, // Check 7 days before due date
		})
		if err != nil {
			logger.Error("Failed to marshal payload for tenant", "tenant_id", tenant.ID, "error", err)
			continue
		}

		task := asynq.NewTask(TaskCheckPurchaseDueDates, payload, asynq.Queue("default"))
		_, err = client.EnqueueContext(ctx, task)
		if err != nil {
			logger.Error("Failed to enqueue purchase due date check for tenant",
				"tenant_id", tenant.ID,
				"tenant_name", tenant.Name,
				"error", err,
			)
			continue
		}

		enqueuedCount++
		logger.Debug("Enqueued purchase due date check",
			"tenant_id", tenant.ID,
			"tenant_name", tenant.Name,
			"schema", tenant.SchemaName,
		)
	}

	logger.Info("Scheduled purchase due date checks",
		"total_tenants", len(activeTenants),
		"enqueued_jobs", enqueuedCount,
	)

	return nil
}
