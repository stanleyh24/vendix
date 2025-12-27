package invoices

import (
	"fmt"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(db *database.DB, cfg *config.Config) (*Handler, error) {
	service, err := NewService(db, cfg)
	if err != nil {
		return nil, err
	}
	return &Handler{service: service}, nil
}

func RegisterRoutes(router fiber.Router, db *database.DB, cfg *config.Config) {
	h, err := NewHandler(db, cfg)
	if err != nil {
		// Log error but don't fail route registration
		// This allows the API to start even if storage is unavailable
		// In production, you might want to panic here
	}
	invoices := router.Group("/invoices")
	invoices.Post("/", h.Create)
	invoices.Get("/", h.List)
	// Specific routes must be registered before parameterized routes
	invoices.Get("/:id/pdf", h.DownloadPDF)
	invoices.Get("/:id/xml", h.DownloadXML)
	invoices.Post("/:id/send", h.SendInvoice)
	invoices.Post("/:id/cancel", h.CancelInvoice)
	invoices.Get("/:id", h.Get)
	invoices.Put("/:id", h.Update)
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

// DownloadPDF retrieves and returns the PDF for the invoice
func (h *Handler) DownloadPDF(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	invoiceID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid invoice ID"})
	}

	// Get PDF bytes from service (will generate if it doesn't exist)
	pdfBytes, err := h.service.GetPDFBytes(c.Context(), schema, invoiceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Get invoice for filename
	inv, err := h.service.GetByID(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invoice not found"})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=invoice-%s.pdf", inv.InvoiceNumber))
	return c.Send(pdfBytes)
}

// DownloadXML retrieves and returns the XML for the invoice
func (h *Handler) DownloadXML(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	invoiceID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid invoice ID"})
	}

	// Get XML bytes from service
	xmlBytes, err := h.service.GetXMLBytes(c.Context(), schema, invoiceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Get invoice for filename
	inv, err := h.service.GetByID(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invoice not found"})
	}

	c.Set("Content-Type", "application/xml")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=invoice-%s.xml", inv.InvoiceNumber))
	return c.Send(xmlBytes)
}
