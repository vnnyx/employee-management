package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/internal/reimbursement"
	"github.com/vnnyx/employee-management/internal/reimbursement/entity"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
)

type reimbursementUseCase struct {
	reimbursementRepo reimbursement.Repository
}

func NewReimbursementUseCase(reimbursementRepo reimbursement.Repository) reimbursement.UseCase {
	return &reimbursementUseCase{
		reimbursementRepo: reimbursementRepo,
	}
}

func (u *reimbursementUseCase) SubmitReimbursement(ctx context.Context, authCredential authCredential.Credential, payload entity.SubmitReimbursement) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"ReimbursementUseCase.SubmitReimbursement()",
	)
	defer span.End()

	err := database.WithAuditContext(ctx, authCredential, pgx.TxOptions{}, func(tx database.DBTx) error {
		reimbursementRepoTx := u.reimbursementRepo.WithTx(tx)

		timeNow := time.Now()
		err := reimbursementRepoTx.StoreNewReimbursement(ctx, entity.Reimbursement{
			ID:                uuid.NewString(),
			UserID:            authCredential.UserID,
			Amount:            payload.Amount,
			Description:       payload.Description,
			ReimbursementDate: payload.Date,
			CreatedAt:         timeNow,
			UpdatedAt:         timeNow,
			CreatedBy:         authCredential.UserID,
			UpdatedBy:         authCredential.UserID,
			IPAddress:         authCredential.IPAddress,
		})
		if err != nil {
			return errors.Wrap(err, "ReimbursementUseCase.SubmitReimbursement().StoreNewReimbursement()")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "ReimbursementUseCase.SubmitReimbursement().WithAuditContext()")
	}

	return nil
}
