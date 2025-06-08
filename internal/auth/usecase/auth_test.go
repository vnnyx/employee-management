package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vnnyx/employee-management/internal/auth/entity"
	mockauth "github.com/vnnyx/employee-management/internal/auth/mock"
	"github.com/vnnyx/employee-management/internal/auth/usecase"
	userEntity "github.com/vnnyx/employee-management/internal/users/entity"
	"go.uber.org/mock/gomock"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		password    string
		mockUser    *userEntity.User
		mockError   error
		expectedErr error
		expectJWT   bool
	}{
		{
			name:     "success",
			username: "admin",
			password: "adminpass",
			mockUser: &userEntity.User{
				ID:       "user-1",
				Username: "admin",
				IsAdmin:  true,
			},
			mockError:   nil,
			expectedErr: nil,
			expectJWT:   true,
		},
		{
			name:        "user not found",
			username:    "notfound",
			password:    "password",
			mockUser:    nil,
			mockError:   nil,
			expectedErr: entity.ErrUserNotFound,
			expectJWT:   false,
		},
		{
			name:        "repo returns error",
			username:    "someuser",
			password:    "somepass",
			mockUser:    nil,
			mockError:   errors.New("db failure"),
			expectedErr: errors.New("db failure"),
			expectJWT:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthRepo := mockauth.NewMockRepository(ctrl)
			useCase := usecase.NewAuthUseCase(mockAuthRepo, usecase.AuthConfig{Key: "test-secret"})

			mockAuthRepo.
				EXPECT().
				GetUserByUsernamePassword(gomock.Any(), tt.username, tt.password).
				Return(tt.mockUser, tt.mockError)

			token, err := useCase.Login(context.Background(), tt.username, tt.password)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			if tt.expectJWT {
				assert.NotEmpty(t, token)

				// Optional: decode token to ensure correctness
				parsed, _ := jwt.ParseWithClaims(token, &entity.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte("test-secret"), nil
				})
				assert.True(t, parsed.Valid)
			} else {
				assert.Empty(t, token)
			}
		})
	}
}
