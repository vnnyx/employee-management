package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/internal/overtime"
	"github.com/vnnyx/employee-management/internal/overtime/entity"
	"github.com/vnnyx/employee-management/pkg/apperror"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/iso8601"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
)

type overtimeUseCase struct {
	overtimeRepo overtime.Repository
}

func NewOvertimeUseCase(overtimeRepo overtime.Repository) overtime.UseCase {
	return &overtimeUseCase{
		overtimeRepo: overtimeRepo,
	}
}

func (u *overtimeUseCase) SubmitOvertime(ctx context.Context, authCredential authCredential.Credential, payload entity.SubmitOvertime) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"OvertimeUseCase.SubmitOvertime()",
	)
	defer span.End()

	timeNow := time.Now()
	// Check if request overtime outside working hours
	if timeNow.Hour() >= 9 && timeNow.Hour() <= 17 && timeNow.Weekday() != time.Saturday && timeNow.Weekday() != time.Sunday {
		return apperror.BadRequest(
			apperror.AppError{
				IssueCode: entity.OvertimeInvalidTimeRequest,
				Message:   entity.GetErrorMessageByIssueCode(entity.OvertimeInvalidTimeRequest),
			},
		)
	}

	err := database.WithAuditContext(ctx, authCredential, pgx.TxOptions{}, func(tx database.DBTx) error {
		overtimeRepoTx := u.overtimeRepo.WithTx(tx)

		overtime, err := overtimeRepoTx.FindOvertimeByUserIDDate(ctx, authCredential.UserID, payload.OvertimeDate)
		if err != nil {
			return errors.Wrap(err, "OvertimeUseCase.SubmitOvertime().FindOvertimeByUserIDDate()")
		}

		newDuration := iso8601.MustParse(payload.Overtime)
		totalDuration := newDuration

		if newDuration > 3*time.Hour {
			return apperror.BadRequest(
				apperror.AppError{
					IssueCode: entity.OvertimeExceedsLimit,
					Message:   entity.GetErrorMessageByIssueCode(entity.OvertimeExceedsLimit),
					Path:      []string{"overtime"},
				},
			)
		}

		if overtime != nil {
			currentDuration := overtime.OvertimeHours
			totalDuration = currentDuration + newDuration
			if totalDuration > 3*time.Hour {
				return apperror.BadRequest(
					apperror.AppError{
						IssueCode: entity.OvertimeExceedsLimit,
						Message:   entity.GetErrorMessageByIssueCode(entity.OvertimeExceedsLimit),
						Path:      []string{"overtime"},
					},
				)
			}
		}

		err = overtimeRepoTx.UpsertOvertime(ctx, entity.Overtime{
			ID:            uuid.NewString(),
			UserID:        authCredential.UserID,
			OverTimeDate:  payload.OvertimeDate,
			OvertimeHours: totalDuration,
			CreatedAt:     timeNow,
			UpdatedAt:     timeNow,
			CreatedBy:     authCredential.UserID,
			UpdatedBy:     authCredential.UserID,
			IPAddress:     authCredential.IPAddress,
		})
		if err != nil {
			return errors.Wrap(err, "OvertimeUseCase.SubmitOvertime().UpsertOvertime()")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "OvertimeUseCase.SubmitOvertime().WithAuditContext()")
	}

	return nil
}
