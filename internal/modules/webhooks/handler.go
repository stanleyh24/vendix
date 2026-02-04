package webhooks

import (
	"github.com/stanleyh24/vendix/internal/config"
	"github.com/stanleyh24/vendix/internal/database"
	"github.com/stanleyh24/vendix/internal/logger"

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
	webhooks := router.Group("/webhooks")
	webhooks.Post("/stripe", h.HandleStripeWebhook)
}

func (h *Handler) HandleStripeWebhook(c *fiber.Ctx) error {
	logger.Info("Received Stripe webhook")
	// TODO: Implement Stripe webhook signature verification
	// TODO: Handle different event types
	return c.JSON(fiber.Map{"received": true})
}
