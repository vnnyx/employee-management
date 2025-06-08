package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	_ "github.com/vnnyx/employee-management/docs/helper"
	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/internal/constants"
	"github.com/vnnyx/employee-management/internal/dtos"
	"github.com/vnnyx/employee-management/internal/payroll"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
	"github.com/vnnyx/employee-management/pkg/resourceful"
)

type PayrollHandler struct {
	uc payroll.UseCase
}

func NewPayrollHandler(uc payroll.UseCase) *PayrollHandler {
	return &PayrollHandler{
		uc: uc,
	}
}

// @Summary      Generate Payroll
// @Description  Generate payroll for a specific period
// @Tags         Payroll
// @Accept       json
// @Produce      json
// @Param        request body dtos.GeneratePayrollRequest true "Generate Payroll Request"
// @Success      201 {object} dtos.Response{request,data=dtos.GeneratedPayrollResponse} "Generated Payroll Response"
// @Failure      400 {object} apperror.Error "Bad Request"
// @Router        /v1/payroll [POST]
// @Security     BearerAuth
func (h *PayrollHandler) GeneratePayroll(c *fiber.Ctx) error {
	ctx, span := instrumentation.NewTraceSpan(
		c.UserContext(),
		"PayrollHandler.GeneratePayroll()",
	)
	defer span.End()

	authCredential := c.Locals(constants.KeyAuthCredential).(authCredential.Credential)

	var req dtos.GeneratePayrollRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.Wrap(err, "PayrollHandler().GeneratePayroll().c.BodyParser()")
	}

	data, err := h.uc.GeneratePayroll(ctx, authCredential, req.PeriodID)
	if err != nil {
		return errors.Wrap(err, "PayrollHandler().GeneratePayroll().uc.GeneratePayroll()")
	}

	return c.Status(fiber.StatusCreated).JSON(
		dtos.Response{
			RequestID: authCredential.RequestID,
			Data:      dtos.GeneratedPayrollResponse(data),
		},
	)
}

// @Summary      Show Payslip
// @Description  Show payslip for a specific payroll
// @Tags         Payroll
// @Accept       json
// @Produce      json
// @Param        payrollId path string true "Payroll ID"
// @Success      200 {object} dtos.Response{data=dtos.PayslipDataResponse} "Payslip Response"
// @Failure      400 {object} apperror.Error "Bad Request"
// @Router       /v1/payroll/{payrollId}/payslip [GET]
// @Security     BearerAuth
func (h *PayrollHandler) ShowPayslip(c *fiber.Ctx) error {
	ctx, span := instrumentation.NewTraceSpan(
		c.UserContext(),
		"PayrollHandler.ShowPayslip()",
	)
	defer span.End()

	authCredential := c.Locals(constants.KeyAuthCredential).(authCredential.Credential)

	var param struct {
		PayrollID uuid.UUID `params:"payrollId"`
	}
	err := c.ParamsParser(&param)
	if err != nil {
		return errors.Wrap(err, "PayrollHandler().ShowPayslip().c.ParamsParser()")
	}

	data, err := h.uc.ShowPayslip(ctx, authCredential, param.PayrollID.String())
	if err != nil {
		return errors.Wrap(err, "PayrollHandler().ShowPayslip().uc.ShowPayslip()")
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.Response{
			RequestID: authCredential.RequestID,
			Data:      dtos.NewShowPayslipResponse(data),
		},
	)
}

// @Summary      List Payslips
// @Description  List payslips for a specific payroll
// @Tags         Payroll
// @Accept       json
// @Produce      json
// @Param        payrollId path string true "Payroll ID"
// @Param        limit query int false "Limit" default(10)
// @Param        page query int false "Page" default(1)
// @Param        mode query string false "Mode" default(offset)
// @Param        cursor query string false "Cursor"
// @Success      200 {object} docshelper.Response[string, dtos.PayslipDataResponse, entity.ListPayslipMetadata] "List of Payslips"
// @Failure      400 {object} apperror.Error "Bad Request"
// @Router       /v1/payroll/{payrollId}/payslips [GET]
func (h *PayrollHandler) ListPayslips(c *fiber.Ctx) error {
	ctx, span := instrumentation.NewTraceSpan(
		c.UserContext(),
		"PayrollHandler.ListPayslips()",
	)
	defer span.End()

	authCredential := c.Locals(constants.KeyAuthCredential).(authCredential.Credential)

	var param struct {
		PayrollID uuid.UUID `params:"payrollId"`
	}
	err := c.ParamsParser(&param)
	if err != nil {
		return errors.Wrap(err, "PayrollHandler().ListPayslips().c.ParamsParser()")
	}

	var request dtos.ListPayslipsRequest
	if err := c.QueryParser(&request); err != nil {
		return errors.Wrap(err, "PayrollHandler().ListPayslips().c.QueryParser()")
	}

	decodedCursor, err := resourceful.DecodeCursor(request.Cursor)
	if err != nil {
		return errors.Wrap(err, "PayrollHandler().ListPayslips().resourceful.DecodeCursor()")
	}

	resourceful := resourceful.NewResource[string, dtos.PayslipDataResponse](&resourceful.Parameter{
		Limit:  request.Limit,
		Page:   request.Page,
		Mode:   resourceful.Mode(request.Mode.GetOrDefault("offset")),
		Cursor: decodedCursor,
	})

	data, err := h.uc.ListPayslips(ctx, authCredential, param.PayrollID.String(), resourceful)
	if err != nil {
		return errors.Wrap(err, "PayrollHandler().ListPayslips().uc.ListPayslips()")
	}

	return c.Status(fiber.StatusOK).JSON(data.Response(authCredential.RequestID))
}
