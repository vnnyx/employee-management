package v1

import "github.com/gofiber/fiber/v2"

func MapReimbursement(routes fiber.Router, h *ReimbursementHandler) {
	reimbursement := routes.Group("/reimbursement")

	reimbursement.Post("/", h.SubmitReimbursement)
}
