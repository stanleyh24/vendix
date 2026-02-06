package tenantconfig

import (
	"github.com/stanleyh24/vendix/internal/config"
	"github.com/stanleyh24/vendix/internal/database"
	"github.com/stanleyh24/vendix/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

// NewHandler creates a new tenant config handler
func NewHandler(db *database.DB, cfg *config.Config) *Handler {
	return &Handler{
		service: NewService(db, cfg),
	}
}

// RegisterRoutes registers tenant config routes (requires tenant middleware)
func RegisterRoutes(router fiber.Router, db *database.DB, cfg *config.Config) {
	h := NewHandler(db, cfg)

	config := router.Group("/config")
	config.Get("/", h.Get)
	config.Put("/", h.Update)
	config.Post("/logo", h.UploadLogo)
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

// UploadLogo handles the logo file upload
// @Summary Upload tenant logo
// @Description Upload a logo image for the current tenant
// @Tags config
// @Accept multipart/form-data
// @Produce json
// @Param logo formData file true "Logo image file"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tenant/config/logo [post]
func (h *Handler) UploadLogo(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	if schema == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant schema not found in context",
		})
	}

	file, err := c.FormFile("logo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No logo file provided",
		})
	}

	f, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open file",
		})
	}
	defer f.Close()

	url, err := h.service.UploadLogo(c.Context(), schema, f, file.Filename, file.Header.Get("Content-Type"), file.Size)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"logo_url": url,
	})
}
