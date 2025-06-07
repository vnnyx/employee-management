package repository

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/vnnyx/employee-management/internal/auth"
	"github.com/vnnyx/employee-management/internal/constants"
	"github.com/vnnyx/employee-management/internal/users/entity"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
	"golang.org/x/crypto/bcrypt"
)

type authRepo struct {
	db database.Queryer
}

func NewAuthRepository(db database.Queryer) auth.Repository {
	return &authRepo{
		db: db,
	}
}

func (r *authRepo) GetUserByUsernamePassword(ctx context.Context, username, password string) (*entity.User, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"AuthRepository.GetUserByUsernamePassword()",
	)
	defer span.End()

	user := new(entity.User)
	err := pgxscan.Get(ctx, r.db, user, findUserByUsername, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, constants.ErrWrapPgxscanGet)
	}

	// Verify the password
	valid, err := verifyPassword(user.Password, password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify password")
	}
	if !valid {
		return nil, nil
	}

	return user, nil
}

func verifyPassword(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, errors.Wrap(err, "failed to compare password")
	}

	return true, nil
}
