package auth

import "context"

type UseCase interface {
	Login(ctx context.Context, username, password string) (string, error)
}
