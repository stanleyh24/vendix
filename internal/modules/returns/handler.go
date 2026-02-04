package returns

import (
	"time"

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

// RegisterRoutes registers returns routes
func RegisterRoutes(router fiber.Router, db *database.DB, cfg *config.Config) {
	h := NewHandler(db, cfg)
	returns := router.Group("/returns")

	// Returns routes
	returns.Get("/", h.List)
	returns.Post("/", h.Create)
	returns.Get("/:id", h.GetByID)
	returns.Patch("/:id/approve", h.Approve)
	returns.Patch("/:id/reject", h.Reject)
	returns.Post("/:id/process", h.Process)
}

// List lists all returns
// @Summary List returns
// @Tags returns
// @Produce json
// @Param status query string false "Filter by status"
// @Param sale_id query string false "Filter by sale ID"
// @Param invoice_id query string false "Filter by invoice ID"
// @Param customer_id query string false "Filter by customer ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} Return
// @Router /api/v1/tenant/returns [get]
func (h *Handler) List(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	if schema == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant schema not found in context",
		})
	}

	// Parse filters
	filters := &ReturnFilters{}

	if status := c.Query("status"); status != "" {
		filters.Status = &status
	}

	if saleID := c.Query("sale_id"); saleID != "" {
		id, err := uuid.Parse(saleID)
		if err == nil {
			filters.SaleID = &id
		}
	}

	if invoiceID := c.Query("invoice_id"); invoiceID != "" {
		id, err := uuid.Parse(invoiceID)
		if err == nil {
			filters.InvoiceID = &id
		}
	}

	if customerID := c.Query("customer_id"); customerID != "" {
		id, err := uuid.Parse(customerID)
		if err == nil {
			filters.CustomerID = &id
		}
	}

	if startDate := c.Query("start_date"); startDate != "" {
		// Parse date
		if date, err := time.Parse("2006-01-02", startDate); err == nil {
			filters.StartDate = &date
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		// Parse date
		if date, err := time.Parse("2006-01-02", endDate); err == nil {
			filters.EndDate = &date
		}
	}

	returns, err := h.service.List(c.Context(), schema, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(returns)
}

// Create creates a new return
// @Summary Create return
// @Tags returns
// @Accept json
// @Produce json
// @Param return body CreateReturnRequest true "Return data"
// @Success 201 {object} Return
// @Router /api/v1/tenant/returns [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	if schema == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant schema not found in context",
		})
	}

	var req CreateReturnRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user ID from context
	var userID *uuid.UUID
	if userIDStr, ok := c.Locals("user_id").(string); ok && userIDStr != "" {
		if id, err := uuid.Parse(userIDStr); err == nil {
			userID = &id
		}
	}

	ret, err := h.service.Create(c.Context(), schema, &req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(ret)
}

// GetByID retrieves a return by ID
// @Summary Get return
// @Tags returns
// @Produce json
// @Param id path string true "Return ID"
// @Success 200 {object} Return
// @Router /api/v1/tenant/returns/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	returnID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid return ID",
		})
	}

	ret, err := h.service.GetByID(c.Context(), schema, returnID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Return not found",
		})
	}

	return c.JSON(ret)
}

// Approve approves a return
// @Summary Approve return
// @Tags returns
// @Success 200 {object} map[string]string
// @Router /api/v1/tenant/returns/{id}/approve [patch]
func (h *Handler) Approve(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	returnID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid return ID",
		})
	}

	// Get user ID from context
	var userID uuid.UUID
	if userIDStr, ok := c.Locals("user_id").(string); ok && userIDStr != "" {
		if id, err := uuid.Parse(userIDStr); err == nil {
			userID = id
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	if err := h.service.Approve(c.Context(), schema, returnID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Return approved successfully",
	})
}

// Reject rejects a return
// @Summary Reject return
// @Tags returns
// @Accept json
// @Produce json
// @Param id path string true "Return ID"
// @Param request body RejectReturnRequest true "Rejection reason"
// @Success 200 {object} map[string]string
// @Router /api/v1/tenant/returns/{id}/reject [patch]
func (h *Handler) Reject(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	returnID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid return ID",
		})
	}

	var req RejectReturnRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user ID from context
	var userID uuid.UUID
	if userIDStr, ok := c.Locals("user_id").(string); ok && userIDStr != "" {
		if id, err := uuid.Parse(userIDStr); err == nil {
			userID = id
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	if err := h.service.Reject(c.Context(), schema, returnID, &req, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Return rejected successfully",
	})
}

// Process processes an approved return
// @Summary Process return
// @Tags returns
// @Accept json
// @Produce json
// @Param id path string true "Return ID"
// @Param request body ProcessReturnRequest true "Process request"
// @Success 200 {object} map[string]string
// @Router /api/v1/tenant/returns/{id}/process [post]
func (h *Handler) Process(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	returnID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid return ID",
		})
	}

	var req ProcessReturnRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user ID from context
	var userID *uuid.UUID
	if userIDStr, ok := c.Locals("user_id").(string); ok && userIDStr != "" {
		if id, err := uuid.Parse(userIDStr); err == nil {
			userID = &id
		}
	}

	if err := h.service.Process(c.Context(), schema, returnID, &req, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Return processed successfully",
	})
}

