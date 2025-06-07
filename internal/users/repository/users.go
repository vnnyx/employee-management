package repository

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/vnnyx/employee-management/internal/constants"
	"github.com/vnnyx/employee-management/internal/users"
	"github.com/vnnyx/employee-management/internal/users/entity"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
)

type userRepo struct {
	db database.Queryer
}

func NewUserRepository(db database.Queryer) users.Repository {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) WithTx(tx database.DBTx) users.Repository {
	return &userRepo{
		db: tx,
	}
}

func (r *userRepo) FindAllUsers(ctx context.Context, opts ...entity.FindUserOptions) (entity.FindUserResult, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"UserRepository.FindUsers()",
	)
	defer span.End()

	var result entity.FindUserResult

	query := findUsersQuery
	if len(opts) > 0 && opts[0].PessimisticLock {
		query += " FOR UPDATE"
	}

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return result, errors.Wrap(err, constants.ErrWrapDbQuery)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.IsAdmin,
			&user.Salary,
		)
		if err != nil {
			return result, errors.Wrap(err, constants.ErrWrapDbQueryRowScan)
		}

		users = append(users, user)
	}

	result.List = users

	if len(opts) > 0 && opts[0].MappedOptions != nil {
		result.IsMapped = true
		result.MappedBy = opts[0].MappedOptions.MappedBy

		mappedUsers := make(map[any][]entity.User)
		var keyFunc func(user entity.User) any

		switch opts[0].MappedOptions.MappedBy {
		case entity.MappedByUserID:
			keyFunc = func(user entity.User) any {
				return user.ID
			}
		default:
			return result, errors.New("unsupported mapped by option")
		}

		for _, user := range users {
			key := keyFunc(user)
			mappedUsers[key] = append(mappedUsers[key], user)
		}

		result.Mapped = mappedUsers
	}

	return result, nil
}

func (r *userRepo) FindUserByID(ctx context.Context, userID string) (*entity.User, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"UserRepository.FindUserByID()",
	)
	defer span.End()

	query := findUserByIDQuery

	var user entity.User
	err := pgxscan.Get(ctx, r.db, &user, query, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, constants.ErrWrapPgxscanGet)
	}

	return &user, nil
}
