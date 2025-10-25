package payments

import (
	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(db *database.DB, cfg *config.Config) *Handler {
	return &Handler{service: NewService(db, cfg)}
}

func RegisterRoutes(router fiber.Router, db *database.DB, cfg *config.Config) {
	h := NewHandler(db, cfg)
	payments := router.Group("/payments")
	payments.Post("/", h.Create)
	payments.Get("/", h.List)
	payments.Get("/:id", h.Get)
	payments.Put("/:id", h.Update)
	payments.Delete("/:id", h.Delete)

	// Alias for expenses
	expenses := router.Group("/expenses")
	expenses.Post("/", h.Create)
	expenses.Get("/", h.List)
	expenses.Get("/:id", h.Get)
	expenses.Put("/:id", h.Update)
	expenses.Delete("/:id", h.Delete)
}

// Create creates a new payment or expense
// @Summary Create payment
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body CreatePaymentRequest true "Payment data"
// @Success 201 {object} Payment
// @Router /api/v1/tenant/payments [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	var req CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	payment, err := h.service.Create(c.Context(), schema, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(payment)
}

// List lists all payments/expenses
// @Summary List payments
// @Tags payments
// @Produce json
// @Success 200 {array} Payment
// @Router /api/v1/tenant/payments [get]
func (h *Handler) List(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	payments, err := h.service.List(c.Context(), schema)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(payments)
}

// Get retrieves a payment by ID
// @Summary Get payment
// @Tags payments
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} Payment
// @Router /api/v1/tenant/payments/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	payment, err := h.service.GetByID(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Payment not found",
		})
	}

	return c.JSON(payment)
}

// Update updates a payment
// @Summary Update payment
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Param payment body UpdatePaymentRequest true "Payment data"
// @Success 200 {object} Payment
// @Router /api/v1/tenant/payments/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	var req UpdatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	payment, err := h.service.Update(c.Context(), schema, id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(payment)
}

// Delete deletes a payment
// @Summary Delete payment
// @Tags payments
// @Param id path string true "Payment ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/tenant/payments/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	if err := h.service.Delete(c.Context(), schema, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "Payment deleted successfully"})
}
