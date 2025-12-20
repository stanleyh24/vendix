package invoices

import (
	"bytes"
	"fmt"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf"
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
	// Specific routes must be registered before parameterized routes
	invoices.Get("/:id/pdf", h.DownloadPDF)
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

// DownloadPDF generates and returns a simple PDF for the invoice
func (h *Handler) DownloadPDF(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	inv, err := h.service.GetByID(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invoice not found"})
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Factura "+inv.InvoiceNumber)
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Cliente: %s", inv.CustomerName))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Fecha: %s", inv.IssueDate.Format("2006-01-02")))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Vence: %s", inv.DueDate.Format("2006-01-02")))
	pdf.Ln(10)

	// Table header
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(100, 8, "Descripción", "B", 0, "L", false, 0, "")
	pdf.CellFormat(20, 8, "Cant.", "B", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, "P. Unit.", "B", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, "Total", "B", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	for _, line := range inv.Lines {
		pdf.CellFormat(100, 7, line.Description, "", 0, "L", false, 0, "")
		pdf.CellFormat(20, 7, fmt.Sprintf("%.2f", line.Quantity), "", 0, "R", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", line.UnitPrice), "", 0, "R", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", line.LineTotal), "", 1, "R", false, 0, "")
	}

	pdf.Ln(4)
	pdf.Cell(0, 0, "")
	pdf.Ln(2)
	pdf.CellFormat(150, 7, "Subtotal:", "", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", inv.Subtotal), "", 1, "R", false, 0, "")
	pdf.CellFormat(150, 7, "ITBIS:", "", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", inv.TaxAmount), "", 1, "R", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(150, 8, "TOTAL:", "", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, fmt.Sprintf("%.2f", inv.Total), "", 1, "R", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate PDF"})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=invoice-%s.pdf", inv.InvoiceNumber))
	return c.SendStream(bytes.NewReader(buf.Bytes()))
}
