package suppliers

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

// RegisterRoutes registers supplier routes under /suppliers
func RegisterRoutes(router fiber.Router, db *database.DB, cfg *config.Config) {
	h := NewHandler(db, cfg)
	sup := router.Group("/suppliers")
	sup.Get("/", h.List)
	sup.Post("/", h.Create)
	sup.Get("/:id", h.Get)
	sup.Put("/:id", h.Update)
	sup.Delete("/:id", h.Delete)
}

func (h *Handler) List(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	var onlyActive *bool
	if v := c.Query("active"); v == "true" {
		b := true
		onlyActive = &b
	}
	list, err := h.service.List(c.Context(), schema, onlyActive)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if list == nil {
		list = []*Supplier{}
	}
	return c.JSON(list)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	var req CreateSupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	sup, err := h.service.Create(c.Context(), schema, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(sup)
}

func (h *Handler) Get(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")
	sup, err := h.service.GetByID(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Supplier not found"})
	}
	return c.JSON(sup)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")
	var req UpdateSupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	sup, err := h.service.Update(c.Context(), schema, id, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(sup)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")
	if err := h.service.Delete(c.Context(), schema, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
