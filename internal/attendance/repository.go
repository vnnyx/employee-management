package attendance

import (
	"context"
	"time"

	"github.com/vnnyx/employee-management/internal/attendance/entity"
	"github.com/vnnyx/employee-management/pkg/database"
)

type Repository interface {
	WithTx(tx database.DBTx) Repository

	StoreNewAttendance(ctx context.Context, attendance entity.Attendance) error
	UpsertAttendance(ctx context.Context, attendance entity.Attendance) error
	StoreNewAttendancePeriod(ctx context.Context, period entity.AttendancePeriod) error
	FindPeriodByID(ctx context.Context, periodID string) (*entity.AttendancePeriod, error)
	FindAttendanceByPeriod(ctx context.Context, startDate, endDate time.Time, opts ...entity.FindAttendanceOptions) (entity.FindAttendanceResult, error)
	FindAttendancePeriodByPayrollID(ctx context.Context, payrollID string) (*entity.AttendancePeriod, error)
}
