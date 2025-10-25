package middleware

import (
	"vendix/internal/database"
	"vendix/internal/logger"

	"github.com/gofiber/fiber/v2"
)

// ModuleMiddleware checks if tenant has access to a specific module based on their plan
func ModuleMiddleware(db *database.DB, moduleName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := GetTenantID(c)
		if tenantID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Tenant not identified",
			})
		}

		// Check if tenant has active subscription with this module enabled
		var isEnabled bool
		query := `
			SELECT EXISTS(
				SELECT 1 
				FROM public.subscriptions s
				JOIN public.plan_modules pm ON s.plan_id = pm.plan_id
				WHERE s.tenant_id = $1
				AND s.status = 'active'
				AND pm.module_name = $2
				AND pm.is_enabled = true
				AND s.current_period_end > NOW()
			)
		`

		err := db.Get(&isEnabled, query, tenantID, moduleName)
		if err != nil {
			logger.Error("Failed to check module access",
				"tenant_id", tenantID,
				"module", moduleName,
				"error", err,
			)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to verify module access",
			})
		}

		if !isEnabled {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Module not available in your plan. Please upgrade to access this feature.",
			})
		}

		return c.Next()
	}
}
