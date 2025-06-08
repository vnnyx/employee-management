package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	attEntity "github.com/vnnyx/employee-management/internal/attendance/entity"
	mockAtt "github.com/vnnyx/employee-management/internal/attendance/mock"
	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
	overtimeEntity "github.com/vnnyx/employee-management/internal/overtime/entity"
	mockOvertime "github.com/vnnyx/employee-management/internal/overtime/mock"
	"github.com/vnnyx/employee-management/internal/payroll/entity"
	mockPayroll "github.com/vnnyx/employee-management/internal/payroll/mock"
	"github.com/vnnyx/employee-management/internal/payroll/usecase"
	reimbursementEntity "github.com/vnnyx/employee-management/internal/reimbursement/entity"
	mockReimbursement "github.com/vnnyx/employee-management/internal/reimbursement/mock"
	userEntity "github.com/vnnyx/employee-management/internal/users/entity"
	mockUser "github.com/vnnyx/employee-management/internal/users/mock"
	"github.com/vnnyx/employee-management/pkg/apperror"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/optional"
	"github.com/vnnyx/employee-management/pkg/testutil"
	"go.uber.org/mock/gomock"
)

func TestGeneratePayroll(t *testing.T) {
	type mockParams struct {
		payrollRepo         *mockPayroll.MockRepository
		payrollRepoTx       *mockPayroll.MockRepository
		userRepo            *mockUser.MockRepository
		userRepoTx          *mockUser.MockRepository
		attRepo             *mockAtt.MockRepository
		attRepoTx           *mockAtt.MockRepository
		overTimeRepo        *mockOvertime.MockRepository
		overTimeRepoTx      *mockOvertime.MockRepository
		reimbursementRepo   *mockReimbursement.MockRepository
		reimbursementRepoTx *mockReimbursement.MockRepository
	}

	type setupMockFunc func(mockParams)

	type testCase struct {
		name             string
		authCredential   authCredential.Credential
		periodID         string
		generatedPayroll entity.GeneratedPayroll
		expectedErr      error
		setupMock        setupMockFunc
		patched          func()
	}

	tests := []testCase{
		{
			name: "success - payroll generated",
			authCredential: authCredential.Credential{
				UserID:    "admin-1",
				IPAddress: "127.0.0.1",
				Username:  "admin",
				IsAdmin:   func(b bool) *bool { return &b }(true),
				RequestID: "req-123",
			},
			periodID: "period-1",
			generatedPayroll: entity.GeneratedPayroll{
				PeriodID:      "period-1",
				PayrollID:     "payroll-1",
				TotalTakeHome: 199,
				TotalEmployee: 1,
				TotalPayslip:  1,
				GeneratedBy:   "admin-1",
				GeneratedAt:   "2023-10-01T00:00:00Z",
			},
			expectedErr: nil,
			setupMock: func(m mockParams) {
				m.payrollRepo.EXPECT().WithTx(gomock.Any()).Return(m.payrollRepoTx)
				m.userRepo.EXPECT().WithTx(gomock.Any()).Return(m.userRepoTx)
				m.attRepo.EXPECT().WithTx(gomock.Any()).Return(m.attRepoTx)
				m.overTimeRepo.EXPECT().WithTx(gomock.Any()).Return(m.overTimeRepoTx)
				m.reimbursementRepo.EXPECT().WithTx(gomock.Any()).Return(m.reimbursementRepoTx)

				m.payrollRepoTx.EXPECT().FindPayrollByPeriodID(gomock.Any(), "period-1", entity.FindPayrollOptions{
					PessimisticLock: true,
				}).Return(nil, nil)

				m.attRepoTx.EXPECT().FindPeriodByID(gomock.Any(), "period-1").Return(&attEntity.AttendancePeriod{
					ID:        "period-1",
					StartDate: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC),
				}, nil)

				m.userRepoTx.EXPECT().FindAllUsers(gomock.Any(), userEntity.FindUserOptions{
					PessimisticLock: true,
					MappedOptions: &userEntity.MappedOptions{
						MappedBy: userEntity.MappedByUserID,
					},
				}).Return(userEntity.FindUserResult{
					List: []userEntity.User{
						{
							ID:       "user-1",
							Username: "testuser",
							IsAdmin:  false,
							Salary:   1000.00,
						},
					},
					Mapped: map[any][]userEntity.User{
						"user-1": {
							{
								ID:       "user-1",
								Username: "testuser",
								IsAdmin:  false,
								Salary:   1000.00,
							},
						},
					},
					IsMapped: true,
					MappedBy: userEntity.MappedByUserID,
				}, nil)

				m.attRepoTx.EXPECT().FindAttendanceByPeriod(gomock.Any(), time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC), attEntity.FindAttendanceOptions{
					PessimisticLock: true,
					MappedOptions: &attEntity.MappedOptions{
						MappedBy: attEntity.MappedByUserID,
					},
				}).Return(attEntity.FindAttendanceResult{
					List: []attEntity.Attendance{
						{
							ID:             "att-1",
							UserID:         "user-1",
							AttendanceDate: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					Mapped: map[any][]attEntity.Attendance{
						"user-1": {
							{
								ID:             "att-1",
								UserID:         "user-1",
								AttendanceDate: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
							},
						},
					},
					IsMapped: true,
					MappedBy: attEntity.MappedByUserID,
				}, nil)

				m.overTimeRepoTx.EXPECT().FindOvertimeByPeriod(gomock.Any(), time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC), overtimeEntity.FindOvertimeOptions{
					PessimisticLock: true,
					MappedOptions: &overtimeEntity.MappedOptions{
						MappedBy: overtimeEntity.MappedByUserID,
					},
				}).Return(overtimeEntity.FindOvertimeResult{
					List: []overtimeEntity.Overtime{
						{
							ID:            "overtime-1",
							UserID:        "user-1",
							OverTimeDate:  time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
							OvertimeHours: 2 * time.Hour,
						},
						{
							ID:            "overtime-2",
							UserID:        "user-1",
							OverTimeDate:  time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
							OvertimeHours: 3 * time.Hour,
						},
					},
					Mapped: map[any][]overtimeEntity.Overtime{
						"user-1": {
							{
								ID:            "overtime-1",
								UserID:        "user-1",
								OverTimeDate:  time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
								OvertimeHours: 2 * time.Hour,
							},
							{
								ID:            "overtime-2",
								UserID:        "user-1",
								OverTimeDate:  time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
								OvertimeHours: 3 * time.Hour,
							},
						},
					},
					IsMapped: true,
					MappedBy: overtimeEntity.MappedByUserID,
				}, nil)

				m.reimbursementRepoTx.EXPECT().FindReimbursementByPeriod(gomock.Any(), time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC), reimbursementEntity.FindReimbursementOptions{
					PessimisticLock: true,
					MappedOptions: &reimbursementEntity.MappedOptions{
						MappedBy: reimbursementEntity.MappedByUserID,
					},
				}).Return(reimbursementEntity.FindReimbursementResult{
					List: []reimbursementEntity.Reimbursement{
						{
							ID:                "reimbursement-1",
							UserID:            "user-1",
							Amount:            100.00,
							Description:       optional.NewString("Travel Expenses"),
							ReimbursementDate: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							ID:                "reimbursement-2",
							UserID:            "user-1",
							Amount:            50.00,
							Description:       optional.NewString("Meal Expenses"),
							ReimbursementDate: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
						},
					},
					Mapped: map[any][]reimbursementEntity.Reimbursement{
						"user-1": {
							{
								ID:                "reimbursement-1",
								UserID:            "user-1",
								Amount:            100.00,
								Description:       optional.NewString("Travel Expenses"),
								ReimbursementDate: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
							},
							{
								ID:                "reimbursement-2",
								UserID:            "user-1",
								Amount:            50.00,
								Description:       optional.NewString("Meal Expenses"),
								ReimbursementDate: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
							},
						},
					},
					IsMapped: true,
					MappedBy: reimbursementEntity.MappedByUserID,
				}, nil)

				m.payrollRepoTx.EXPECT().StoreNewPayslips(gomock.Any(), mock.MatchedBy(func(args []entity.Payslip) bool {
					return testutil.EqualVerbose(
						[]entity.Payslip{
							{
								ID:                 "payslip-1",
								UserID:             "user-1",
								PayrollID:          "payroll-1",
								BaseSalary:         1000.00,
								AttendanceDays:     1,
								OvertimeHours:      optional.NewDuration(5 * time.Hour),
								ReimbursementTotal: 150.00,
							},
						},
						args,
						cmpopts.IgnoreFields(entity.Payslip{},
							"ID",
							"PayrollID",
							"OvertimePay",
							"TotalTakeHome",
							"CreatedAt",
							"UpdatedAt",
							"CreatedBy",
							"UpdatedBy",
							"IPAddress",
						),
					)
				}))

				m.payrollRepoTx.EXPECT().StoreNewPayroll(gomock.Any(), mock.MatchedBy(func(args entity.Payroll) bool {
					return testutil.EqualVerbose(
						entity.Payroll{
							ID:       "payroll-1",
							PeriodID: "period-1",
							RunBy:    "admin-1",
							RunAt:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
						args,
						cmpopts.IgnoreFields(entity.Payroll{},
							"ID",
							"CreatedAt",
							"UpdatedAt",
							"CreatedBy",
							"UpdatedBy",
							"IPAddress",
						),
					)
				}))

				m.payrollRepoTx.EXPECT().StoreNewPayrollSummary(gomock.Any(), mock.MatchedBy(func(args entity.PayrollSummary) bool {
					return testutil.EqualVerbose(
						entity.PayrollSummary{
							ID:            "payroll-summary-1",
							PayrollID:     "payroll-1",
							TotalTakeHome: 199,
							GeneratedBy:   "admin-1",
							GeneratedAt:   time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						},
						args,
						cmpopts.IgnoreFields(entity.PayrollSummary{},
							"ID",
							"CreatedAt",
							"UpdatedAt",
							"CreatedBy",
							"UpdatedBy",
							"IPAddress",
						),
					)
				}))
			},
			patched: func() {
				gomonkey.ApplyFunc(uuid.NewString, func() string {
					return "payroll-1"
				})
				gomonkey.ApplyFunc(time.Now, func() time.Time {
					return time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
				},
				)
			},
		},
		{
			name: "error - payroll already exists",
			authCredential: authCredential.Credential{
				UserID:    "admin-1",
				IPAddress: "127.0.0.1",
				Username:  "admin",
				IsAdmin:   func(b bool) *bool { return &b }(true),
				RequestID: "req-123",
			},
			periodID: "period-1",
			expectedErr: apperror.BadRequest(
				apperror.AppError{
					IssueCode: entity.PayrollAlreadyGenerated,
					Message:   entity.GetErrorMessageByIssueCode(entity.PayrollAlreadyGenerated),
				},
			),
			setupMock: func(m mockParams) {
				m.payrollRepo.EXPECT().WithTx(gomock.Any()).Return(m.payrollRepoTx)
				m.userRepo.EXPECT().WithTx(gomock.Any()).Return(m.userRepoTx)
				m.attRepo.EXPECT().WithTx(gomock.Any()).Return(m.attRepoTx)
				m.overTimeRepo.EXPECT().WithTx(gomock.Any()).Return(m.overTimeRepoTx)
				m.reimbursementRepo.EXPECT().WithTx(gomock.Any()).Return(m.reimbursementRepoTx)

				m.payrollRepoTx.EXPECT().FindPayrollByPeriodID(gomock.Any(), "period-1", entity.FindPayrollOptions{
					PessimisticLock: true,
				}).Return(&entity.Payroll{
					ID:       "payroll-1",
					PeriodID: "period-1",
					RunBy:    "admin-1",
					RunAt:    time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
		},
		{
			name: "error - period not found",
			authCredential: authCredential.Credential{
				UserID:    "admin-1",
				IPAddress: "127.0.0.1",
				Username:  "admin",
				IsAdmin:   func(b bool) *bool { return &b }(true),
				RequestID: "req-123",
			},
			periodID: "invalid-period",
			expectedErr: apperror.NotFound(
				apperror.AppError{
					IssueCode: entity.AttendancePeriodNotFound,
					Message:   entity.GetErrorMessageByIssueCode(entity.AttendancePeriodNotFound),
					Received:  "invalid-period",
				},
			),
			setupMock: func(m mockParams) {
				m.payrollRepo.EXPECT().WithTx(gomock.Any()).Return(m.payrollRepoTx)
				m.userRepo.EXPECT().WithTx(gomock.Any()).Return(m.userRepoTx)
				m.attRepo.EXPECT().WithTx(gomock.Any()).Return(m.attRepoTx)
				m.overTimeRepo.EXPECT().WithTx(gomock.Any()).Return(m.overTimeRepoTx)
				m.reimbursementRepo.EXPECT().WithTx(gomock.Any()).Return(m.reimbursementRepoTx)

				m.payrollRepoTx.EXPECT().FindPayrollByPeriodID(gomock.Any(), "invalid-period", entity.FindPayrollOptions{
					PessimisticLock: true,
				}).Return(nil, nil)
				m.attRepoTx.EXPECT().FindPeriodByID(gomock.Any(), "invalid-period").Return(nil, nil)
			},
		},
		{
			name: "error - user not admin",
			authCredential: authCredential.Credential{
				UserID:    "user-1",
				IPAddress: "127.0.0.1",
				Username:  "testuser",
				IsAdmin:   func(b bool) *bool { return &b }(false),
				RequestID: "req-123",
			},
			periodID:         "period-1",
			generatedPayroll: entity.GeneratedPayroll{},
			expectedErr: apperror.Forbidden(
				apperror.AppError{
					IssueCode: entity.PayrollNotAuthorized,
					Message:   entity.GetErrorMessageByIssueCode(entity.PayrollNotAuthorized),
				},
			),
			setupMock: func(m mockParams) {},
			patched:   func() {},
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

			mockParams := mockParams{
				payrollRepo:         mockPayroll.NewMockRepository(ctrl),
				payrollRepoTx:       mockPayroll.NewMockRepository(ctrl),
				userRepo:            mockUser.NewMockRepository(ctrl),
				userRepoTx:          mockUser.NewMockRepository(ctrl),
				attRepo:             mockAtt.NewMockRepository(ctrl),
				attRepoTx:           mockAtt.NewMockRepository(ctrl),
				overTimeRepo:        mockOvertime.NewMockRepository(ctrl),
				overTimeRepoTx:      mockOvertime.NewMockRepository(ctrl),
				reimbursementRepo:   mockReimbursement.NewMockRepository(ctrl),
				reimbursementRepoTx: mockReimbursement.NewMockRepository(ctrl),
			}
			if tt.setupMock != nil {
				tt.setupMock(mockParams)
			}

			if tt.patched != nil {
				tt.patched()
			}

			useCase := usecase.NewPayrollUseCase(
				mockParams.payrollRepo,
				mockParams.userRepo,
				mockParams.attRepo,
				mockParams.overTimeRepo,
				mockParams.reimbursementRepo,
			)
			result, err := useCase.GeneratePayroll(context.Background(), tt.authCredential, tt.periodID)

			if tt.expectedErr != nil {
				assert.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.generatedPayroll.PeriodID, result.PeriodID)
			}
		})
	}
}

