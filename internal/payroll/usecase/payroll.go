package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/vnnyx/employee-management/internal/attendance"
	attendanceEntity "github.com/vnnyx/employee-management/internal/attendance/entity"
	authCredential "github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/internal/overtime"
	overtimeEntity "github.com/vnnyx/employee-management/internal/overtime/entity"
	"github.com/vnnyx/employee-management/internal/payroll"
	"github.com/vnnyx/employee-management/internal/payroll/entity"
	"github.com/vnnyx/employee-management/internal/reimbursement"
	reimbursementEntity "github.com/vnnyx/employee-management/internal/reimbursement/entity"
	"github.com/vnnyx/employee-management/internal/users"
	userEntity "github.com/vnnyx/employee-management/internal/users/entity"
	"github.com/vnnyx/employee-management/pkg/apperror"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/iso8601"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
	"github.com/vnnyx/employee-management/pkg/optional"
)

type payrollUseCase struct {
	payrollRepo       payroll.Repository
	userRepo          users.Repository
	attendanceRepo    attendance.Repository
	overtimeRepo      overtime.Repository
	reimbursementRepo reimbursement.Repository
}

func NewPayrollUseCase(payrollRepo payroll.Repository, userRepo users.Repository, attendanceRepo attendance.Repository, overtimeRepo overtime.Repository, reimbursementRepo reimbursement.Repository) payroll.UseCase {
	return &payrollUseCase{
		payrollRepo:       payrollRepo,
		userRepo:          userRepo,
		attendanceRepo:    attendanceRepo,
		overtimeRepo:      overtimeRepo,
		reimbursementRepo: reimbursementRepo,
	}
}

