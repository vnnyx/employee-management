package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/vnnyx/employee-management/internal/attendance"
	"github.com/vnnyx/employee-management/internal/attendance/entity"
	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/pkg/apperror"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
)

type attendanceUseCase struct {
	attendanceRepo attendance.Repository
}

func NewAttendanceUseCase(attendanceRepo attendance.Repository) attendance.UseCase {
	return &attendanceUseCase{
		attendanceRepo: attendanceRepo,
	}
}

func (u *attendanceUseCase) SubmitAttendance(ctx context.Context, authCredential authCredential.Credential) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"AttendanceUseCase.SubmitAttendance()",
	)
	defer span.End()

	timeNow := time.Now()

	// Validate if current date is not weekend
	if timeNow.Weekday() == time.Saturday || timeNow.Weekday() == time.Sunday {
		return apperror.BadRequest(
			apperror.AppError{
				IssueCode: entity.AttendanceInvalidDay,
				Message:   entity.GetErrorMessageByIssueCode(entity.AttendanceInvalidDay),
			},
		)
	}

	err := database.WithAuditContext(ctx, authCredential, pgx.TxOptions{}, func(tx database.DBTx) error {
		attendanceRepoTx := u.attendanceRepo.WithTx(tx)

		err := attendanceRepoTx.UpsertAttendance(ctx, entity.Attendance{
			ID:             uuid.NewString(),
			UserID:         authCredential.UserID,
			AttendanceDate: timeNow,
			UpdatedAt:      timeNow,
			CreatedAt:      timeNow,
			IPAddress:      authCredential.IPAddress,
			CreatedBy:      authCredential.UserID,
			UpdatedBy:      authCredential.UserID,
		})
		if err != nil {
			return errors.Wrap(err, "AttendanceUseCase.SubmitAttendance().UpsertAttendance()")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "AttendanceUseCase.SubmitAttendance().WithAuditContext()")
	}

	return nil
}

func (u *attendanceUseCase) CreateAttendancePeriod(ctx context.Context, authCredential authCredential.Credential, payload entity.CreateAttendancePeriod) (string, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"AttendanceUseCase.CreateAttendancePeriod()",
	)
	defer span.End()

	if !*authCredential.IsAdmin {
		return "", apperror.Forbidden(
			apperror.AppError{
				IssueCode: entity.AttendanceNotAuthorized,
				Message:   entity.GetErrorMessageByIssueCode(entity.AttendanceNotAuthorized),
			},
		)
	}

	if payload.StartDate.After(payload.EndDate) {
		return "", apperror.BadRequest(
			apperror.AppError{
				IssueCode: entity.AttendanceInvalidPeriod,
				Message:   entity.GetErrorMessageByIssueCode(entity.AttendanceInvalidPeriod),
			},
		)
	}

	timeNow := time.Now()
	uuidString := uuid.NewString()
	err := database.WithAuditContext(ctx, authCredential, pgx.TxOptions{}, func(tx database.DBTx) error {
		attendanceRepoTx := u.attendanceRepo.WithTx(tx)

		err := attendanceRepoTx.StoreNewAttendancePeriod(ctx, entity.AttendancePeriod{
			ID:        uuidString,
			StartDate: payload.StartDate,
			EndDate:   payload.EndDate,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			CreatedBy: authCredential.UserID,
			UpdatedBy: authCredential.UserID,
			IPAddress: authCredential.IPAddress,
		})
		if err != nil {
			if database.IsUniqueViolation(err, "attendance_periods_start_date_end_date_key") {
				return apperror.BadRequest(
					apperror.AppError{
						IssueCode: entity.AttendancePeriodAlreadyExists,
						Message:   entity.GetErrorMessageByIssueCode(entity.AttendancePeriodAlreadyExists),
					},
				)
			}
			return errors.Wrap(err, "AttendanceUseCase.CreateAttendancePeriod().StoreNewAttendancePeriod()")
		}

		return nil
	})
	if err != nil {
		return "", errors.Wrap(err, "AttendanceUseCase.CreateAttendancePeriod().WithAuditContext()")
	}

	return uuidString, nil
}
