package payroll

import (
	"context"

	"github.com/vnnyx/employee-management/internal/payroll/entity"
	"github.com/vnnyx/employee-management/pkg/database"
)

type Repository interface {
	WithTx(tx database.DBTx) Repository

	StoreNewPayroll(ctx context.Context, payroll entity.Payroll) error
	StoreNewPayslips(ctx context.Context, payslips []entity.Payslip) error
	StoreNewPayrollSummary(ctx context.Context, summary entity.PayrollSummary) error
	FindPayrollByPeriodID(ctx context.Context, periodID string, opts ...entity.FindPayrollOptions) (*entity.Payroll, error)
	FindPayslipByUserIDPeriod(ctx context.Context, userID, periodID string) (*entity.Payslip, error)
	FindPayrollByID(ctx context.Context, payrollID string) (*entity.Payroll, error)
	FindPayslipByPayrollID(ctx context.Context, payrollID string, opts ...entity.FindPayslipOptions) (entity.FindPayslipResult, error)
}