func (u *payrollUseCase) GeneratePayroll(ctx context.Context, authCredential authCredential.Credential, periodID string) (entity.GeneratedPayroll, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"PayrollUseCase.GeneratePayroll()",
	)
	defer span.End()

	var generatedPayroll entity.GeneratedPayroll

	if !*authCredential.IsAdmin {
		return entity.GeneratedPayroll{}, apperror.Forbidden(
			apperror.AppError{
				IssueCode: entity.PayrollNotAuthorized,
				Message:   entity.GetErrorMessageByIssueCode(entity.PayrollNotAuthorized),
			},
		)
	}

	err := database.WithAuditContext(ctx, authCredential, pgx.TxOptions{}, func(tx database.DBTx) error {
		payrollRepoTx := u.payrollRepo.WithTx(tx)
		userRepoTx := u.userRepo.WithTx(tx)
		attendanceRepoTx := u.attendanceRepo.WithTx(tx)
		overtimeRepoTx := u.overtimeRepo.WithTx(tx)
		reimbursementRepoTx := u.reimbursementRepo.WithTx(tx)

		// Check if payroll for the period already exists
		payroll, err := payrollRepoTx.FindPayrollByPeriodID(ctx, periodID, entity.FindPayrollOptions{
			PessimisticLock: true,
		})
		if err != nil {
			return errors.Wrap(err, "PayrollUseCase.GeneratePayroll().FindPayrollByPeriodID()")
		}

		if payroll != nil {
			return apperror.BadRequest(
				apperror.AppError{
					IssueCode: entity.PayrollAlreadyGenerated,
					Message:   entity.GetErrorMessageByIssueCode(entity.PayrollAlreadyGenerated),
				},
			)
		}

		period, err := attendanceRepoTx.FindPeriodByID(ctx, periodID)
		if err != nil {
			return errors.Wrap(err, "PayrollUseCase.GeneratePayroll().FindPeriodByID()")
		}
		if period == nil {
			return apperror.NotFound(
				apperror.AppError{
					IssueCode: entity.AttendancePeriodNotFound,
					Message:   entity.GetErrorMessageByIssueCode(entity.AttendancePeriodNotFound),
					Received:  periodID,
				},
			)
		}

		users, err := userRepoTx.FindAllUsers(ctx, userEntity.FindUserOptions{
			PessimisticLock: true,
			MappedOptions: &userEntity.MappedOptions{
				MappedBy: userEntity.MappedByUserID,
			},
		})
		if err != nil {
			return errors.Wrap(err, "PayrollUseCase.GeneratePayroll().FindAllUsers()")
		}

		attendances, err := attendanceRepoTx.FindAttendanceByPeriod(ctx, period.StartDate, period.EndDate, attendanceEntity.FindAttendanceOptions{
			PessimisticLock: true,
			MappedOptions: &attendanceEntity.MappedOptions{
				MappedBy: attendanceEntity.MappedByUserID,
			},
		})
		if err != nil {
			return errors.Wrap(err, "PayrollUseCase.GeneratePayroll().FindAttendanceByPeriod()")
		}

		overtimes, err := overtimeRepoTx.FindOvertimeByPeriod(ctx, period.StartDate, period.EndDate, overtimeEntity.FindOvertimeOptions{
			PessimisticLock: true,
			MappedOptions: &overtimeEntity.MappedOptions{
				MappedBy: overtimeEntity.MappedByUserID,
			},
		})
		if err != nil {
			return errors.Wrap(err, "PayrollUseCase.GeneratePayroll().FindOvertimeByPeriod()")
		}

		reimbursements, err := reimbursementRepoTx.FindReimbursementByPeriod(ctx, period.StartDate, period.EndDate, reimbursementEntity.FindReimbursementOptions{
			PessimisticLock: true,
			MappedOptions: &reimbursementEntity.MappedOptions{
				MappedBy: reimbursementEntity.MappedByUserID,
			},
		})
		if err != nil {
			return errors.Wrap(err, "PayrollUseCase.GeneratePayroll().FindReimbursementByPeriod()")
		}

		var (
			payslips                []entity.Payslip
			totalTakeHomePayPayroll int64
		)

		payrollID := uuid.NewString()
		timeNow := time.Now()
		for _, user := range users.List {
			var totalAttendanceDays int64
			if attendances.IsMapped {
				if attendanceList, ok := attendances.Mapped[user.ID]; ok {
					totalAttendanceDays = int64(len(attendanceList))
				}
			}

			var totalOvertimeHours time.Duration
			if overtimes.IsMapped {
				if overtimeList, ok := overtimes.Mapped[user.ID]; ok {
					for _, overtime := range overtimeList {
						totalOvertimeHours += overtime.OvertimeHours
					}
				}
			}

			var totalReimbursementAmount int64
			if reimbursements.IsMapped {
				if reimbursementList, ok := reimbursements.Mapped[user.ID]; ok {
					for _, reimbursement := range reimbursementList {
						totalReimbursementAmount += reimbursement.Amount
					}
				}
			}

			attendancePay := int64(float64(user.Salary) * float64(totalAttendanceDays) / float64(period.EndDate.Sub(period.StartDate).Hours()/24))
			hourlyRate := float64(user.Salary) / (period.EndDate.Sub(period.StartDate).Hours() / 24)
			overtimePayMultiplier := 1.0 // This can be adjusted based on business rules
			totalOvertimePay := int64(totalOvertimeHours.Hours() * hourlyRate * overtimePayMultiplier)
			totalTakeHomePay := attendancePay + totalOvertimePay

			payslips = append(payslips, entity.Payslip{
				ID:                 uuid.NewString(),
				UserID:             user.ID,
				PayrollID:          payrollID,
				BaseSalary:         user.Salary,
				AttendanceDays:     totalAttendanceDays,
				OvertimeHours:      optional.NewDuration(totalOvertimeHours),
				OvertimePay:        int64(totalOvertimePay),
				ReimbursementTotal: totalReimbursementAmount,
				TotalTakeHome:      totalTakeHomePay,
				CreatedAt:          timeNow,
				UpdatedAt:          timeNow,
				CreatedBy:          authCredential.UserID,
				UpdatedBy:          authCredential.UserID,
				IPAddress:          authCredential.IPAddress,
			})

			totalTakeHomePayPayroll += totalTakeHomePay
		}

		// Store the generated payroll
		err = payrollRepoTx.StoreNewPayroll(ctx, entity.Payroll{
			ID:        payrollID,
			PeriodID:  period.ID,
			RunBy:     authCredential.UserID,
			RunAt:     timeNow,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			CreatedBy: authCredential.UserID,
			UpdatedBy: authCredential.UserID,
			IPAddress: authCredential.IPAddress,
		})
		if err != nil {
			return errors.Wrap(err, "PayrollUseCase.GeneratePayroll().StoreNewPayroll()")
		}

		// Store the payslips
		if len(payslips) > 0 {
			err := payrollRepoTx.StoreNewPayslips(ctx, payslips)
			if err != nil {
				return errors.Wrap(err, "PayrollUseCase.GeneratePayroll().StoreNewPayslips()")
			}
		}

		// Store the payroll summary
		err = payrollRepoTx.StoreNewPayrollSummary(ctx, entity.PayrollSummary{
			ID:            uuid.NewString(),
			PayrollID:     payrollID,
			TotalTakeHome: totalTakeHomePayPayroll,
			GeneratedBy:   authCredential.UserID,
			GeneratedAt:   timeNow,
			CreatedAt:     timeNow,
			UpdatedAt:     timeNow,
			CreatedBy:     authCredential.UserID,
			UpdatedBy:     authCredential.UserID,
			IPAddress:     authCredential.IPAddress,
		})
		if err != nil {
			return errors.Wrap(err, "PayrollUseCase.GeneratePayroll().StoreNewPayrollSummary()")
		}

		generatedPayroll = entity.GeneratedPayroll{
			PeriodID:      period.ID,
			PayrollID:     payrollID,
			TotalTakeHome: totalTakeHomePayPayroll,
			TotalEmployee: int64(len(users.List)),
			TotalPayslip:  int64(len(payslips)),
			GeneratedBy:   authCredential.UserID,
			GeneratedAt:   timeNow.Format(time.RFC3339),
		}

		return nil
	})
	if err != nil {
		return entity.GeneratedPayroll{}, errors.Wrap(err, "PayrollUseCase.GeneratePayroll().WithAuditContext()")
	}

	return generatedPayroll, nil
}

