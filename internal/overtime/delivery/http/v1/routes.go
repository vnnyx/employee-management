package v1

import "github.com/gofiber/fiber/v2"

func MapOvertime(routes fiber.Router, h *OvertimeHandler) {
	overtime := routes.Group("/overtime")

	overtime.Post("/", h.SubmitOvertime)
}
