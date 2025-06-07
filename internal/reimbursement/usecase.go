package reimbursement

import (
	"context"

	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/internal/reimbursement/entity"
)

type UseCase interface {
	SubmitReimbursement(ctx context.Context, authCredential authCredential.Credential, payload entity.SubmitReimbursement) error
}
