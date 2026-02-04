package purchases

import (
	"github.com/stanleyh24/vendix/internal/config"
	"github.com/stanleyh24/vendix/internal/database"
	"github.com/stanleyh24/vendix/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(db *database.DB, cfg *config.Config) *Handler {
	return &Handler{service: NewService(db, cfg)}
}

// RegisterRoutes registers purchase routes
func RegisterRoutes(router fiber.Router, db *database.DB, cfg *config.Config) {
	h := NewHandler(db, cfg)
	purchases := router.Group("/purchases")

	purchases.Get("/", h.List)
	purchases.Post("/", h.Create)
	purchases.Get("/:id", h.GetByID)
	purchases.Put("/:id", h.Update)
	purchases.Delete("/:id", h.Delete)
	purchases.Post("/:id/receive", h.Receive)
	purchases.Get("/:id/allocations", h.GetPurchaseAllocations)

	// Supplier payments routes
	supplierPayments := router.Group("/supplier-payments")
	supplierPayments.Get("/", h.ListSupplierPayments)
	supplierPayments.Post("/", h.CreateSupplierPayment)
	supplierPayments.Get("/:id", h.GetSupplierPaymentByID)
}

// List gets all purchases
// @Summary List purchases
// @Tags purchases
// @Produce json
// @Param supplier_id query string false "Filter by supplier ID"
// @Param status query string false "Filter by status (draft, pending, received, cancelled)"
// @Param limit query int false "Limit results"
// @Param offset query int false "Offset results"
// @Success 200 {array} PurchaseResponse
// @Router /api/v1/tenant/purchases [get]
func (h *Handler) List(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	supplierID := c.Query("supplier_id")
	status := c.Query("status")
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	var supplierIDPtr *string
	if supplierID != "" {
		supplierIDPtr = &supplierID
	}

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	purchases, err := h.service.List(c.Context(), schema, supplierIDPtr, statusPtr, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(purchases)
}

// Create creates a new purchase
// @Summary Create purchase
// @Tags purchases
// @Accept json
// @Produce json
// @Param purchase body CreatePurchaseRequest true "Purchase data"
// @Success 201 {object} PurchaseResponse
// @Router /api/v1/tenant/purchases [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	var req CreatePurchaseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user ID from context (if available)
	var userID *uuid.UUID
	if userIDStr := c.Get("user_id"); userIDStr != "" {
		if uid, err := uuid.Parse(userIDStr); err == nil {
			userID = &uid
		}
	}

	purchase, err := h.service.Create(c.Context(), schema, &req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(purchase)
}

// GetByID gets a purchase by ID
// @Summary Get purchase by ID
// @Tags purchases
// @Produce json
// @Param id path string true "Purchase ID"
// @Success 200 {object} PurchaseResponse
// @Router /api/v1/tenant/purchases/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	purchase, err := h.service.GetByID(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Purchase not found",
		})
	}

	return c.JSON(purchase)
}

// Update updates a purchase
// @Summary Update purchase
// @Tags purchases
// @Accept json
// @Produce json
// @Param id path string true "Purchase ID"
// @Param purchase body UpdatePurchaseRequest true "Purchase data"
// @Success 200 {object} PurchaseResponse
// @Router /api/v1/tenant/purchases/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	var req UpdatePurchaseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	purchase, err := h.service.Update(c.Context(), schema, id, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(purchase)
}

// Delete deletes a purchase
// @Summary Delete purchase
// @Tags purchases
// @Param id path string true "Purchase ID"
// @Success 204 "No Content"
// @Router /api/v1/tenant/purchases/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	if err := h.service.Delete(c.Context(), schema, id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Receive marks a purchase as received
// @Summary Receive purchase
// @Tags purchases
// @Accept json
// @Produce json
// @Param id path string true "Purchase ID"
// @Param request body ReceivePurchaseRequest true "Receive data"
// @Success 200 {object} map[string]string
// @Router /api/v1/tenant/purchases/{id}/receive [post]
func (h *Handler) Receive(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	var req ReceivePurchaseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.service.Receive(c.Context(), schema, id, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Purchase received successfully",
	})
}

// GetPurchaseAllocations gets all payment allocations for a purchase
// @Summary Get purchase payment allocations
// @Tags purchases
// @Produce json
// @Param id path string true "Purchase ID"
// @Success 200 {array} SupplierPaymentAllocation
// @Router /api/v1/tenant/purchases/{id}/allocations [get]
func (h *Handler) GetPurchaseAllocations(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	allocations, err := h.service.GetPurchaseAllocations(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(allocations)
}

// ========== Supplier Payments Handlers ==========

// ListSupplierPayments lists supplier payments
// @Summary List supplier payments
// @Tags supplier-payments
// @Produce json
// @Param supplier_id query string false "Filter by supplier ID"
// @Param limit query int false "Limit results"
// @Param offset query int false "Offset results"
// @Success 200 {array} SupplierPaymentResponse
// @Router /api/v1/tenant/supplier-payments [get]
func (h *Handler) ListSupplierPayments(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	supplierID := c.Query("supplier_id")
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	var supplierIDPtr *string
	if supplierID != "" {
		supplierIDPtr = &supplierID
	}

	payments, err := h.service.ListSupplierPayments(c.Context(), schema, supplierIDPtr, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(payments)
}

// CreateSupplierPayment creates a new supplier payment
// @Summary Create supplier payment
// @Tags supplier-payments
// @Accept json
// @Produce json
// @Param payment body CreateSupplierPaymentRequest true "Payment data"
// @Success 201 {object} SupplierPaymentResponse
// @Router /api/v1/tenant/supplier-payments [post]
func (h *Handler) CreateSupplierPayment(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)

	var req CreateSupplierPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user ID from context (if available)
	var userID *uuid.UUID
	if userIDStr := c.Get("user_id"); userIDStr != "" {
		if uid, err := uuid.Parse(userIDStr); err == nil {
			userID = &uid
		}
	}

	payment, err := h.service.CreateSupplierPayment(c.Context(), schema, &req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(payment)
}

// GetSupplierPaymentByID gets a supplier payment by ID
// @Summary Get supplier payment by ID
// @Tags supplier-payments
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} SupplierPaymentResponse
// @Router /api/v1/tenant/supplier-payments/{id} [get]
func (h *Handler) GetSupplierPaymentByID(c *fiber.Ctx) error {
	schema := middleware.GetTenantSchema(c)
	id := c.Params("id")

	payment, err := h.service.GetSupplierPaymentByID(c.Context(), schema, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Payment not found",
		})
	}

	return c.JSON(payment)
}

