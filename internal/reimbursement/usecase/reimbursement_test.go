package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/internal/reimbursement/entity"
	mockreimbursement "github.com/vnnyx/employee-management/internal/reimbursement/mock"
	"github.com/vnnyx/employee-management/internal/reimbursement/usecase"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/optional"
	"go.uber.org/mock/gomock"
)

func TestSubmitReimbursement(t *testing.T) {
	tests := []struct {
		name           string
		authCredential authCredential.Credential
		payload        entity.SubmitReimbursement
		mockError      error
		expectedErr    error
	}{
		{
			name: "success - reimbursement submitted",
			authCredential: authCredential.Credential{
				UserID:    "user-1",
				IPAddress: "127.0.0.1",
				Username:  "tester",
			},
			payload: entity.SubmitReimbursement{
				Amount:      100000,
				Date:        time.Date(2025, 6, 7, 0, 0, 0, 0, time.UTC),
				Description: optional.NewString("Lunch with client"),
			},
			mockError:   nil,
			expectedErr: nil,
		},
		{
			name: "error - failed to store reimbursement",
			authCredential: authCredential.Credential{
				UserID:    "user-2",
				IPAddress: "127.0.0.1",
				Username:  "tester2",
			},
			payload: entity.SubmitReimbursement{
				Amount:      200000,
				Date:        time.Date(2025, 6, 7, 0, 0, 0, 0, time.UTC),
				Description: optional.NewString("Taxi to airport"),
			},
			mockError: errors.New("db error"),
			expectedErr: errors.New(
				"ReimbursementUseCase.SubmitReimbursement().WithAuditContext(): ReimbursementUseCase.SubmitReimbursement().StoreNewReimbursement(): db error",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.ApplyFunc(database.WithAuditContext, func(
				ctx context.Context,
				cred authCredential.Credential,
				txOpt pgx.TxOptions,
				fn func(tx database.DBTx) error,
			) error {
				return fn(nil)
			})
			defer patches.Reset()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockreimbursement.NewMockRepository(ctrl)
			useCase := usecase.NewReimbursementUseCase(mockRepo)

			mockRepoTx := mockreimbursement.NewMockRepository(ctrl)
			mockRepo.EXPECT().WithTx(gomock.Any()).Return(mockRepoTx)

			mockRepoTx.
				EXPECT().
				StoreNewReimbursement(gomock.Any(), gomock.Any()).
				Return(tt.mockError)

			err := useCase.SubmitReimbursement(context.Background(), tt.authCredential, tt.payload)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
