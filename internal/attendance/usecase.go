package attendance

import (
	"context"

	"github.com/vnnyx/employee-management/internal/attendance/entity"
	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
)

type UseCase interface {
	SubmitAttendance(ctx context.Context, authCredential authCredential.Credential) error
	CreateAttendancePeriod(ctx context.Context, authCredential authCredential.Credential, payload entity.CreateAttendancePeriod) (string, error)
}
