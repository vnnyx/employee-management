package repository

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/vnnyx/employee-management/internal/constants"
	"github.com/vnnyx/employee-management/internal/reimbursement"
	"github.com/vnnyx/employee-management/internal/reimbursement/entity"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
)

type reimbursementRepo struct {
	db database.Queryer
}

func NewReimbursementRepository(db database.Queryer) reimbursement.Repository {
	return &reimbursementRepo{
		db: db,
	}
}

func (r *reimbursementRepo) WithTx(tx database.DBTx) reimbursement.Repository {
	return &reimbursementRepo{
		db: tx,
	}
}

func (r *reimbursementRepo) StoreNewReimbursement(ctx context.Context, reimbursement entity.Reimbursement) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"ReimbursementRepository.StoreNewReimbursement()",
	)
	defer span.End()

	query, args, err := sqlx.Named(insertReimbursementQuery, reimbursement)
	if err != nil {
		return errors.Wrap(err, constants.ErrWrapSqlxNamed)
	}
	query = database.Rebind(query)

	var returnedID string
	err = pgxscan.Get(ctx, r.db, &returnedID, query, args...)
	if err != nil {
		return errors.Wrap(err, constants.ErrWrapPgxscanGet)
	}

	if returnedID == "" {
		return errors.Wrap(errors.New("failed to insert reimbursement"), constants.ErrWrapPgxscanGet)
	}

	return nil
}

func (r *reimbursementRepo) FindReimbursementByPeriod(ctx context.Context, startDate, endDate time.Time, opts ...entity.FindReimbursementOptions) (entity.FindReimbursementResult, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"ReimbursementRepository.FindReimbursementByPeriod()",
	)
	defer span.End()

	var result entity.FindReimbursementResult

	query := findReimbursementByPeriodQuery
	if len(opts) > 0 && opts[0].PessimisticLock {
		query += " FOR UPDATE"
	}

	rows, err := r.db.Query(ctx, query, startDate, endDate)
	if err != nil {
		return result, errors.Wrap(err, constants.ErrWrapPgxQuery)
	}
	defer rows.Close()

	var reimbursements []entity.Reimbursement
	for rows.Next() {
		var reimbursement entity.Reimbursement
		err := rows.Scan(
			&reimbursement.ID,
			&reimbursement.UserID,
			&reimbursement.Amount,
			&reimbursement.Description,
			&reimbursement.ReimbursementDate,
			&reimbursement.CreatedAt,
			&reimbursement.UpdatedAt,
			&reimbursement.CreatedBy,
			&reimbursement.UpdatedBy,
			&reimbursement.IPAddress,
		)
		if err != nil {
			return result, errors.Wrap(err, constants.ErrWrapDbQueryRowScan)
		}

		reimbursements = append(reimbursements, reimbursement)
	}

	result.List = reimbursements

	if len(opts) > 0 && opts[0].MappedOptions != nil {
		result.IsMapped = true
		result.MappedBy = opts[0].MappedOptions.MappedBy

		mappedReimbursements := make(map[any][]entity.Reimbursement)
		var keyFunc func(reimbursement entity.Reimbursement) any

		switch opts[0].MappedOptions.MappedBy {
		case entity.MappedByUserID:
			keyFunc = func(reimbursement entity.Reimbursement) any {
				return reimbursement.UserID
			}
		case entity.MappedByDate:
			keyFunc = func(reimbursement entity.Reimbursement) any {
				return reimbursement.ReimbursementDate.Format("2006-01-02")
			}
		default:
			return result, errors.New("unsupported mapped by option")
		}

		for _, reimbursement := range reimbursements {
			key := keyFunc(reimbursement)
			mappedReimbursements[key] = append(mappedReimbursements[key], reimbursement)
		}

		result.Mapped = mappedReimbursements
	}

	return result, nil
}

func (r *reimbursementRepo) FindReimbursementByUserIDPeriod(ctx context.Context, userID string, startDate, endDate time.Time) ([]entity.Reimbursement, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"ReimbursementRepository.FindReimbursementByUserIDPeriod()",
	)
	defer span.End()

	var reimbursements []entity.Reimbursement

	query := findReimbursementByUserIDPeriodQuery
	rows, err := r.db.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		return reimbursements, errors.Wrap(err, constants.ErrWrapPgxQuery)
	}
	defer rows.Close()

	for rows.Next() {
		var reimbursement entity.Reimbursement
		err := rows.Scan(
			&reimbursement.ID,
			&reimbursement.UserID,
			&reimbursement.Amount,
			&reimbursement.Description,
			&reimbursement.ReimbursementDate,
			&reimbursement.CreatedAt,
			&reimbursement.UpdatedAt,
			&reimbursement.CreatedBy,
			&reimbursement.UpdatedBy,
			&reimbursement.IPAddress,
		)
		if err != nil {
			return reimbursements, errors.Wrap(err, constants.ErrWrapDbQueryRowScan)
		}

		reimbursements = append(reimbursements, reimbursement)
	}

	return reimbursements, nil
}
