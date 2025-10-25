package invoices

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
	invoices := router.Group("/invoices")
	invoices.Post("/", h.Create)
	invoices.Get("/", h.List)
	invoices.Get("/:id", h.Get)
	invoices.Put("/:id", h.Update)
	invoices.Post("/:id/send", h.SendInvoice)
	invoices.Post("/:id/cancel", h.CancelInvoice)
}

// Create creates a new invoice
// @Summary Create invoice
// @Tags invoices
// @Accept json
// @Produce json
// @Param invoice body CreateInvoiceRequest true "Invoice data"
// @Success 201 {object} Invoice
// @Router /api/v1/tenant/invoices [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	var req CreateInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	invoice, err := h.service.Create(c.Context(), schema, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(invoice)
}

// List lists all invoices
// @Summary List invoices
// @Tags invoices
// @Produce json
// @Param status query string false "Filter by status"
// @Success 200 {array} Invoice
// @Router /api/v1/tenant/invoices [get]
func (h *Handler) List(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	status := c.Query("status")

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	invoices, err := h.service.List(c.Context(), schema, statusPtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(invoices)
}

// Get retrieves an invoice by ID
// @Summary Get invoice
// @Tags invoices
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} Invoice
// @Router /api/v1/tenant/invoices/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	invoice, err := h.service.GetByID(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Invoice not found",
		})
	}

	return c.JSON(invoice)
}

// Update updates an invoice
// @Summary Update invoice
// @Tags invoices
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Param invoice body UpdateInvoiceRequest true "Invoice data"
// @Success 200 {object} Invoice
// @Router /api/v1/tenant/invoices/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	var req UpdateInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	invoice, err := h.service.Update(c.Context(), schema, id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(invoice)
}

// SendInvoice sends an invoice to DGII
// @Summary Send invoice
// @Tags invoices
// @Param id path string true "Invoice ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/tenant/invoices/{id}/send [post]
func (h *Handler) SendInvoice(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	if err := h.service.SendInvoice(c.Context(), schema, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "Invoice sent successfully"})
}

// CancelInvoice cancels an invoice
// @Summary Cancel invoice
// @Tags invoices
// @Param id path string true "Invoice ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/tenant/invoices/{id}/cancel [post]
func (h *Handler) CancelInvoice(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	if err := h.service.CancelInvoice(c.Context(), schema, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "Invoice cancelled successfully"})
}
