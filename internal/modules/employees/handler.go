package employees

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
	return &Handler{
		service: NewService(db, cfg),
	}
}

// RegisterRoutes registers employee routes
func RegisterRoutes(router fiber.Router, db *database.DB, cfg *config.Config) {
	h := NewHandler(db, cfg)

	employees := router.Group("/employees")

	employees.Post("/", h.Create)
	employees.Get("/", h.List)
	employees.Get("/active", h.ListActive)
	employees.Get("/:id", h.Get)
	employees.Put("/:id", h.Update)
	employees.Delete("/:id", h.Delete)
}

// Create creates a new employee
// @Summary Create employee
// @Tags employees
// @Accept json
// @Produce json
// @Param employee body CreateEmployeeRequest true "Employee data"
// @Success 201 {object} EmployeeResponse
// @Security BearerAuth
// @Router /api/v1/tenant/employees [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	userID := middleware.GetUserID(c)

	var req CreateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	employee, err := h.service.Create(c.Context(), schema, userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(employee)
}

// List lists all employees
// @Summary List employees
// @Tags employees
// @Produce json
// @Success 200 {array} EmployeeResponse
// @Security BearerAuth
// @Router /api/v1/tenant/employees [get]
func (h *Handler) List(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	employees, err := h.service.List(c.Context(), schema)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(employees)
}

// ListActive lists all active employees
// @Summary List active employees
// @Tags employees
// @Produce json
// @Success 200 {array} EmployeeResponse
// @Security BearerAuth
// @Router /api/v1/tenant/employees/active [get]
func (h *Handler) ListActive(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	employees, err := h.service.ListActive(c.Context(), schema)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(employees)
}

// Get retrieves an employee by ID
// @Summary Get employee
// @Tags employees
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} EmployeeResponse
// @Security BearerAuth
// @Router /api/v1/tenant/employees/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	employee, err := h.service.GetByID(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Employee not found",
		})
	}

	return c.JSON(employee)
}

// Update updates an employee
// @Summary Update employee
// @Tags employees
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Param employee body UpdateEmployeeRequest true "Employee data"
// @Success 200 {object} EmployeeResponse
// @Security BearerAuth
// @Router /api/v1/tenant/employees/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	var req UpdateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	employee, err := h.service.Update(c.Context(), schema, id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(employee)
}

// Delete deletes an employee
// @Summary Delete employee
// @Tags employees
// @Param id path string true "Employee ID"
// @Success 204
// @Security BearerAuth
// @Router /api/v1/tenant/employees/{id} [delete]
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

