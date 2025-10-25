package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimitMiddleware creates a rate limiter per tenant or user
func RateLimitMiddleware(maxRequests int, window time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxRequests,
		Expiration: window,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Rate limit by tenant + user
			tenantID := GetTenantID(c)
			userID := GetUserID(c)

			if tenantID != "" && userID != "" {
				return fmt.Sprintf("ratelimit:%s:%s", tenantID, userID)
			} else if tenantID != "" {
				return fmt.Sprintf("ratelimit:%s", tenantID)
			}

			// Fallback to IP if no tenant/user identified
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded. Please try again later.",
			})
		},
	})
}
