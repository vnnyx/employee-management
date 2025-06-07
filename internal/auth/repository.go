package auth

import (
	"context"

	"github.com/vnnyx/employee-management/internal/users/entity"
)

type Repository interface {
	GetUserByUsernamePassword(ctx context.Context, username, password string) (*entity.User, error)
}
