package creditnotes

import (
	"fmt"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/logger"
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
		// Log error and return early - routes won't be registered
		logger.Error("Failed to initialize credit notes handler", "error", err)
		return
	}
	creditNotes := router.Group("/credit-notes")
	creditNotes.Post("/", h.Create)
	creditNotes.Get("/", h.List)
	creditNotes.Get("/:id/pdf", h.DownloadPDF)
	creditNotes.Get("/:id/xml", h.DownloadXML)
	creditNotes.Post("/:id/send", h.SendCreditNote)
	creditNotes.Post("/:id/cancel", h.CancelCreditNote)
	creditNotes.Get("/:id", h.Get)
}

// Create creates a new credit note
func (h *Handler) Create(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	var req CreateCreditNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user ID from context if available
	var createdBy *uuid.UUID
	if userIDStr, ok := c.Locals("user_id").(string); ok && userIDStr != "" {
		if id, err := uuid.Parse(userIDStr); err == nil {
			createdBy = &id
		}
	}

	creditNote, err := h.service.Create(c.Context(), schema, &req, createdBy)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(creditNote)
}

// List lists all credit notes
func (h *Handler) List(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	filters := &CreditNoteFilters{}

	// Parse query parameters
	if status := c.Query("status"); status != "" {
		filters.Status = &status
	}
	if invoiceID := c.Query("invoice_id"); invoiceID != "" {
		if id, err := uuid.Parse(invoiceID); err == nil {
			filters.OriginalInvoiceID = &id
		}
	}
	if returnID := c.Query("return_id"); returnID != "" {
		if id, err := uuid.Parse(returnID); err == nil {
			filters.ReturnID = &id
		}
	}
	if customerID := c.Query("customer_id"); customerID != "" {
		if id, err := uuid.Parse(customerID); err == nil {
			filters.CustomerID = &id
		}
	}

	creditNotes, err := h.service.List(c.Context(), schema, filters)
	if err != nil {
		// Log error but return empty list instead of error to frontend
		logger.Error("Error listing credit notes", "error", err, "schema", schema)
		// Return empty array instead of error - no credit notes is not an error
		return c.JSON([]*CreditNote{})
	}

	// Always return an array, even if empty
	if creditNotes == nil {
		creditNotes = []*CreditNote{}
	}

	return c.JSON(creditNotes)
}

// Get retrieves a credit note by ID
func (h *Handler) Get(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid credit note ID",
		})
	}

	creditNote, err := h.service.GetByID(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Credit note not found",
		})
	}

	return c.JSON(creditNote)
}

// SendCreditNote sends a credit note to DGII
func (h *Handler) SendCreditNote(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid credit note ID",
		})
	}

	if err := h.service.SendCreditNote(c.Context(), schema, id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Credit note sent to DGII successfully",
	})
}

// CancelCreditNote cancels a credit note
func (h *Handler) CancelCreditNote(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid credit note ID",
		})
	}

	if err := h.service.repo.UpdateStatus(c.Context(), schema, id, CreditNoteStatusCancelled); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Credit note cancelled successfully",
	})
}

// DownloadPDF downloads the PDF of a credit note
func (h *Handler) DownloadPDF(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid credit note ID",
		})
	}

	pdfBytes, err := h.service.GeneratePDF(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to generate PDF: %v", err),
		})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=credit-note-%s.pdf", idParam))
	return c.Send(pdfBytes)
}

// DownloadXML downloads the XML of a credit note
func (h *Handler) DownloadXML(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid credit note ID",
		})
	}

	creditNote, err := h.service.GetByID(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Credit note not found",
		})
	}

	if creditNote.XMLURL == nil || *creditNote.XMLURL == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "XML not available: credit note has not been sent to DGII",
		})
	}

	// Extract object key from URL and download
	// For now, return error indicating XML needs to be downloaded from storage
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error":   "XML download not yet implemented. Use the XML URL from the credit note.",
		"xml_url": creditNote.XMLURL,
	})
}
