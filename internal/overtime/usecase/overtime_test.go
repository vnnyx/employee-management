package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/internal/overtime/entity"
	mockOvertime "github.com/vnnyx/employee-management/internal/overtime/mock"
	"github.com/vnnyx/employee-management/internal/overtime/usecase"
	"github.com/vnnyx/employee-management/pkg/apperror"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/testutil"
	"go.uber.org/mock/gomock"
)

func TestSubmitOvertime(t *testing.T) {
	type testCase struct {
		name           string
		authCredential authCredential.Credential
		payload        entity.SubmitOvertime
		overtime       entity.Overtime
		expectedErr    error
		mockNow        *time.Time
		setupMock      func(repo *mockOvertime.MockRepository, txRepo *mockOvertime.MockRepository)
	}

	tests := []testCase{
		{
			name: "success - overtime submitted",
			authCredential: authCredential.Credential{
				UserID:    "user-1",
				IPAddress: "127.0.0.1",
				Username:  "testuser",
				IsAdmin:   func(b bool) *bool { return &b }(false),
				RequestID: "req-123",
			},
			payload: entity.SubmitOvertime{
				OvertimeDate: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
				Overtime:     "PT3H",
			},
			overtime: entity.Overtime{
				UserID:        "user-1",
				OverTimeDate:  time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
				OvertimeHours: 3 * time.Hour,
			},
			expectedErr: nil,
			setupMock: func(repo, txRepo *mockOvertime.MockRepository) {
				repo.EXPECT().WithTx(gomock.Any()).Return(txRepo)
				txRepo.EXPECT().
					FindOvertimeByUserIDDate(gomock.Any(), "user-1", time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)).
					Return(nil, nil)
				txRepo.EXPECT().
					UpsertOvertime(gomock.Any(), mock.MatchedBy(func(args entity.Overtime) bool {
						return testutil.EqualVerbose(
							entity.Overtime{
								UserID:        "user-1",
								OverTimeDate:  time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
								OvertimeHours: 3 * time.Hour,
							},
							args,
							cmpopts.IgnoreFields(entity.Overtime{},
								"ID", "CreatedAt", "UpdatedAt", "CreatedBy", "UpdatedBy", "IPAddress",
							))
					})).
					Return(nil)
			},
		},
		{
			name: "error - overtime exceeds limit",
			authCredential: authCredential.Credential{
				UserID:    "user-1",
				IPAddress: "127.0.0.1",
				Username:  "testuser",
				IsAdmin:   func(b bool) *bool { return &b }(false),
				RequestID: "req-123",
			},
			payload: entity.SubmitOvertime{
				OvertimeDate: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
				Overtime:     "PT4H",
			},
			expectedErr: apperror.BadRequest(
				apperror.AppError{
					IssueCode: entity.OvertimeExceedsLimit,
					Message:   entity.GetErrorMessageByIssueCode(entity.OvertimeExceedsLimit),
				}),
			setupMock: func(repo, txRepo *mockOvertime.MockRepository) {
				repo.EXPECT().WithTx(gomock.Any()).Return(txRepo)
				txRepo.EXPECT().
					FindOvertimeByUserIDDate(gomock.Any(), "user-1", time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)).
					Return(nil, nil)
				// No Upsert expected
			},
		},
		{
			name: "error - invalid overtime time request",
			authCredential: authCredential.Credential{
				UserID:    "user-1",
				IPAddress: "127.0.0.1",
				Username:  "testuser",
				IsAdmin:   func(b bool) *bool { return &b }(false),
				RequestID: "req-123",
			},
			payload: entity.SubmitOvertime{
				OvertimeDate: time.Date(2023, 10, 1, 9, 0, 0, 0, time.UTC),
				Overtime:     "PT2H",
			},
			mockNow: ptrTime(time.Date(2023, 10, 2, 10, 0, 0, 0, time.UTC)), // Monday 10 AM
			expectedErr: apperror.BadRequest(
				apperror.AppError{
					IssueCode: entity.OvertimeInvalidTimeRequest,
					Message:   entity.GetErrorMessageByIssueCode(entity.OvertimeInvalidTimeRequest),
				}),
			setupMock: func(repo, txRepo *mockOvertime.MockRepository) {
			},
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

			if tt.mockNow != nil {
				patches.ApplyFunc(time.Now, func() time.Time {
					return *tt.mockNow
				})
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockOvertime.NewMockRepository(ctrl)
			mockRepoTx := mockOvertime.NewMockRepository(ctrl)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo, mockRepoTx)
			}

			useCase := usecase.NewOvertimeUseCase(mockRepo)

			err := useCase.SubmitOvertime(context.Background(), tt.authCredential, tt.payload)

			if tt.expectedErr != nil {
				assert.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
