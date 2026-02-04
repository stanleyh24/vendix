package tenantconfig

import (
	"github.com/stanleyh24/vendix/internal/database"
	"github.com/stanleyh24/vendix/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

// NewHandler creates a new tenant config handler
func NewHandler(db *database.DB) *Handler {
	return &Handler{
		service: NewService(db),
	}
}

// RegisterRoutes registers tenant config routes (requires tenant middleware)
func RegisterRoutes(router fiber.Router, db *database.DB) {
	h := NewHandler(db)

	config := router.Group("/config")
	config.Get("/", h.Get)
	config.Put("/", h.Update)
}

// Get retrieves the tenant configuration
// @Summary Get tenant configuration
// @Description Get the configuration settings for the current tenant
// @Tags config
// @Produce json
// @Success 200 {object} TenantConfig
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tenant/config [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	if schema == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant schema not found in context",
		})
	}

	config, err := h.service.Get(c.Context(), schema)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve configuration",
		})
	}

	return c.JSON(config)
}

// Update updates the tenant configuration
// @Summary Update tenant configuration
// @Description Update the configuration settings for the current tenant
// @Tags config
// @Accept json
// @Produce json
// @Param config body UpdateTenantConfigRequest true "Configuration data"
// @Success 200 {object} TenantConfig
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tenant/config [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	if schema == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant schema not found in context",
		})
	}

	var req UpdateTenantConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	config, err := h.service.Update(c.Context(), schema, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update configuration",
		})
	}

	return c.JSON(config)
}