func (u *payrollUseCase) ShowPayslip(ctx context.Context, authCredential authCredential.Credential, payrollID string) (*entity.PayslipData, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"PayrollUseCase.ShowPayslip()",
	)
	defer span.End()

	var payslipData entity.PayslipData

	user, err := u.userRepo.FindUserByID(ctx, authCredential.UserID)
	if err != nil {
		return &payslipData, errors.Wrap(err, "PayrollUseCase.ShowPayslip().FindUserByID()")
	}
	if user == nil {
		return &payslipData, apperror.NotFound(
			apperror.AppError{
				IssueCode: entity.UserNotFound,
				Message:   entity.GetErrorMessageByIssueCode(entity.UserNotFound),
				Received:  authCredential.UserID,
			},
		)
	}

	period, err := u.attendanceRepo.FindAttendancePeriodByPayrollID(ctx, payrollID)
	if err != nil {
		return &payslipData, errors.Wrap(err, "PayrollUseCase.ShowPayslip().FindAttendancePeriodByPayrollID()")
	}
	if period == nil {
		return &payslipData, apperror.NotFound(
			apperror.AppError{
				IssueCode: entity.AttendancePeriodNotFound,
				Message:   entity.GetErrorMessageByIssueCode(entity.AttendancePeriodNotFound),
				Received:  payrollID,
			},
		)
	}

	payslip, err := u.payrollRepo.FindPayslipByUserIDPeriod(ctx, authCredential.UserID, payrollID)
	if err != nil {
		return &payslipData, errors.Wrap(err, "PayrollUseCase.ShowPayslip().FindPayslipByUserIDPeriod()")
	}
	if payslip == nil {
		return &payslipData, apperror.NotFound(
			apperror.AppError{
				IssueCode: entity.PayslipNotFound,
				Message:   entity.GetErrorMessageByIssueCode(entity.PayslipNotFound),
				Received: map[string]string{
					"user_id":    authCredential.UserID,
					"payroll_id": payrollID,
				},
			},
		)
	}

	reimbursements, err := u.reimbursementRepo.FindReimbursementByUserIDPeriod(ctx, authCredential.UserID, period.StartDate, period.EndDate)
	if err != nil {
		return &payslipData, errors.Wrap(err, "PayrollUseCase.ShowPayslip().FindReimbursementByUserIDPeriod()")
	}

	reimbursementData := make([]entity.ReimbursementData, 0, len(reimbursements))
	for _, reimbursement := range reimbursements {
		reimbursementData = append(reimbursementData, entity.ReimbursementData{
			Description:       reimbursement.Description,
			Amount:            reimbursement.Amount,
			ReimbursementDate: reimbursement.ReimbursementDate.Format(time.RFC3339),
		})
	}

	return &entity.PayslipData{
		ID: payslip.ID,
		User: entity.UserData{
			ID:       user.ID,
			Username: user.Username,
		},
		AttendancePeriod: entity.AttendancePeriodData{
			StartDate: period.StartDate.Format(time.RFC3339),
			EndDate:   period.EndDate.Format(time.RFC3339),
		},
		BaseSalary:     user.Salary,
		WorkingDays:    int64(period.EndDate.Sub(period.StartDate).Hours() / 24),
		AttendanceDays: payslip.AttendanceDays,
		AttendancePay:  int64(float64(user.Salary) * float64(payslip.AttendanceDays) / float64(period.EndDate.Sub(period.StartDate).Hours()/24)),
		Overtime: entity.OvertimeData{
			OvertimeHours: iso8601.ToString(payslip.OvertimeHours.MustGet()),
			RatePerHour:   int64(float64(user.Salary) / (period.EndDate.Sub(period.StartDate).Hours() / 24)),
			Multiplier:    1,
			OvertimePay:   payslip.OvertimePay,
		},
		Reimbursements:     reimbursementData,
		ReimbursementTotal: payslip.ReimbursementTotal,
		TotalTakeHome:      payslip.TotalTakeHome,
	}, nil
}

