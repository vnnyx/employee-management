package repository_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/vnnyx/employee-management/internal/users/entity"
	"github.com/vnnyx/employee-management/internal/users/repository"
)

func TestFindUserByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewUserRepository(mock)

	tests := []struct {
		name      string
		setupMock func()
		userID    string
		wantUser  *entity.User
		expectErr bool
	}{
		{
			name:   "success - user found",
			userID: "123",
			setupMock: func() {
				rows := pgxmock.NewRows([]string{"id", "username", "is_admin", "salary"}).
					AddRow("123", "john_doe", true, 5000)
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id =").
					WithArgs("123").
					WillReturnRows(rows)
			},
			wantUser: &entity.User{
				ID:       "123",
				Username: "john_doe",
				IsAdmin:  true,
				Salary:   5000,
			},
			expectErr: false,
		},
		{
			name:   "user not found",
			userID: "456",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id =").
					WithArgs("456").
					WillReturnError(pgx.ErrNoRows)
			},
			wantUser:  nil,
			expectErr: false,
		},
		{
			name:   "db error",
			userID: "789",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id =").
					WithArgs("789").
					WillReturnError(errors.New("db failed"))
			},
			wantUser:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			user, err := repo.FindUserByID(context.Background(), tt.userID)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser, user)
			}
		})
	}
}

func TestFindAllUsers(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewUserRepository(mock)

	tests := []struct {
		name       string
		setupMock  func()
		opts       []entity.FindUserOptions
		wantCount  int
		wantMapped bool
		expectErr  bool
	}{
		{
			name: "success - no options",
			setupMock: func() {
				rows := pgxmock.NewRows([]string{"id", "username", "is_admin", "salary"}).
					AddRow("1", "user1", false, 1000).
					AddRow("2", "user2", true, 2000)
				mock.ExpectQuery("SELECT (.+) FROM users$").
					WillReturnRows(rows)
			},
			opts:       nil,
			wantCount:  2,
			wantMapped: false,
			expectErr:  false,
		},
		{
			name: "success - with mapped by user id",
			setupMock: func() {
				rows := pgxmock.NewRows([]string{"id", "username", "is_admin", "salary"}).
					AddRow("1", "user1", false, 1000).
					AddRow("2", "user2", true, 2000).
					AddRow("1", "user3", false, 1500)
				mock.ExpectQuery("SELECT (.+) FROM users$").
					WillReturnRows(rows)
			},
			opts: []entity.FindUserOptions{
				{
					MappedOptions: &entity.MappedOptions{
						MappedBy: entity.MappedByUserID,
					},
				},
			},
			wantCount:  3,
			wantMapped: true,
			expectErr:  false,
		},
		{
			name: "error on query",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users$").
					WillReturnError(errors.New("db failure"))
			},
			opts:       nil,
			wantCount:  0,
			wantMapped: false,
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			result, err := repo.FindAllUsers(context.Background(), tt.opts...)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result.List, tt.wantCount)

				if tt.wantMapped {
					assert.True(t, result.IsMapped)
					assert.Equal(t, entity.MappedByUserID, result.MappedBy)
					assert.NotEmpty(t, result.Mapped)

					if usersByID, ok := result.Mapped["1"]; ok {
						assert.GreaterOrEqual(t, len(usersByID), 1)
					} else {
						t.Errorf("expected mapped users to include key '1'")
					}
				} else {
					assert.False(t, result.IsMapped)
					assert.Empty(t, result.Mapped)
				}
			}
		})
	}
}
