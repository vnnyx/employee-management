package users

import (
	"context"

	"github.com/vnnyx/employee-management/internal/users/entity"
	"github.com/vnnyx/employee-management/pkg/database"
)

type Repository interface {
	WithTx(tx database.DBTx) Repository

	FindAllUsers(ctx context.Context, opts ...entity.FindUserOptions) (entity.FindUserResult, error)
	FindUserByID(ctx context.Context, userID string) (*entity.User, error)
}
