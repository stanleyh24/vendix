package notifications

import (
	"github.com/stanleyh24/vendix/internal/config"
	"github.com/stanleyh24/vendix/internal/database"
	"github.com/stanleyh24/vendix/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(db *database.DB, cfg *config.Config) *Handler {
	return &Handler{service: NewService(db, cfg)}
}

// RegisterRoutes registers notification routes
func RegisterRoutes(router fiber.Router, db *database.DB, cfg *config.Config) {
	h := NewHandler(db, cfg)
	notifications := router.Group("/notifications")

	notifications.Get("/", h.List)
	notifications.Get("/unread-count", h.GetUnreadCount)
	notifications.Get("/:id", h.GetByID)
	notifications.Post("/:id/read", h.MarkAsRead)
	notifications.Post("/read-all", h.MarkAllAsRead)
	notifications.Post("/check-purchase-due", h.CheckPurchaseDue)
	notifications.Delete("/:id", h.Delete)
}

// List gets all notifications for the tenant
// @Summary List notifications
// @Tags notifications
// @Produce json
// @Param unread_only query bool false "Filter by unread only"
// @Param type query string false "Filter by notification type"
// @Param limit query int false "Limit results"
// @Param offset query int false "Offset results"
// @Success 200 {array} NotificationResponse
// @Router /api/v1/tenant/notifications [get]
func (h *Handler) List(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	tenantIDStr := middleware.GetTenantID(c)
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID",
		})
	}

	unreadOnly := c.Query("unread_only") == "true"
	notificationType := c.Query("type")
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	var typePtr *string
	if notificationType != "" {
		typePtr = &notificationType
	}

	notifications, err := h.service.List(c.Context(), schema, tenantID, unreadOnly, typePtr, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(notifications)
}

// GetUnreadCount gets the count of unread notifications
// @Summary Get unread notifications count
// @Tags notifications
// @Produce json
// @Success 200 {object} map[string]int
// @Router /api/v1/tenant/notifications/unread-count [get]
func (h *Handler) GetUnreadCount(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	tenantIDStr := middleware.GetTenantID(c)
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID",
		})
	}

	count, err := h.service.GetUnreadCount(c.Context(), schema, tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"count": count,
	})
}

// GetByID gets a notification by ID
// @Summary Get notification by ID
// @Tags notifications
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} NotificationResponse
// @Router /api/v1/tenant/notifications/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	notification, err := h.service.GetByID(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Notification not found",
		})
	}

	return c.JSON(notification)
}

// MarkAsRead marks a notification as read
// @Summary Mark notification as read
// @Tags notifications
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/tenant/notifications/{id}/read [post]
func (h *Handler) MarkAsRead(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	if err := h.service.MarkAsRead(c.Context(), schema, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Notification marked as read",
	})
}

// MarkAllAsRead marks all notifications as read
// @Summary Mark all notifications as read
// @Tags notifications
// @Success 200 {object} map[string]string
// @Router /api/v1/tenant/notifications/read-all [post]
func (h *Handler) MarkAllAsRead(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	tenantIDStr := middleware.GetTenantID(c)
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID",
		})
	}

	if err := h.service.MarkAllAsRead(c.Context(), schema, tenantID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "All notifications marked as read",
	})
}

// Delete deletes a notification
// @Summary Delete notification
// @Tags notifications
// @Param id path string true "Notification ID"
// @Success 204 "No Content"
// @Router /api/v1/tenant/notifications/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	if err := h.service.Delete(c.Context(), schema, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// CheckPurchaseDue manually triggers checking for purchases due soon or overdue
// @Summary Check purchase due dates
// @Tags notifications
// @Produce json
// @Param days_before query int false "Days before due date to check (default: 7)"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tenant/notifications/check-purchase-due [post]
func (h *Handler) CheckPurchaseDue(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	tenantIDStr := middleware.GetTenantID(c)
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID",
		})
	}

	daysBefore := c.QueryInt("days_before", 7)
	if daysBefore <= 0 {
		daysBefore = 7
	}

	if err := h.service.CheckPurchaseDueNotifications(c.Context(), schema, tenantID, daysBefore); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Purchase due dates checked successfully",
		"days_before": daysBefore,
	})
}

