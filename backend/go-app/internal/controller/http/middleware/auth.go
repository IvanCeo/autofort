package middleware

import (
	"autofort/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

func New(s *usecase.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		access := c.Get(fiber.HeaderAuthorization)
		if err := s.AuthCheck(access); err != nil {
			return c.Redirect("/refresh", fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}