func (u *payrollUseCase) ListPayslips(ctx context.Context, authCredential authCredential.Credential, payrollID string) (entity.ListPayslips, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"PayrollUseCase.ListPayslips()",
	)
	defer span.End()

	var payslipsData entity.ListPayslips

	if !*authCredential.IsAdmin {
		return payslipsData, apperror.Forbidden(
			apperror.AppError{
				IssueCode: entity.PayrollNotAuthorized,
				Message:   entity.GetErrorMessageByIssueCode(entity.PayrollNotAuthorized),
			},
		)
	}

	users, err := u.userRepo.FindAllUsers(ctx, userEntity.FindUserOptions{
		PessimisticLock: true,
	})
	if err != nil {
		return payslipsData, errors.Wrap(err, "PayrollUseCase.IndexPayslips().FindAllUsers()")
	}

	payslips, err := u.payrollRepo.FindPayslipByPayrollID(ctx, payrollID, entity.FindPayslipOptions{
		PessimisticLock: true,
		MappedOptions: &entity.MappedOptions{
			MappedBy: entity.MappedByUserID,
		},
	})
	if err != nil {
		return payslipsData, errors.Wrap(err, "PayrollUseCase.IndexPayslips().FindPayslipByPayrollID()")
	}

	period, err := u.attendanceRepo.FindAttendancePeriodByPayrollID(ctx, payrollID)
	if err != nil {
		return payslipsData, errors.Wrap(err, "PayrollUseCase.IndexPayslips().FindAttendancePeriodByPayrollID()")
	}
	if period == nil {
		return payslipsData, apperror.NotFound(
			apperror.AppError{
				IssueCode: entity.AttendancePeriodNotFound,
				Message:   entity.GetErrorMessageByIssueCode(entity.AttendancePeriodNotFound),
				Received:  payrollID,
			},
		)
	}

	reimbursements, err := u.reimbursementRepo.FindReimbursementByPeriod(ctx, period.StartDate, period.EndDate, reimbursementEntity.FindReimbursementOptions{
		PessimisticLock: true,
		MappedOptions: &reimbursementEntity.MappedOptions{
			MappedBy: reimbursementEntity.MappedByUserID,
		},
	})
	if err != nil {
		return payslipsData, errors.Wrap(err, "PayrollUseCase.IndexPayslips().FindReimbursementByPeriod()")
	}

	var (
		payslipDataList  []entity.PayslipData
		totalTakeHomePay int64
	)
	for _, user := range users.List {
		payslips, found := payslips.Mapped[user.ID]
		if !found {
			continue // Skip users without payslips
		}

		// Each period should have only one payslip per user
		payslip := payslips[0]

		reimbursementData := make([]entity.ReimbursementData, 0)
		if reimbursements.IsMapped {
			if reimbursementList, ok := reimbursements.Mapped[user.ID]; ok {
				for _, reimbursement := range reimbursementList {
					reimbursementData = append(reimbursementData, entity.ReimbursementData{
						Description:       reimbursement.Description,
						Amount:            reimbursement.Amount,
						ReimbursementDate: reimbursement.ReimbursementDate.Format(time.RFC3339),
					})
				}
			}
		}

		payslipData := entity.PayslipData{
			ID: payslip.ID,
			User: entity.UserData{
				ID:       user.ID,
				Username: user.Username,
			},
			AttendancePeriod: entity.AttendancePeriodData{
				StartDate: period.StartDate.Format(time.RFC3339),
				EndDate:   period.EndDate.Format(time.RFC3339),
			},
			BaseSalary:     user.Salary,
			WorkingDays:    int64(period.EndDate.Sub(period.StartDate).Hours() / 24),
			AttendanceDays: payslip.AttendanceDays,
			AttendancePay:  int64(float64(user.Salary) * float64(payslip.AttendanceDays) / float64(period.EndDate.Sub(period.StartDate).Hours()/24)),
			Overtime: entity.OvertimeData{
				OvertimeHours: iso8601.ToString(payslip.OvertimeHours.MustGet()),
				RatePerHour:   int64(float64(user.Salary) / (period.EndDate.Sub(period.StartDate).Hours() / 24)),
				Multiplier:    1,
				OvertimePay:   payslip.OvertimePay,
			},
			Reimbursements:     reimbursementData,
			ReimbursementTotal: payslip.ReimbursementTotal,
			TotalTakeHome:      payslip.TotalTakeHome,
		}

		payslipDataList = append(payslipDataList, payslipData)
		totalTakeHomePay += payslip.TotalTakeHome
	}

	payslipsData.TotalTakeHome = totalTakeHomePay
	payslipsData.PayslipsData = payslipDataList

	return payslipsData, nil
}
