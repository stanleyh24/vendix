package server

import (
	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/middleware"
	"vendix/internal/modules/accounting"
	"vendix/internal/modules/auth"
	"vendix/internal/modules/customers"
	"vendix/internal/modules/invoices"
	"vendix/internal/modules/payments"
	"vendix/internal/modules/products"
	"vendix/internal/modules/reports"
	"vendix/internal/modules/sales"
	"vendix/internal/modules/tenantconfig"
	"vendix/internal/modules/tenants"
	"vendix/internal/modules/webhooks"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	swagger "github.com/swaggo/fiber-swagger"
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

	// Swagger documentation
	s.Get("/swagger/*", swagger.WrapHandler)

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
	tenantconfig.RegisterRoutes(protected, s.DB)

	// Customer management
	customers.RegisterRoutes(protected, s.DB, s.Config)

	// Product management
	products.RegisterRoutes(protected, s.DB, s.Config)

	// Invoice management
	invoices.RegisterRoutes(protected, s.DB, s.Config)

	// Sales management (POS)
	sales.RegisterRoutes(protected, s.DB, s.Config)

	// Payment management
	payments.RegisterRoutes(protected, s.DB, s.Config)

	// Reports
	reports.RegisterRoutes(protected, s.DB, s.Config)

	// Accounting
	accounting.RegisterRoutes(protected, s.DB, s.Config)
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
