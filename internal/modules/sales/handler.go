package sales

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
	return &Handler{service: NewService(db, cfg)}
}

func RegisterRoutes(router fiber.Router, db *database.DB, cfg *config.Config) {
	h := NewHandler(db, cfg)
	sales := router.Group("/sales")
	sales.Post("/", h.Create)
	sales.Get("/", h.List)
	sales.Get("/:id", h.Get)
	sales.Put("/:id", h.Update)
}

// Create creates a new sale
// @Summary Create sale
// @Tags sales
// @Accept json
// @Produce json
// @Param sale body CreateSaleRequest true "Sale data"
// @Success 201 {object} Sale
// @Router /api/v1/tenant/sales [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	var req CreateSaleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	sale, err := h.service.Create(c.Context(), schema, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(sale)
}

// List lists all sales
// @Summary List sales
// @Tags sales
// @Produce json
// @Param status query string false "Filter by status"
// @Success 200 {array} Sale
// @Router /api/v1/tenant/sales [get]
func (h *Handler) List(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	status := c.Query("status")

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	sales, err := h.service.List(c.Context(), schema, statusPtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(sales)
}

// Get retrieves a sale by ID
// @Summary Get sale
// @Tags sales
// @Produce json
// @Param id path string true "Sale ID"
// @Success 200 {object} Sale
// @Router /api/v1/tenant/sales/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	sale, err := h.service.GetByID(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Sale not found",
		})
	}

	return c.JSON(sale)
}

// Update updates a sale
// @Summary Update sale
// @Tags sales
// @Accept json
// @Produce json
// @Param id path string true "Sale ID"
// @Param sale body UpdateSaleRequest true "Sale data"
// @Success 200 {object} Sale
// @Router /api/v1/tenant/sales/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	var req UpdateSaleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	sale, err := h.service.Update(c.Context(), schema, id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(sale)
}
