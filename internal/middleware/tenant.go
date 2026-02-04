package middleware

import (
	"context"
	"strings"

	"github.com/stanleyh24/vendix/internal/database"
	"github.com/stanleyh24/vendix/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TenantContextKey string

const (
	TenantIDKey     TenantContextKey = "tenant_id"
	TenantSchemaKey TenantContextKey = "tenant_schema"
)

type Tenant struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Slug       string    `db:"slug"`
	SchemaName string    `db:"schema_name"`
	Status     string    `db:"status"`
}

// TenantMiddleware identifies and validates the tenant for each request
func TenantMiddleware(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tenantSlug string

		// Try to get tenant from X-Tenant-ID header
		if headerTenant := c.Get("X-Tenant-ID"); headerTenant != "" {
			tenantSlug = headerTenant
		} else {
			// Try to extract from subdomain
			host := c.Hostname()
			parts := strings.Split(host, ".")
			if len(parts) > 2 {
				tenantSlug = parts[0]
			}
		}

		if tenantSlug == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Tenant not specified. Use X-Tenant-ID header or subdomain",
			})
		}

		// Fetch tenant from database
		var tenant Tenant
		query := `SELECT id, name, slug, schema_name, status FROM public.tenants WHERE slug = $1 AND deleted_at IS NULL`

		err := db.Get(&tenant, query, tenantSlug)
		if err != nil {
			logger.Error("Tenant not found", "slug", tenantSlug, "error", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Tenant not found",
			})
		}

		// Check if tenant is active
		if tenant.Status != "active" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Tenant is not active",
			})
		}

		// Set tenant in context
		c.Locals("tenant_id", tenant.ID.String())
		c.Locals("tenant_schema", tenant.SchemaName)

		// Set search path for this request (if using connection from pool)
		ctx := context.Background()
		if err := db.SetSearchPath(ctx, tenant.SchemaName); err != nil {
			logger.Error("Failed to set search path", "schema", tenant.SchemaName, "error", err)
		}

		return c.Next()
	}
}

// GetTenantID retrieves tenant ID from fiber context
func GetTenantID(c *fiber.Ctx) string {
	if id, ok := c.Locals("tenant_id").(string); ok {
		return id
	}
	return ""
}

// GetTenantSchema retrieves tenant schema from fiber context
func GetTenantSchema(c *fiber.Ctx) string {
	if schema, ok := c.Locals("tenant_schema").(string); ok {
		return schema
	}
	return ""
}
