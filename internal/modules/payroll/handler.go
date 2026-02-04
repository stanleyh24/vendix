package payroll

import (
	"github.com/stanleyh24/vendix/internal/config"
	"github.com/stanleyh24/vendix/internal/database"
	"github.com/stanleyh24/vendix/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(db *database.DB, cfg *config.Config) *Handler {
	return &Handler{
		service: NewService(db, cfg),
	}
}

// RegisterRoutes registers payroll routes
func RegisterRoutes(router fiber.Router, db *database.DB, cfg *config.Config) {
	h := NewHandler(db, cfg)

	payroll := router.Group("/payroll")

	// Period routes
	payroll.Post("/periods", h.CreatePeriod)
	payroll.Get("/periods", h.ListPeriods)
	payroll.Get("/periods/:id", h.GetPeriod)
	payroll.Put("/periods/:id", h.UpdatePeriod)
	payroll.Delete("/periods/:id", h.DeletePeriod)

	// Entry routes
	payroll.Post("/periods/:periodId/entries", h.CreateEntry)
	payroll.Get("/periods/:periodId/entries", h.ListEntries)
	payroll.Put("/entries/:id", h.UpdateEntry)
	payroll.Delete("/entries/:id", h.DeleteEntry)
}

// CreatePeriod creates a new payroll period
// @Summary Create payroll period
// @Tags payroll
// @Accept json
// @Produce json
// @Param period body CreatePayrollPeriodRequest true "Payroll period data"
// @Success 201 {object} PayrollPeriodResponse
// @Security BearerAuth
// @Router /api/v1/tenant/payroll/periods [post]
func (h *Handler) CreatePeriod(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	userID := middleware.GetUserID(c)

	var req CreatePayrollPeriodRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	period, err := h.service.CreatePeriod(c.Context(), schema, userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(period)
}

// ListPeriods lists all payroll periods
// @Summary List payroll periods
// @Tags payroll
// @Produce json
// @Success 200 {array} PayrollPeriodResponse
// @Security BearerAuth
// @Router /api/v1/tenant/payroll/periods [get]
func (h *Handler) ListPeriods(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	periods, err := h.service.ListPeriods(c.Context(), schema)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(periods)
}

// GetPeriod retrieves a payroll period by ID
// @Summary Get payroll period
// @Tags payroll
// @Produce json
// @Param id path string true "Period ID"
// @Success 200 {object} PayrollPeriodResponse
// @Security BearerAuth
// @Router /api/v1/tenant/payroll/periods/{id} [get]
func (h *Handler) GetPeriod(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	period, err := h.service.GetPeriodByID(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Payroll period not found",
		})
	}

	return c.JSON(period)
}

// UpdatePeriod updates a payroll period
// @Summary Update payroll period
// @Tags payroll
// @Accept json
// @Produce json
// @Param id path string true "Period ID"
// @Param period body UpdatePayrollPeriodRequest true "Payroll period data"
// @Success 200 {object} PayrollPeriodResponse
// @Security BearerAuth
// @Router /api/v1/tenant/payroll/periods/{id} [put]
func (h *Handler) UpdatePeriod(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	var req UpdatePayrollPeriodRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	period, err := h.service.UpdatePeriod(c.Context(), schema, id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(period)
}

// DeletePeriod deletes a payroll period
// @Summary Delete payroll period
// @Tags payroll
// @Param id path string true "Period ID"
// @Success 204
// @Security BearerAuth
// @Router /api/v1/tenant/payroll/periods/{id} [delete]
func (h *Handler) DeletePeriod(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	if err := h.service.DeletePeriod(c.Context(), schema, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// CreateEntry creates a new payroll entry
// @Summary Create payroll entry
// @Tags payroll
// @Accept json
// @Produce json
// @Param periodId path string true "Period ID"
// @Param entry body CreatePayrollEntryRequest true "Payroll entry data"
// @Success 201 {object} PayrollEntryResponse
// @Security BearerAuth
// @Router /api/v1/tenant/payroll/periods/{periodId}/entries [post]
func (h *Handler) CreateEntry(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	periodID := c.Params("periodId")

	var req CreatePayrollEntryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	entry, err := h.service.CreateEntry(c.Context(), schema, periodID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(entry)
}

// ListEntries lists all entries for a payroll period
// @Summary List payroll entries
// @Tags payroll
// @Produce json
// @Param periodId path string true "Period ID"
// @Success 200 {array} PayrollEntryResponse
// @Security BearerAuth
// @Router /api/v1/tenant/payroll/periods/{periodId}/entries [get]
func (h *Handler) ListEntries(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	periodID := c.Params("periodId")

	entries, err := h.service.ListEntriesByPeriod(c.Context(), schema, periodID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(entries)
}

// UpdateEntry updates a payroll entry
// @Summary Update payroll entry
// @Tags payroll
// @Accept json
// @Produce json
// @Param id path string true "Entry ID"
// @Param entry body UpdatePayrollEntryRequest true "Payroll entry data"
// @Success 200 {object} PayrollEntryResponse
// @Security BearerAuth
// @Router /api/v1/tenant/payroll/entries/{id} [put]
func (h *Handler) UpdateEntry(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	var req UpdatePayrollEntryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	entry, err := h.service.UpdateEntry(c.Context(), schema, id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(entry)
}

// DeleteEntry deletes a payroll entry
// @Summary Delete payroll entry
// @Tags payroll
// @Param id path string true "Entry ID"
// @Success 204
// @Security BearerAuth
// @Router /api/v1/tenant/payroll/entries/{id} [delete]
func (h *Handler) DeleteEntry(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	if err := h.service.DeleteEntry(c.Context(), schema, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

