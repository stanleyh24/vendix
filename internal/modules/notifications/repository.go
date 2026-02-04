package notifications

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/stanleyh24/vendix/internal/database"

	"github.com/google/uuid"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new notification
func (r *Repository) Create(ctx context.Context, schema string, notification *Notification) error {
	metadataJSON, err := json.Marshal(notification.Metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	query := fmt.Sprintf(`
		INSERT INTO %s.notifications (
			id, tenant_id, type, title, message, entity_type, entity_id,
			priority, is_read, read_at, metadata, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, schema)

	var entityID interface{}
	if notification.EntityID != nil {
		entityID = *notification.EntityID
	} else {
		entityID = nil
	}

	_, err = r.db.ExecContext(ctx, query,
		notification.ID,
		notification.TenantID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.EntityType,
		entityID,
		notification.Priority,
		notification.IsRead,
		notification.ReadAt,
		metadataJSON,
		notification.CreatedAt,
	)
	return err
}

// GetByID gets a notification by ID
func (r *Repository) GetByID(ctx context.Context, schema string, id uuid.UUID) (*Notification, error) {
	var notification Notification
	query := fmt.Sprintf(`
		SELECT id, tenant_id, type, title, message, entity_type, entity_id,
			priority, is_read, read_at, metadata, created_at
		FROM %s.notifications WHERE id = $1
	`, schema)

	var metadataJSON []byte
	var entityID sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&notification.ID,
		&notification.TenantID,
		&notification.Type,
		&notification.Title,
		&notification.Message,
		&notification.EntityType,
		&entityID,
		&notification.Priority,
		&notification.IsRead,
		&notification.ReadAt,
		&metadataJSON,
		&notification.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if entityID.Valid {
		if id, err := uuid.Parse(entityID.String); err == nil {
			notification.EntityID = &id
		}
	}

	if len(metadataJSON) > 0 {
		json.Unmarshal(metadataJSON, &notification.Metadata)
	}

	return &notification, nil
}

// List gets notifications with optional filters
func (r *Repository) List(ctx context.Context, schema string, tenantID uuid.UUID, unreadOnly bool, notificationType *string, limit, offset int) ([]*Notification, error) {
	query := fmt.Sprintf(`
		SELECT id, tenant_id, type, title, message, entity_type, entity_id,
			priority, is_read, read_at, metadata, created_at
		FROM %s.notifications
		WHERE tenant_id = $1
	`, schema)

	args := []interface{}{tenantID}
	argPos := 2

	if unreadOnly {
		query += fmt.Sprintf(" AND is_read = false")
	}

	if notificationType != nil {
		query += fmt.Sprintf(" AND type = $%d", argPos)
		args = append(args, *notificationType)
		argPos++
	}

	query += " ORDER BY created_at DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, limit)
		argPos++
		if offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argPos)
			args = append(args, offset)
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*Notification
	for rows.Next() {
		var notification Notification
		var metadataJSON []byte
		var entityID sql.NullString

		err := rows.Scan(
			&notification.ID,
			&notification.TenantID,
			&notification.Type,
			&notification.Title,
			&notification.Message,
			&notification.EntityType,
			&entityID,
			&notification.Priority,
			&notification.IsRead,
			&notification.ReadAt,
			&metadataJSON,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if entityID.Valid {
			if id, err := uuid.Parse(entityID.String); err == nil {
				notification.EntityID = &id
			}
		}

		if len(metadataJSON) > 0 {
			json.Unmarshal(metadataJSON, &notification.Metadata)
		}

		notifications = append(notifications, &notification)
	}

	return notifications, nil
}

// MarkAsRead marks a notification as read
func (r *Repository) MarkAsRead(ctx context.Context, schema string, id uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.notifications
		SET is_read = true, read_at = $1
		WHERE id = $2
	`, schema)

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// MarkAllAsRead marks all notifications for a tenant as read
func (r *Repository) MarkAllAsRead(ctx context.Context, schema string, tenantID uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.notifications
		SET is_read = true, read_at = $1
		WHERE tenant_id = $2 AND is_read = false
	`, schema)

	_, err := r.db.ExecContext(ctx, query, time.Now(), tenantID)
	return err
}

// GetUnreadCount gets the count of unread notifications
func (r *Repository) GetUnreadCount(ctx context.Context, schema string, tenantID uuid.UUID) (int, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) FROM %s.notifications
		WHERE tenant_id = $1 AND is_read = false
	`, schema)

	var count int
	err := r.db.GetContext(ctx, &count, query, tenantID)
	return count, err
}

// Delete deletes a notification
func (r *Repository) Delete(ctx context.Context, schema string, id uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s.notifications WHERE id = $1`, schema)
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetPurchasesDueSoon gets purchases that are due within the specified days
func (r *Repository) GetPurchasesDueSoon(ctx context.Context, schema string, daysBefore int) ([]*PurchaseDueInfo, error) {
	query := fmt.Sprintf(`
		SELECT 
			p.id,
			p.purchase_number,
			p.supplier_id,
			s.name as supplier_name,
			p.payment_due_date,
			p.total,
			COALESCE(SUM(spa.amount), 0) as paid_amount,
			(p.total - COALESCE(SUM(spa.amount), 0)) as balance,
			EXTRACT(DAY FROM (p.payment_due_date - CURRENT_DATE))::int as days_until_due
		FROM %s.purchases p
		INNER JOIN %s.suppliers s ON p.supplier_id = s.id
		LEFT JOIN %s.supplier_payment_allocations spa ON p.id = spa.purchase_id
		WHERE p.payment_method = 'credit'
		  AND p.payment_status IN ('pending', 'partial')
		  AND p.payment_due_date IS NOT NULL
		  AND p.payment_due_date >= CURRENT_DATE
		  AND p.payment_due_date <= CURRENT_DATE + INTERVAL '%d days'
		GROUP BY p.id, p.purchase_number, p.supplier_id, s.name, p.payment_due_date, p.total
		HAVING (p.total - COALESCE(SUM(spa.amount), 0)) > 0
		ORDER BY p.payment_due_date ASC
	`, schema, schema, schema, daysBefore)

	type row struct {
		ID            uuid.UUID  `db:"id"`
		PurchaseNumber string    `db:"purchase_number"`
		SupplierID     uuid.UUID `db:"supplier_id"`
		SupplierName   string    `db:"supplier_name"`
		PaymentDueDate time.Time `db:"payment_due_date"`
		Total          float64   `db:"total"`
		PaidAmount     float64   `db:"paid_amount"`
		Balance        float64   `db:"balance"`
		DaysUntilDue   int       `db:"days_until_due"`
	}

	var rows []row
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, err
	}

	result := make([]*PurchaseDueInfo, len(rows))
	for i, r := range rows {
		result[i] = &PurchaseDueInfo{
			PurchaseID:     r.ID,
			PurchaseNumber: r.PurchaseNumber,
			SupplierID:     r.SupplierID,
			SupplierName:   r.SupplierName,
			PaymentDueDate: r.PaymentDueDate,
			Total:          r.Total,
			PaidAmount:     r.PaidAmount,
			Balance:        r.Balance,
			DaysUntilDue:   r.DaysUntilDue,
		}
	}

	return result, nil
}

// GetPurchasesOverdue gets purchases that are overdue
func (r *Repository) GetPurchasesOverdue(ctx context.Context, schema string) ([]*PurchaseDueInfo, error) {
	query := fmt.Sprintf(`
		SELECT 
			p.id,
			p.purchase_number,
			p.supplier_id,
			s.name as supplier_name,
			p.payment_due_date,
			p.total,
			COALESCE(SUM(spa.amount), 0) as paid_amount,
			(p.total - COALESCE(SUM(spa.amount), 0)) as balance,
			EXTRACT(DAY FROM (CURRENT_DATE - p.payment_due_date))::int as days_overdue
		FROM %s.purchases p
		INNER JOIN %s.suppliers s ON p.supplier_id = s.id
		LEFT JOIN %s.supplier_payment_allocations spa ON p.id = spa.purchase_id
		WHERE p.payment_method = 'credit'
		  AND p.payment_status IN ('pending', 'partial')
		  AND p.payment_due_date IS NOT NULL
		  AND p.payment_due_date < CURRENT_DATE
		GROUP BY p.id, p.purchase_number, p.supplier_id, s.name, p.payment_due_date, p.total
		HAVING (p.total - COALESCE(SUM(spa.amount), 0)) > 0
		ORDER BY p.payment_due_date ASC
	`, schema, schema, schema)

	type row struct {
		ID            uuid.UUID  `db:"id"`
		PurchaseNumber string    `db:"purchase_number"`
		SupplierID     uuid.UUID `db:"supplier_id"`
		SupplierName   string    `db:"supplier_name"`
		PaymentDueDate time.Time `db:"payment_due_date"`
		Total          float64   `db:"total"`
		PaidAmount     float64   `db:"paid_amount"`
		Balance        float64   `db:"balance"`
		DaysOverdue    int       `db:"days_overdue"`
	}

	var rows []row
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, err
	}

	result := make([]*PurchaseDueInfo, len(rows))
	for i, r := range rows {
		result[i] = &PurchaseDueInfo{
			PurchaseID:     r.ID,
			PurchaseNumber: r.PurchaseNumber,
			SupplierID:     r.SupplierID,
			SupplierName:   r.SupplierName,
			PaymentDueDate: r.PaymentDueDate,
			Total:          r.Total,
			PaidAmount:     r.PaidAmount,
			Balance:        r.Balance,
			DaysOverdue:    r.DaysOverdue,
		}
	}

	return result, nil
}

// PurchaseDueInfo represents purchase information for notifications
type PurchaseDueInfo struct {
	PurchaseID     uuid.UUID
	PurchaseNumber string
	SupplierID     uuid.UUID
	SupplierName   string
	PaymentDueDate time.Time
	Total          float64
	PaidAmount     float64
	Balance        float64
	DaysUntilDue   int
	DaysOverdue    int
}

