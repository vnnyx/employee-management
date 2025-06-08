package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vnnyx/employee-management/internal/attendance/entity"
	mockAttendance "github.com/vnnyx/employee-management/internal/attendance/mock"
	"github.com/vnnyx/employee-management/internal/attendance/usecase"
	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/pkg/apperror"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/testutil"
	"go.uber.org/mock/gomock"
)

func TestSubmitAttendance(t *testing.T) {
	type testCase struct {
		name           string
		authCredential authCredential.Credential
		mockNow        time.Time
		expectedErr    error
		setupMock      func(repo *mockAttendance.MockRepository, txRepo *mockAttendance.MockRepository)
	}

	tests := []testCase{
		{
			name: "success - weekday attendance",
			authCredential: authCredential.Credential{
				UserID:    "user-1",
				IPAddress: "127.0.0.1",
			},
			mockNow:     time.Date(2025, 6, 4, 9, 0, 0, 0, time.UTC), // Wednesday
			expectedErr: nil,
			setupMock: func(repo, txRepo *mockAttendance.MockRepository) {
				repo.EXPECT().WithTx(gomock.Any()).Return(txRepo)

				txRepo.EXPECT().UpsertAttendance(gomock.Any(), mock.MatchedBy(func(att entity.Attendance) bool {
					expected := entity.Attendance{
						UserID:    "user-1",
						IPAddress: "127.0.0.1",
					}
					return testutil.EqualVerbose(expected, att,
						cmpopts.IgnoreFields(entity.Attendance{},
							"ID", "AttendanceDate", "CreatedAt", "UpdatedAt", "CreatedBy", "UpdatedBy"),
					)
				})).Return(nil)
			},
		},
		{
			name: "error - weekend attendance",
			authCredential: authCredential.Credential{
				UserID:    "user-2",
				IPAddress: "192.168.0.1",
			},
			mockNow: time.Date(2025, 6, 7, 10, 0, 0, 0, time.UTC), // Saturday
			expectedErr: apperror.BadRequest(
				apperror.AppError{
					IssueCode: entity.AttendanceInvalidDay,
					Message:   entity.GetErrorMessageByIssueCode(entity.AttendanceInvalidDay),
				}),
			setupMock: func(repo, txRepo *mockAttendance.MockRepository) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patch := gomonkey.NewPatches()
			defer patch.Reset()

			patch.ApplyFunc(time.Now, func() time.Time {
				return tt.mockNow
			})

			patch.ApplyFunc(database.WithAuditContext, func(
				ctx context.Context,
				cred authCredential.Credential,
				txOpt pgx.TxOptions,
				fn func(tx database.DBTx) error,
			) error {
				return fn(nil)
			})

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockAttendance.NewMockRepository(ctrl)
			mockRepoTx := mockAttendance.NewMockRepository(ctrl)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo, mockRepoTx)
			}

			useCase := usecase.NewAttendanceUseCase(mockRepo)

			err := useCase.SubmitAttendance(context.Background(), tt.authCredential)

			if tt.expectedErr != nil {
				assert.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateAttendancePeriod(t *testing.T) {
	type testCase struct {
		name           string
		authCredential authCredential.Credential
		payload        entity.CreateAttendancePeriod
		expectedID     string
		expectedErr    error
		setupMock      func(repo *mockAttendance.MockRepository, txRepo *mockAttendance.MockRepository)
	}

	tests := []testCase{
		{
			name: "success - create attendance period",
			authCredential: authCredential.Credential{
				UserID:    "user-1",
				IPAddress: "127.0.0.1",
				Username:  "testuser",
				IsAdmin:   func(b bool) *bool { return &b }(true),
				RequestID: "req-123",
			},
			payload: entity.CreateAttendancePeriod{
				StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 6, 30, 23, 59, 59, 999999999, time.UTC),
			},
			expectedID:  "period-123",
			expectedErr: nil,
			setupMock: func(repo, txRepo *mockAttendance.MockRepository) {
				repo.EXPECT().WithTx(gomock.Any()).Return(txRepo)

				txRepo.EXPECT().StoreNewAttendancePeriod(gomock.Any(), mock.MatchedBy(func(period entity.AttendancePeriod) bool {
					expected := entity.AttendancePeriod{
						StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
						EndDate:   time.Date(2025, 6, 30, 23, 59, 59, 999999999, time.UTC),
						CreatedBy: "user-1",
						UpdatedBy: "user-1",
						IPAddress: "127.0.0.1",
					}
					return testutil.EqualVerbose(expected, period,
						cmpopts.IgnoreFields(entity.AttendancePeriod{},
							"ID", "CreatedAt", "UpdatedAt"),
					)
				})).Return(nil)

			},
		},
		{
			name: "error - non-admin user",
			authCredential: authCredential.Credential{
				UserID:    "user-2",
				IPAddress: "127.0.0.1",
				Username:  "testuser2",
				IsAdmin:   func(b bool) *bool { return &b }(false),
			},
			payload: entity.CreateAttendancePeriod{
				StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 6, 30, 23, 59, 59, 999999999, time.UTC),
			},
			expectedID: "",
			expectedErr: apperror.Forbidden(
				apperror.AppError{
					IssueCode: entity.AttendanceNotAuthorized,
					Message:   entity.GetErrorMessageByIssueCode(entity.AttendanceNotAuthorized),
				}),
			setupMock: func(repo, txRepo *mockAttendance.MockRepository) {},
		},
		{
			name: "error - invalid period dates",
			authCredential: authCredential.Credential{
				UserID:    "user-3",
				IPAddress: "127.0.0.1",
				Username:  "testuser3",
				IsAdmin:   func(b bool) *bool { return &b }(true),
			},
			payload: entity.CreateAttendancePeriod{
				StartDate: time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 6, 1, 23, 59, 59, 999999999, time.UTC),
			},
			expectedID: "",
			expectedErr: apperror.BadRequest(
				apperror.AppError{
					IssueCode: entity.AttendanceInvalidPeriod,
					Message:   entity.GetErrorMessageByIssueCode(entity.AttendanceInvalidPeriod),
				}),
			setupMock: func(repo, txRepo *mockAttendance.MockRepository) {},
		},
		{
			name: "error - unique constraint violation",
			authCredential: authCredential.Credential{
				UserID:    "user-4",
				IPAddress: "127.0.0.1",
				Username:  "testuser4",
				IsAdmin:   func(b bool) *bool { return &b }(true),
			},
			payload: entity.CreateAttendancePeriod{
				StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 6, 30, 23, 59, 59, 999999999, time.UTC),
			},
			expectedID: "",
			expectedErr: apperror.BadRequest(
				apperror.AppError{
					IssueCode: entity.AttendancePeriodAlreadyExists,
					Message:   entity.GetErrorMessageByIssueCode(entity.AttendancePeriodAlreadyExists),
				}),
			setupMock: func(repo, txRepo *mockAttendance.MockRepository) {
				repo.EXPECT().WithTx(gomock.Any()).Return(txRepo)

				txRepo.EXPECT().StoreNewAttendancePeriod(gomock.Any(), mock.MatchedBy(func(period entity.AttendancePeriod) bool {
					expected := entity.AttendancePeriod{
						StartDate: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
						EndDate:   time.Date(2025, 6, 30, 23, 59, 59, 999999999, time.UTC),
						CreatedBy: "user-4",
						UpdatedBy: "user-4",
						IPAddress: "127.0.0.1",
					}
					return testutil.EqualVerbose(expected, period,
						cmpopts.IgnoreFields(entity.AttendancePeriod{},
							"ID", "CreatedAt", "UpdatedAt"),
					)
				})).Return(errors.New("attendance_periods_start_date_end_date_key"))
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

			if tt.expectedID != "" {
				patches.ApplyFunc(uuid.NewString, func() string {
					return tt.expectedID
				})
			}

			if tt.name == "error - unique constraint violation" {
				patches.ApplyFunc(database.IsUniqueViolation, func(err error, constraintName string) bool {
					return constraintName == "attendance_periods_start_date_end_date_key"
				})
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockAttendance.NewMockRepository(ctrl)
			mockRepoTx := mockAttendance.NewMockRepository(ctrl)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo, mockRepoTx)
			}

			useCase := usecase.NewAttendanceUseCase(mockRepo)

			id, err := useCase.CreateAttendancePeriod(context.Background(), tt.authCredential, tt.payload)

			if tt.expectedErr != nil {
				assert.ErrorContains(t, err, tt.expectedErr.Error())
				assert.Empty(t, id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}
