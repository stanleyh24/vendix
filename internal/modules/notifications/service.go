package notifications

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

// Create creates a new notification
func (s *Service) Create(ctx context.Context, schema string, tenantID uuid.UUID, req *CreateNotificationRequest) (*Notification, error) {
	priority := string(PriorityNormal)
	if req.Priority != "" {
		priority = req.Priority
	}

	var entityID *uuid.UUID
	if req.EntityID != nil {
		id, err := uuid.Parse(*req.EntityID)
		if err != nil {
			return nil, fmt.Errorf("invalid entity ID: %w", err)
		}
		entityID = &id
	}

	notification := &Notification{
		ID:         uuid.New(),
		TenantID:   tenantID,
		Type:       req.Type,
		Title:      req.Title,
		Message:    req.Message,
		EntityType: req.EntityType,
		EntityID:   entityID,
		Priority:   priority,
		IsRead:     false,
		Metadata:   req.Metadata,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, schema, notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return notification, nil
}

// List gets notifications for a tenant
func (s *Service) List(ctx context.Context, schema string, tenantID uuid.UUID, unreadOnly bool, notificationType *string, limit, offset int) ([]*NotificationResponse, error) {
	notifications, err := s.repo.List(ctx, schema, tenantID, unreadOnly, notificationType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}

	responses := make([]*NotificationResponse, len(notifications))
	for i, n := range notifications {
		response := &NotificationResponse{
			Notification: *n,
		}

		// Calculate days until due or days overdue from metadata
		if daysUntilDue, ok := n.Metadata["days_until_due"].(float64); ok {
			days := int(daysUntilDue)
			response.DaysUntilDue = &days
		}
		if daysOverdue, ok := n.Metadata["days_overdue"].(float64); ok {
			days := int(daysOverdue)
			response.DaysOverdue = &days
		}

		responses[i] = response
	}

	return responses, nil
}

// GetByID gets a notification by ID
func (s *Service) GetByID(ctx context.Context, schema string, id string) (*NotificationResponse, error) {
	notificationID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid notification ID: %w", err)
	}

	notification, err := s.repo.GetByID(ctx, schema, notificationID)
	if err != nil {
		return nil, fmt.Errorf("notification not found: %w", err)
	}

	response := &NotificationResponse{
		Notification: *notification,
	}

	// Calculate days until due or days overdue from metadata
	if daysUntilDue, ok := notification.Metadata["days_until_due"].(float64); ok {
		days := int(daysUntilDue)
		response.DaysUntilDue = &days
	}
	if daysOverdue, ok := notification.Metadata["days_overdue"].(float64); ok {
		days := int(daysOverdue)
		response.DaysOverdue = &days
	}

	return response, nil
}

// MarkAsRead marks a notification as read
func (s *Service) MarkAsRead(ctx context.Context, schema string, id string) error {
	notificationID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid notification ID: %w", err)
	}

	return s.repo.MarkAsRead(ctx, schema, notificationID)
}

// MarkAllAsRead marks all notifications as read for a tenant
func (s *Service) MarkAllAsRead(ctx context.Context, schema string, tenantID uuid.UUID) error {
	return s.repo.MarkAllAsRead(ctx, schema, tenantID)
}

// GetUnreadCount gets the count of unread notifications
func (s *Service) GetUnreadCount(ctx context.Context, schema string, tenantID uuid.UUID) (int, error) {
	return s.repo.GetUnreadCount(ctx, schema, tenantID)
}

// Delete deletes a notification
func (s *Service) Delete(ctx context.Context, schema string, id string) error {
	notificationID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid notification ID: %w", err)
	}

	return s.repo.Delete(ctx, schema, notificationID)
}

