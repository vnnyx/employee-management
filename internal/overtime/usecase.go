package overtime

import (
	"context"

	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/internal/overtime/entity"
)

type UseCase interface {
	SubmitOvertime(ctx context.Context, authCredential authCredential.Credential, payload entity.SubmitOvertime) error
}
