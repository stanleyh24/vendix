package notifications

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypePurchaseDueSoon  NotificationType = "purchase_due_soon"
	NotificationTypePurchaseOverdue NotificationType = "purchase_overdue"
	NotificationTypeInvoiceDueSoon   NotificationType = "invoice_due_soon"
	NotificationTypeInvoiceOverdue   NotificationType = "invoice_overdue"
)

// NotificationPriority represents the priority level
type NotificationPriority string

const (
	PriorityLow    NotificationPriority = "low"
	PriorityNormal NotificationPriority = "normal"
	PriorityHigh   NotificationPriority = "high"
	PriorityUrgent NotificationPriority = "urgent"
)

// Notification represents a notification
type Notification struct {
	ID         uuid.UUID              `db:"id" json:"id"`
	TenantID   uuid.UUID              `db:"tenant_id" json:"tenant_id"`
	Type       string                 `db:"type" json:"type"`
	Title      string                 `db:"title" json:"title"`
	Message    string                 `db:"message" json:"message"`
	EntityType *string                `db:"entity_type" json:"entity_type,omitempty"`
	EntityID   *uuid.UUID             `db:"entity_id" json:"entity_id,omitempty"`
	Priority   string                 `db:"priority" json:"priority"`
	IsRead     bool                   `db:"is_read" json:"is_read"`
	ReadAt     *time.Time             `db:"read_at" json:"read_at,omitempty"`
	Metadata   map[string]interface{} `db:"metadata" json:"metadata,omitempty"`
	CreatedAt  time.Time              `db:"created_at" json:"created_at"`
}

// CreateNotificationRequest DTO
type CreateNotificationRequest struct {
	Type       string                 `json:"type" validate:"required"`
	Title      string                 `json:"title" validate:"required"`
	Message    string                 `json:"message" validate:"required"`
	EntityType *string                `json:"entity_type,omitempty"`
	EntityID   *string                `json:"entity_id,omitempty"`
	Priority   string                 `json:"priority"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// NotificationResponse includes additional computed fields
type NotificationResponse struct {
	Notification
	DaysUntilDue *int `json:"days_until_due,omitempty"`
	DaysOverdue  *int `json:"days_overdue,omitempty"`
}

