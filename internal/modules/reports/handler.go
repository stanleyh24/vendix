package reports

import (
	"vendix/internal/config"
	"vendix/internal/database"

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
	reports.Get("/sales", h.SalesReport)
	reports.Get("/taxes", h.TaxReport)
}

func (h *Handler) SalesReport(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"report": "sales"})
}

func (h *Handler) TaxReport(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"report": "taxes"})
}
