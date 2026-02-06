package server

import (
	"github.com/stanleyh24/vendix/internal/config"
	"github.com/stanleyh24/vendix/internal/database"
	"github.com/stanleyh24/vendix/internal/middleware"
	"github.com/stanleyh24/vendix/internal/modules/accounting"
	"github.com/stanleyh24/vendix/internal/modules/auth"
	"github.com/stanleyh24/vendix/internal/modules/cashregisters"
	"github.com/stanleyh24/vendix/internal/modules/creditnotes"
	"github.com/stanleyh24/vendix/internal/modules/customers"
	"github.com/stanleyh24/vendix/internal/modules/employees"
	"github.com/stanleyh24/vendix/internal/modules/invoices"
	"github.com/stanleyh24/vendix/internal/modules/notifications"
	"github.com/stanleyh24/vendix/internal/modules/payments"
	"github.com/stanleyh24/vendix/internal/modules/payroll"
	"github.com/stanleyh24/vendix/internal/modules/products"
	"github.com/stanleyh24/vendix/internal/modules/purchases"
	"github.com/stanleyh24/vendix/internal/modules/reports"
	"github.com/stanleyh24/vendix/internal/modules/returns"
	"github.com/stanleyh24/vendix/internal/modules/sales"
	"github.com/stanleyh24/vendix/internal/modules/suppliers"
	"github.com/stanleyh24/vendix/internal/modules/tenantconfig"
	"github.com/stanleyh24/vendix/internal/modules/tenants"
	"github.com/stanleyh24/vendix/internal/modules/webhooks"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Server struct {
	*fiber.App
	Config *config.Config
	DB     *database.DB
}

func NewServer(cfg *config.Config, db *database.DB) *Server {
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Tenant-ID",
	}))

	server := &Server{
		App:    app,
		Config: cfg,
		DB:     db,
	}

	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	// Health check
	s.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"app":    s.Config.AppName,
		})
	})

	// API v1
	api := s.Group("/api/v1")

	// Public routes (no authentication required)
	public := api.Group("/public")

	// Tenant management (superadmin only, no tenant context)
	tenants.RegisterRoutes(public, s.DB, s.Config)

	// Webhook routes (no authentication, validated by signature)
	webhooks.RegisterRoutes(public, s.DB, s.Config)

	// Tenant-specific routes (require tenant identification)
	tenant := api.Group("/tenant")
	tenant.Use(middleware.TenantMiddleware(s.DB))

	// Auth routes (within tenant context, but before auth)
	auth.RegisterRoutes(tenant, s.DB, s.Config)

	// Protected routes (require authentication)
	protected := tenant.Group("")
	protected.Use(middleware.AuthMiddleware(s.Config))

	// Apply rate limiting to protected routes
	if s.Config.RateLimitEnabled {
		protected.Use(middleware.RateLimitMiddleware(
			s.Config.RateLimitMaxRequests,
			s.Config.RateLimitWindow,
		))
	}

	// Tenant configuration
	tenantconfig.RegisterRoutes(protected, s.DB, s.Config)

	// Customer management
	customers.RegisterRoutes(protected, s.DB, s.Config)

	// Product management
	products.RegisterRoutes(protected, s.DB, s.Config)

	// Suppliers management
	suppliers.RegisterRoutes(protected, s.DB, s.Config)

	// Purchases management
	purchases.RegisterRoutes(protected, s.DB, s.Config)

	// Invoice management
	invoices.RegisterRoutes(protected, s.DB, s.Config)

	// Credit Notes management (DGII NCF tipo 04)
	creditnotes.RegisterRoutes(protected, s.DB, s.Config)

	// Sales management (POS)
	sales.RegisterRoutes(protected, s.DB, s.Config)

	// Returns management
	returns.RegisterRoutes(protected, s.DB, s.Config)

	// Payment management
	payments.RegisterRoutes(protected, s.DB, s.Config)

	// Reports
	reports.RegisterRoutes(protected, s.DB, s.Config)

	// Accounting
	accounting.RegisterRoutes(protected, s.DB, s.Config)

	// Cash Registers
	cashregisters.RegisterRoutes(protected, s.DB, s.Config)

	// Employee management
	employees.RegisterRoutes(protected, s.DB, s.Config)

	// Payroll management
	payroll.RegisterRoutes(protected, s.DB, s.Config)

	// Notifications
	notifications.RegisterRoutes(protected, s.DB, s.Config)
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
		"code":  code,
	})
}
