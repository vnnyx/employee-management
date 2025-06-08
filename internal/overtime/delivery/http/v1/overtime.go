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

// @Summary      Submit Overtime
// @Description  Submit an overtime request
// @Tags         Overtime
// @Accept       json
// @Produce      json
// @Param        request body dtos.OvertimeRequest true "Overtime Request"
// @Success      200 {object} dtos.Response "Success"
// @Failure      400 {object} apperror.Error "Bad Request"
// @Router       /v1/overtime [POST]
// @Security     BearerAuth
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
