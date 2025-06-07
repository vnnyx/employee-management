package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/vnnyx/employee-management/internal/attendance"
	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/internal/constants"
	"github.com/vnnyx/employee-management/internal/dtos"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
)

type AttendanceHandler struct {
	uc attendance.UseCase
}

func NewAttendanceHandler(uc attendance.UseCase) *AttendanceHandler {
	return &AttendanceHandler{
		uc: uc,
	}
}

func (h *AttendanceHandler) SubmitAttendance(c *fiber.Ctx) error {
	ctx, span := instrumentation.NewTraceSpan(
		c.UserContext(),
		"AttendanceHandler.SubmitAttendance()",
	)
	defer span.End()

	authCredential := c.Locals(constants.KeyAuthCredential).(authCredential.Credential)
	err := h.uc.SubmitAttendance(ctx, authCredential)
	if err != nil {
		return errors.Wrap(err, "AttendanceHandler().SubmitAttendance().uc.SubmitAttendance()")
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.Response{
			RequestID: authCredential.RequestID,
		},
	)
}

func (h *AttendanceHandler) CreateAttendancePeriod(c *fiber.Ctx) error {
	ctx, span := instrumentation.NewTraceSpan(
		c.UserContext(),
		"AttendanceHandler.CreateAttendancePeriod()",
	)
	defer span.End()

	authCredential := c.Locals(constants.KeyAuthCredential).(authCredential.Credential)

	var req dtos.AttendancePeriodRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.Wrap(err, "AttendanceHandler().CreateAttendancePeriod().c.BodyParser()")
	}

	id, err := h.uc.CreateAttendancePeriod(ctx, authCredential, req.ToRequestEntity())
	if err != nil {
		return errors.Wrap(err, "AttendanceHandler().CreateAttendancePeriod().uc.CreateAttendancePeriod()")
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.Response{
			RequestID: authCredential.RequestID,
			Data: map[string]string{
				"id": id,
			},
		},
	)
}
