package usecase

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/vnnyx/employee-management/internal/auth"
	"github.com/vnnyx/employee-management/internal/auth/entity"
	userEntity "github.com/vnnyx/employee-management/internal/users/entity"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
)

type authUseCase struct {
	authRepo auth.Repository
	key      string
}

type AuthConfig struct {
	Key string
}

func NewAuthUseCase(authRepo auth.Repository, authConfig AuthConfig) auth.UseCase {
	return &authUseCase{
		authRepo: authRepo,
		key:      authConfig.Key,
	}
}

func (u *authUseCase) Login(ctx context.Context, username, password string) (string, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"AuthUseCase.Login()",
	)
	defer span.End()

	user, err := u.authRepo.GetUserByUsernamePassword(ctx, username, password)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", entity.ErrUserNotFound
	}

	token, err := u.generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *authUseCase) generateJWT(user *userEntity.User) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := entity.AccessTokenClaims{
		UserID:   user.ID,
		IsAdmin:  user.IsAdmin,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.key))
	if err != nil {
		return "", errors.Wrap(err, "AuthUseCase.generateJWT().SignedString()")
	}

	return tokenString, nil
}
