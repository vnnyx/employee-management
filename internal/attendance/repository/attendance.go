package repository

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/vnnyx/employee-management/internal/attendance"
	"github.com/vnnyx/employee-management/internal/attendance/entity"
	"github.com/vnnyx/employee-management/internal/constants"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
)

type attendanceRepo struct {
	db database.Queryer
}

func NewAttendanceRepository(db database.Queryer) attendance.Repository {
	return &attendanceRepo{
		db: db,
	}
}

func (r *attendanceRepo) WithTx(tx database.DBTx) attendance.Repository {
	return &attendanceRepo{
		db: tx,
	}
}

func (r *attendanceRepo) StoreNewAttendance(ctx context.Context, attendance entity.Attendance) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"AttendanceRepository.StoreNewAttendance()",
	)
	defer span.End()

	query, args, err := sqlx.Named(insertAttendanceQuery, attendance)
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
		return errors.Wrap(errors.New("failed to insert attendance"), constants.ErrWrapPgxscanGet)
	}

	return nil
}

func (r *attendanceRepo) UpsertAttendance(ctx context.Context, attendance entity.Attendance) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"AttendanceRepository.UpsertAttendance()",
	)
	defer span.End()

	query, args, err := sqlx.Named(upsertAttendanceQuery, attendance)
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
		return errors.Wrap(errors.New("failed to upsert attendance"), constants.ErrWrapPgxscanGet)
	}

	return nil
}

func (r *attendanceRepo) StoreNewAttendancePeriod(ctx context.Context, period entity.AttendancePeriod) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"AttendanceRepository.StoreNewAttendancePeriod()",
	)
	defer span.End()

	query, args, err := sqlx.Named(insertAttendancePeriodQuery, period)
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
		return errors.Wrap(errors.New("failed to insert attendance period"), constants.ErrWrapPgxscanGet)
	}

	return nil
}

func (r *attendanceRepo) FindPeriodByID(ctx context.Context, periodID string) (*entity.AttendancePeriod, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"AttendanceRepository.FindPeriodByID()",
	)
	defer span.End()

	var period entity.AttendancePeriod
	err := pgxscan.Get(ctx, r.db, &period, findAttendancePeriodByIDQuery, periodID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, constants.ErrWrapPgxscanGet)
	}

	return &period, nil
}

func (r *attendanceRepo) FindAttendanceByPeriod(ctx context.Context, startDate, endDate time.Time, opts ...entity.FindAttendanceOptions) (entity.FindAttendanceResult, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"AttendanceRepository.FindAttendanceByPeriod()",
	)
	defer span.End()

	var result entity.FindAttendanceResult

	query := findAttendanceByPeriodQuery
	if len(opts) > 0 && opts[0].PessimisticLock {
		query += " FOR UPDATE"
	}

	rows, err := r.db.Query(ctx, query, startDate, endDate)
	if err != nil {
		return result, errors.Wrap(err, constants.ErrWrapPgxQuery)
	}
	defer rows.Close()

	var attendances []entity.Attendance
	for rows.Next() {
		var attendance entity.Attendance
		err := rows.Scan(
			&attendance.ID,
			&attendance.UserID,
			&attendance.AttendanceDate,
			&attendance.CreatedAt,
			&attendance.UpdatedAt,
			&attendance.CreatedBy,
			&attendance.UpdatedBy,
			&attendance.IPAddress,
		)
		if err != nil {
			return result, errors.Wrap(err, constants.ErrWrapDbQueryRowScan)
		}

		attendances = append(attendances, attendance)
	}

	result.List = attendances

	if len(opts) > 0 && opts[0].MappedOptions != nil {
		result.IsMapped = true
		result.MappedBy = opts[0].MappedOptions.MappedBy

		mappedAttendances := make(map[any][]entity.Attendance)
		var keyFunc func(attendance entity.Attendance) any

		switch opts[0].MappedOptions.MappedBy {
		case entity.MappedByUserID:
			keyFunc = func(attendance entity.Attendance) any {
				return attendance.UserID
			}
		case entity.MappedByAttendanceDate:
			keyFunc = func(attendance entity.Attendance) any {
				return attendance.AttendanceDate.Format("2006-01-02")
			}
		default:
			return result, errors.New("unsupported mapped by option")
		}

		for _, attendance := range attendances {
			key := keyFunc(attendance)
			mappedAttendances[key] = append(mappedAttendances[key], attendance)
		}

		result.Mapped = mappedAttendances
	}

	return result, nil
}

func (r *attendanceRepo) FindAttendancePeriodByPayrollID(ctx context.Context, payrollID string) (*entity.AttendancePeriod, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"AttendanceRepository.FindAttendancePeriodByPayrollID()",
	)
	defer span.End()

	var period entity.AttendancePeriod
	err := pgxscan.Get(ctx, r.db, &period, findAttendancePeriodByPayrollIDQuery, payrollID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, constants.ErrWrapPgxscanGet)
	}

	return &period, nil
}
