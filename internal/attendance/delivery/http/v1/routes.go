package v1

import "github.com/gofiber/fiber/v2"

func MapAttendance(routes fiber.Router, h *AttendanceHandler) {
	attendance := routes.Group("/attendance")

	attendance.Post("/", h.SubmitAttendance)
	attendance.Post("/period", h.CreateAttendancePeriod)
}
