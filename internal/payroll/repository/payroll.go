package repository

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/vnnyx/employee-management/internal/constants"
	"github.com/vnnyx/employee-management/internal/payroll"
	"github.com/vnnyx/employee-management/internal/payroll/entity"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
)

type payrollRepository struct {
	db database.Queryer
}

func NewPayrollRepository(db database.Queryer) payroll.Repository {
	return &payrollRepository{
		db: db,
	}
}

func (r *payrollRepository) WithTx(tx database.DBTx) payroll.Repository {
	return &payrollRepository{
		db: tx,
	}
}

func (r *payrollRepository) StoreNewPayroll(ctx context.Context, payroll entity.Payroll) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"PayrollRepository.StoreNewPayroll()",
	)
	defer span.End()

	query, args, err := sqlx.Named(insertPayrollQuery, payroll)
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
		return errors.Wrap(errors.New("failed to insert payroll"), constants.ErrWrapPgxscanGet)
	}

	return nil
}

func (r *payrollRepository) StoreNewPayslips(ctx context.Context, payslips []entity.Payslip) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"PayrollRepository.StoreNewPayslips()",
	)
	defer span.End()

	if len(payslips) == 0 {
		return nil
	}

	query, args, err := sqlx.Named(insertPayslipQuery, payslips)
	if err != nil {
		return errors.Wrap(err, constants.ErrWrapSqlxNamed)
	}
	query = database.Rebind(query)

	var returnedIDs []string
	err = pgxscan.Select(ctx, r.db, &returnedIDs, query, args...)
	if err != nil {
		return errors.Wrap(err, constants.ErrWrapPgxscanSelect)
	}

	if len(returnedIDs) != len(payslips) {
		return errors.Wrap(errors.New("failed to insert payslips"), constants.ErrWrapPgxscanSelect)
	}

	return nil
}

func (r *payrollRepository) StoreNewPayrollSummary(ctx context.Context, summary entity.PayrollSummary) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"PayrollRepository.StoreNewPayrollSummary()",
	)
	defer span.End()

	query, args, err := sqlx.Named(insertPayrollSummaryQuery, summary)
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
		return errors.Wrap(errors.New("failed to insert payroll summary"), constants.ErrWrapPgxscanGet)
	}

	return nil
}

func (r *payrollRepository) FindPayrollByPeriodID(ctx context.Context, periodID string, opts ...entity.FindPayrollOptions) (*entity.Payroll, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"PayrollRepository.FindPayrollByPeriodID()",
	)
	defer span.End()

	query := findPayrollByPeriodIDQuery
	if len(opts) > 0 && opts[0].PessimisticLock {
		query += " FOR UPDATE"
	}

	var payroll entity.Payroll
	err := pgxscan.Get(ctx, r.db, &payroll, query, periodID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, constants.ErrWrapPgxscanGet)
	}

	return &payroll, nil
}

func (r *payrollRepository) FindPayslipByUserIDPeriod(ctx context.Context, userID, periodID string) (*entity.Payslip, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"PayrollRepository.FindPayslipByUserIDPeriod()",
	)
	defer span.End()

	var payslip entity.Payslip
	err := pgxscan.Get(ctx, r.db, &payslip, findPayslipByUserIDPeriodQuery, userID, periodID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, constants.ErrWrapPgxscanGet)
	}

	return &payslip, nil
}

func (r *payrollRepository) FindPayrollByID(ctx context.Context, payrollID string) (*entity.Payroll, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"PayrollRepository.FindPayrollByID()",
	)
	defer span.End()

	var payroll entity.Payroll
	err := pgxscan.Get(ctx, r.db, &payroll, findPayrollByIDQuery, payrollID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, constants.ErrWrapPgxscanGet)
	}

	return &payroll, nil
}

func (r *payrollRepository) FindPayslipByPayrollID(ctx context.Context, payrollID string, opts ...entity.FindPayslipOptions) (entity.FindPayslipResult, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"PayrollRepository.FindPayslipByPayrollID()",
	)
	defer span.End()

	var result entity.FindPayslipResult

	query := findPayslipByPayrollIDQuery
	if len(opts) > 0 && opts[0].PessimisticLock {
		query += " FOR UPDATE"
	}

	rows, err := r.db.Query(ctx, query, payrollID)
	if err != nil {
		return result, errors.Wrap(err, constants.ErrWrapPgxQuery)
	}
	defer rows.Close()

	var payslips []entity.Payslip
	for rows.Next() {
		var payslip entity.Payslip
		err := rows.Scan(
			&payslip.ID,
			&payslip.UserID,
			&payslip.PayrollID,
			&payslip.BaseSalary,
			&payslip.AttendanceDays,
			&payslip.OvertimeHours,
			&payslip.OvertimePay,
			&payslip.ReimbursementTotal,
			&payslip.TotalTakeHome,
			&payslip.CreatedAt,
			&payslip.UpdatedAt,
			&payslip.CreatedBy,
			&payslip.UpdatedBy,
			&payslip.IPAddress,
		)
		if err != nil {
			return result, errors.Wrap(err, constants.ErrWrapDbQueryRowScan)
		}

		payslips = append(payslips, payslip)
	}

	result.List = payslips

	if len(opts) > 0 && opts[0].MappedOptions != nil {
		result.IsMapped = true
		result.MappedBy = opts[0].MappedOptions.MappedBy

		mappedPayslips := make(map[any][]entity.Payslip)
		var keyFunc func(payslip entity.Payslip) any

		switch opts[0].MappedOptions.MappedBy {
		case entity.MappedByUserID:
			keyFunc = func(payslip entity.Payslip) any {
				return payslip.UserID
			}
		case entity.MappedByPayslipID:
			keyFunc = func(payslip entity.Payslip) any {
				return payslip.ID
			}
		case entity.MappedByPayrollID:
			keyFunc = func(payslip entity.Payslip) any {
				return payslip.PayrollID
			}
		default:
			return result, errors.New("unsupported mapped by option")
		}

		for _, payslip := range payslips {
			key := keyFunc(payslip)
			mappedPayslips[key] = append(mappedPayslips[key], payslip)
		}

		result.Mapped = mappedPayslips
	}

	return result, nil
}
