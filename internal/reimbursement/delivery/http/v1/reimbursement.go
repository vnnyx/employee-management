package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/internal/constants"
	"github.com/vnnyx/employee-management/internal/dtos"
	"github.com/vnnyx/employee-management/internal/reimbursement"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
)

type ReimbursementHandler struct {
	reimbursementUC reimbursement.UseCase
}

func NewReimbursementHandler(reimbursementUC reimbursement.UseCase) *ReimbursementHandler {
	return &ReimbursementHandler{
		reimbursementUC: reimbursementUC,
	}
}

// @Summary      Submit Reimbursement
// @Description  Submit a reimbursement request
// @Tags         Reimbursement
// @Accept       json
// @Produce      json
// @Param        request body dtos.ReimbursementRequest true "Reimbursement Request"
// @Success      200 {object} dtos.Response "Success"
// @Failure      400 {object} apperror.Error "Bad Request"
// @Router       /v1/reimbursement [POST]
// @Security     BearerAuth
func (h *ReimbursementHandler) SubmitReimbursement(c *fiber.Ctx) error {
	ctx, span := instrumentation.NewTraceSpan(
		c.UserContext(),
		"ReimbursementHandler.SubmitReimbursement()",
	)
	defer span.End()

	authCredential := c.Locals(constants.KeyAuthCredential).(authCredential.Credential)
	var reimbursementRequest dtos.ReimbursementRequest
	if err := c.BodyParser(&reimbursementRequest); err != nil {
		return errors.Wrap(err, "ReimbursementHandler().SubmitReimbursement().c.BodyParser()")
	}

	if err := reimbursementRequest.Validate(); err != nil {
		return errors.Wrap(err, "ReimbursementHandler().SubmitReimbursement().reimbursementRequest.Validate()")
	}

	err := h.reimbursementUC.SubmitReimbursement(ctx, authCredential, reimbursementRequest.ToRequestEntity())
	if err != nil {
		return errors.Wrap(err, "ReimbursementHandler().SubmitReimbursement().uc.SubmitReimbursement()")
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.Response{
			RequestID: authCredential.RequestID,
		},
	)
}
