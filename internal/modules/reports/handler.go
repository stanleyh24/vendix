package reports

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
	return &Handler{service: NewService(db, cfg)}
}

func RegisterRoutes(router fiber.Router, db *database.DB, cfg *config.Config) {
	h := NewHandler(db, cfg)
	reports := router.Group("/reports")
	reports.Get("/dashboard", h.Dashboard)
	reports.Get("/sales", h.SalesReport)
	reports.Get("/taxes", h.TaxReport)
	reports.Get("/inventory-valuation", h.InventoryValuation)
	reports.Get("/accounts-receivable", h.AccountsReceivable)
	reports.Get("/accounts-payable", h.AccountsPayable)
	reports.Get("/returns", h.ReturnsReport)
}

// SalesReport generates a comprehensive sales report
// @Summary Get sales report
// @Tags reports
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD), defaults to 30 days ago"
// @Param end_date query string false "End date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} SalesReport
// @Router /api/v1/tenant/reports/sales [get]
func (h *Handler) SalesReport(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var startDatePtr, endDatePtr *string
	if startDate != "" {
		startDatePtr = &startDate
	}
	if endDate != "" {
		endDatePtr = &endDate
	}

	report, err := h.service.GetSalesReport(c.Context(), schema, startDatePtr, endDatePtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(report)
}

// TaxReport generates a comprehensive tax report (ITBIS)
// @Summary Get tax report
// @Tags reports
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD), defaults to 30 days ago"
// @Param end_date query string false "End date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} TaxReport
// @Router /api/v1/tenant/reports/taxes [get]
func (h *Handler) TaxReport(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var startDatePtr, endDatePtr *string
	if startDate != "" {
		startDatePtr = &startDate
	}
	if endDate != "" {
		endDatePtr = &endDate
	}

	report, err := h.service.GetTaxReport(c.Context(), schema, startDatePtr, endDatePtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(report)
}

// InventoryValuation generates an inventory valuation report
// @Summary Get inventory valuation report
// @Tags reports
// @Produce json
// @Param include_zero_stock query bool false "Include products with zero stock"
// @Success 200 {object} InventoryValuationReport
// @Router /api/v1/tenant/reports/inventory-valuation [get]
func (h *Handler) InventoryValuation(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	// Get query parameter for including zero stock items
	includeZeroStock := c.Query("include_zero_stock") == "true"

	report, err := h.service.GetInventoryValuation(c.Context(), schema, includeZeroStock)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(report)
}

// AccountsReceivable generates an accounts receivable aging report
// @Summary Get accounts receivable report
// @Tags reports
// @Produce json
// @Param as_of_date query string false "Report date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} AccountsReceivableReport
// @Router /api/v1/tenant/reports/accounts-receivable [get]
func (h *Handler) AccountsReceivable(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	asOfDate := c.Query("as_of_date")

	var asOfDatePtr *string
	if asOfDate != "" {
		asOfDatePtr = &asOfDate
	}

	report, err := h.service.GetAccountsReceivableReport(c.Context(), schema, asOfDatePtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(report)
}

// AccountsPayable generates an accounts payable aging report
// @Summary Get accounts payable report
// @Tags reports
// @Produce json
// @Param as_of_date query string false "Report date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} AccountsPayableReport
// @Router /api/v1/tenant/reports/accounts-payable [get]
func (h *Handler) AccountsPayable(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	asOfDate := c.Query("as_of_date")

	var asOfDatePtr *string
	if asOfDate != "" {
		asOfDatePtr = &asOfDate
	}

	report, err := h.service.GetAccountsPayableReport(c.Context(), schema, asOfDatePtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(report)
}

// ReturnsReport generates a comprehensive returns report
// @Summary Get returns report
// @Tags reports
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD), defaults to 30 days ago"
// @Param end_date query string false "End date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} ReturnsReport
// @Router /api/v1/tenant/reports/returns [get]
func (h *Handler) ReturnsReport(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var startDatePtr, endDatePtr *string
	if startDate != "" {
		startDatePtr = &startDate
	}
	if endDate != "" {
		endDatePtr = &endDate
	}

	report, err := h.service.GetReturnsReport(c.Context(), schema, startDatePtr, endDatePtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(report)
}

// Dashboard generates comprehensive dashboard summary
// @Summary Get dashboard summary
// @Tags reports
// @Produce json
// @Success 200 {object} DashboardSummary
// @Router /api/v1/tenant/reports/dashboard [get]
func (h *Handler) Dashboard(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	summary, err := h.service.GetDashboardSummary(c.Context(), schema)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(summary)
}
