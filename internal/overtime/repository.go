package overtime

import (
	"context"
	"time"

	"github.com/vnnyx/employee-management/internal/overtime/entity"
	"github.com/vnnyx/employee-management/pkg/database"
)

type Repository interface {
	WithTx(tx database.DBTx) Repository

	StoreNewOvertime(ctx context.Context, overtime entity.Overtime) error
	FindOvertimeByUserIDDate(ctx context.Context, userID string, date time.Time) (*entity.Overtime, error)
	UpsertOvertime(ctx context.Context, overtime entity.Overtime) error
	FindOvertimeByPeriod(ctx context.Context, startDate, endDate time.Time, opts ...entity.FindOvertimeOptions) (entity.FindOvertimeResult, error)
}
