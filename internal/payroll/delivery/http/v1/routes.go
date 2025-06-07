package v1

import "github.com/gofiber/fiber/v2"

func MapPayroll(routes fiber.Router, h *PayrollHandler) {
	payroll := routes.Group("/payroll")

	payroll.Post("/", h.GeneratePayroll)
	payroll.Get("/:payrollId/payslip", h.ShowPayslip)
	payroll.Get("/:payrollId/payslips", h.ListPayslips)
}