func TestShowPay(t *testing.T) {
	type mockParams struct {
		payrollRepo       *mockPayroll.MockRepository
		userRepo          *mockUser.MockRepository
		attendanceRepo    *mockAtt.MockRepository
		overtimeRepo      *mockOvertime.MockRepository
		reimbursementRepo *mockReimbursement.MockRepository
	}

	type setupMockFunc func(mockParams)

	type testCase struct {
		name            string
		authCredential  authCredential.Credential
		payrollID       string
		expectedPayslip *entity.Payslip
		expectedErr     error
		setupMock       setupMockFunc
		patched         func()
	}

	tests := []testCase{
		{
			name: "success - payslip found",
			authCredential: authCredential.Credential{
				UserID:    "user-1",
				IPAddress: "127.0.0.1",
				Username:  "testuser",
				IsAdmin:   func(b bool) *bool { return &b }(false),
				RequestID: "req-123",
			},
			payrollID: "payroll-1",
			expectedPayslip: &entity.Payslip{
				ID:                 "payslip-1",
				UserID:             "user-1",
				PayrollID:          "payroll-1",
				BaseSalary:         1000.00,
				AttendanceDays:     20,
				OvertimeHours:      optional.NewDuration(5 * time.Hour),
				OvertimePay:        100.00,
				ReimbursementTotal: 150.00,
				TotalTakeHome:      1150.00,
			},
			expectedErr: nil,
			setupMock: func(m mockParams) {
				m.userRepo.EXPECT().FindUserByID(gomock.Any(), "user-1").Return(&userEntity.User{
					ID:       "user-1",
					Username: "testuser",
					IsAdmin:  false,
					Salary:   1000.00,
				}, nil)

				m.attendanceRepo.EXPECT().FindAttendancePeriodByPayrollID(gomock.Any(), "payroll-1").Return(&attEntity.AttendancePeriod{
					ID:        "period-1",
					StartDate: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC),
				}, nil)

				m.payrollRepo.EXPECT().FindPayslipByUserIDPeriod(gomock.Any(), "user-1", "payroll-1").Return(&entity.Payslip{
					ID:                 "payslip-1",
					UserID:             "user-1",
					PayrollID:          "payroll-1",
					BaseSalary:         1000.00,
					AttendanceDays:     20,
					OvertimeHours:      optional.NewDuration(5 * time.Hour),
					OvertimePay:        100.00,
					ReimbursementTotal: 150.00,
					TotalTakeHome:      1150.00,
				}, nil)

				m.reimbursementRepo.EXPECT().FindReimbursementByUserIDPeriod(gomock.Any(), "user-1", time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC)).Return([]reimbursementEntity.Reimbursement{
					{
						ID:          "reimbursement-1",
						UserID:      "user-1",
						Amount:      150.00,
						Description: optional.NewString("Travel Expenses"),
					},
					{
						ID:          "reimbursement-2",
						UserID:      "user-1",
						Amount:      50.00,
						Description: optional.NewString("Meal Expenses"),
					},
				}, nil)
			},
		},
		{
			name: "error - user not found",
			authCredential: authCredential.Credential{
				UserID:    "user-1",
				IPAddress: "127.0.0.1",
				Username:  "testuser",
				IsAdmin:   func(b bool) *bool { return &b }(false),
				RequestID: "req-123",
			},
			payrollID:       "payroll-1",
			expectedPayslip: nil,
			expectedErr: apperror.NotFound(
				apperror.AppError{
					IssueCode: entity.UserNotFound,
					Message:   entity.GetErrorMessageByIssueCode(entity.UserNotFound),
					Received:  "user-1",
				},
			),
			setupMock: func(m mockParams) {
				m.userRepo.EXPECT().FindUserByID(gomock.Any(), "user-1").Return(nil, nil)
			},
			patched: func() {},
		},
		{
			name: "error - period not found",
			authCredential: authCredential.Credential{
				UserID:    "user-1",
				IPAddress: "127.0.0.1",
				Username:  "testuser",
				IsAdmin:   func(b bool) *bool { return &b }(false),
				RequestID: "req-123",
			},
			payrollID:       "payroll-1",
			expectedPayslip: nil,
			expectedErr: apperror.NotFound(
				apperror.AppError{
					IssueCode: entity.AttendancePeriodNotFound,
					Message:   entity.GetErrorMessageByIssueCode(entity.AttendancePeriodNotFound),
					Received:  "payroll-1",
				},
			),
			setupMock: func(m mockParams) {
				m.userRepo.EXPECT().FindUserByID(gomock.Any(), "user-1").Return(&userEntity.User{
					ID:       "user-1",
					Username: "testuser",
					IsAdmin:  false,
					Salary:   1000.00,
				}, nil)

				m.attendanceRepo.EXPECT().FindAttendancePeriodByPayrollID(gomock.Any(), "payroll-1").Return(nil, nil)
			},
			patched: func() {},
		},
		{
			name: "error - payslip not found",
			authCredential: authCredential.Credential{
				UserID:    "user-1",
				IPAddress: "127.0.0.1",
				Username:  "testuser",
				IsAdmin:   func(b bool) *bool { return &b }(false),
				RequestID: "req-123",
			},
			payrollID:       "payroll-1",
			expectedPayslip: nil,
			expectedErr: apperror.NotFound(
				apperror.AppError{
					IssueCode: entity.PayslipNotFound,
					Message:   entity.GetErrorMessageByIssueCode(entity.PayslipNotFound),
					Received:  "payroll-1",
				},
			),
			setupMock: func(m mockParams) {
				m.userRepo.EXPECT().FindUserByID(gomock.Any(), "user-1").Return(&userEntity.User{
					ID:       "user-1",
					Username: "testuser",
					IsAdmin:  false,
					Salary:   1000.00,
				}, nil)

				m.attendanceRepo.EXPECT().FindAttendancePeriodByPayrollID(gomock.Any(), "payroll-1").Return(&attEntity.AttendancePeriod{
					ID:        "period-1",
					StartDate: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC),
				}, nil)

				m.payrollRepo.EXPECT().FindPayslipByUserIDPeriod(gomock.Any(), "user-1", "payroll-1").Return(nil, nil)
			},
			patched: func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockParams := mockParams{
				payrollRepo:       mockPayroll.NewMockRepository(ctrl),
				userRepo:          mockUser.NewMockRepository(ctrl),
				attendanceRepo:    mockAtt.NewMockRepository(ctrl),
				overtimeRepo:      mockOvertime.NewMockRepository(ctrl),
				reimbursementRepo: mockReimbursement.NewMockRepository(ctrl),
			}

			if tt.setupMock != nil {
				tt.setupMock(mockParams)
			}

			if tt.patched != nil {
				tt.patched()
			}

			useCase := usecase.NewPayrollUseCase(
				mockParams.payrollRepo,
				mockParams.userRepo,
				mockParams.attendanceRepo,
				mockParams.overtimeRepo,
				mockParams.reimbursementRepo,
			)
			payslip, err := useCase.ShowPayslip(context.Background(), tt.authCredential, tt.payrollID)
			if tt.expectedErr != nil {
				assert.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, payslip)
				assert.Equal(t, tt.expectedPayslip.ID, payslip.ID)
				assert.Equal(t, tt.expectedPayslip.BaseSalary, payslip.BaseSalary)
				assert.Equal(t, tt.expectedPayslip.AttendanceDays, payslip.AttendanceDays)
				assert.Equal(t, tt.expectedPayslip.ReimbursementTotal, payslip.ReimbursementTotal)
				assert.Equal(t, tt.expectedPayslip.TotalTakeHome, payslip.TotalTakeHome)
			}
		})
	}
}
