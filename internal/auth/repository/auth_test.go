package repository_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vnnyx/employee-management/internal/auth/repository"
	"golang.org/x/crypto/bcrypt"
)

func TestGetUserByUsernamePassword(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewAuthRepository(mock)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		username    string
		password    string
		setupMock   func()
		expectedNil bool
		expectErr   bool
	}{
		{
			name:     "success - user found and password matches",
			username: "john",
			password: "correct-password",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("john").
					WillReturnRows(pgxmock.NewRows([]string{"id", "username", "is_admin", "password"}).
						AddRow("user-1", "john", false, string(hashedPassword)))
			},
			expectedNil: false,
			expectErr:   false,
		},
		{
			name:     "wrong password",
			username: "john",
			password: "wrong-password",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("john").
					WillReturnRows(pgxmock.NewRows([]string{"id", "username", "is_admin", "password"}).
						AddRow("user-1", "john", false, string(hashedPassword)))
			},
			expectedNil: true,
			expectErr:   false,
		},
		{
			name:     "user not found",
			username: "missing-user",
			password: "whatever",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("missing-user").
					WillReturnError(pgx.ErrNoRows)
			},
			expectedNil: true,
			expectErr:   false,
		},
		{
			name:     "query error",
			username: "error-user",
			password: "irrelevant",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("error-user").
					WillReturnError(assert.AnError)
			},
			expectedNil: true,
			expectErr:   true,
		},
		{
			name:     "invalid hash format",
			username: "broken-hash",
			password: "somepass",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("broken-hash").
					WillReturnRows(pgxmock.NewRows([]string{"id", "username", "is_admin", "password"}).
						AddRow("user-99", "broken-hash", false, "not-a-valid-bcrypt-hash"))
			},
			expectedNil: true,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			user, err := repo.GetUserByUsernamePassword(context.Background(), tt.username, tt.password)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectedNil {
					assert.Nil(t, user)
				} else {
					assert.NotNil(t, user)
					assert.Equal(t, tt.username, user.Username)
				}
			}
		})
	}
}
