package customers

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

// RegisterRoutes registers customer routes
func RegisterRoutes(router fiber.Router, db *database.DB, cfg *config.Config) {
	h := NewHandler(db, cfg)

	customers := router.Group("/customers")

	customers.Post("/", h.Create)
	customers.Get("/", h.List)
	customers.Get("/:id", h.Get)
	customers.Put("/:id", h.Update)
	customers.Delete("/:id", h.Delete)
}

// Create creates a new customer
// @Summary Create customer
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body CreateCustomerRequest true "Customer data"
// @Success 201 {object} CustomerResponse
// @Security BearerAuth
// @Router /api/v1/tenant/customers [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	userID := middleware.GetUserID(c)

	var req CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	customer, err := h.service.Create(c.Context(), schema, userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(customer)
}

// List lists all customers
// @Summary List customers
// @Tags customers
// @Produce json
// @Success 200 {array} CustomerResponse
// @Security BearerAuth
// @Router /api/v1/tenant/customers [get]
func (h *Handler) List(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	customers, err := h.service.List(c.Context(), schema)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(customers)
}

// Get retrieves a customer by ID
// @Summary Get customer
// @Tags customers
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} CustomerResponse
// @Security BearerAuth
// @Router /api/v1/tenant/customers/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	customer, err := h.service.GetByID(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Customer not found",
		})
	}

	return c.JSON(customer)
}

// Update updates a customer
// @Summary Update customer
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body UpdateCustomerRequest true "Customer data"
// @Success 200 {object} CustomerResponse
// @Security BearerAuth
// @Router /api/v1/tenant/customers/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	var req UpdateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	customer, err := h.service.Update(c.Context(), schema, id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(customer)
}

// Delete deletes a customer
// @Summary Delete customer
// @Tags customers
// @Param id path string true "Customer ID"
// @Success 204
// @Security BearerAuth
// @Router /api/v1/tenant/customers/{id} [delete]
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
