package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/internal/constants"
	"github.com/vnnyx/employee-management/internal/dtos"
	"github.com/vnnyx/employee-management/internal/overtime"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
)

type OvertimeHandler struct {
	overtimeUC overtime.UseCase
}

func NewOvertimeHandler(overtimeUC overtime.UseCase) *OvertimeHandler {
	return &OvertimeHandler{
		overtimeUC: overtimeUC,
	}
}

func (h *OvertimeHandler) SubmitOvertime(c *fiber.Ctx) error {
	ctx, span := instrumentation.NewTraceSpan(
		c.UserContext(),
		"OvertimeHandler.SubmitOvertime()",
	)
	defer span.End()

	authCredential := c.Locals(constants.KeyAuthCredential).(authCredential.Credential)
	var overtimeRequest dtos.OvertimeRequest
	if err := c.BodyParser(&overtimeRequest); err != nil {
		return errors.Wrap(err, "OvertimeHandler().SubmitOvertime().c.BodyParser()")
	}

	if err := overtimeRequest.Validate(); err != nil {
		return errors.Wrap(err, "OvertimeHandler().SubmitOvertime().overtimeRequest.Validate()")
	}

	err := h.overtimeUC.SubmitOvertime(ctx, authCredential, overtimeRequest.ToRequestEntity())
	if err != nil {
		return errors.Wrap(err, "OvertimeHandler().SubmitOvertime().uc.SubmitOvertime()")
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.Response{
			RequestID: authCredential.RequestID,
		},
	)
}
