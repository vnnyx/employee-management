package reimbursement

import (
	"context"
	"time"

	"github.com/vnnyx/employee-management/internal/reimbursement/entity"
	"github.com/vnnyx/employee-management/pkg/database"
)

type Repository interface {
	WithTx(tx database.DBTx) Repository

	StoreNewReimbursement(ctx context.Context, reimbursement entity.Reimbursement) error
	FindReimbursementByPeriod(ctx context.Context, startDate, endDate time.Time, opts ...entity.FindReimbursementOptions) (entity.FindReimbursementResult, error)
	FindReimbursementByUserIDPeriod(ctx context.Context, userID string, startDate, endDate time.Time) ([]entity.Reimbursement, error)
}
