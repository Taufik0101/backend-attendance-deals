package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func (m *CustomMiddleware) FiberContextToContextMiddleware(c *fiber.Ctx) error {
	c.Locals("fiberContext", c)
	return c.Next()
}
