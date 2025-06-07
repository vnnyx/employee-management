package v1

import "github.com/gofiber/fiber/v2"

func MapAuth(routes fiber.Router, h *AuthHandler) {
	auth := routes.Group("/auth")

	auth.Post("/login", h.Login)
}
