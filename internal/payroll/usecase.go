package payroll

import (
	"context"

	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/internal/dtos"
	"github.com/vnnyx/employee-management/internal/payroll/entity"
	"github.com/vnnyx/employee-management/pkg/resourceful"
)

type UseCase interface {
	GeneratePayroll(ctx context.Context, authCredential authCredential.Credential, periodID string) (entity.GeneratedPayroll, error)
	ShowPayslip(ctx context.Context, authCredential authCredential.Credential, payrollID string) (*entity.PayslipData, error)
	ListPayslips(ctx context.Context, authCredential authCredential.Credential, payrollID string, resource *resourceful.Resource[string, dtos.PayslipDataResponse]) (*resourceful.Resource[string, dtos.PayslipDataResponse], error)
}
