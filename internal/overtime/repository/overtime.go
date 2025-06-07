package repository

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/vnnyx/employee-management/internal/constants"
	"github.com/vnnyx/employee-management/internal/overtime"
	"github.com/vnnyx/employee-management/internal/overtime/entity"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
)

type overtimeRepo struct {
	db database.Queryer
}

func NewOvertimeRepository(db database.Queryer) overtime.Repository {
	return &overtimeRepo{
		db: db,
	}
}

func (r *overtimeRepo) WithTx(tx database.DBTx) overtime.Repository {
	return &overtimeRepo{
		db: tx,
	}
}

func (r *overtimeRepo) StoreNewOvertime(ctx context.Context, overtime entity.Overtime) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"OvertimeRepository.StoreNewOvertime()",
	)
	defer span.End()

	query, args, err := sqlx.Named(insertOvertimeQuery, overtime)
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
		return errors.Wrap(errors.New("failed to insert overtime"), constants.ErrWrapPgxscanGet)
	}

	return nil
}

func (r *overtimeRepo) FindOvertimeByUserIDDate(ctx context.Context, userID string, date time.Time) (*entity.Overtime, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"OvertimeRepository.GetOvertimeByUserIDDate()",
	)
	defer span.End()

	var overtime entity.Overtime
	err := pgxscan.Get(ctx, r.db, &overtime, findOvertimeByUserIDDate, userID, date)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, constants.ErrWrapPgxscanGet)
	}

	return &overtime, nil
}

func (r *overtimeRepo) UpsertOvertime(ctx context.Context, overtime entity.Overtime) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"OvertimeRepository.UpsertOvertime()",
	)
	defer span.End()

	query, args, err := sqlx.Named(upsertOvertimeQuery, overtime)
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
		return errors.Wrap(errors.New("failed to upsert overtime"), constants.ErrWrapPgxscanGet)
	}

	return nil
}

func (r *overtimeRepo) FindOvertimeByPeriod(ctx context.Context, startDate, endDate time.Time, opts ...entity.FindOvertimeOptions) (entity.FindOvertimeResult, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"OvertimeRepository.FindOvertimeByPeriod()",
	)
	defer span.End()

	var result entity.FindOvertimeResult

	query := findOvertimeByPeriodQuery
	if len(opts) > 0 && opts[0].PessimisticLock {
		query += " FOR UPDATE"
	}

	rows, err := r.db.Query(ctx, query, startDate, endDate)
	if err != nil {
		return result, errors.Wrap(err, constants.ErrWrapPgxQuery)
	}
	defer rows.Close()

	var overtimes []entity.Overtime
	for rows.Next() {
		var overtime entity.Overtime
		err := rows.Scan(
			&overtime.ID,
			&overtime.UserID,
			&overtime.OverTimeDate,
			&overtime.OvertimeHours,
			&overtime.CreatedAt,
			&overtime.UpdatedAt,
			&overtime.CreatedBy,
			&overtime.UpdatedBy,
			&overtime.IPAddress,
		)
		if err != nil {
			return result, errors.Wrap(err, constants.ErrWrapDbQueryRowScan)
		}

		overtimes = append(overtimes, overtime)
	}

	result.List = overtimes

	if len(opts) > 0 && opts[0].MappedOptions != nil {
		result.IsMapped = true
		result.MappedBy = opts[0].MappedOptions.MappedBy

		mappedOvertimes := make(map[any][]entity.Overtime)
		var keyFunc func(overtime entity.Overtime) any

		switch opts[0].MappedOptions.MappedBy {
		case entity.MappedByUserID:
			keyFunc = func(overtime entity.Overtime) any {
				return overtime.UserID
			}
		case entity.MappedByAttendanceDate:
			keyFunc = func(overtime entity.Overtime) any {
				return overtime.OverTimeDate.Format("2006-01-02")
			}
		default:
			return result, errors.New("unsupported mapped by option")
		}

		for _, overtime := range overtimes {
			key := keyFunc(overtime)
			mappedOvertimes[key] = append(mappedOvertimes[key], overtime)
		}

		result.Mapped = mappedOvertimes
	}

	return result, nil
}
