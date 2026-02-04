package tenants

import (
	"github.com/stanleyh24/vendix/internal/config"
	"github.com/stanleyh24/vendix/internal/database"

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

// RegisterRoutes registers tenant management routes
func RegisterRoutes(router fiber.Router, db *database.DB, cfg *config.Config) {
	h := NewHandler(db, cfg)

	tenants := router.Group("/tenants")

	tenants.Post("/", h.Create)
	tenants.Get("/", h.List)
	tenants.Get("/:id", h.Get)
	tenants.Put("/:id", h.Update)
	tenants.Delete("/:id", h.Delete)
	tenants.Post("/:id/init", h.InitializeSchema)
}

// Create creates a new tenant
// @Summary Create tenant
// @Tags tenants
// @Accept json
// @Produce json
// @Param tenant body CreateTenantRequest true "Tenant data"
// @Success 201 {object} TenantResponse
// @Router /api/v1/public/tenants [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	tenant, err := h.service.Create(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(tenant)
}

// List lists all tenants
// @Summary List tenants
// @Tags tenants
// @Produce json
// @Success 200 {array} TenantResponse
// @Router /api/v1/public/tenants [get]
func (h *Handler) List(c *fiber.Ctx) error {
	tenants, err := h.service.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(tenants)
}

// Get retrieves a tenant by ID
// @Summary Get tenant
// @Tags tenants
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} TenantResponse
// @Router /api/v1/public/tenants/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	id := c.Params("id")

	tenant, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tenant not found",
		})
	}

	return c.JSON(tenant)
}

// Update updates a tenant
// @Summary Update tenant
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path string true "Tenant ID"
// @Param tenant body UpdateTenantRequest true "Tenant data"
// @Success 200 {object} TenantResponse
// @Router /api/v1/public/tenants/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req UpdateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	tenant, err := h.service.Update(c.Context(), id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(tenant)
}

// Delete soft deletes a tenant
// @Summary Delete tenant
// @Tags tenants
// @Param id path string true "Tenant ID"
// @Success 204
// @Router /api/v1/public/tenants/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// InitializeSchema initializes the tenant's database schema
// @Summary Initialize tenant schema
// @Tags tenants
// @Param id path string true "Tenant ID"
// @Success 200
// @Router /api/v1/public/tenants/{id}/init [post]
func (h *Handler) InitializeSchema(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.InitializeSchema(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Tenant schema initialized successfully",
	})
}
