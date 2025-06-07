package payroll

import (
	"context"

	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/internal/payroll/entity"
)

type UseCase interface {
	GeneratePayroll(ctx context.Context, authCredential authCredential.Credential, periodID string) (entity.GeneratedPayroll, error)
	ShowPayslip(ctx context.Context, authCredential authCredential.Credential, payrollID string) (*entity.PayslipData, error)
	ListPayslips(ctx context.Context, authCredential authCredential.Credential, payrollID string) (entity.ListPayslips, error)
}
