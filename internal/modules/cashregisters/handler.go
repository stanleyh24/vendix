package cashregisters

import (
	"time"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(db *database.DB, cfg *config.Config) *Handler {
	return &Handler{
		service: NewService(db, cfg),
	}
}

func RegisterRoutes(router fiber.Router, db *database.DB, cfg *config.Config) {
	h := NewHandler(db, cfg)
	cash := router.Group("/cash-registers")

	// Sessions routes - register all session routes in a separate group first
	// This ensures they are registered before any dynamic /:id routes
	sessions := cash.Group("/sessions")
	sessions.Post("/open", h.OpenSession)
	sessions.Get("", h.ListSessions)
	sessions.Get("/:id/summary", h.GetSessionSummary)
	sessions.Post("/:id/close", h.CloseSession)

	// Cash Registers CRUD - base routes before parameterized routes
	cash.Get("", h.ListCashRegisters)
	cash.Post("", h.CreateCashRegister)

	// Dynamic routes with specific suffixes must come before generic /:id
	cash.Get("/:id/open-session", h.GetOpenSession)

	// Generic dynamic route last
	cash.Get("/:id", h.GetCashRegisterByID)
}

// ListCashRegisters lists all cash registers
// @Summary List cash registers
// @Tags cash-registers
// @Produce json
// @Param active_only query bool false "Filter by active status"
// @Success 200 {array} CashRegister
// @Router /api/v1/tenant/cash-registers [get]
func (h *Handler) ListCashRegisters(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	if schema == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant schema not found in context",
		})
	}

	activeOnly := c.Query("active_only", "false") == "true"

	registers, err := h.service.ListCashRegisters(c.Context(), schema, activeOnly)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(registers)
}

// CreateCashRegister creates a new cash register
// @Summary Create cash register
// @Tags cash-registers
// @Accept json
// @Produce json
// @Param request body CreateCashRegisterRequest true "Cash register data"
// @Success 201 {object} CashRegister
// @Router /api/v1/tenant/cash-registers [post]
func (h *Handler) CreateCashRegister(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	if schema == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant schema not found in context",
		})
	}

	// Get user ID
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req CreateCashRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	register, err := h.service.CreateCashRegister(c.Context(), schema, userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(register)
}

// GetCashRegisterByID gets a cash register by ID
// @Summary Get cash register
// @Tags cash-registers
// @Produce json
// @Param id path string true "Cash register ID"
// @Success 200 {object} CashRegister
// @Router /api/v1/tenant/cash-registers/{id} [get]
func (h *Handler) GetCashRegisterByID(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	if schema == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant schema not found in context",
		})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid cash register ID",
		})
	}

	register, err := h.service.GetCashRegisterByID(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(register)
}

// OpenSession opens a new cash register session
// @Summary Open cash register session
// @Tags cash-registers
// @Accept json
// @Produce json
// @Param request body OpenSessionRequest true "Session data"
// @Success 201 {object} CashRegisterSession
// @Router /api/v1/tenant/cash-registers/sessions/open [post]
func (h *Handler) OpenSession(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	if schema == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant schema not found in context",
		})
	}

	// Get user ID
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req OpenSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	session, err := h.service.OpenSession(c.Context(), schema, userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(session)
}

// CloseSession closes a cash register session and generates journal entry
// @Summary Close cash register session
// @Tags cash-registers
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body CloseSessionRequest true "Close session data"
// @Success 200 {object} SessionSummary
// @Router /api/v1/tenant/cash-registers/sessions/{id}/close [post]
func (h *Handler) CloseSession(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	if schema == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant schema not found in context",
		})
	}

	// Get user ID
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	sessionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid session ID",
		})
	}

	var req CloseSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	summary, err := h.service.CloseSession(c.Context(), schema, sessionID, userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(summary)
}

// GetSessionSummary gets the summary of a session
// @Summary Get session summary
// @Tags cash-registers
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} SessionSummary
// @Router /api/v1/tenant/cash-registers/sessions/{id}/summary [get]
func (h *Handler) GetSessionSummary(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	if schema == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant schema not found in context",
		})
	}

	sessionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid session ID",
		})
	}

	summary, err := h.service.GetSessionSummary(c.Context(), schema, sessionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(summary)
}

// GetOpenSession gets the open session for a cash register
// @Summary Get open session
// @Tags cash-registers
// @Produce json
// @Param id path string true "Cash register ID"
// @Success 200 {object} CashRegisterSession
// @Router /api/v1/tenant/cash-registers/{id}/open-session [get]
func (h *Handler) GetOpenSession(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	if schema == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant schema not found in context",
		})
	}

	cashRegisterID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid cash register ID",
		})
	}

	session, err := h.service.GetOpenSession(c.Context(), schema, cashRegisterID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(session)
}

// ListSessions lists sessions with filters
// @Summary List sessions
// @Tags cash-registers
// @Produce json
// @Param cash_register_id query string false "Filter by cash register ID"
// @Param status query string false "Filter by status (open, closed, cancelled)"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} CashRegisterSession
// @Router /api/v1/tenant/cash-registers/sessions [get]
func (h *Handler) ListSessions(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	if schema == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant schema not found in context",
		})
	}

	var cashRegisterID *uuid.UUID
	if idStr := c.Query("cash_register_id"); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			cashRegisterID = &id
		}
	}

	var status *string
	if s := c.Query("status"); s != "" {
		status = &s
	}

	var startDate, endDate *time.Time
	if sd := c.Query("start_date"); sd != "" {
		if t, err := time.Parse("2006-01-02", sd); err == nil {
			startDate = &t
		}
	}
	if ed := c.Query("end_date"); ed != "" {
		if t, err := time.Parse("2006-01-02", ed); err == nil {
			endDate = &t
		}
	}

	sessions, err := h.service.ListSessions(c.Context(), schema, cashRegisterID, status, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(sessions)
}