// CheckPurchaseDueNotifications checks for purchases due soon or overdue and creates notifications
func (s *Service) CheckPurchaseDueNotifications(ctx context.Context, schema string, tenantID uuid.UUID, daysBefore int) error {
	// Get purchases due soon
	dueSoon, err := s.repo.GetPurchasesDueSoon(ctx, schema, daysBefore)
	if err != nil {
		return fmt.Errorf("failed to get purchases due soon: %w", err)
	}

	// Create notifications for purchases due soon
	for _, purchase := range dueSoon {
		// Check if notification already exists
		existing, _ := s.repo.List(ctx, schema, tenantID, false, stringPtr(string(NotificationTypePurchaseDueSoon)), 1, 0)
		alreadyNotified := false
		for _, n := range existing {
			if n.EntityID != nil && *n.EntityID == purchase.PurchaseID {
				alreadyNotified = true
				break
			}
		}

		if !alreadyNotified {
			priority := string(PriorityNormal)
			if purchase.DaysUntilDue <= 3 {
				priority = string(PriorityHigh)
			} else if purchase.DaysUntilDue <= 7 {
				priority = string(PriorityNormal)
			}

			req := &CreateNotificationRequest{
				Type:       string(NotificationTypePurchaseDueSoon),
				Title:      fmt.Sprintf("Compra próxima a vencer: %s", purchase.PurchaseNumber),
				Message:    fmt.Sprintf("La compra %s de %s vence en %d días. Monto pendiente: RD$%.2f", purchase.PurchaseNumber, purchase.SupplierName, purchase.DaysUntilDue, purchase.Balance),
				EntityType: stringPtr("purchase"),
				EntityID:   stringPtr(purchase.PurchaseID.String()),
				Priority:   priority,
				Metadata: map[string]interface{}{
					"purchase_number": purchase.PurchaseNumber,
					"supplier_name":    purchase.SupplierName,
					"balance":          purchase.Balance,
					"days_until_due":   purchase.DaysUntilDue,
					"due_date":         purchase.PaymentDueDate.Format("2006-01-02"),
				},
			}

			if _, err := s.Create(ctx, schema, tenantID, req); err != nil {
				logger.Warn("Failed to create notification for purchase due soon", "error", err, "purchase_id", purchase.PurchaseID)
			}
		}
	}

	// Get overdue purchases
	overdue, err := s.repo.GetPurchasesOverdue(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to get overdue purchases: %w", err)
	}

	// Create notifications for overdue purchases
	for _, purchase := range overdue {
		// Check if notification already exists (only create one per purchase)
		existing, _ := s.repo.List(ctx, schema, tenantID, false, stringPtr(string(NotificationTypePurchaseOverdue)), 1, 0)
		alreadyNotified := false
		for _, n := range existing {
			if n.EntityID != nil && *n.EntityID == purchase.PurchaseID {
				// Update existing notification if it's been more than 7 days
				if daysSinceNotif := int(time.Since(n.CreatedAt).Hours() / 24); daysSinceNotif < 7 {
					alreadyNotified = true
				}
				break
			}
		}

		if !alreadyNotified {
			priority := string(PriorityUrgent)
			if purchase.DaysOverdue > 30 {
				priority = string(PriorityUrgent)
			} else if purchase.DaysOverdue > 15 {
				priority = string(PriorityHigh)
			}

			req := &CreateNotificationRequest{
				Type:       string(NotificationTypePurchaseOverdue),
				Title:      fmt.Sprintf("Compra vencida: %s", purchase.PurchaseNumber),
				Message:    fmt.Sprintf("La compra %s de %s está vencida hace %d días. Monto pendiente: RD$%.2f", purchase.PurchaseNumber, purchase.SupplierName, purchase.DaysOverdue, purchase.Balance),
				EntityType: stringPtr("purchase"),
				EntityID:   stringPtr(purchase.PurchaseID.String()),
				Priority:   priority,
				Metadata: map[string]interface{}{
					"purchase_number": purchase.PurchaseNumber,
					"supplier_name":    purchase.SupplierName,
					"balance":          purchase.Balance,
					"days_overdue":     purchase.DaysOverdue,
					"due_date":         purchase.PaymentDueDate.Format("2006-01-02"),
				},
			}

			if _, err := s.Create(ctx, schema, tenantID, req); err != nil {
				logger.Warn("Failed to create notification for overdue purchase", "error", err, "purchase_id", purchase.PurchaseID)
			}
		}
	}

	return nil
}

func stringPtr(s string) *string {
	return &s
}

